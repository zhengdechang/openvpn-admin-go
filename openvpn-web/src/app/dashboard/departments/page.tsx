"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { UserRole } from "@/types/types";
import { useTranslation } from "react-i18next";
import MainLayout from "@/components/layout/main-layout";
import { departmentAPI, userManagementAPI } from "@/services/api";
import type { Department, AdminUser } from "@/types/types";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogClose,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

type DepartmentTree = Department & { children: DepartmentTree[] };

export default function DepartmentsPage() {
  const { user: currentUser } = useAuth();
  const { t } = useTranslation();
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
      const [dList, uList] = await Promise.all([
        departmentAPI.list(),
        userManagementAPI.list(),
      ]);
      const userMap = new Map(uList.map((u) => [u.id, u]));
      const deptsWithHead = dList.map((dept) => ({
        ...dept,
        head: dept.headId ? userMap.get(dept.headId) : undefined,
      }));
      console.log("dList:", dList);
      console.log("deptsWithHead:", deptsWithHead);
      setDepts(deptsWithHead);
      setUsers(uList);
      console.log("depts after setDepts:", deptsWithHead);
      console.log("tree after setDepts:", buildTree(deptsWithHead));
    } catch {
      toast.error(t("dashboard.departments.loadError"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);
  // 构建部门树
  const buildTree = (list: Department[]): DepartmentTree[] => {
    console.log("buildTree list:", list);
    const map = new Map<string, DepartmentTree>();
    list.forEach((item) => {
      map.set(item.id, { ...item, children: [] });
    });
    console.log("map:", map);
    const roots: DepartmentTree[] = [];
    map.forEach((item) => {
      if (
        item.parentId &&
        item.parentId !== item.id &&
        map.has(item.parentId)
      ) {
        const parent = map.get(item.parentId);
        parent?.children.push(item);
      } else {
        roots.push(item);
      }
    });
    console.log("roots:", roots);
    return roots;
  };
  const tree: DepartmentTree[] = buildTree(depts);

  const handleCreate = async () => {
    if (!form.name) {
      toast.error(t("dashboard.departments.nameRequired"));
      return;
    }
    try {
      await departmentAPI.create(form);
      toast.success(t("dashboard.departments.createSuccess"));
      setForm({ name: "", headId: "", parentId: "" });
      fetchData();
      setOpen(false);
    } catch {
      toast.error(t("dashboard.departments.createError"));
    }
  };
  // 编辑部门
  const handleEdit = async () => {
    if (!editingDept) return;
    if (!form.name) {
      toast.error(t("dashboard.departments.nameRequired"));
      return;
    }
    try {
      await departmentAPI.update(editingDept.id, form);
      toast.success(t("dashboard.departments.updateSuccess"));
      setForm({ name: "", headId: "", parentId: "" });
      setEditingDept(null);
      setEditOpen(false);
      fetchData();
    } catch {
      toast.error(t("dashboard.departments.updateError"));
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm(t("dashboard.departments.deleteConfirm"))) return;
    try {
      await departmentAPI.delete(id);
      toast.success(t("dashboard.departments.deleteSuccess"));
      fetchData();
    } catch {
      toast.error(t("dashboard.departments.deleteError"));
    }
  };

  // 切换部门展开/收起
  const toggleExpand = (id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev); // Create a new Set based on the previous state
      if (next.has(id)) {
        next.delete(id); // If the id is already in the Set, remove it (collapse)
      } else {
        next.add(id); // Otherwise, add it to the Set (expand)
      }
      return next;
    });
  };
  // 递归渲染树形列表，支持多节点展开
  const renderRows = (
    nodes: DepartmentTree[],
    level: number = 0
  ): React.ReactNode[] => {
    console.log("renderRows nodes:", nodes);
    return nodes.flatMap((node) => {
      const hasChildren = node.children && node.children.length > 0;
      const isExpanded = expandedIds.has(node.id);
      return [
        <TableRow key={node.id}>
          <TableCell
            style={{
              paddingLeft: level === 0 ? undefined : level * 20,
              position:"relative"
            }}
          >
            {hasChildren && (
              <span
                className="cursor-pointer select-none mr-1 flex items-center "
                onClick={() => toggleExpand(node.id)}
                style={{
                  width: 20,
                  position:"absolute"
                }}
              >
                {isExpanded ? (
                  <svg
                    width="16"
                    height="16"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                ) : (
                  <svg
                    width="16"
                    height="16"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                )}
              </span>
            )}
            <span style={{ flex: 1 }}>{node.name}</span>
          </TableCell>
          <TableCell>
            {node.head?.name || t("dashboard.departments.emptyData")}
          </TableCell>
          <TableCell className="space-x-2">
            {(currentUser?.role === UserRole.ADMIN ||
              currentUser?.role === UserRole.SUPERADMIN) && (
              <>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => {
                    setEditingDept(node);
                    setForm({
                      name: node.name,
                      headId: node.headId || "",
                      parentId: node.parentId || "",
                    });
                    setEditOpen(true);
                  }}
                >
                  {t("dashboard.departments.edit")}
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={() => handleDelete(node.id)}
                >
                  {t("dashboard.departments.delete")}
                </Button>
              </>
            )}
          </TableCell>
        </TableRow>,
        // 如果展开，则渲染子节点
        ...(hasChildren && isExpanded
          ? renderRows(node.children, level + 1)
          : []),
      ];
    });
  };
  return (
    <MainLayout className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">
          {t("dashboard.departments.pageTitle")}
        </h2>
        {(currentUser?.role === UserRole.ADMIN ||
          currentUser?.role === UserRole.SUPERADMIN) && (
          <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
              <Button>{t("dashboard.departments.addDepartment")}</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>
                  {t("dashboard.departments.addDepartmentDialogTitle")}
                </DialogTitle>
              </DialogHeader>
              <div className="space-y-2 pt-2">
                <Input
                  placeholder={t(
                    "dashboard.departments.departmentNamePlaceholder"
                  )}
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                />
                <select
                  className="border px-2 w-full py-1"
                  value={form.parentId}
                  onChange={(e) =>
                    setForm({ ...form, parentId: e.target.value })
                  }
                >
                  <option value="">
                    {t("dashboard.departments.selectParentDepartment")}
                  </option>
                  {depts.map((d) => (
                    <option key={d.id} value={d.id}>
                      {d.name}
                    </option>
                  ))}
                </select>
                <select
                  className="border px-2 w-full py-1"
                  value={form.headId}
                  onChange={(e) => setForm({ ...form, headId: e.target.value })}
                >
                  <option value="">
                    {t("dashboard.departments.selectHead")}
                  </option>
                  {users.map((u) => (
                    <option key={u.id} value={u.id}>
                      {u.name}
                    </option>
                  ))}
                </select>
              </div>
              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline">
                    {t("dashboard.departments.cancel")}
                  </Button>
                </DialogClose>
                <Button onClick={handleCreate}>
                  {t("dashboard.departments.create")}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        )}
        {/* Edit Department Dialog */}
        {(currentUser?.role === UserRole.ADMIN ||
          currentUser?.role === UserRole.SUPERADMIN) &&
          editingDept && ( // Ensure editingDept is not null to render
            <Dialog
              open={editOpen}
              onOpenChange={(isOpen) => {
                setEditOpen(isOpen);
                if (!isOpen) {
                  setEditingDept(null);
                  setForm({ name: "", headId: "", parentId: "" }); // Reset form on close
                }
              }}
            >
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>
                    {t("dashboard.departments.editDepartmentDialogTitle")}
                  </DialogTitle>
                </DialogHeader>
                <div className="space-y-2 pt-2">
                  <Input
                    placeholder={t(
                      "dashboard.departments.departmentNamePlaceholder"
                    )}
                    value={form.name}
                    onChange={(e) => setForm({ ...form, name: e.target.value })}
                  />
                  <select
                    className="border px-2 w-full py-1"
                    value={form.parentId}
                    onChange={(e) =>
                      setForm({ ...form, parentId: e.target.value })
                    }
                  >
                    <option value="">
                      {t("dashboard.departments.selectParentDepartment")}
                    </option>
                    {/* Filter out the current department being edited from the parent list */}
                    {depts
                      .filter((d) => d.id !== editingDept?.id)
                      .map((d) => (
                        <option key={d.id} value={d.id}>
                          {d.name}
                        </option>
                      ))}
                  </select>
                  <select
                    className="border px-2 w-full py-1"
                    value={form.headId}
                    onChange={(e) =>
                      setForm({ ...form, headId: e.target.value })
                    }
                  >
                    <option value="">
                      {t("dashboard.departments.selectHead")}
                    </option>
                    {users.map((u) => (
                      <option key={u.id} value={u.id}>
                        {u.name}
                      </option>
                    ))}
                  </select>
                </div>
                <DialogFooter>
                  <DialogClose asChild>
                    <Button
                      variant="outline"
                      onClick={() => {
                        setEditOpen(false);
                        setEditingDept(null);
                        setForm({ name: "", headId: "", parentId: "" }); // Reset form on cancel
                      }}
                    >
                      {t("dashboard.departments.cancel")}
                    </Button>
                  </DialogClose>
                  <Button onClick={handleEdit}>
                    {t("dashboard.departments.saveChangesButton")}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          )}
      </div>
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.departments.listTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>{t("dashboard.departments.loading")}</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t("dashboard.departments.columnName")}</TableHead>
                  <TableHead>{t("dashboard.departments.columnHead")}</TableHead>
                  <TableHead>
                    {t("dashboard.departments.columnActions")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>{renderRows(tree)}</TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
