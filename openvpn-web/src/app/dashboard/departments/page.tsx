"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { UserRole } from "@/types/types";
import { useTranslation } from "react-i18next";
import MainLayout from "@/components/layout/main-layout";
import { departmentAPI, userManagementAPI } from "@/services/api";
import type { Department, AdminUser } from "@/types/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";
import { toast } from "sonner";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";

type DepartmentTree = Department & { children: DepartmentTree[] };

export default function DepartmentsPage() {
  const { user: currentUser } = useAuth();
  const { t } = useTranslation();
  const [depts, setDepts] = useState<Department[]>([]);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());
  const [form, setForm] = useState({ name: "", headId: "", parentId: "" });
  const [deptPage, setDeptPage] = useState(0);
  const deptPageSize = 20;
  const [editOpen, setEditOpen] = useState(false);
  const [editingDept, setEditingDept] = useState<Department | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<{ open: boolean; id: string }>({ open: false, id: "" });

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
      setDepts(deptsWithHead);
      setUsers(uList);
    } catch {
      toast.error(t("dashboard.departments.loadError"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const buildTree = (list: Department[]): DepartmentTree[] => {
    const map = new Map<string, DepartmentTree>();
    list.forEach((item) => {
      map.set(item.id, { ...item, children: [] });
    });
    const roots: DepartmentTree[] = [];
    map.forEach((item) => {
      if (item.parentId && item.parentId !== item.id && map.has(item.parentId)) {
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

  const handleDelete = (id: string) => {
    setConfirmDelete({ open: true, id });
  };

  const doDelete = async () => {
    try {
      await departmentAPI.delete(confirmDelete.id);
      toast.success(t("dashboard.departments.deleteSuccess"));
      fetchData();
    } catch {
      toast.error(t("dashboard.departments.deleteError"));
    }
  };

  const toggleExpand = (id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const renderRows = (nodes: DepartmentTree[], level: number = 0): React.ReactNode[] => {
    return nodes.flatMap((node) => {
      const hasChildren = node.children && node.children.length > 0;
      const isExpanded = expandedIds.has(node.id);
      return [
        <TableRow key={node.id}>
          <TableCell
            style={{
              paddingLeft: level === 0 ? undefined : level * 20,
              position: "relative",
            }}
          >
            {hasChildren && (
              <span
                className="cursor-pointer select-none mr-1 flex items-center"
                onClick={() => toggleExpand(node.id)}
                style={{ width: 20, position: "absolute" }}
              >
                {isExpanded ? (
                  <svg width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                ) : (
                  <svg width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
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
            {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (
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
        ...(hasChildren && isExpanded ? renderRows(node.children, level + 1) : []),
      ];
    });
  };

  return (
    <MainLayout className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">
          {t("dashboard.departments.pageTitle")}
        </h2>
        {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (
          <Button onClick={() => setOpen(true)}>
            {t("dashboard.departments.addDepartment")}
          </Button>
        )}
      </div>

      {/* Add Department Dialog */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("dashboard.departments.addDepartmentDialogTitle")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 pt-2">
            <div className="space-y-1">
              <Label htmlFor="add-dept-name">{t("dashboard.departments.departmentNamePlaceholder")}</Label>
              <Input
                id="add-dept-name"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder={t("dashboard.departments.departmentNamePlaceholder")}
              />
            </div>
            <div className="space-y-1">
              <Label>{t("dashboard.departments.selectParentDepartment")}</Label>
              <Select
                value={form.parentId || "__none__"}
                onValueChange={(val) => setForm({ ...form, parentId: val === "__none__" ? "" : val })}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder={t("dashboard.departments.selectParentDepartment")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__none__">{t("dashboard.departments.selectParentDepartment")}</SelectItem>
                  {depts.map((d) => (
                    <SelectItem key={d.id} value={d.id}>{d.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1">
              <Label>{t("dashboard.departments.selectHead")}</Label>
              <Select
                value={form.headId || "__none__"}
                onValueChange={(val) => setForm({ ...form, headId: val === "__none__" ? "" : val })}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder={t("dashboard.departments.selectHead")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__none__">{t("dashboard.departments.selectHead")}</SelectItem>
                  {users.map((u) => (
                    <SelectItem key={u.id} value={u.id}>{u.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>
              {t("dashboard.departments.cancel")}
            </Button>
            <Button onClick={handleCreate}>
              {t("dashboard.departments.create")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Department Dialog */}
      {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (
        <Dialog
          open={editOpen}
          onOpenChange={(open) => {
            if (!open) {
              setEditOpen(false);
              setEditingDept(null);
              setForm({ name: "", headId: "", parentId: "" });
            } else {
              setEditOpen(true);
            }
          }}
        >
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{t("dashboard.departments.editDepartmentDialogTitle")}</DialogTitle>
            </DialogHeader>
            <div className="space-y-4 pt-2">
              <div className="space-y-1">
                <Label htmlFor="edit-dept-name">{t("dashboard.departments.departmentNamePlaceholder")}</Label>
                <Input
                  id="edit-dept-name"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder={t("dashboard.departments.departmentNamePlaceholder")}
                />
              </div>
              <div className="space-y-1">
                <Label>{t("dashboard.departments.selectParentDepartment")}</Label>
                <Select
                  value={form.parentId || "__none__"}
                  onValueChange={(val) => setForm({ ...form, parentId: val === "__none__" ? "" : val })}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("dashboard.departments.selectParentDepartment")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none__">{t("dashboard.departments.selectParentDepartment")}</SelectItem>
                    {depts.filter((d) => d.id !== editingDept?.id).map((d) => (
                      <SelectItem key={d.id} value={d.id}>{d.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1">
                <Label>{t("dashboard.departments.selectHead")}</Label>
                <Select
                  value={form.headId || "__none__"}
                  onValueChange={(val) => setForm({ ...form, headId: val === "__none__" ? "" : val })}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("dashboard.departments.selectHead")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none__">{t("dashboard.departments.selectHead")}</SelectItem>
                    {users.map((u) => (
                      <SelectItem key={u.id} value={u.id}>{u.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => {
                  setEditOpen(false);
                  setEditingDept(null);
                  setForm({ name: "", headId: "", parentId: "" });
                }}
              >
                {t("dashboard.departments.cancel")}
              </Button>
              <Button onClick={handleEdit}>
                {t("dashboard.departments.saveChangesButton")}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}

      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.departments.listTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>{t("dashboard.departments.loading")}</p>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t("dashboard.departments.columnName")}</TableHead>
                    <TableHead>{t("dashboard.departments.columnHead")}</TableHead>
                    <TableHead>{t("dashboard.departments.columnActions")}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {tree.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={3} className="py-10 text-center text-muted-foreground">
                        {t("dashboard.departments.emptyData")}
                      </TableCell>
                    </TableRow>
                  ) : (
                    renderRows(tree.slice(deptPage * deptPageSize, (deptPage + 1) * deptPageSize))
                  )}
                </TableBody>
              </Table>
              {/* Pagination */}
              {tree.length > deptPageSize && (
                <div className="flex items-center justify-between pt-4 border-t mt-2">
                  <span className="text-sm text-muted-foreground">
                    {deptPage * deptPageSize + 1}–{Math.min((deptPage + 1) * deptPageSize, tree.length)} / {tree.length}
                  </span>
                  <div className="flex items-center gap-1">
                    <Button size="sm" variant="outline" disabled={deptPage === 0} onClick={() => setDeptPage(0)} className="min-w-8 px-2">«</Button>
                    <Button size="sm" variant="outline" disabled={deptPage === 0} onClick={() => setDeptPage((p) => p - 1)} className="min-w-8 px-2">‹</Button>
                    <span className="px-3 text-sm">{deptPage + 1} / {Math.max(1, Math.ceil(tree.length / deptPageSize))}</span>
                    <Button size="sm" variant="outline" disabled={(deptPage + 1) * deptPageSize >= tree.length} onClick={() => setDeptPage((p) => p + 1)} className="min-w-8 px-2">›</Button>
                    <Button size="sm" variant="outline" disabled={(deptPage + 1) * deptPageSize >= tree.length} onClick={() => setDeptPage(Math.ceil(tree.length / deptPageSize) - 1)} className="min-w-8 px-2">»</Button>
                  </div>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>

      <AlertDialog
        open={confirmDelete.open}
        onOpenChange={(o) => setConfirmDelete((prev) => ({ ...prev, open: o }))}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("dashboard.departments.deleteConfirmTitle", "Delete Department")}</AlertDialogTitle>
            <AlertDialogDescription>{t("dashboard.departments.deleteConfirm")}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("dashboard.departments.cancel")}</AlertDialogCancel>
            <AlertDialogAction
              className={buttonVariants({ variant: "destructive" })}
              onClick={doDelete}
            >
              {t("dashboard.departments.delete")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </MainLayout>
  );
}
