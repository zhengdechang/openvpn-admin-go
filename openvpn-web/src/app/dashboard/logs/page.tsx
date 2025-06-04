"use client";

import React, { useEffect, useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

export default function LogsPage() {
  const [serverLogs, setServerLogs] = useState<string>("");
  const [clientUsername, setClientUsername] = useState<string>("");
  const [clientLogs, setClientLogs] = useState<string[] | null>(null);
  const [loadingServer, setLoadingServer] = useState(true);

  useEffect(() => {
    const fetchServerLogs = async () => {
      try {
        const logs = await openvpnAPI.getServerLogs();
        setServerLogs(logs);
      } catch (error) {
        toast.error("获取服务器日志失败");
      } finally {
        setLoadingServer(false);
      }
    };
    fetchServerLogs();
  }, []);

  const handleFetchClientLogs = async () => {
    if (!clientUsername) {
      toast.error("请输入用户名");
      return;
    }
    try {
      const logs = await openvpnAPI.getClientLogs(clientUsername);
      setClientLogs(logs);
    } catch (error) {
      toast.error("获取客户端日志失败");
    }
  };

  return (
    <MainLayout className="p-4 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>服务器日志</CardTitle>
        </CardHeader>
        <CardContent>
          {loadingServer ? (
            <p>加载中...</p>
          ) : (
            <pre className="whitespace-pre-wrap">{serverLogs || "暂无日志"}</pre>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>客户端日志查询</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center space-x-2">
            <Input
              placeholder="用户名"
              value={clientUsername}
              onChange={(e) => setClientUsername(e.target.value)}
            />
            <Button onClick={handleFetchClientLogs}>查询</Button>
          </div>
          {clientLogs && (
            <pre className="whitespace-pre-wrap">{clientLogs.join("\n") || "暂无日志"}</pre>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}