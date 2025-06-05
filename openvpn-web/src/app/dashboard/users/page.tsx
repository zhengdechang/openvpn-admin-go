"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { Dialog, DialogTrigger, DialogContent, DialogHeader, DialogFooter, DialogTitle, DialogClose } from "@/components/ui/dialog";
import MainLayout from "@/components/layout/main-layout";
import { userManagementAPI, departmentAPI, openvpnAPI } from "@/services/api";
import { AdminUser, Department, UserRole } from "@/types/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

export default function UsersPage() {
  const { user: currentUser } = useAuth();
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [depts, setDepts] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState({ name: "", email: "", password: "", role: UserRole.USER, departmentId: "" });
  const [open, setOpen] = useState(false);

  const fetchAll = async () => {
    setLoading(true);
    try {
      const [u, d] = await Promise.all([userManagementAPI.list(), departmentAPI.list()]);
      setUsers(u);
      setDepts(d);
    } catch {
      toast.error("加载用户或部门失败");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchAll(); }, []);
  // 根据角色过滤用户
  let visibleUsers = users;
  // 下载配置
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
    if (!form.name || !form.email || !form.password) {
      toast.error("请填写必填项");
      return;
    }
    try {
      await userManagementAPI.create(form);
      toast.success("创建成功");
      setForm({ name: "", email: "", password: "", role: UserRole.USER, departmentId: "" });
      fetchAll();
    } catch {
      toast.error("创建失败");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("确定删除此用户?")) return;
    try {
      await userManagementAPI.delete(id);
      toast.success("删除成功");
      fetchAll();
    } catch {
      toast.error("删除失败");
    }
  };

  return (
    <MainLayout className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">用户管理</h1>
        {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.MANAGER || currentUser?.role === UserRole.SUPERADMIN) && (
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button>新增用户</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>新增用户</DialogTitle>
            </DialogHeader>
            <div className="space-y-2 pt-2">
              <Input
                placeholder="姓名"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
              />
              <Input
                placeholder="邮箱"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
              />
              <Input
                type="password"
                placeholder="密码"
                value={form.password}
                onChange={(e) => setForm({ ...form, password: e.target.value })}
              />
              <div className="flex space-x-2">
                <select
                  className="border px-2"
                  value={form.role}
                  onChange={(e) => setForm({ ...form, role: e.target.value as UserRole })}
                >
                  <option value={UserRole.USER}>User</option>
                  {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (
                    <> 
                      <option value={UserRole.MANAGER}>Manager</option>
                      <option value={UserRole.ADMIN}>Admin</option>
                      <option value={UserRole.SUPERADMIN}>Superadmin</option>
                    </>
                  )}
                </select>
                <select
                  className="border px-2"
                  value={form.departmentId}
                  onChange={(e) => setForm({ ...form, departmentId: e.target.value })}
                >
                  <option value="">-- 选择部门 --</option>
                  {depts.map((d) => (
                    <option key={d.id} value={d.id}>{d.name}</option>
                  ))}
                </select>
              </div>
            </div>
            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline">取消</Button>
              </DialogClose>
              <Button onClick={() => { handleCreate(); setOpen(false); }}>创建</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
        )}
      </div>

      <Card>
        <CardHeader><CardTitle>用户列表</CardTitle></CardHeader>
        <CardContent>
          {loading ? (
            <p>加载中...</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>姓名</TableHead>
                  <TableHead>邮箱</TableHead>
                  <TableHead>角色</TableHead>
                  <TableHead>部门</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {visibleUsers.map((u) => (
                  <TableRow key={u.id}>
                    <TableCell>{u.name}</TableCell>
                    <TableCell>{u.email}</TableCell>
                    <TableCell>{u.role}</TableCell>
                    <TableCell>{depts.find((d) => d.id === u.departmentId)?.name || '-'}</TableCell>
                    <TableCell className="space-x-2">
                      <select
                        className="border px-2 py-1"
                        defaultValue=""
                        onChange={(e) => handleDownload(u.id, e.target.value)}
                      >
                        <option value="" disabled>下载配置</option>
                        <option value="windows">Windows</option>
                        <option value="macos">macOS</option>
                        <option value="linux">Linux</option>
                      </select>
                      {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.MANAGER || currentUser?.role === UserRole.SUPERADMIN) && (
                        <Button size="sm" variant="destructive" onClick={() => handleDelete(u.id)}>删除</Button>
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