"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/services/api";
import type { OpenVPNClient } from "@/types/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function ClientsPage() {
  const [clients, setClients] = useState<OpenVPNClient[]>([]);
  const { user } = useAuth();
  const [loading, setLoading] = useState<boolean>(true);

  const fetchClients = async () => {
    try {
      const data = await openvpnAPI.getClientList();
      setClients(data);
    } catch {
      toast.error("加载客户端列表失败");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchClients(); }, []);
  
  // 根据角色过滤客户端: 普通用户仅能查看自己的客户端
  const visibleClients = user?.role === 'user'
    ? clients.filter(c => c.username === user.id)
    : clients;

  const handleAdd = async () => {
    const username = prompt("请输入客户端用户名");
    if (!username) return;
    try {
      await openvpnAPI.addClient(username);
      toast.success("添加成功");
      fetchClients();
    } catch {
      toast.error("添加失败");
    }
  };

  const handleDelete = async (username: string) => {
    if (!confirm(`确定删除客户端 ${username} ?`)) return;
    try {
      await openvpnAPI.deleteClient(username);
      toast.success("删除成功");
      fetchClients();
    } catch {
      toast.error("删除失败");
    }
  };

  const handleDownload = async (username: string, os: string) => {
    try {
      const { config } = await openvpnAPI.getClientConfig(username, os);
      const blob = new Blob([config], { type: "application/x-openvpn-profile" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      // 根据操作系统选择文件扩展名
      const ext = os === 'linux' ? 'conf' : 'ovpn';
      a.download = `${username}.${ext}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      toast.success("下载成功");
    } catch {
      toast.error("下载失败");
    }
  };

  return (
    <MainLayout className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">客户端管理</h2>
        {user?.role !== 'user' && <Button onClick={handleAdd}>添加客户端</Button>}
      </div>
      <Card>
        <CardHeader><CardTitle>客户端列表</CardTitle></CardHeader>
        <CardContent>
          {loading ? (
            <p>加载中...</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>用户名</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {visibleClients.map((c) => (
                  <TableRow key={c.username}>
                    <TableCell>{c.username}</TableCell>
                    <TableCell className="space-x-2">
                      <select
                        className="border px-2 py-1"
                        defaultValue=""
                        onChange={(e) => handleDownload(c.username, e.target.value)}
                      >
                        <option value="" disabled>下载配置</option>
                        <option value="windows">Windows</option>
                        <option value="macos">macOS</option>
                        <option value="linux">Linux</option>
                      </select>
                      {user?.role !== 'user' && (
                        <Button size="sm" variant="destructive" onClick={() => handleDelete(c.username)}>删除</Button>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}