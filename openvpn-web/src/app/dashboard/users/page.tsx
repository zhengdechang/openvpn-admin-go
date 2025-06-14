// In openvpn-web/src/app/dashboard/users/page.tsx
"use client";

import React, { useState, useEffect, useCallback } from "react"; // Added useCallback
import { useAuth } from "@/lib/auth-context";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogDescription, // Added DialogDescription
  DialogClose,
} from "@/components/ui/dialog";
import { useTranslation } from "react-i18next";
import MainLayout from "@/components/layout/main-layout";
import { userManagementAPI, departmentAPI, openvpnAPI } from "@/services/api";
// Ensure UserUpdateRequest is imported if defined and used
import {
  AdminUser,
  Department,
  UserRole,
  UserUpdateRequest,
} from "@/types/types";
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
import { Label } from "@/components/ui/label"; // Added Label
import { toast } from "sonner";

// Helper function to format bytes into a readable string
const formatBytes = (bytes?: number, decimals = 2): string => {
  if (bytes === undefined || bytes === null || bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
};

// Define initial state for the edit form
const initialEditFormState: UserUpdateRequest = {
  name: "",
  email: "",
  role: UserRole.USER,
  departmentId: "",
  fixedIp: "", // Initialize with empty string
  subnet: "", // Initialize with empty string
  password: "", // For password changes
};

export default function UsersPage() {
  const { user: currentUser } = useAuth();
  const { t } = useTranslation();
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [depts, setDepts] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);

  // Form state for adding a new user
  const [addUserForm, setAddUserForm] = useState<UserUpdateRequest>({
    // Using UserUpdateRequest for consistency
    name: "",
    email: "",
    password: "",
    role: UserRole.USER,
    departmentId: "",
    fixedIp: "",
    subnet: "", // Initialize with empty string
  });
  const [addUserDialogOpen, setAddUserDialogOpen] = useState(false);

  // State for Edit User Dialog
  const [editUserDialogOpen, setEditUserDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<AdminUser | null>(null);
  const [editForm, setEditForm] =
    useState<UserUpdateRequest>(initialEditFormState);

  // Define handlePauseUser function
  const handlePauseUser = async (username: string) => {
    if (!confirm(t("dashboard.users.pauseConfirm", `Are you sure you want to pause user ${username}?`))) return;
    try {
      await userManagementAPI.pauseUser(username);
      toast.success(t("dashboard.users.pauseSuccess", `User ${username} paused successfully.`));
      fetchAll(); // Refresh the user list to show updated status
    } catch (error: any) {
      toast.error(error?.response?.data?.error || t("dashboard.users.pauseError", `Failed to pause user ${username}.`));
    }
  };

  // Define handleResumeUser function
  const handleResumeUser = async (username: string) => {
    if (!confirm(t("dashboard.users.resumeConfirm", `Are you sure you want to resume user ${username}?`))) return;
    try {
      await userManagementAPI.resumeUser(username);
      toast.success(t("dashboard.users.resumeSuccess", `User ${username} resumed successfully.`));
      fetchAll(); // Refresh the user list
    } catch (error: any) {
      toast.error(error?.response?.data?.error || t("dashboard.users.resumeError", `Failed to resume user ${username}.`));
    }
  };

  const fetchAll = useCallback(async () => {
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
  }, [t]); // Added t

  useEffect(() => {
    fetchAll();
  }, [fetchAll]);

  const handleDownload = async (username: string, os: string) => {
    try {
      const data = await openvpnAPI.getClientConfig(username, os);
      const config = data.config;
      const blob = new Blob([config], {
        type: "application/x-openvpn-profile",
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      const ext = os === "linux" ? "conf" : "ovpn";
      a.download = `${username}.${ext}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      toast.success(t("dashboard.users.downloadConfigSuccess"));
    } catch {
      toast.error(t("dashboard.users.downloadConfigError"));
    }
  };

  const handleCreateUser = async () => {
    if (!addUserForm.name || !addUserForm.email || !addUserForm.password) {
      toast.error(t("dashboard.users.formRequiredFields"));
      return;
    }
    try {
      const payload: UserUpdateRequest = { ...addUserForm };
      if (
        !(
          currentUser?.role === UserRole.ADMIN ||
          currentUser?.role === UserRole.SUPERADMIN
        ) &&
        (payload.fixedIp || payload.subnet)
      ) {
        // Non-admins cannot set fixed IP or subnet on creation, clear them if set by mistake in form state
        payload.fixedIp = "";
        payload.subnet = "";
      }
      // Ensure departmentId is set if manager is creating user
      if (currentUser?.role === UserRole.MANAGER && !payload.departmentId) {
        payload.departmentId = currentUser.departmentId || "";
      }

      // If fixedIp is empty string, set it to null
      if (payload.fixedIp === "") {
        payload.fixedIp = null;
      }
      // If subnet is empty string, set it to null
      if (payload.subnet === "") {
        payload.subnet = null;
      }

      await userManagementAPI.create(
        payload as Partial<AdminUser> & { password: string }
      );
      toast.success(t("dashboard.users.createSuccess"));
      setAddUserDialogOpen(false);
      fetchAll();
    } catch (error: any) {
      toast.error(
        error?.response?.data?.error || t("dashboard.users.createError")
      );
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

  const handleEditClick = (userToEdit: AdminUser) => {
    setEditingUser(userToEdit);
    setEditForm({
      name: userToEdit.name,
      email: userToEdit.email,
      role: userToEdit.role,
      departmentId: userToEdit.departmentId || "",
      fixedIp: userToEdit.fixedIp || "",
      subnet: userToEdit.subnet || "",
      password: "",
    });
    setEditUserDialogOpen(true);
  };

  const handleUpdateUser = async () => {
    if (!editingUser) return;

    const updatePayload: UserUpdateRequest = { ...editForm };

    if (!updatePayload.password?.trim()) {
      delete updatePayload.password;
    }

    if (
      !(
        currentUser?.role === UserRole.ADMIN ||
        currentUser?.role === UserRole.SUPERADMIN
      )
    ) {
      // If user is not admin/superadmin, don't send fixedIp or subnet
      delete updatePayload.fixedIp;
      delete updatePayload.subnet;
    } else {
      // If fixedIp is empty string, set it to null
      if (updatePayload.fixedIp === "") {
        updatePayload.fixedIp = null;
      }
      // If subnet is empty string, set it to null
      if (updatePayload.subnet === "") {
        updatePayload.subnet = null;
      }
    }

    try {
      await userManagementAPI.update(editingUser.id, updatePayload);
      toast.success(
        t("dashboard.users.editUserSuccess", "User updated successfully!")
      );
      setEditUserDialogOpen(false);
      fetchAll();
    } catch (error: any) {
      toast.error(
        error?.response?.data?.error ||
          t("dashboard.users.editUserError", "Failed to update user.")
      );
    }
  };

  const canEditFixedIp =
    currentUser?.role === UserRole.ADMIN ||
    currentUser?.role === UserRole.SUPERADMIN;
  const canManageUsers =
    currentUser?.role === UserRole.ADMIN ||
    currentUser?.role === UserRole.SUPERADMIN ||
    currentUser?.role === UserRole.MANAGER;

  return (
    <MainLayout className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">{t("dashboard.users.pageTitle")}</h1>
        {canManageUsers && (
          <Dialog open={addUserDialogOpen} onOpenChange={setAddUserDialogOpen}>
            <DialogTrigger asChild>
              <Button
                onClick={() => {
                  let initialDeptId = "";
                  let initialRole = UserRole.USER;
                  if (currentUser?.role === UserRole.MANAGER) {
                    initialDeptId = currentUser.departmentId || "";
                    // Managers can only create users
                  }
                  setAddUserForm({
                    name: "",
                    email: "",
                    password: "",
                    role: initialRole,
                    departmentId: initialDeptId,
                    fixedIp: "",
                    subnet: "", // Initialize with empty string
                  });
                  setAddUserDialogOpen(true);
                }}
              >
                {t("dashboard.users.addUserButton")}
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>
                  {t("dashboard.users.addUserDialogTitle")}
                </DialogTitle>
              </DialogHeader>
              <div className="space-y-3 py-3">
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="add-name" className="text-right">
                    {t("dashboard.users.namePlaceholder")}
                  </Label>
                  <Input
                    id="add-name"
                    value={addUserForm.name}
                    onChange={(e) =>
                      setAddUserForm({ ...addUserForm, name: e.target.value })
                    }
                    className="col-span-3"
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="add-email" className="text-right">
                    {t("dashboard.users.emailPlaceholder")}
                  </Label>
                  <Input
                    id="add-email"
                    type="email"
                    value={addUserForm.email}
                    onChange={(e) =>
                      setAddUserForm({ ...addUserForm, email: e.target.value })
                    }
                    className="col-span-3"
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="add-password" className="text-right">
                    {t("dashboard.users.passwordPlaceholder")}
                  </Label>
                  <Input
                    id="add-password"
                    type="password"
                    value={addUserForm.password}
                    onChange={(e) =>
                      setAddUserForm({
                        ...addUserForm,
                        password: e.target.value,
                      })
                    }
                    className="col-span-3"
                  />
                </div>

                {canEditFixedIp && (
                  <>
                    <div className="grid grid-cols-4 items-center gap-4">
                      <Label htmlFor="add-fixedIp" className="text-right">
                        {t("dashboard.users.fixedIpLabel", "Fixed VPN IP (Optional)")}
                      </Label>
                      <Input
                        id="add-fixedIp"
                        value={addUserForm.fixedIp || ""}
                        onChange={(e) =>
                          setAddUserForm({
                            ...addUserForm,
                            fixedIp: e.target.value,
                          })
                        }
                        placeholder={t(
                          "dashboard.users.fixedIpPlaceholder",
                          "e.g., 10.8.0.100 or empty"
                        )}
                        className="col-span-3"
                      />
                    </div>
                    <div className="grid grid-cols-4 items-center gap-4">
                      <Label htmlFor="add-subnet" className="text-right">
                        {t("dashboard.users.subnetLabel", "Subnet (Optional)")}
                      </Label>
                      <Input
                        id="add-subnet"
                        value={addUserForm.subnet || ""}
                        onChange={(e) =>
                          setAddUserForm({
                            ...addUserForm,
                            subnet: e.target.value,
                          })
                        }
                        placeholder={t(
                          "dashboard.users.subnetPlaceholder",
                          "e.g., 10.10.120.0/23 or empty"
                        )}
                        className="col-span-3"
                      />
                    </div>
                  </>
                )}

                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="add-role" className="text-right">
                    {t("dashboard.users.roleLabel", "Role")}
                  </Label>
                  <select
                    id="add-role"
                    value={addUserForm.role}
                    onChange={(e) =>
                      setAddUserForm({
                        ...addUserForm,
                        role: e.target.value as UserRole,
                      })
                    }
                    className="col-span-3 border px-2 py-2 rounded-md"
                    disabled={
                      currentUser?.role === UserRole.MANAGER
                    } /* Managers can only create 'user' role */
                  >
                    {Object.values(UserRole)
                      .filter(
                        (role) =>
                          currentUser?.role === UserRole.SUPERADMIN || // Superadmin can assign any role
                          (currentUser?.role === UserRole.ADMIN &&
                            role !== UserRole.SUPERADMIN) || // Admin can assign any role except Superadmin
                          (currentUser?.role === UserRole.MANAGER &&
                            role === UserRole.USER) // Manager can only assign User
                      )
                      .map((role) => (
                        <option key={role} value={role}>
                          {t(
                            `dashboard.users.role${
                              role.charAt(0).toUpperCase() + role.slice(1)
                            }`,
                            role
                          )}
                        </option>
                      ))}
                  </select>
                </div>

                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="add-department" className="text-right">
                    {t("dashboard.users.departmentLabel")}
                  </Label>
                  <select
                    id="add-department"
                    value={addUserForm.departmentId || ""}
                    onChange={(e) =>
                      setAddUserForm({
                        ...addUserForm,
                        departmentId: e.target.value,
                      })
                    }
                    className="col-span-3 border px-2 py-2 rounded-md"
                    disabled={
                      currentUser?.role === UserRole.MANAGER &&
                      !!currentUser.departmentId
                    } /* Manager cannot change their own department selection */
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
                </div>
              </div>
              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline">{t("common.cancel")}</Button>
                </DialogClose>
                <Button onClick={handleCreateUser}>{t("common.create")}</Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        )}
      </div>

      <Dialog open={editUserDialogOpen} onOpenChange={setEditUserDialogOpen}>
        <DialogContent className="sm:max-w-[525px]">
          <DialogHeader>
            <DialogTitle>
              {t("dashboard.users.editUserDialogTitle", "Edit User")}
            </DialogTitle>
            <DialogDescription>
              {t(
                "dashboard.users.editUserDescription",
                "Make changes to the user profile here. Click save when you're done."
              )}
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit-name" className="text-right">
                {t("dashboard.users.namePlaceholder")}
              </Label>
              <Input
                id="edit-name"
                value={editForm.name || ""}
                onChange={(e) =>
                  setEditForm({ ...editForm, name: e.target.value })
                }
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit-email" className="text-right">
                {t("dashboard.users.emailPlaceholder")}
              </Label>
              <Input
                id="edit-email"
                type="email"
                value={editForm.email || ""}
                onChange={(e) =>
                  setEditForm({ ...editForm, email: e.target.value })
                }
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit-password" className="text-right">
                {t(
                  "dashboard.users.passwordOptionalPlaceholder",
                  "Password (optional)"
                )}
              </Label>
              <Input
                id="edit-password"
                type="password"
                value={editForm.password || ""}
                onChange={(e) =>
                  setEditForm({ ...editForm, password: e.target.value })
                }
                className="col-span-3"
                placeholder={t(
                  "dashboard.users.passwordLeaveBlankPlaceholder",
                  "Leave blank to keep current"
                )}
              />
            </div>

            {canEditFixedIp && (
              <>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="edit-fixedIp" className="text-right">
                    {t("dashboard.users.fixedIpLabel", "Fixed VPN IP (Optional)")}
                  </Label>
                  <Input
                    id="edit-fixedIp"
                    value={editForm.fixedIp || ""}
                    onChange={(e) =>
                      setEditForm({
                        ...editForm,
                        fixedIp: e.target.value,
                      })
                    }
                    placeholder={t(
                      "dashboard.users.fixedIpPlaceholder",
                      "e.g., 10.8.0.100 or empty"
                    )}
                    className="col-span-3"
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="edit-subnet" className="text-right">
                    {t("dashboard.users.subnetLabel", "Subnet (Optional)")}
                  </Label>
                  <Input
                    id="edit-subnet"
                    value={editForm.subnet || ""}
                    onChange={(e) =>
                      setEditForm({
                        ...editForm,
                        subnet: e.target.value,
                      })
                    }
                    placeholder={t(
                      "dashboard.users.subnetPlaceholder",
                      "e.g., 10.10.120.0/23 or empty"
                    )}
                    className="col-span-3"
                  />
                </div>
              </>
            )}
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit-department" className="text-right">
                {t("dashboard.users.departmentLabel")}
              </Label>
              <select
                id="edit-department"
                value={editForm.departmentId || ""}
                onChange={(e) =>
                  setEditForm({ ...editForm, departmentId: e.target.value })
                }
                className="col-span-3 border px-2 py-2 rounded-md"
                disabled={
                  currentUser?.role === UserRole.MANAGER &&
                  editingUser?.departmentId !== currentUser.departmentId &&
                  editingUser?.id !== currentUser.id
                }
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
            </div>
            {(currentUser?.role === UserRole.SUPERADMIN ||
              (currentUser?.role === UserRole.ADMIN &&
                editingUser?.role !== UserRole.SUPERADMIN)) &&
              (editingUser?.id !== currentUser?.id ||
                currentUser?.role ===
                  UserRole.SUPERADMIN) /* Can't change own role unless superadmin */ && (
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="edit-role" className="text-right">
                    {t("dashboard.users.roleLabel", "Role")}
                  </Label>
                  <select
                    id="edit-role"
                    value={editForm.role || ""}
                    onChange={(e) =>
                      setEditForm({
                        ...editForm,
                        role: e.target.value as UserRole,
                      })
                    }
                    className="col-span-3 border px-2 py-2 rounded-md"
                  >
                    {Object.values(UserRole)
                      .filter(
                        (role) =>
                          currentUser?.role === UserRole.SUPERADMIN ||
                          role !== UserRole.SUPERADMIN
                      )
                      .map((role) => (
                        <option key={role} value={role}>
                          {t(
                            `dashboard.users.role${
                              role.charAt(0).toUpperCase() + role.slice(1)
                            }`,
                            role
                          )}
                        </option>
                      ))}
                  </select>
                </div>
              )}
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">{t("common.cancel")}</Button>
            </DialogClose>
            <Button onClick={handleUpdateUser}>
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
            <div className="relative">
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-[150px]">{t("dashboard.users.columnName")}</TableHead>
                      <TableHead className="w-[200px]">{t("dashboard.users.columnEmail")}</TableHead>
                      <TableHead className="w-[100px]">{t("dashboard.users.columnRole")}</TableHead>
                      <TableHead className="w-[150px]">{t("dashboard.users.columnDepartment")}</TableHead>
                      <TableHead className="w-[120px]">
                        {t("dashboard.users.columnFixedIp", "Fixed IP")}
                      </TableHead>
                      <TableHead className="w-[120px]">
                        {t("dashboard.users.columnSubnet", "Subnet")}
                      </TableHead>
                      <TableHead className="w-[120px]">
                        {t("dashboard.users.columnConnectionIp", "Connection IP")}
                      </TableHead>
                      <TableHead className="w-[120px]">
                        {t("dashboard.users.columnAllocatedVpnIp", "VPN IP")}
                      </TableHead>
                      <TableHead className="w-[180px]">
                        {t("dashboard.users.columnLastConnection")}
                      </TableHead>
                      <TableHead className="w-[100px]">
                        {t("dashboard.users.columnOnlineStatus")}
                      </TableHead>
                      <TableHead className="w-[100px]">
                        {t("dashboard.users.columnAccessState", "Access State")}
                      </TableHead>
                      <TableHead className="w-[120px]">{t("dashboard.users.columnCreator")}</TableHead>
                      <TableHead className="w-[120px]">{t("dashboard.users.columnBytesReceived", "Bytes Received")}</TableHead>
                      <TableHead className="w-[120px]">{t("dashboard.users.columnBytesSent", "Bytes Sent")}</TableHead>
                      <TableHead className="w-[300px] sticky right-0 bg-background shadow-[-4px_0_8px_rgba(0,0,0,0.2)]"> {/* Increased width for new buttons */}
                        {t("dashboard.users.columnActions")}
                      </TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {users.map((u: AdminUser) => (
                      <TableRow key={u.id}>
                        <TableCell>{u.name}</TableCell>
                        <TableCell>{u.email}</TableCell>
                        <TableCell>{u.role}</TableCell>
                        <TableCell>
                          {depts.find((d) => d.id === u.departmentId)?.name ||
                            t("dashboard.users.emptyDepartment")}
                        </TableCell>
                        <TableCell>{u.fixedIp || "-"} </TableCell>
                        <TableCell>
                          {u.subnet || "-"}
                        </TableCell>
                        <TableCell>
                          {u.connectionIp || t("common.na")}
                        </TableCell>
                        <TableCell>
                          {u.allocatedVpnIp || t("common.na")}
                        </TableCell>
                        <TableCell>
                          {u.lastConnectionTime
                            ? new Date(u.lastConnectionTime).toLocaleString()
                            : t("common.na")}
                        </TableCell>
                        <TableCell>
                          {typeof u.isOnline === "boolean"
                            ? u.isOnline
                              ? t("dashboard.users.statusOnline")
                              : t("dashboard.users.statusOffline")
                            : t("common.na")}
                        </TableCell>
                        <TableCell>
                          {users.find((creator) => creator.id === u.creatorId)
                            ?.name || t("common.na")}
                        </TableCell>
                        <TableCell>
                          {u.isPaused ? t("dashboard.users.statusPaused", "Paused") : t("dashboard.users.statusActive", "Active")}
                        </TableCell>
                        <TableCell>{formatBytes(u.bytesReceived)}</TableCell>
                        <TableCell>{formatBytes(u.bytesSent)}</TableCell>
                        <TableCell className="sticky right-0 bg-background shadow-[-4px_0_8px_rgba(0,0,0,0.1)]">
                          <div className="flex items-center justify-center gap-1">
                            {canManageUsers && (
                              <>
                                {u.isPaused ? (
                                  <Button
                                    size="sm"
                                    variant="outline"
                                    className="h-8 px-2"
                                    onClick={() => handleResumeUser(u.name)}
                                  >
                                    {t("dashboard.users.resumeButton", "Resume")}
                                  </Button>
                                ) : (
                                  <Button
                                    size="sm"
                                    variant="outline"
                                    className="h-8 px-2"
                                    onClick={() => handlePauseUser(u.name)}
                                  >
                                    {t("dashboard.users.pauseButton", "Pause")}
                                  </Button>
                                )}
                              </>
                            )}
                            {(currentUser?.role === UserRole.ADMIN ||
                              currentUser?.role === UserRole.SUPERADMIN ||
                              (currentUser?.role === UserRole.MANAGER &&
                                currentUser.departmentId === u.departmentId)) && (
                              <Button
                                size="sm"
                                variant="outline"
                                className="h-8 px-2"
                                onClick={() => handleEditClick(u)}
                              >
                                {t("common.edit")}
                              </Button>
                            )}
                            {(currentUser?.role === UserRole.SUPERADMIN ||
                              (currentUser?.role === UserRole.ADMIN &&
                                u.role !== UserRole.SUPERADMIN) ||
                              (currentUser?.role === UserRole.MANAGER &&
                                u.departmentId === currentUser.departmentId &&
                                u.role !== UserRole.SUPERADMIN &&
                                u.role !== UserRole.ADMIN)) &&
                              u.id !== currentUser?.id && (
                                <Button
                                  size="sm"
                                  variant="destructive"
                                  className="h-8 px-2"
                                  onClick={() => handleDelete(u.id)}
                                >
                                  {t("dashboard.users.deleteButton")}
                                </Button>
                              )}
                            <select
                              className="border px-1 py-1 rounded-md text-sm h-8"
                              defaultValue=""
                              onChange={(e) => handleDownload(u.name, e.target.value)}
                            >
                              <option value="" disabled>
                                {t("dashboard.users.downloadConfigButtonShort", "DL")}
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
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
