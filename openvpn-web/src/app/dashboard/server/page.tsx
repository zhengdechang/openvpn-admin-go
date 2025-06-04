"use client";

import React, { useEffect, useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/lib/api";
import { ServerStatus } from "@/lib/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";

export default function ServerPage() {
  const [status, setStatus] = useState<ServerStatus | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const data = await openvpnAPI.getServerStatus();
        setStatus(data);
      } catch (error) {
        toast.error("获取服务器状态失败");
      } finally {
        setLoading(false);
      }
    };
    fetchStatus();
  }, []);

  return (
    <MainLayout className="p-4">
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
    </MainLayout>
  );
}