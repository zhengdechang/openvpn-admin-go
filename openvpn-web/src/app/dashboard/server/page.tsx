"use client";

"use client";
import React, { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import MainLayout from "@/components/layout/main-layout";
import { serverAPI } from "@/services/api";
import type { ServerStatus } from "@/types/types";
import { UserRole } from "@/types/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function ServerPage() {
  const { user: currentUser } = useAuth();
  const [status, setStatus] = useState<ServerStatus | null>(null);
  const [config, setConfig] = useState<string>("");
  const [loading, setLoading] = useState(true);

  const fetchStatus = async () => {
    setLoading(true);
    try {
      const data = await serverAPI.getStatus();
      setStatus(data);
      const tpl = await serverAPI.getConfigTemplate();
      setConfig(tpl.template);
    } catch (error) {
      toast.error("获取服务器信息失败");
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => { fetchStatus(); }, []);

  if (!currentUser || currentUser.role !== UserRole.SUPERADMIN) {
    return (
      <MainLayout className="p-4">
        <p className="text-center mt-10">无权限访问此页面</p>
      </MainLayout>
    );
  }
  return (
    <MainLayout className="p-4 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>服务器状态</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>加载中...</p>
          ) : status ? (
            <div className="space-y-2">
              <p>名称: {status.name}</p>
              <p>状态: {status.status}</p>
              <p>运行时长: {status.uptime}</p>
              <p>当前连接: {status.connected}</p>
              <p>历史总数: {status.total}</p>
              <p>最后更新时间: {new Date(status.lastUpdated).toLocaleString()}</p>
            </div>
          ) : (
            <p>暂无数据</p>
          )}
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>服务器配置管理</CardTitle>
        </CardHeader>
        <CardContent>
          <textarea
            className="w-full h-64 border p-2"
            value={config}
            onChange={(e) => setConfig(e.target.value)}
          />
          <div className="mt-4 space-x-2">
            <Button onClick={async () => {
              try {
                await serverAPI.updateConfig(config);
                toast.success('配置更新成功');
              } catch {
                toast.error('配置更新失败');
              }
            }}>
              保存配置
            </Button>
            <Button onClick={async () => {
              try { await serverAPI.start(); toast.success('启动成功'); fetchStatus(); } catch { toast.error('启动失败'); }
            }}>
              启动
            </Button>
            <Button onClick={async () => {
              try { await serverAPI.stop(); toast.success('停止成功'); fetchStatus(); } catch { toast.error('停止失败'); }
            }}>
              停止
            </Button>
            <Button onClick={async () => {
              try { await serverAPI.restart(); toast.success('重启成功'); fetchStatus(); } catch { toast.error('重启失败'); }
            }}>
              重启
            </Button>
          </div>
        </CardContent>
      </Card>
    </MainLayout>
  );
}
