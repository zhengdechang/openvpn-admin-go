// In openvpn-web/src/app/dashboard/users/page.tsx
"use client";

import React, { useState, useEffect, useCallback, useMemo } from "react";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import MainLayout from "@/components/layout/main-layout";
import { userManagementAPI, departmentAPI, openvpnAPI } from "@/services/api";
import {
  AdminUser,
  Department,
  UserRole,
  UserUpdateRequest,
} from "@/types/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { toast } from "sonner";
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  flexRender,
  type ColumnDef,
  type SortingState,
} from "@tanstack/react-table";
import MuiButton from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Chip from "@mui/material/Chip";
import { Dialog as MuiDialog, DialogTitle, DialogContent, DialogActions, DialogContentText } from "@mui/material";
import { FormControl, InputLabel, Select as MuiSelect, MenuItem } from "@mui/material";

// Helper function to format bytes into a readable string
const formatBytes = (bytes?: number, decimals = 2): string => {
  if (bytes === undefined || bytes === null || bytes === 0) return "0 Bytes";
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
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
  const [searchTerm, setSearchTerm] = useState("");
  const [departmentFilter, setDepartmentFilter] = useState<string>("all");
  const [statusFilter, setStatusFilter] = useState<string>("all");

  // Form state for adding a new user
  const [addUserForm, setAddUserForm] = useState<UserUpdateRequest>({
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
    if (
      !confirm(
        t(
          "dashboard.users.pauseConfirm",
          `Are you sure you want to pause user ${username}?`
        )
      )
    )
      return;
    try {
      await userManagementAPI.pauseUser(username);
      toast.success(
        t(
          "dashboard.users.pauseSuccess",
          `User ${username} paused successfully.`
        )
      );
      fetchAll(); // Refresh the user list to show updated status
    } catch (error: any) {
      toast.error(
        error?.response?.data?.error ||
          t("dashboard.users.pauseError", `Failed to pause user ${username}.`)
      );
    }
  };

  // Define handleResumeUser function
  const handleResumeUser = async (username: string) => {
    if (
      !confirm(
        t(
          "dashboard.users.resumeConfirm",
          `Are you sure you want to resume user ${username}?`
        )
      )
    )
      return;
    try {
      await userManagementAPI.resumeUser(username);
      toast.success(
        t(
          "dashboard.users.resumeSuccess",
          `User ${username} resumed successfully.`
        )
      );
      fetchAll(); // Refresh the user list
    } catch (error: any) {
      toast.error(
        error?.response?.data?.error ||
          t("dashboard.users.resumeError", `Failed to resume user ${username}.`)
      );
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

  const departmentNameById = useMemo(() => {
    return depts.reduce<Record<string, string>>((acc, dept) => {
      acc[dept.id] = dept.name;
      return acc;
    }, {});
  }, [depts]);

  const stats = useMemo(() => {
    const total = users.length;
    const online = users.filter((u) => u.isOnline).length;
    const paused = users.filter((u) => u.isPaused).length;
    const totalTraffic = users.reduce((sum, u) => {
      return sum + (u.bytesReceived || 0) + (u.bytesSent || 0);
    }, 0);
    return { total, online, paused, totalTraffic };
  }, [users]);

  const normalizedSearch = searchTerm.trim().toLowerCase();

  const filteredUsers = useMemo(() => {
    return users
      .filter((user) => {
        const matchesSearch =
          normalizedSearch.length === 0 ||
          user.name?.toLowerCase().includes(normalizedSearch) ||
          user.email?.toLowerCase().includes(normalizedSearch);
        const matchesDepartment =
          departmentFilter === "all" ||
          (departmentFilter === "none" && !user.departmentId) ||
          (user.departmentId || "") === departmentFilter;
        const matchesStatus =
          statusFilter === "all" ||
          (statusFilter === "paused" && user.isPaused) ||
          (statusFilter === "online" && !user.isPaused && user.isOnline) ||
          (statusFilter === "offline" &&
            !user.isPaused &&
            user.isOnline === false);
        return matchesSearch && matchesDepartment && matchesStatus;
      })
      .sort((a, b) => (a.name || "").localeCompare(b.name || ""));
  }, [
    users,
    normalizedSearch,
    departmentFilter,
    statusFilter,
  ]);

  const hasFilters =
    normalizedSearch.length > 0 ||
    departmentFilter !== "all" ||
    statusFilter !== "all";

  const handleClearFilters = () => {
    setSearchTerm("");
    setDepartmentFilter("all");
    setStatusFilter("all");
  };

  // ── TanStack React Table ────────────────────────────────────────────
  const [sorting, setSorting] = useState<SortingState>([]);

  const columns = useMemo<ColumnDef<AdminUser>[]>(
    () => [
      {
        accessorKey: "name",
        header: () => t("dashboard.users.columnName"),
        cell: ({ row }) => <span className="font-medium">{row.original.name}</span>,
      },
      {
        accessorKey: "email",
        header: () => t("dashboard.users.columnEmail"),
      },
      {
        accessorKey: "role",
        header: () => t("dashboard.users.columnRole"),
        cell: ({ row }) => (
          <Chip
            label={t(`dashboard.users.role${row.original.role.charAt(0).toUpperCase() + row.original.role.slice(1)}`, row.original.role)}
            size="small"
            color="default"
            sx={{ textTransform: "capitalize" }}
          />
        ),
      },
      {
        accessorKey: "departmentId",
        header: () => t("dashboard.users.columnDepartment"),
        cell: ({ row }) =>
          row.original.departmentId
            ? departmentNameById[row.original.departmentId] || t("dashboard.users.emptyDepartment")
            : t("dashboard.users.emptyDepartment"),
      },
      {
        accessorKey: "fixedIp",
        header: () => t("dashboard.users.columnFixedIp", "Fixed IP"),
        cell: ({ row }) => row.original.fixedIp || "-",
      },
      {
        accessorKey: "subnet",
        header: () => t("dashboard.users.columnSubnet", "Subnet"),
        cell: ({ row }) => row.original.subnet || "-",
      },
      {
        accessorKey: "connectionIp",
        header: () => t("dashboard.users.columnConnectionIp", "Connection IP"),
        cell: ({ row }) => row.original.connectionIp || t("common.na"),
      },
      {
        accessorKey: "allocatedVpnIp",
        header: () => t("dashboard.users.columnAllocatedVpnIp", "VPN IP"),
        cell: ({ row }) => row.original.allocatedVpnIp || t("common.na"),
      },
      {
        accessorKey: "lastConnectionTime",
        header: () => t("dashboard.users.columnLastConnection"),
        cell: ({ row }) =>
          row.original.lastConnectionTime
            ? new Date(row.original.lastConnectionTime).toLocaleString()
            : t("common.na"),
      },
      {
        accessorKey: "isOnline",
        header: () => t("dashboard.users.columnOnlineStatus"),
        cell: ({ row }) =>
          typeof row.original.isOnline === "boolean" ? (
            <Chip
              label={row.original.isOnline ? t("dashboard.users.statusOnline") : t("dashboard.users.statusOffline")}
              color={row.original.isOnline ? "success" : "default"}
              size="small"
            />
          ) : (
            t("common.na")
          ),
      },
      {
        accessorKey: "isPaused",
        header: () => t("dashboard.users.columnAccessState", "Access State"),
        cell: ({ row }) => (
          <Chip
            label={row.original.isPaused ? t("dashboard.users.statusPaused", "Paused") : t("dashboard.users.statusActive", "Active")}
            color={row.original.isPaused ? "error" : "success"}
            size="small"
          />
        ),
      },
      {
        accessorKey: "bytesReceived",
        header: () => t("dashboard.users.columnBytesReceived", "Bytes Received"),
        cell: ({ row }) => formatBytes(row.original.bytesReceived),
      },
      {
        accessorKey: "bytesSent",
        header: () => t("dashboard.users.columnBytesSent", "Bytes Sent"),
        cell: ({ row }) => formatBytes(row.original.bytesSent),
      },
      {
        id: "actions",
        header: () => t("dashboard.users.columnActions"),
        cell: ({ row }) => {
          const u = row.original;
          return (
            <div className="flex flex-wrap items-center justify-center gap-1">
              {canManageUsers && (
                u.isPaused ? (
                  <MuiButton size="small" variant="outlined" sx={{ height: 32, px: 1, minWidth: "auto" }} onClick={() => handleResumeUser(u.name)}>
                    {t("dashboard.users.resumeButton", "Resume")}
                  </MuiButton>
                ) : (
                  <MuiButton size="small" variant="outlined" sx={{ height: 32, px: 1, minWidth: "auto" }} onClick={() => handlePauseUser(u.name)}>
                    {t("dashboard.users.pauseButton", "Pause")}
                  </MuiButton>
                )
              )}
              {(currentUser?.role === UserRole.ADMIN ||
                currentUser?.role === UserRole.SUPERADMIN ||
                (currentUser?.role === UserRole.MANAGER && currentUser.departmentId === u.departmentId)) && (
                <MuiButton size="small" variant="outlined" sx={{ height: 32, px: 1, minWidth: "auto" }} onClick={() => handleEditClick(u)}>
                  {t("common.edit")}
                </MuiButton>
              )}
              {(currentUser?.role === UserRole.SUPERADMIN ||
                (currentUser?.role === UserRole.ADMIN && u.role !== UserRole.SUPERADMIN) ||
                (currentUser?.role === UserRole.MANAGER && u.departmentId === currentUser.departmentId && u.role !== UserRole.SUPERADMIN && u.role !== UserRole.ADMIN)) &&
                u.id !== currentUser?.id && (
                  <MuiButton size="small" variant="contained" color="error" sx={{ height: 32, px: 1, minWidth: "auto" }} onClick={() => handleDelete(u.id)}>
                    {t("dashboard.users.deleteButton")}
                  </MuiButton>
                )}
              <FormControl size="small" sx={{ minWidth: 80 }}>
                <MuiSelect
                  displayEmpty
                  value=""
                  onChange={(e) => handleDownload(u.name, e.target.value)}
                  sx={{ height: 32, fontSize: "0.8rem" }}
                  renderValue={() => t("dashboard.users.downloadConfigButtonShort", "DL Cfg")}
                >
                  <MenuItem value="windows">{t("dashboard.users.osWindows")}</MenuItem>
                  <MenuItem value="macos">{t("dashboard.users.osMacOS")}</MenuItem>
                  <MenuItem value="linux">{t("dashboard.users.osLinux")}</MenuItem>
                </MuiSelect>
              </FormControl>
            </div>
          );
        },
      },
    ],
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [t, departmentNameById, canManageUsers, currentUser, users]
  );

  const table = useReactTable({
    data: filteredUsers,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  return (
    <MainLayout className="p-4 space-y-6">
      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <h1 className="text-2xl font-bold">{t("dashboard.users.pageTitle")}</h1>
        {canManageUsers && (
          <MuiButton
            variant="contained"
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
          </MuiButton>
        )}
      </div>

      {/* Add User Dialog */}
      <MuiDialog open={addUserDialogOpen} onClose={() => setAddUserDialogOpen(false)}>
        <DialogTitle>
          {t("dashboard.users.addUserDialogTitle")}
        </DialogTitle>
        <DialogContent>
          <div className="space-y-4 py-2">
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.namePlaceholder")}</span>
              <TextField
                size="small"
                value={addUserForm.name}
                onChange={(e) => setAddUserForm({ ...addUserForm, name: e.target.value })}
                className="col-span-3"
                fullWidth
                sx={{ gridColumn: "span 3" }}
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.emailPlaceholder")}</span>
              <TextField
                size="small"
                type="email"
                value={addUserForm.email}
                onChange={(e) => setAddUserForm({ ...addUserForm, email: e.target.value })}
                fullWidth
                sx={{ gridColumn: "span 3" }}
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.passwordPlaceholder")}</span>
              <TextField
                size="small"
                type="password"
                value={addUserForm.password}
                onChange={(e) => setAddUserForm({ ...addUserForm, password: e.target.value })}
                fullWidth
                sx={{ gridColumn: "span 3" }}
              />
            </div>

            {canEditFixedIp && (
              <>
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">
                    {t("dashboard.users.fixedIpLabel", "Fixed VPN IP (Optional)")}
                  </span>
                  <TextField
                    size="small"
                    value={addUserForm.fixedIp || ""}
                    onChange={(e) => setAddUserForm({ ...addUserForm, fixedIp: e.target.value })}
                    placeholder={t("dashboard.users.fixedIpPlaceholder", "e.g., 10.8.0.100 or empty")}
                    fullWidth
                    sx={{ gridColumn: "span 3" }}
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">
                    {t("dashboard.users.subnetLabel", "Subnet (Optional)")}
                  </span>
                  <TextField
                    size="small"
                    value={addUserForm.subnet || ""}
                    onChange={(e) => setAddUserForm({ ...addUserForm, subnet: e.target.value })}
                    placeholder={t("dashboard.users.subnetPlaceholder", "e.g., 10.10.120.0/23 or empty")}
                    fullWidth
                    sx={{ gridColumn: "span 3" }}
                  />
                </div>
              </>
            )}

            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.roleLabel", "Role")}</span>
              <FormControl size="small" sx={{ gridColumn: "span 3", width: "100%" }}>
                <InputLabel>{t("dashboard.users.roleLabel", "Role")}</InputLabel>
                <MuiSelect
                  value={addUserForm.role}
                  label={t("dashboard.users.roleLabel", "Role")}
                  onChange={(e) => setAddUserForm({ ...addUserForm, role: e.target.value as UserRole })}
                  disabled={currentUser?.role === UserRole.MANAGER}
                >
                  {Object.values(UserRole)
                    .filter(
                      (role) =>
                        currentUser?.role === UserRole.SUPERADMIN ||
                        (currentUser?.role === UserRole.ADMIN && role !== UserRole.SUPERADMIN) ||
                        (currentUser?.role === UserRole.MANAGER && role === UserRole.USER)
                    )
                    .map((role) => (
                      <MenuItem key={role} value={role}>
                        {t(`dashboard.users.role${role.charAt(0).toUpperCase() + role.slice(1)}`, role)}
                      </MenuItem>
                    ))}
                </MuiSelect>
              </FormControl>
            </div>

            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.departmentLabel")}</span>
              <FormControl size="small" sx={{ gridColumn: "span 3", width: "100%" }}>
                <InputLabel>{t("dashboard.users.selectDepartmentPlaceholder")}</InputLabel>
                <MuiSelect
                  value={addUserForm.departmentId || ""}
                  label={t("dashboard.users.selectDepartmentPlaceholder")}
                  onChange={(e) => setAddUserForm({ ...addUserForm, departmentId: e.target.value })}
                  disabled={currentUser?.role === UserRole.MANAGER && !!currentUser.departmentId}
                >
                  <MenuItem value="">
                    {t("dashboard.users.selectDepartmentPlaceholder")}
                  </MenuItem>
                  {depts.map((d) => (
                    <MenuItem key={d.id} value={d.id}>
                      {d.name}
                    </MenuItem>
                  ))}
                </MuiSelect>
              </FormControl>
            </div>
          </div>
        </DialogContent>
        <DialogActions>
          <MuiButton variant="outlined" onClick={() => setAddUserDialogOpen(false)}>{t("common.cancel")}</MuiButton>
          <MuiButton variant="contained" onClick={handleCreateUser}>{t("common.create")}</MuiButton>
        </DialogActions>
      </MuiDialog>

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("dashboard.users.statsTotalUsers")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-semibold">{stats.total}</p>
            <p className="text-sm text-muted-foreground">
              {t("dashboard.users.statsTotalUsersHint", {
                count: stats.total,
              })}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("dashboard.users.statsOnline")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-semibold">{stats.online}</p>
            <p className="text-sm text-muted-foreground">
              {t("dashboard.users.statsOnlineHint", {
                count: stats.online,
              })}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("dashboard.users.statsPaused")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-semibold">{stats.paused}</p>
            <p className="text-sm text-muted-foreground">
              {t("dashboard.users.statsPausedHint", {
                count: stats.paused,
              })}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("dashboard.users.statsTraffic")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-semibold">
              {formatBytes(stats.totalTraffic, 1)}
            </p>
            <p className="text-sm text-muted-foreground">
              {t("dashboard.users.statsTrafficHint")}
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div className="grid w-full gap-3 md:grid-cols-2 xl:grid-cols-4 xl:max-w-4xl">
          <div className="flex flex-col gap-1">
            <TextField
              label={t("dashboard.users.searchLabel", "Search")}
              size="small"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder={t("dashboard.users.searchPlaceholder", "Search by name or email")}
              fullWidth
            />
          </div>
          <div className="flex flex-col gap-1">
            <FormControl size="small" fullWidth>
              <InputLabel>{t("dashboard.users.departmentFilterLabel")}</InputLabel>
              <MuiSelect
                value={departmentFilter}
                label={t("dashboard.users.departmentFilterLabel")}
                onChange={(e) => setDepartmentFilter(e.target.value)}
              >
                <MenuItem value="all">{t("dashboard.users.filterDepartmentAll")}</MenuItem>
                <MenuItem value="none">{t("dashboard.users.filterDepartmentNone")}</MenuItem>
                {depts.map((d) => (
                  <MenuItem key={d.id} value={d.id}>
                    {d.name}
                  </MenuItem>
                ))}
              </MuiSelect>
            </FormControl>
          </div>
          <div className="flex flex-col gap-1">
            <FormControl size="small" fullWidth>
              <InputLabel>{t("dashboard.users.statusFilterLabel")}</InputLabel>
              <MuiSelect
                value={statusFilter}
                label={t("dashboard.users.statusFilterLabel")}
                onChange={(e) => setStatusFilter(e.target.value)}
              >
                <MenuItem value="all">{t("dashboard.users.filterStatusAll")}</MenuItem>
                <MenuItem value="online">{t("dashboard.users.filterStatusOnline")}</MenuItem>
                <MenuItem value="offline">{t("dashboard.users.filterStatusOffline")}</MenuItem>
                <MenuItem value="paused">{t("dashboard.users.filterStatusPaused")}</MenuItem>
              </MuiSelect>
            </FormControl>
          </div>
        </div>
        {hasFilters && (
          <MuiButton
            variant="text"
            className="self-start lg:self-auto"
            onClick={handleClearFilters}
          >
            {t("common.clearFilters", "Clear filters")}
          </MuiButton>
        )}
      </div>

      <div className="flex items-center justify-between text-sm text-muted-foreground">
        <span>
          {t("dashboard.users.resultsSummary", {
            count: filteredUsers.length,
            total: users.length,
          })}
        </span>
      </div>

      {/* Edit User Dialog */}
      <MuiDialog open={editUserDialogOpen} onClose={() => setEditUserDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {t("dashboard.users.editUserDialogTitle", "Edit User")}
        </DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ mb: 2 }}>
            {t(
              "dashboard.users.editUserDescription",
              "Make changes to the user profile here. Click save when you're done."
            )}
          </DialogContentText>
          <div className="grid gap-4 py-2">
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.namePlaceholder")}</span>
              <TextField
                size="small"
                value={editForm.name || ""}
                onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                fullWidth
                sx={{ gridColumn: "span 3" }}
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.emailPlaceholder")}</span>
              <TextField
                size="small"
                type="email"
                value={editForm.email || ""}
                onChange={(e) => setEditForm({ ...editForm, email: e.target.value })}
                fullWidth
                sx={{ gridColumn: "span 3" }}
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">
                {t("dashboard.users.passwordOptionalPlaceholder", "Password (optional)")}
              </span>
              <TextField
                size="small"
                type="password"
                value={editForm.password || ""}
                onChange={(e) => setEditForm({ ...editForm, password: e.target.value })}
                fullWidth
                placeholder={t("dashboard.users.passwordLeaveBlankPlaceholder", "Leave blank to keep current")}
                sx={{ gridColumn: "span 3" }}
              />
            </div>

            {canEditFixedIp && (
              <>
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">
                    {t("dashboard.users.fixedIpLabel", "Fixed VPN IP (Optional)")}
                  </span>
                  <TextField
                    size="small"
                    value={editForm.fixedIp || ""}
                    onChange={(e) => setEditForm({ ...editForm, fixedIp: e.target.value })}
                    placeholder={t("dashboard.users.fixedIpPlaceholder", "e.g., 10.8.0.100 or empty")}
                    fullWidth
                    sx={{ gridColumn: "span 3" }}
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">
                    {t("dashboard.users.subnetLabel", "Subnet (Optional)")}
                  </span>
                  <TextField
                    size="small"
                    value={editForm.subnet || ""}
                    onChange={(e) => setEditForm({ ...editForm, subnet: e.target.value })}
                    placeholder={t("dashboard.users.subnetPlaceholder", "e.g., 10.10.120.0/23 or empty")}
                    fullWidth
                    sx={{ gridColumn: "span 3" }}
                  />
                </div>
              </>
            )}
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.departmentLabel")}</span>
              <FormControl size="small" sx={{ gridColumn: "span 3", width: "100%" }}>
                <InputLabel>{t("dashboard.users.selectDepartmentPlaceholder")}</InputLabel>
                <MuiSelect
                  value={editForm.departmentId || ""}
                  label={t("dashboard.users.selectDepartmentPlaceholder")}
                  onChange={(e) => setEditForm({ ...editForm, departmentId: e.target.value })}
                  disabled={
                    currentUser?.role === UserRole.MANAGER &&
                    editingUser?.departmentId !== currentUser.departmentId &&
                    editingUser?.id !== currentUser.id
                  }
                >
                  <MenuItem value="">
                    {t("dashboard.users.selectDepartmentPlaceholder")}
                  </MenuItem>
                  {depts.map((d) => (
                    <MenuItem key={d.id} value={d.id}>
                      {d.name}
                    </MenuItem>
                  ))}
                </MuiSelect>
              </FormControl>
            </div>
            {(currentUser?.role === UserRole.SUPERADMIN ||
              (currentUser?.role === UserRole.ADMIN &&
                editingUser?.role !== UserRole.SUPERADMIN)) &&
              (editingUser?.id !== currentUser?.id ||
                currentUser?.role === UserRole.SUPERADMIN) && (
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">{t("dashboard.users.roleLabel", "Role")}</span>
                  <FormControl size="small" sx={{ gridColumn: "span 3", width: "100%" }}>
                    <InputLabel>{t("dashboard.users.roleLabel", "Role")}</InputLabel>
                    <MuiSelect
                      value={editForm.role || ""}
                      label={t("dashboard.users.roleLabel", "Role")}
                      onChange={(e) => setEditForm({ ...editForm, role: e.target.value as UserRole })}
                    >
                      {Object.values(UserRole)
                        .filter(
                          (role) =>
                            currentUser?.role === UserRole.SUPERADMIN ||
                            role !== UserRole.SUPERADMIN
                        )
                        .map((role) => (
                          <MenuItem key={role} value={role}>
                            {t(`dashboard.users.role${role.charAt(0).toUpperCase() + role.slice(1)}`, role)}
                          </MenuItem>
                        ))}
                    </MuiSelect>
                  </FormControl>
                </div>
              )}
          </div>
        </DialogContent>
        <DialogActions>
          <MuiButton variant="outlined" onClick={() => setEditUserDialogOpen(false)}>{t("common.cancel")}</MuiButton>
          <MuiButton variant="contained" onClick={handleUpdateUser}>
            {t("common.saveChanges")}
          </MuiButton>
        </DialogActions>
      </MuiDialog>

      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.users.listTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>{t("common.loading")}</p>
          ) : filteredUsers.length === 0 ? (
            <p className="py-10 text-center text-muted-foreground">{t("dashboard.users.noResults")}</p>
          ) : (
            <>
              {/* ── Mobile: card list (hidden on md+) ───────────────────── */}
              <div className="md:hidden space-y-3">
                {table.getRowModel().rows.map((row) => {
                  const u = row.original;
                  return (
                    <div key={u.id} className="border rounded-xl p-4 space-y-3 bg-card shadow-sm">
                      {/* Header row */}
                      <div className="flex items-start justify-between gap-2">
                        <div>
                          <div className="font-semibold text-sm">{u.name}</div>
                          <div className="text-xs text-muted-foreground mt-0.5">{u.email}</div>
                        </div>
                        <div className="flex flex-col items-end gap-1 shrink-0">
                          <Chip
                            label={u.isOnline ? t("dashboard.users.statusOnline") : t("dashboard.users.statusOffline")}
                            color={u.isOnline ? "success" : "default"}
                            size="small"
                          />
                          <Chip
                            label={u.isPaused ? t("dashboard.users.statusPaused", "Paused") : t("dashboard.users.statusActive", "Active")}
                            color={u.isPaused ? "error" : "default"}
                            variant="outlined"
                            size="small"
                          />
                        </div>
                      </div>
                      {/* Meta row */}
                      <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-muted-foreground">
                        <span><span className="font-medium text-foreground">{t("dashboard.users.columnRole")}:</span> {u.role}</span>
                        <span><span className="font-medium text-foreground">{t("dashboard.users.columnDepartment")}:</span> {u.departmentId ? departmentNameById[u.departmentId] || "-" : "-"}</span>
                        {u.fixedIp && <span><span className="font-medium text-foreground">IP:</span> {u.fixedIp}</span>}
                        {u.allocatedVpnIp && <span><span className="font-medium text-foreground">VPN IP:</span> {u.allocatedVpnIp}</span>}
                        <span><span className="font-medium text-foreground">↓</span> {formatBytes(u.bytesReceived)}</span>
                        <span><span className="font-medium text-foreground">↑</span> {formatBytes(u.bytesSent)}</span>
                      </div>
                      {/* Actions */}
                      <div className="flex flex-wrap gap-2 pt-1 border-t">
                        {canManageUsers && (
                          u.isPaused ? (
                            <MuiButton size="small" variant="outlined" sx={{ height: 28, fontSize: "0.75rem", px: 1 }} onClick={() => handleResumeUser(u.name)}>
                              {t("dashboard.users.resumeButton", "Resume")}
                            </MuiButton>
                          ) : (
                            <MuiButton size="small" variant="outlined" sx={{ height: 28, fontSize: "0.75rem", px: 1 }} onClick={() => handlePauseUser(u.name)}>
                              {t("dashboard.users.pauseButton", "Pause")}
                            </MuiButton>
                          )
                        )}
                        {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN ||
                          (currentUser?.role === UserRole.MANAGER && currentUser.departmentId === u.departmentId)) && (
                          <MuiButton size="small" variant="outlined" sx={{ height: 28, fontSize: "0.75rem", px: 1 }} onClick={() => handleEditClick(u)}>
                            {t("common.edit")}
                          </MuiButton>
                        )}
                        {(currentUser?.role === UserRole.SUPERADMIN ||
                          (currentUser?.role === UserRole.ADMIN && u.role !== UserRole.SUPERADMIN) ||
                          (currentUser?.role === UserRole.MANAGER && u.departmentId === currentUser.departmentId && u.role !== UserRole.SUPERADMIN && u.role !== UserRole.ADMIN)) &&
                          u.id !== currentUser?.id && (
                            <MuiButton size="small" variant="contained" color="error" sx={{ height: 28, fontSize: "0.75rem", px: 1 }} onClick={() => handleDelete(u.id)}>
                              {t("dashboard.users.deleteButton")}
                            </MuiButton>
                          )}
                        <FormControl size="small" sx={{ minWidth: 72 }}>
                          <MuiSelect
                            displayEmpty
                            value=""
                            onChange={(e) => handleDownload(u.name, e.target.value)}
                            sx={{ height: 28, fontSize: "0.75rem" }}
                            renderValue={() => t("dashboard.users.downloadConfigButtonShort", "DL Cfg")}
                          >
                            <MenuItem value="windows">{t("dashboard.users.osWindows")}</MenuItem>
                            <MenuItem value="macos">{t("dashboard.users.osMacOS")}</MenuItem>
                            <MenuItem value="linux">{t("dashboard.users.osLinux")}</MenuItem>
                          </MuiSelect>
                        </FormControl>
                      </div>
                    </div>
                  );
                })}
              </div>

              {/* ── Desktop: TanStack table (hidden below md) ────────────── */}
              <div className="hidden md:block overflow-x-auto">
                <table className="w-full caption-bottom text-sm border-collapse">
                  <thead>
                    {table.getHeaderGroups().map((hg) => (
                      <tr key={hg.id} className="border-b">
                        {hg.headers.map((header) => (
                          <th
                            key={header.id}
                            className={[
                              "h-12 px-4 text-left align-middle font-medium text-muted-foreground whitespace-nowrap",
                              header.column.getCanSort() ? "cursor-pointer select-none" : "",
                              header.id === "actions" ? "sticky right-0 min-w-[260px]" : "min-w-[120px]",
                            ].join(" ")}
                            style={header.id === "actions" ? { backgroundColor: "hsl(var(--card))", boxShadow: "inset 10px 0 0px -9px #0505050f" } : undefined}
                            onClick={header.column.getToggleSortingHandler()}
                          >
                            {flexRender(header.column.columnDef.header, header.getContext())}
                            {header.column.getIsSorted() === "asc" ? " ↑" : header.column.getIsSorted() === "desc" ? " ↓" : ""}
                          </th>
                        ))}
                      </tr>
                    ))}
                  </thead>
                  <tbody>
                    {table.getRowModel().rows.map((row) => (
                      <tr key={row.id} className="border-b transition-colors hover:bg-muted/50">
                        {row.getVisibleCells().map((cell) => (
                          <td
                            key={cell.id}
                            className={[
                              "px-4 py-3 align-middle",
                              cell.column.id === "actions" ? "sticky right-0" : "",
                            ].join(" ")}
                            style={cell.column.id === "actions" ? { backgroundColor: "hsl(var(--card))", boxShadow: "inset 10px 0 0px -9px #0505050f" } : undefined}
                          >
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
