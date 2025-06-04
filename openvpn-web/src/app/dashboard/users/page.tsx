"use client";

import React, { useState, useEffect } from "react";
import { Dialog, DialogTrigger, DialogContent, DialogHeader, DialogFooter, DialogTitle, DialogClose } from "@/components/ui/dialog";
import MainLayout from "@/components/layout/main-layout";
import { userManagementAPI, departmentAPI } from "@/services/api";
import type { AdminUser, Department } from "@/types/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

export default function UsersPage() {
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [depts, setDepts] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState({ name: "", email: "", password: "", role: "user", departmentId: "" });
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

  const handleCreate = async () => {
    if (!form.name || !form.email || !form.password) {
      toast.error("请填写必填项");
      return;
    }
    try {
      await userManagementAPI.create(form);
      toast.success("创建成功");
      setForm({ name: "", email: "", password: "", role: "user", departmentId: "" });
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
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button>新增用户</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>新增用户</n              </DialogTitle>
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
                  onChange={(e) => setForm({ ...form, role: e.target.value })}
                >
                  <option value="user">User</option>
                  <option value="manager">Manager</option>
                  <option value="admin">Admin</option>
                  <option value="superadmin">SuperAdmin</option>
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
                {users.map((u) => (
                  <TableRow key={u.id}>
                    <TableCell>{u.name}</TableCell>
                    <TableCell>{u.email}</TableCell>
                    <TableCell>{u.role}</TableCell>
                    <TableCell>{depts.find((d) => d.id === u.departmentId)?.name || '-'}</TableCell>
                    <TableCell>
                      <Button size="sm" variant="destructive" onClick={() => handleDelete(u.id)}>删除</Button>
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