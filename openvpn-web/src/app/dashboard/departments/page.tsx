"use client";

import React, { useState, useEffect } from "react"; // Added useCallback if needed later, for now just useEffect
import { useAuth } from "@/lib/auth-context";
import { UserRole } from "@/types/types";
import { useTranslation } from "react-i18next";
import MainLayout from "@/components/layout/main-layout";
import { departmentAPI, userManagementAPI } from "@/services/api";
import type { Department, AdminUser } from "@/types/types"; // Department type should already support 'children?: Department[]'
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
import { Label } from "@/components/ui/label"; // For pagination UI
import { toast } from "sonner";

// DepartmentTree is used for rendering, ensuring children property
type DepartmentTree = Department & { children: DepartmentTree[] };

export default function DepartmentsPage() {
  const { user: currentUser } = useAuth();
  const { t } = useTranslation();

  // State for displayed departments (paginated top-level trees from backend)
  const [departments, setDepartments] = useState<DepartmentTree[]>([]);
  // State for all users (for 'Select Head' dropdown) - fetched once or with large pagination
  const [allUsersForSelect, setAllUsersForSelect] = useState<AdminUser[]>([]);
  // State for all departments (flat list for 'Select Parent' dropdown) - fetched once or with large pagination
  const [allDeptsForSelect, setAllDeptsForSelect] = useState<Department[]>([]);

  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false); // For Add dialog
  const [editOpen, setEditOpen] = useState(false); // For Edit dialog
  const [editingDept, setEditingDept] = useState<Department | null>(null);

  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());
  const [form, setForm] = useState({ name: "", headId: "", parentId: "" });

  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const [pageSize, setPageSize] = useState(10); // Default page size

  const fetchData = async (page: number, size: number) => {
    setLoading(true);
    try {
      // Fetch paginated departments (these are top-level, with children pre-loaded by backend)
      const deptResponse = await departmentAPI.list(page, size);
      setDepartments(deptResponse.departments as DepartmentTree[]); // Children are already part of the Department type from API
      setTotalItems(deptResponse.totalItems);
      setTotalPages(deptResponse.totalPages > 0 ? deptResponse.totalPages : 1);
      setCurrentPage(deptResponse.currentPage);

      // Fetch users for dropdowns (e.g., first 100 users) - adjust as needed
      if (allUsersForSelect.length === 0) { // Fetch only if not already populated
        const usersResponse = await userManagementAPI.list(1, 100);
        setAllUsersForSelect(usersResponse.users);
      }

      // Fetch all departments for parent selection dropdown - adjust as needed
      if (allDeptsForSelect.length === 0) { // Fetch only if not already populated
        // This assumes departmentAPI.list can fetch a large number or all if needed for parent selection
        // A better approach might be a dedicated endpoint for selectable parents.
        const allDeptsResponse = await departmentAPI.list(1, 500); // Fetch up to 500 depts for selection
        setAllDeptsForSelect(allDeptsResponse.departments);
      }

    } catch (error: any) {
      toast.error(error?.response?.data?.error || t("dashboard.departments.loadError"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData(currentPage, pageSize);
  }, [currentPage, pageSize]);

  // buildTree function is no longer needed as hierarchy comes from backend.

  const handleCreate = async () => {
    if (!form.name) {
      toast.error(t("dashboard.departments.nameRequired"));
      return;
    }
    try {
      await departmentAPI.create(form);
      toast.success(t("dashboard.departments.createSuccess"));
      setForm({ name: "", headId: "", parentId: "" });
      fetchData(1, pageSize); // Refresh list and go to first page
      setOpen(false);
    } catch (error: any) {
      toast.error(error?.response?.data?.error || t("dashboard.departments.createError"));
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
      fetchData(currentPage, pageSize); // Refresh current page
    } catch (error: any) {
      toast.error(error?.response?.data?.error || t("dashboard.departments.updateError"));
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm(t("dashboard.departments.deleteConfirm"))) return;
    try {
      await departmentAPI.delete(id);
      toast.success(t("dashboard.departments.deleteSuccess"));
      // After deleting, it's common to go to the first page or refetch current.
      // If current page becomes empty, it might be better to go to page 1 or previous page.
      fetchData(1, pageSize); // For simplicity, go to page 1.
    } catch (error: any) {
      toast.error(error?.response?.data?.error || t("dashboard.departments.deleteError"));
    }
  };

  const toggleExpand = (id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const renderRows = (
    nodes: DepartmentTree[],
    level: number = 0
  ): React.ReactNode[] => {
    return nodes.flatMap((node) => {
      const hasChildren = node.children && node.children.length > 0;
      const isExpanded = expandedIds.has(node.id);
      return [
        <TableRow key={node.id}>
          <TableCell style={{ paddingLeft: level * 20, position: "relative" }}>
            {hasChildren && (
              <span
                className="cursor-pointer select-none mr-1 flex items-center"
                onClick={() => toggleExpand(node.id)}
                style={{ width: 20, display: "inline-block", verticalAlign: "middle" }} // Adjusted style
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
            <span style={{ marginLeft: hasChildren ? 0 : 20 }}>{node.name}</span> {/* Indent text if no expand icon */}
          </TableCell>
          <TableCell>
            {allUsersForSelect.find(u => u.id === node.headId)?.name || t("dashboard.departments.emptyData")}
          </TableCell>
          <TableCell className="space-x-2 text-right"> {/* Ensure actions are right-aligned */}
            {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (
              <>
                <Button
                  size="sm" variant="outline"
                  onClick={() => {
                    setEditingDept(node);
                    setForm({ name: node.name, headId: node.headId || "", parentId: node.parentId || "" });
                    setEditOpen(true);
                  }}
                >
                  {t("dashboard.departments.edit")}
                </Button>
                <Button size="sm" variant="destructive" onClick={() => handleDelete(node.id)}>
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
        <h2 className="text-xl font-semibold">{t("dashboard.departments.pageTitle")}</h2>
        {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (
          <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => {
                setForm({ name: "", headId: "", parentId: "" }); // Reset form for add
                setOpen(true);
              }}>{t("dashboard.departments.addDepartment")}</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader><DialogTitle>{t("dashboard.departments.addDepartmentDialogTitle")}</DialogTitle></DialogHeader>
              <div className="space-y-2 pt-2">
                <Input placeholder={t("dashboard.departments.departmentNamePlaceholder")} value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
                <select className="border px-2 w-full py-1.5 rounded-md" value={form.parentId} onChange={(e) => setForm({ ...form, parentId: e.target.value })}>
                  <option value="">{t("dashboard.departments.selectParentDepartment")}</option>
                  {/* Use allDeptsForSelect, which should be a flat list of all departments */}
                  {allDeptsForSelect.map((d) => ( <option key={d.id} value={d.id}>{d.name}</option> ))}
                </select>
                <select className="border px-2 w-full py-1.5 rounded-md" value={form.headId} onChange={(e) => setForm({ ...form, headId: e.target.value })}>
                  <option value="">{t("dashboard.departments.selectHead")}</option>
                  {allUsersForSelect.map((u) => ( <option key={u.id} value={u.id}>{u.name}</option> ))}
                </select>
              </div>
              <DialogFooter>
                <DialogClose asChild><Button variant="outline">{t("dashboard.departments.cancel")}</Button></DialogClose>
                <Button onClick={handleCreate}>{t("dashboard.departments.create")}</Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        )}
      </div>

      {/* Edit Department Dialog - Placed outside the header's conditional rendering for clarity */}
      {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && editingDept && (
        <Dialog open={editOpen} onOpenChange={(isOpen) => {
          setEditOpen(isOpen);
          if (!isOpen) setEditingDept(null); // Reset editingDept when dialog closes
        }}>
          <DialogContent>
            <DialogHeader><DialogTitle>{t("dashboard.departments.editDepartmentDialogTitle")}</DialogTitle></DialogHeader>
            <div className="space-y-2 pt-2">
              <Input placeholder={t("dashboard.departments.departmentNamePlaceholder")} value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
              <select className="border px-2 w-full py-1.5 rounded-md" value={form.parentId} onChange={(e) => setForm({ ...form, parentId: e.target.value })}>
                <option value="">{t("dashboard.departments.selectParentDepartment")}</option>
                {allDeptsForSelect.filter(d => d.id !== editingDept?.id).map((d) => ( <option key={d.id} value={d.id}>{d.name}</option> ))}
              </select>
              <select className="border px-2 w-full py-1.5 rounded-md" value={form.headId} onChange={(e) => setForm({ ...form, headId: e.target.value })}>
                <option value="">{t("dashboard.departments.selectHead")}</option>
                {allUsersForSelect.map((u) => ( <option key={u.id} value={u.id}>{u.name}</option> ))}
              </select>
            </div>
            <DialogFooter>
              <DialogClose asChild><Button variant="outline" onClick={() => { setEditOpen(false); setEditingDept(null); }}>{t("dashboard.departments.cancel")}</Button></DialogClose>
              <Button onClick={handleEdit}>{t("dashboard.departments.saveChangesButton")}</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}

      <Card>
        <CardHeader><CardTitle>{t("dashboard.departments.listTitle")}</CardTitle></CardHeader>
        <CardContent>
          {loading ? ( <p>{t("dashboard.departments.loading")}</p> ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[50%]">{t("dashboard.departments.columnName")}</TableHead>
                    <TableHead className="w-[30%]">{t("dashboard.departments.columnHead")}</TableHead>
                    <TableHead className="w-[20%] text-right">{t("dashboard.departments.columnActions")}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {departments && departments.length > 0 ? renderRows(departments, 0) : (
                    <TableRow><TableCell colSpan={3} className="text-center">{t("dashboard.departments.noDepartmentsFound", "No departments found.")}</TableCell></TableRow>
                  )}
                </TableBody>
              </Table>
              {/* Pagination Controls */}
              {!loading && totalItems > 0 && (
                <div className="flex items-center justify-between mt-4 pt-2 border-t">
                  <div className="text-sm text-muted-foreground space-x-2">
                    <span>{t("dashboard.pagination.totalItems", { count: totalItems })}</span>
                    <span>|</span>
                    <span>{t("dashboard.pagination.pageInfo", { currentPage, totalPages })}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Label htmlFor="pageSizeSelectDept" className="sr-only">{t("dashboard.pagination.pageSizeLabel")}</Label>
                    <select
                      id="pageSizeSelectDept"
                      value={pageSize}
                      onChange={(e) => { setPageSize(Number(e.target.value)); setCurrentPage(1); }}
                      className="border px-2 py-1.5 rounded-md text-sm h-8 bg-background focus:ring-ring focus:border-input"
                    >
                      {[10, 20, 30, 50].map(size => ( <option key={size} value={size}>{t("dashboard.pagination.pageSizeOption", { count: size })}</option> ))}
                    </select>
                    <Button variant="outline" size="sm" onClick={() => setCurrentPage(p => Math.max(1, p - 1))} disabled={currentPage === 1} className="h-8 px-2">
                      {t("dashboard.pagination.previousButton", "Previous")}
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))} disabled={currentPage === totalPages || totalPages === 0} className="h-8 px-2">
                      {t("dashboard.pagination.nextButton", "Next")}
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
