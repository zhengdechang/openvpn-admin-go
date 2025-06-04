"use client";
import React, { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI, departmentAPI } from "@/services/api";
import type { OpenVPNClient, Department } from "@/types/types";
import { UserRole } from "@/types/types";
import { Dialog, DialogTrigger, DialogContent, DialogHeader, DialogFooter, DialogTitle, DialogClose } from "@/components/ui/dialog";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

export default function ClientsPage() {
  const { user: currentUser } = useAuth();
  const [clients, setClients] = useState<OpenVPNClient[]>([]);
  const [depts, setDepts] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [form, setForm] = useState({ name: "", email: "", departmentId: "" });

  const fetchData = async () => {
    setLoading(true);
    try {
      const [cList, dList] = await Promise.all([
        openvpnAPI.getClientList(),
        departmentAPI.list(),
      ]);
      setClients(cList);
      setDepts(dList);
    } catch {
      toast.error("加载客户端或部门列表失败");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchData(); }, []);
  // 根据角色过滤客户端
  let visibleClients = clients;
  if (currentUser?.role === UserRole.USER) {
    visibleClients = clients.filter(c => c.email === currentUser.email);
  } else if (currentUser?.role === UserRole.MANAGER) {
    visibleClients = clients.filter(c => c.departmentId === currentUser.departmentId);
  }

  const handleDownload = async (id: string, os: string) => {
    try {
      const { config } = await openvpnAPI.getClientConfig(id, os);
      const blob = new Blob([config], { type: 'application/x-openvpn-profile' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      const ext = os === 'linux' ? 'conf' : 'ovpn';
      a.download = `${id}.${ext}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      toast.success('下载成功');
    } catch {
      toast.error('下载失败');
    }
  };

  const handleCreate = async () => {
    if (!form.name || !form.email || !form.departmentId) {
      toast.error("请填写必填项");
      return;
    }
    try {
      await openvpnAPI.addClient(form.name, form.departmentId);
      toast.success("创建成功");
      setForm({ name: "", email: "", departmentId: "" });
      setDialogOpen(false);
      fetchData();
    } catch {
      toast.error("创建失败");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm(`确定删除客户端 ${id} ?`)) return;
    try {
      await openvpnAPI.deleteClient(id);
      toast.success("删除成功");
      fetchData();
    } catch {
      toast.error("删除失败");
    }
  };

  return (
    <MainLayout className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">客户端管理</h2>
        {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.MANAGER) && (
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button>新增客户端</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>新增客户端</DialogTitle>
            </DialogHeader>
            <div className="space-y-2 pt-2">
              <Input
                placeholder="客户端名称"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
              />
              <Input
                placeholder="邮箱"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
              />
              <select
                className="border px-2 w-full py-1"
                value={form.departmentId}
                onChange={(e) => setForm({ ...form, departmentId: e.target.value })}
              >
                <option value="">-- 选择部门 --</option>
                {depts.map((d) => (
                  <option key={d.id} value={d.id}>{d.name}</option>
                ))}
              </select>
            </div>
            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline">取消</Button>
              </DialogClose>
              <Button onClick={handleCreate}>创建</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
        )}
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
                  <TableHead>名称</TableHead>
                  <TableHead>邮箱</TableHead>
                  <TableHead>部门</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {visibleClients.map((c) => (
                  <TableRow key={c.id}>
                    <TableCell>{c.name}</TableCell>
                    <TableCell>{c.email}</TableCell>
                    <TableCell>{depts.find((d) => d.id === c.departmentId)?.name || '-'}</TableCell>
                    <TableCell className="space-x-2">
                      <select
                        className="border px-2 py-1"
                        defaultValue=""
                        onChange={(e) => handleDownload(c.id, e.target.value)}
                      >
                        <option value="" disabled>下载配置</option>
                        <option value="windows">Windows</option>
                        <option value="macos">macOS</option>
                        <option value="linux">Linux</option>
                      </select>
                      {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.MANAGER) && (
                        <Button size="sm" variant="destructive" onClick={() => handleDelete(c.id)}>删除</Button>
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