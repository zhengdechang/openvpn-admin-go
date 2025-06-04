"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { UserRole } from "@/types/types";
import MainLayout from "@/components/layout/main-layout";
import { departmentAPI, userManagementAPI } from "@/services/api";
import type { Department, AdminUser } from "@/types/types";
import { Dialog, DialogTrigger, DialogContent, DialogHeader, DialogFooter, DialogTitle, DialogClose } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

type DepartmentTree = Department & { children: DepartmentTree[] };

export default function DepartmentsPage() {
  const { user: currentUser } = useAuth();
  const [depts, setDepts] = useState<Department[]>([]);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  // 控制树形展开的部门ID集合
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());
  const [form, setForm] = useState({ name: "", headId: "", parentId: "" });
  // 编辑模式状态
  const [editOpen, setEditOpen] = useState(false);
  const [editingDept, setEditingDept] = useState<Department | null>(null);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [dList, uList] = await Promise.all([departmentAPI.list(), userManagementAPI.list()]);
      setDepts(dList);
      setUsers(uList);
    } catch {
      toast.error("加载部门或用户失败");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchData(); }, []);
  // 构建部门树
  const buildTree = (list: Department[]): DepartmentTree[] => {
    const map = new Map<string, DepartmentTree>();
    list.forEach(item => {
      map.set(item.id, { ...item, children: [] });
    });
    const roots: DepartmentTree[] = [];
    map.forEach(item => {
      if (item.parentId) {
        const parent = map.get(item.parentId);
        parent?.children.push(item);
      } else {
        roots.push(item);
      }
    });
    return roots;
  };
  const tree: DepartmentTree[] = buildTree(depts);


  const handleCreate = async () => {
    if (!form.name) {
      toast.error("请填写部门名称");
      return;
    }
    try {
      await departmentAPI.create(form);
      toast.success("创建成功");
      setForm({ name: "", headId: "", parentId: "" });
      fetchData();
      setOpen(false);
    } catch {
      toast.error("创建失败");
    }
  };
  // 编辑部门
  const handleEdit = async () => {
    if (!editingDept) return;
    if (!form.name) {
      toast.error("请填写部门名称");
      return;
    }
    try {
      await departmentAPI.update(editingDept.id, form);
      toast.success("更新成功");
      setForm({ name: "", headId: "", parentId: "" });
      setEditingDept(null);
      setEditOpen(false);
      fetchData();
    } catch {
      toast.error("更新失败");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("确定删除此部门？")) return;
    try {
      await departmentAPI.delete(id);
      toast.success("删除成功");
      fetchData();
    } catch {
      toast.error("删除失败");
    }
  };

  // 切换部门展开/收起
  const toggleExpand = (id: string) => {
    setExpandedIds(prev => {
      const next = new Set<string>(); // Always start with a new, empty Set
      if (!prev.has(id)) { // If the clicked item was not already expanded (i.e., it's currently collapsed)
        next.add(id); // Add it to the new Set (it will be the only expanded one)
      }
      // If the clicked item was already expanded (prev.has(id) is true),
      // 'next' will remain empty, effectively collapsing it.
      // This ensures only one item is expanded at a time, or all are collapsed.
      return next;
    });
  };
  // 递归渲染树形列表，支持多节点展开
  const renderRows = (nodes: DepartmentTree[], level: number = 0): React.ReactNode[] => {
    return nodes.flatMap(node => {
      const hasChildren = node.children && node.children.length > 0;
      const isExpanded = expandedIds.has(node.id);
      return [
        <TableRow key={node.id}>
          <TableCell style={{ paddingLeft: level * 20, display: 'flex', alignItems: 'center' }}>
            {hasChildren && (
              <span
                className="cursor-pointer select-none mr-1"
                onClick={() => toggleExpand(node.id)}
              >
                {isExpanded ? '▼' : '▶'}
              </span>
            )}
            {node.name}
          </TableCell>
          <TableCell>{node.head?.name || '-'}</TableCell>
        <TableCell className="space-x-2">
          {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (
            <>
              <Button size="sm" variant="outline" onClick={() => {
                setEditingDept(node);
                setForm({ name: node.name, headId: node.headId || "", parentId: node.parentId || "" });
                setEditOpen(true);
              }}>编辑</Button>
              <Button size="sm" variant="destructive" onClick={() => handleDelete(node.id)}>删除</Button>
            </>
          )}
        </TableCell>
        </TableRow>,
        // 如果展开，则渲染子节点
        ...(hasChildren && isExpanded ? renderRows(node.children, level + 1) : [])
      ];
    });
  };
  return (
    <MainLayout className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">部门管理</h2>
        {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (
          <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
              <Button>新增部门</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>新增部门</DialogTitle>
              </DialogHeader>
              <div className="space-y-2 pt-2">
                <Input
                  placeholder="部门名称"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                />
                <select
                  className="border px-2 w-full py-1"
                  value={form.parentId}
                  onChange={(e) => setForm({ ...form, parentId: e.target.value })}
                >
                  <option value="">-- 选择上级部门 --</option>
                  {depts.map((d) => (
                    <option key={d.id} value={d.id}>{d.name}</option>
                  ))}
                </select>
                <select
                  className="border px-2 w-full py-1"
                  value={form.headId}
                  onChange={(e) => setForm({ ...form, headId: e.target.value })}
                >
                  <option value="">-- 选择负责人 --</option>
                  {users.map((u) => (
                    <option key={u.id} value={u.id}>{u.name}</option>
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
        <CardHeader>
          <CardTitle>部门列表</CardTitle>
        </CardHeader>
      <CardContent>
        {loading ? (
          <p>加载中...</p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>名称</TableHead>
                <TableHead>负责人</TableHead>
                <TableHead>操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {renderRows(tree)}
            </TableBody>
          </Table>
        )}
      </CardContent>
      </Card>
    </MainLayout>
  );
}