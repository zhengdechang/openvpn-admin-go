"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogClose,
} from "@/components/ui/dialog";
import { useTranslation } from "react-i18next";
import MainLayout from "@/components/layout/main-layout";
import { userManagementAPI, departmentAPI, openvpnAPI } from "@/services/api";
import { AdminUser, Department, UserRole } from "@/types/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

export default function UsersPage() {
  const { user: currentUser } = useAuth();
  const { t } = useTranslation();
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [depts, setDepts] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
    role: UserRole.USER,
    departmentId: "",
  });
  const [open, setOpen] = useState(false); // For Add User Dialog

  // State for Edit User Dialog
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<AdminUser | null>(null);
  const [editFormDepartmentId, setEditFormDepartmentId] = useState<string>("");

  const fetchAll = async () => {
    setLoading(true);
    try {
      const [u, d] = await Promise.all([
        userManagementAPI.list(),
        departmentAPI.list(),
      ]);
      setUsers(u);
      setDepts(d);
    } catch {
      toast.error(t("dashboard.users.loadError"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAll();
  }, []);

  // 根据角色过滤用户
  let visibleUsers = users;
  // 下载配置
  const handleDownload = async (id: string, os: string) => {
    try {
      const { config } = await openvpnAPI.getClientConfig(id, os);
      const blob = new Blob([config], {
        type: "application/x-openvpn-profile",
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      const ext = os === "linux" ? "conf" : "ovpn";
      a.download = `${id}.${ext}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      toast.success(t("dashboard.users.downloadConfigSuccess"));
    } catch {
      toast.error(t("dashboard.users.downloadConfigError"));
    }
  };

  const handleCreate = async () => {
    if (!form.name || !form.email || !form.password) {
      toast.error(t("dashboard.users.formRequiredFields"));
      return;
    }
    try {
      await userManagementAPI.create({
        ...form,
      });
      toast.success(t("dashboard.users.createSuccess"));
      setForm({
        name: "",
        email: "",
        password: "",
        role: UserRole.USER,
        departmentId: "",
      });
      fetchAll();
    } catch {
      toast.error(t("dashboard.users.createError"));
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm(t("dashboard.users.deleteConfirm"))) return;
    try {
      await userManagementAPI.delete(id);
      toast.success(t("dashboard.users.deleteSuccess"));
      fetchAll();
    } catch {
      toast.error(t("dashboard.users.deleteError"));
    }
  };

  const handleEditClick = (user: AdminUser) => {
    setEditingUser(user);
    setEditFormDepartmentId(user.departmentId || "");
    setEditDialogOpen(true);
  };

  const handleUpdateUserDepartment = async () => {
    if (!editingUser || !editFormDepartmentId) {
      toast.error(t("dashboard.users.editUserErrorMissingInfo"));
      return;
    }
    try {
      await userManagementAPI.update(editingUser.id, {
        departmentId: editFormDepartmentId,
      });
      toast.success(t("dashboard.users.editUserSuccess"));
      setEditDialogOpen(false);
      fetchAll();
    } catch {
      toast.error(t("dashboard.users.editUserError"));
    } finally {
      setEditingUser(null);
      setEditFormDepartmentId("");
    }
  };

  return (
    <MainLayout className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">{t("dashboard.users.pageTitle")}</h1>
        {(currentUser?.role === UserRole.ADMIN ||
          currentUser?.role === UserRole.MANAGER ||
          currentUser?.role === UserRole.SUPERADMIN) && (
          <Dialog
            open={open}
            onOpenChange={(isOpen) => {
              setOpen(isOpen);
              if (isOpen) {
                console.log("Dialog opened. currentUser:", currentUser);
                console.log(
                  "currentUser.departmentId:",
                  currentUser?.departmentId
                );

                // Reset form fields and set initial values when dialog opens
                let initialDepartmentId = "";
                let initialRole = UserRole.USER; // Default role for new users

                if (currentUser?.role === UserRole.MANAGER) {
                  // Manager adds USER role
                  initialRole = UserRole.USER;
                  // Manager's department is pre-filled and disabled
                  initialDepartmentId = currentUser.departmentId || ""; // Ensure it's a string
                }
                // For other roles (Admin, Superadmin), default role is USER, department is selectable (initialDepartmentId remains "")

                setForm({
                  name: "",
                  email: "",
                  password: "",
                  role: initialRole,
                  departmentId: initialDepartmentId,
                });
              }
            }}
          >
            <DialogTrigger asChild>
              <Button>{t("dashboard.users.addUserButton")}</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>
                  {t("dashboard.users.addUserDialogTitle")}
                </DialogTitle>
              </DialogHeader>
              <div className="space-y-2 pt-2">
                <Input
                  placeholder={t("dashboard.users.namePlaceholder")}
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                />
                <Input
                  placeholder={t("dashboard.users.emailPlaceholder")}
                  value={form.email}
                  onChange={(e) => setForm({ ...form, email: e.target.value })}
                />
                <Input
                  type="password"
                  placeholder={t("dashboard.users.passwordPlaceholder")}
                  value={form.password}
                  onChange={(e) =>
                    setForm({ ...form, password: e.target.value })
                  }
                />
                <div className="flex space-x-2">
                  <select
                    className="border px-2"
                    value={form.role}
                    onChange={(e) =>
                      setForm({ ...form, role: e.target.value as UserRole })
                    }
                    disabled={currentUser?.role === UserRole.MANAGER}
                  >
                    <option value={UserRole.USER}>
                      {t("dashboard.users.roleUser")}
                    </option>
                    {(currentUser?.role === UserRole.ADMIN ||
                      currentUser?.role === UserRole.SUPERADMIN) && (
                      <>
                        <option value={UserRole.MANAGER}>
                          {t("dashboard.users.roleManager")}
                        </option>
                        <option value={UserRole.ADMIN}>
                          {t("dashboard.users.roleAdmin")}
                        </option>
                        <option value={UserRole.SUPERADMIN}>
                          {t("dashboard.users.roleSuperadmin")}
                        </option>
                      </>
                    )}
                  </select>
                  <select
                    className="border px-2 py-2 w-full rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                    value={form.departmentId}
                    onChange={(e) => {
                      if (currentUser?.role !== UserRole.MANAGER) {
                        setForm({ ...form, departmentId: e.target.value });
                      }
                    }}
                    disabled={currentUser?.role === UserRole.MANAGER}
                  >
                    {currentUser?.role === UserRole.MANAGER ? (
                      currentUser.departmentId && (
                        <option value={currentUser.departmentId}>
                          {depts.length > 0
                            ? depts.find(
                                (d) => d.id === currentUser.departmentId
                              )?.name || t("dashboard.users.emptyDepartment")
                            : t("common.loading")}
                        </option>
                      )
                    ) : (
                      <>
                        <option value="">
                          {t("dashboard.users.selectDepartmentPlaceholder")}
                        </option>
                        {depts.map((d) => (
                          <option key={d.id} value={d.id}>
                            {d.name}
                          </option>
                        ))}
                      </>
                    )}
                  </select>
                </div>
              </div>
              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline">{t("common.cancel")}</Button>
                </DialogClose>
                <Button
                  onClick={() => {
                    handleCreate();
                    setOpen(false);
                  }}
                >
                  {t("common.create")}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        )}
      </div>

      {/* Edit User Dialog */}
      <Dialog
        open={editDialogOpen}
        onOpenChange={(isOpen) => {
          setEditDialogOpen(isOpen);
          if (!isOpen) {
            setEditingUser(null);
            setEditFormDepartmentId("");
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {t("dashboard.users.editUserDialogTitle")}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-2 pt-2">
            <label
              htmlFor="edit-department"
              className="block text-sm font-medium text-gray-700"
            >
              {t("dashboard.users.departmentLabel", "Department")}
            </label>
            <select
              id="edit-department"
              className="border px-2 py-2 w-full rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              value={editFormDepartmentId}
              onChange={(e) => setEditFormDepartmentId(e.target.value)}
            >
              <option value="">
                {t("dashboard.users.selectDepartmentPlaceholder")}
              </option>
              {depts.map((d) => (
                <option key={d.id} value={d.id}>
                  {d.name}
                </option>
              ))}
            </select>
            {editingUser && (
              <p className="text-sm text-gray-500 mt-2">
                {t("dashboard.users.editingUserLabel")} {editingUser.name} (
                {editingUser.email})
              </p>
            )}
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button
                variant="outline"
                onClick={() => {
                  setEditingUser(null);
                  setEditFormDepartmentId("");
                }}
              >
                {t("common.cancel")}
              </Button>
            </DialogClose>
            <Button onClick={handleUpdateUserDepartment}>
              {t("common.saveChanges")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.users.listTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>{t("common.loading")}</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t("dashboard.users.columnName")}</TableHead>
                  <TableHead>{t("dashboard.users.columnEmail")}</TableHead>
                  <TableHead>{t("dashboard.users.columnRole")}</TableHead>
                  <TableHead>{t("dashboard.users.columnDepartment")}</TableHead>
                  <TableHead>{t("dashboard.users.columnLastConnection")}</TableHead>
                  <TableHead>{t("dashboard.users.columnOnlineStatus")}</TableHead>
                  <TableHead>{t("dashboard.users.columnCreator")}</TableHead>
                  <TableHead>{t("dashboard.users.columnActions")}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {visibleUsers.map((u: AdminUser) => ( // Explicitly use AdminUser type
                  <TableRow key={u.id}>
                    <TableCell>{u.name}</TableCell>
                    <TableCell>{u.email}</TableCell>
                    <TableCell>{u.role}</TableCell>
                    <TableCell>
                      {depts.find((d) => d.id === u.departmentId)?.name ||
                        t("dashboard.users.emptyDepartment")}
                    </TableCell>
                    <TableCell>
                      {u.lastConnectionTime
                        ? new Date(u.lastConnectionTime).toLocaleString()
                        : t("common.na")}
                    </TableCell>
                    <TableCell>
                      {typeof u.isOnline === 'boolean'
                        ? u.isOnline
                          ? t("dashboard.users.statusOnline")
                          : t("dashboard.users.statusOffline")
                        : t("common.na")}
                    </TableCell>
                    <TableCell>{u.creatorId || t("common.na")}</TableCell>
                    <TableCell className="space-x-2">
                      <div className="flex items-center gap-2">
                        {(currentUser?.role === UserRole.ADMIN ||
                          currentUser?.role === UserRole.SUPERADMIN) && (
                          <Button
                            size="sm"
                            variant="outline"
                            className="h-8 px-3"
                            onClick={() => handleEditClick(u)}
                          >
                            {t("common.edit")}
                          </Button>
                        )}
                        {(currentUser?.role === UserRole.ADMIN ||
                          currentUser?.role === UserRole.MANAGER ||
                          currentUser?.role === UserRole.SUPERADMIN) && (
                          <Button
                            size="sm"
                            variant="destructive"
                            className="h-8 px-3"
                            onClick={() => handleDelete(u.id)}
                          >
                            {t("dashboard.users.deleteButton")}
                          </Button>
                        )}
                        <select
                          className="border px-2 py-1 rounded-md text-sm h-8"
                          defaultValue=""
                          onChange={(e) => handleDownload(u.id, e.target.value)}
                        >
                          <option value="" disabled>
                            {t("dashboard.users.downloadConfigButton")}
                          </option>
                          <option value="windows">
                            {t("dashboard.users.osWindows")}
                          </option>
                          <option value="macos">
                            {t("dashboard.users.osMacOS")}
                          </option>
                          <option value="linux">
                            {t("dashboard.users.osLinux")}
                          </option>
                        </select>
                      </div>
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
