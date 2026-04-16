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
  getPaginationRowModel,
  flexRender,
  type ColumnDef,
  type SortingState,
  type PaginationState,
} from "@tanstack/react-table";
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
import { Badge } from "@/components/ui/badge";
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

// Helper function to format bytes into a readable string
const formatBytes = (bytes?: number, decimals = 2): string => {
  if (bytes === undefined || bytes === null || bytes === 0) return "0 Bytes";
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
};

const initialEditFormState: UserUpdateRequest = {
  name: "",
  email: "",
  role: UserRole.USER,
  departmentId: "",
  fixedIp: "",
  subnet: "",
  password: "",
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

  const [addUserForm, setAddUserForm] = useState<UserUpdateRequest>({
    name: "",
    email: "",
    password: "",
    role: UserRole.USER,
    departmentId: "",
    fixedIp: "",
    subnet: "",
  });
  const [addUserDialogOpen, setAddUserDialogOpen] = useState(false);

  const [editUserDialogOpen, setEditUserDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<AdminUser | null>(null);
  const [editForm, setEditForm] = useState<UserUpdateRequest>(initialEditFormState);
  const [confirmState, setConfirmState] = useState<{
    open: boolean;
    title: string;
    description: string;
    onConfirm: () => void;
    destructive?: boolean;
  }>({ open: false, title: "", description: "", onConfirm: () => {} });

  const handlePauseUser = (username: string) => {
    setConfirmState({
      open: true,
      title: t("dashboard.users.pauseConfirmTitle", "Pause User"),
      description: t("dashboard.users.pauseConfirm", `Are you sure you want to pause user ${username}?`),
      onConfirm: async () => {
        try {
          await userManagementAPI.pauseUser(username);
          toast.success(t("dashboard.users.pauseSuccess", `User ${username} paused successfully.`));
          fetchAll();
        } catch (error: any) {
          toast.error(error?.response?.data?.error || t("dashboard.users.pauseError", `Failed to pause user ${username}.`));
        }
      },
    });
  };

  const handleResumeUser = (username: string) => {
    setConfirmState({
      open: true,
      title: t("dashboard.users.resumeConfirmTitle", "Resume User"),
      description: t("dashboard.users.resumeConfirm", `Are you sure you want to resume user ${username}?`),
      onConfirm: async () => {
        try {
          await userManagementAPI.resumeUser(username);
          toast.success(t("dashboard.users.resumeSuccess", `User ${username} resumed successfully.`));
          fetchAll();
        } catch (error: any) {
          toast.error(error?.response?.data?.error || t("dashboard.users.resumeError", `Failed to resume user ${username}.`));
        }
      },
    });
  };

  const fetchAll = useCallback(async () => {
    setLoading(true);
    try {
      const [u, d] = await Promise.all([userManagementAPI.list(), departmentAPI.list()]);
      setUsers(u);
      setDepts(d);
    } catch {
      toast.error(t("dashboard.users.loadError"));
    } finally {
      setLoading(false);
    }
  }, [t]);

  useEffect(() => {
    fetchAll();
  }, [fetchAll]);

  const handleDownload = async (username: string, os: string) => {
    try {
      const data = await openvpnAPI.getClientConfig(username, os);
      const config = data.config;
      const blob = new Blob([config], { type: "application/x-openvpn-profile" });
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
      if (!(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN) && (payload.fixedIp || payload.subnet)) {
        payload.fixedIp = "";
        payload.subnet = "";
      }
      if (currentUser?.role === UserRole.MANAGER && !payload.departmentId) {
        payload.departmentId = currentUser.departmentId || "";
      }
      if (payload.fixedIp === "") payload.fixedIp = null;
      if (payload.subnet === "") payload.subnet = null;
      await userManagementAPI.create(payload as Partial<AdminUser> & { password: string });
      toast.success(t("dashboard.users.createSuccess"));
      setAddUserDialogOpen(false);
      fetchAll();
    } catch (error: any) {
      toast.error(error?.response?.data?.error || t("dashboard.users.createError"));
    }
  };

  const handleDelete = (id: string) => {
    setConfirmState({
      open: true,
      title: t("dashboard.users.deleteConfirmTitle", "Delete User"),
      description: t("dashboard.users.deleteConfirm"),
      destructive: true,
      onConfirm: async () => {
        try {
          await userManagementAPI.delete(id);
          toast.success(t("dashboard.users.deleteSuccess"));
          fetchAll();
        } catch {
          toast.error(t("dashboard.users.deleteError"));
        }
      },
    });
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
    if (!updatePayload.password?.trim()) delete updatePayload.password;
    if (!(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN)) {
      delete updatePayload.fixedIp;
      delete updatePayload.subnet;
    } else {
      if (updatePayload.fixedIp === "") updatePayload.fixedIp = null;
      if (updatePayload.subnet === "") updatePayload.subnet = null;
    }
    try {
      await userManagementAPI.update(editingUser.id, updatePayload);
      toast.success(t("dashboard.users.editUserSuccess", "User updated successfully!"));
      setEditUserDialogOpen(false);
      fetchAll();
    } catch (error: any) {
      toast.error(error?.response?.data?.error || t("dashboard.users.editUserError", "Failed to update user."));
    }
  };

  const canEditFixedIp = currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN;
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
    const totalTraffic = users.reduce((sum, u) => sum + (u.bytesReceived || 0) + (u.bytesSent || 0), 0);
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
          (statusFilter === "offline" && !user.isPaused && user.isOnline === false);
        return matchesSearch && matchesDepartment && matchesStatus;
      })
      .sort((a, b) => (a.name || "").localeCompare(b.name || ""));
  }, [users, normalizedSearch, departmentFilter, statusFilter]);

  const hasFilters = normalizedSearch.length > 0 || departmentFilter !== "all" || statusFilter !== "all";

  const handleClearFilters = () => {
    setSearchTerm("");
    setDepartmentFilter("all");
    setStatusFilter("all");
  };

  const [sorting, setSorting] = useState<SortingState>([]);
  const [pagination, setPagination] = useState<PaginationState>({ pageIndex: 0, pageSize: 25 });

  const DownloadSelect = ({ username, size = "sm" }: { username: string; size?: "sm" | "xs" }) => (
    <Select value="" onValueChange={(val) => handleDownload(username, val)}>
      <SelectTrigger className={size === "xs" ? "h-7 w-20 text-xs" : "h-8 w-20 text-xs"}>
        <SelectValue placeholder={t("dashboard.users.downloadConfigButtonShort", "DL Cfg")} />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="windows">{t("dashboard.users.osWindows")}</SelectItem>
        <SelectItem value="macos">{t("dashboard.users.osMacOS")}</SelectItem>
        <SelectItem value="linux">{t("dashboard.users.osLinux")}</SelectItem>
      </SelectContent>
    </Select>
  );

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
          <Badge variant="secondary" className="capitalize">
            {t(`dashboard.users.role${row.original.role.charAt(0).toUpperCase() + row.original.role.slice(1)}`, row.original.role)}
          </Badge>
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
            <Badge variant={row.original.isOnline ? "success" : "secondary"}>
              {row.original.isOnline ? t("dashboard.users.statusOnline") : t("dashboard.users.statusOffline")}
            </Badge>
          ) : (
            t("common.na")
          ),
      },
      {
        accessorKey: "isPaused",
        header: () => t("dashboard.users.columnAccessState", "Access State"),
        cell: ({ row }) => (
          <Badge variant={row.original.isPaused ? "destructive" : "success"}>
            {row.original.isPaused ? t("dashboard.users.statusPaused", "Paused") : t("dashboard.users.statusActive", "Active")}
          </Badge>
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
                  <Button size="sm" variant="outline" className="h-8 px-2" onClick={() => handleResumeUser(u.name)}>
                    {t("dashboard.users.resumeButton", "Resume")}
                  </Button>
                ) : (
                  <Button size="sm" variant="outline" className="h-8 px-2" onClick={() => handlePauseUser(u.name)}>
                    {t("dashboard.users.pauseButton", "Pause")}
                  </Button>
                )
              )}
              {(currentUser?.role === UserRole.ADMIN ||
                currentUser?.role === UserRole.SUPERADMIN ||
                (currentUser?.role === UserRole.MANAGER && currentUser.departmentId === u.departmentId)) && (
                <Button size="sm" variant="outline" className="h-8 px-2" onClick={() => handleEditClick(u)}>
                  {t("common.edit")}
                </Button>
              )}
              {(currentUser?.role === UserRole.SUPERADMIN ||
                (currentUser?.role === UserRole.ADMIN && u.role !== UserRole.SUPERADMIN) ||
                (currentUser?.role === UserRole.MANAGER && u.departmentId === currentUser.departmentId && u.role !== UserRole.SUPERADMIN && u.role !== UserRole.ADMIN)) &&
                u.id !== currentUser?.id && (
                  <Button size="sm" variant="destructive" className="h-8 px-2" onClick={() => handleDelete(u.id)}>
                    {t("dashboard.users.deleteButton")}
                  </Button>
                )}
              <DownloadSelect username={u.name} />
            </div>
          );
        },
      },
    ],
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [t, departmentNameById, canManageUsers, currentUser, users]
  );

  useEffect(() => {
    setPagination((p) => ({ ...p, pageIndex: 0 }));
  }, [filteredUsers.length]);

  const table = useReactTable({
    data: filteredUsers,
    columns,
    state: { sorting, pagination },
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    manualPagination: false,
  });

  // ── Regular-user self view ─────────────────────────────────────────────
  if (currentUser?.role === UserRole.USER) {
    const me = users[0];
    return (
      <MainLayout className="p-4 space-y-6">
        <h1 className="text-2xl font-bold">{t("dashboard.users.pageTitle")}</h1>

        {loading || !me ? (
          <p>{t("common.loading")}</p>
        ) : (
          <>
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">{t("dashboard.users.namePlaceholder")}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-2xl font-semibold">{me.name}</p>
                  <p className="text-sm text-muted-foreground">{me.email}</p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">{t("dashboard.users.columnOnlineStatus")}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-2xl font-semibold">
                    {me.isOnline ? t("dashboard.users.statusOnline") : t("dashboard.users.statusOffline")}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {me.lastConnectionTime ? new Date(me.lastConnectionTime).toLocaleString() : t("common.na")}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">VPN IP</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-2xl font-semibold">{me.allocatedVpnIp || t("common.na")}</p>
                  <p className="text-sm text-muted-foreground">
                    {t("dashboard.users.columnConnectionIp", "Connection IP")}: {me.connectionIp || t("common.na")}
                  </p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">{t("dashboard.users.statsTraffic")}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-2xl font-semibold">{formatBytes((me.bytesReceived || 0) + (me.bytesSent || 0), 1)}</p>
                  <p className="text-sm text-muted-foreground">↓ {formatBytes(me.bytesReceived)} · ↑ {formatBytes(me.bytesSent)}</p>
                </CardContent>
              </Card>
            </div>

            <Card>
              <CardHeader>
                <CardTitle>{t("dashboard.users.downloadConfigButton")}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-3">
                  {(["windows", "macos", "linux"] as const).map((os) => (
                    <Button key={os} size="lg" onClick={() => handleDownload(me.name, os)}>
                      {os === "windows" ? (
                        <svg className="mr-2" width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-13.051-1.8"/></svg>
                      ) : os === "macos" ? (
                        <svg className="mr-2" width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12.152 6.896c-.948 0-2.415-1.078-3.96-1.04-2.04.027-3.91 1.183-4.961 3.014-2.117 3.675-.546 9.103 1.519 12.09 1.013 1.454 2.208 3.09 3.792 3.039 1.52-.065 2.09-.987 3.935-.987 1.831 0 2.35.987 3.96.948 1.637-.026 2.676-1.48 3.676-2.948 1.156-1.688 1.636-3.325 1.662-3.415-.039-.013-3.182-1.221-3.22-4.857-.026-3.04 2.48-4.494 2.597-4.559-1.429-2.09-3.623-2.324-4.39-2.376-2-.156-3.675 1.09-4.61 1.09zM15.53 3.83c.843-1.012 1.4-2.427 1.245-3.83-1.207.052-2.662.805-3.532 1.818-.78.896-1.454 2.338-1.273 3.714 1.338.104 2.715-.688 3.559-1.701"/></svg>
                      ) : (
                        <svg className="mr-2" width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12.504 0c-.155 0-.315.008-.48.021-4.226.333-3.105 4.807-3.17 6.298-.076 1.092-.3 1.953-1.05 3.02-.885 1.051-2.127 2.75-2.716 4.521-.278.832-.41 1.684-.287 2.489a5.745 5.745 0 00.762 2.211c1.083 1.67 3.045 2.89 5.266 3.259.073.014.147.024.22.034.111.017.218.03.327.043a12.1 12.1 0 001.17.082c4.886 0 9.23-2.879 10.342-6.717.8-2.772-.267-5.774-2.887-7.664 1.005 2.3.474 4.85-.98 5.822-1.477.983-3.218.37-4.296-.854a9.387 9.387 0 01-.725-1.05 9.08 9.08 0 01-.638-1.491c-.304-.908-.39-1.933-.188-2.9C14.26 1.32 13.578 0 12.504 0zM7.03 4.282c-.53.197-.97.508-1.298.916C4.9 6.22 4.38 7.64 4.38 8.998c0 2.35 1.193 4.432 2.918 5.633-.197-.63-.267-1.29-.209-1.94-.35-.4-.648-.843-.875-1.33C5.793 10.39 5.615 9.24 5.85 8.127c.17-.816.512-1.59 1.01-2.244C7.11 5.6 7.36 5.2 7.03 4.282z"/></svg>
                      )}
                      {os === "windows" ? t("dashboard.users.osWindows") : os === "macos" ? t("dashboard.users.osMacOS") : t("dashboard.users.osLinux")}
                    </Button>
                  ))}
                </div>
              </CardContent>
            </Card>
          </>
        )}
      </MainLayout>
    );
  }

  return (
    <MainLayout className="p-4 space-y-6">
      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <h1 className="text-2xl font-bold">{t("dashboard.users.pageTitle")}</h1>
        {canManageUsers && (
          <Button
            onClick={() => {
              let initialDeptId = "";
              const initialRole = UserRole.USER;
              if (currentUser?.role === UserRole.MANAGER) {
                initialDeptId = currentUser.departmentId || "";
              }
              setAddUserForm({ name: "", email: "", password: "", role: initialRole, departmentId: initialDeptId, fixedIp: "", subnet: "" });
              setAddUserDialogOpen(true);
            }}
          >
            {t("dashboard.users.addUserButton")}
          </Button>
        )}
      </div>

      {/* Add User Dialog */}
      <Dialog open={addUserDialogOpen} onOpenChange={setAddUserDialogOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>{t("dashboard.users.addUserDialogTitle")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-2">
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.namePlaceholder")}</span>
              <div className="col-span-3">
                <Input value={addUserForm.name} onChange={(e) => setAddUserForm({ ...addUserForm, name: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.emailPlaceholder")}</span>
              <div className="col-span-3">
                <Input type="email" value={addUserForm.email} onChange={(e) => setAddUserForm({ ...addUserForm, email: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.passwordPlaceholder")}</span>
              <div className="col-span-3">
                <Input type="password" value={addUserForm.password} onChange={(e) => setAddUserForm({ ...addUserForm, password: e.target.value })} />
              </div>
            </div>
            {canEditFixedIp && (
              <>
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">{t("dashboard.users.fixedIpLabel", "Fixed VPN IP (Optional)")}</span>
                  <div className="col-span-3">
                    <Input value={addUserForm.fixedIp || ""} onChange={(e) => setAddUserForm({ ...addUserForm, fixedIp: e.target.value })} placeholder={t("dashboard.users.fixedIpPlaceholder", "e.g., 10.8.0.100 or empty")} />
                  </div>
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">{t("dashboard.users.subnetLabel", "Subnet (Optional)")}</span>
                  <div className="col-span-3">
                    <Input value={addUserForm.subnet || ""} onChange={(e) => setAddUserForm({ ...addUserForm, subnet: e.target.value })} placeholder={t("dashboard.users.subnetPlaceholder", "e.g., 10.10.120.0/23 or empty")} />
                  </div>
                </div>
              </>
            )}
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.roleLabel", "Role")}</span>
              <div className="col-span-3">
                <Select
                  value={addUserForm.role}
                  onValueChange={(val) => setAddUserForm({ ...addUserForm, role: val as UserRole })}
                  disabled={currentUser?.role === UserRole.MANAGER}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("dashboard.users.roleLabel", "Role")} />
                  </SelectTrigger>
                  <SelectContent>
                    {Object.values(UserRole)
                      .filter((role) =>
                        currentUser?.role === UserRole.SUPERADMIN ||
                        (currentUser?.role === UserRole.ADMIN && role !== UserRole.SUPERADMIN) ||
                        (currentUser?.role === UserRole.MANAGER && role === UserRole.USER)
                      )
                      .map((role) => (
                        <SelectItem key={role} value={role}>
                          {t(`dashboard.users.role${role.charAt(0).toUpperCase() + role.slice(1)}`, role)}
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.departmentLabel")}</span>
              <div className="col-span-3">
                <Select
                  value={addUserForm.departmentId || "__none__"}
                  onValueChange={(val) => setAddUserForm({ ...addUserForm, departmentId: val === "__none__" ? "" : val })}
                  disabled={currentUser?.role === UserRole.MANAGER && !!currentUser.departmentId}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("dashboard.users.selectDepartmentPlaceholder")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none__">{t("dashboard.users.selectDepartmentPlaceholder")}</SelectItem>
                    {depts.map((d) => (
                      <SelectItem key={d.id} value={d.id}>{d.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setAddUserDialogOpen(false)}>{t("common.cancel")}</Button>
            <Button onClick={handleCreateUser}>{t("common.create")}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">{t("dashboard.users.statsTotalUsers")}</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-semibold">{stats.total}</p>
            <p className="text-sm text-muted-foreground">{t("dashboard.users.statsTotalUsersHint", { count: stats.total })}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">{t("dashboard.users.statsOnline")}</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-semibold">{stats.online}</p>
            <p className="text-sm text-muted-foreground">{t("dashboard.users.statsOnlineHint", { count: stats.online })}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">{t("dashboard.users.statsPaused")}</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-semibold">{stats.paused}</p>
            <p className="text-sm text-muted-foreground">{t("dashboard.users.statsPausedHint", { count: stats.paused })}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">{t("dashboard.users.statsTraffic")}</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-semibold">{formatBytes(stats.totalTraffic, 1)}</p>
            <p className="text-sm text-muted-foreground">{t("dashboard.users.statsTrafficHint")}</p>
          </CardContent>
        </Card>
      </div>

      <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div className="grid w-full gap-3 md:grid-cols-2 xl:grid-cols-4 xl:max-w-4xl">
          <Input
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder={t("dashboard.users.searchPlaceholder", "Search by name or email")}
          />
          <Select value={departmentFilter} onValueChange={setDepartmentFilter}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder={t("dashboard.users.departmentFilterLabel")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("dashboard.users.filterDepartmentAll")}</SelectItem>
              <SelectItem value="none">{t("dashboard.users.filterDepartmentNone")}</SelectItem>
              {depts.map((d) => (
                <SelectItem key={d.id} value={d.id}>{d.name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder={t("dashboard.users.statusFilterLabel")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("dashboard.users.filterStatusAll")}</SelectItem>
              <SelectItem value="online">{t("dashboard.users.filterStatusOnline")}</SelectItem>
              <SelectItem value="offline">{t("dashboard.users.filterStatusOffline")}</SelectItem>
              <SelectItem value="paused">{t("dashboard.users.filterStatusPaused")}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        {hasFilters && (
          <Button variant="ghost" className="self-start lg:self-auto" onClick={handleClearFilters}>
            {t("common.clearFilters", "Clear filters")}
          </Button>
        )}
      </div>

      <div className="flex items-center justify-between text-sm text-muted-foreground">
        <span>{t("dashboard.users.resultsSummary", { count: filteredUsers.length, total: users.length })}</span>
      </div>

      {/* Edit User Dialog */}
      <Dialog open={editUserDialogOpen} onOpenChange={setEditUserDialogOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>{t("dashboard.users.editUserDialogTitle", "Edit User")}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            {t("dashboard.users.editUserDescription", "Make changes to the user profile here. Click save when you're done.")}
          </p>
          <div className="grid gap-4 py-2">
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.namePlaceholder")}</span>
              <div className="col-span-3">
                <Input value={editForm.name || ""} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.emailPlaceholder")}</span>
              <div className="col-span-3">
                <Input type="email" value={editForm.email || ""} onChange={(e) => setEditForm({ ...editForm, email: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.passwordOptionalPlaceholder", "Password (optional)")}</span>
              <div className="col-span-3">
                <Input type="password" value={editForm.password || ""} onChange={(e) => setEditForm({ ...editForm, password: e.target.value })} placeholder={t("dashboard.users.passwordLeaveBlankPlaceholder", "Leave blank to keep current")} />
              </div>
            </div>
            {canEditFixedIp && (
              <>
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">{t("dashboard.users.fixedIpLabel", "Fixed VPN IP (Optional)")}</span>
                  <div className="col-span-3">
                    <Input value={editForm.fixedIp || ""} onChange={(e) => setEditForm({ ...editForm, fixedIp: e.target.value })} placeholder={t("dashboard.users.fixedIpPlaceholder", "e.g., 10.8.0.100 or empty")} />
                  </div>
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">{t("dashboard.users.subnetLabel", "Subnet (Optional)")}</span>
                  <div className="col-span-3">
                    <Input value={editForm.subnet || ""} onChange={(e) => setEditForm({ ...editForm, subnet: e.target.value })} placeholder={t("dashboard.users.subnetPlaceholder", "e.g., 10.10.120.0/23 or empty")} />
                  </div>
                </div>
              </>
            )}
            <div className="grid grid-cols-4 items-center gap-4">
              <span className="text-right text-sm font-medium">{t("dashboard.users.departmentLabel")}</span>
              <div className="col-span-3">
                <Select
                  value={editForm.departmentId || "__none__"}
                  onValueChange={(val) => setEditForm({ ...editForm, departmentId: val === "__none__" ? "" : val })}
                  disabled={
                    currentUser?.role === UserRole.MANAGER &&
                    editingUser?.departmentId !== currentUser.departmentId &&
                    editingUser?.id !== currentUser.id
                  }
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={t("dashboard.users.selectDepartmentPlaceholder")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none__">{t("dashboard.users.selectDepartmentPlaceholder")}</SelectItem>
                    {depts.map((d) => (
                      <SelectItem key={d.id} value={d.id}>{d.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            {(currentUser?.role === UserRole.SUPERADMIN ||
              (currentUser?.role === UserRole.ADMIN && editingUser?.role !== UserRole.SUPERADMIN)) &&
              (editingUser?.id !== currentUser?.id || currentUser?.role === UserRole.SUPERADMIN) && (
                <div className="grid grid-cols-4 items-center gap-4">
                  <span className="text-right text-sm font-medium">{t("dashboard.users.roleLabel", "Role")}</span>
                  <div className="col-span-3">
                    <Select
                      value={editForm.role || ""}
                      onValueChange={(val) => setEditForm({ ...editForm, role: val as UserRole })}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder={t("dashboard.users.roleLabel", "Role")} />
                      </SelectTrigger>
                      <SelectContent>
                        {Object.values(UserRole)
                          .filter((role) => currentUser?.role === UserRole.SUPERADMIN || role !== UserRole.SUPERADMIN)
                          .map((role) => (
                            <SelectItem key={role} value={role}>
                              {t(`dashboard.users.role${role.charAt(0).toUpperCase() + role.slice(1)}`, role)}
                            </SelectItem>
                          ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              )}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditUserDialogOpen(false)}>{t("common.cancel")}</Button>
            <Button onClick={handleUpdateUser}>{t("common.saveChanges")}</Button>
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
            <>
              {/* ── Mobile: card list (hidden on md+) ───────────────────── */}
              <div className="md:hidden space-y-3">
                {table.getRowModel().rows.length === 0 ? (
                  <p className="py-10 text-center text-muted-foreground">{t("dashboard.users.noResults")}</p>
                ) : table.getRowModel().rows.map((row) => {
                  const u = row.original;
                  return (
                    <div key={u.id} className="border rounded-xl p-4 space-y-3 bg-card shadow-sm">
                      <div className="flex items-start justify-between gap-2">
                        <div>
                          <div className="font-semibold text-sm">{u.name}</div>
                          <div className="text-xs text-muted-foreground mt-0.5">{u.email}</div>
                        </div>
                        <div className="flex flex-col items-end gap-1 shrink-0">
                          <Badge variant={u.isOnline ? "success" : "secondary"}>
                            {u.isOnline ? t("dashboard.users.statusOnline") : t("dashboard.users.statusOffline")}
                          </Badge>
                          <Badge variant={u.isPaused ? "destructive" : "outline"}>
                            {u.isPaused ? t("dashboard.users.statusPaused", "Paused") : t("dashboard.users.statusActive", "Active")}
                          </Badge>
                        </div>
                      </div>
                      <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-muted-foreground">
                        <span><span className="font-medium text-foreground">{t("dashboard.users.columnRole")}:</span> {u.role}</span>
                        <span><span className="font-medium text-foreground">{t("dashboard.users.columnDepartment")}:</span> {u.departmentId ? departmentNameById[u.departmentId] || "-" : "-"}</span>
                        {u.fixedIp && <span><span className="font-medium text-foreground">IP:</span> {u.fixedIp}</span>}
                        {u.allocatedVpnIp && <span><span className="font-medium text-foreground">VPN IP:</span> {u.allocatedVpnIp}</span>}
                        <span><span className="font-medium text-foreground">↓</span> {formatBytes(u.bytesReceived)}</span>
                        <span><span className="font-medium text-foreground">↑</span> {formatBytes(u.bytesSent)}</span>
                      </div>
                      <div className="flex flex-wrap gap-2 pt-1 border-t">
                        {canManageUsers && (
                          u.isPaused ? (
                            <Button size="sm" variant="outline" className="h-7 text-xs px-2" onClick={() => handleResumeUser(u.name)}>
                              {t("dashboard.users.resumeButton", "Resume")}
                            </Button>
                          ) : (
                            <Button size="sm" variant="outline" className="h-7 text-xs px-2" onClick={() => handlePauseUser(u.name)}>
                              {t("dashboard.users.pauseButton", "Pause")}
                            </Button>
                          )
                        )}
                        {(currentUser?.role === UserRole.ADMIN || currentUser?.role === UserRole.SUPERADMIN ||
                          (currentUser?.role === UserRole.MANAGER && currentUser.departmentId === u.departmentId)) && (
                          <Button size="sm" variant="outline" className="h-7 text-xs px-2" onClick={() => handleEditClick(u)}>
                            {t("common.edit")}
                          </Button>
                        )}
                        {(currentUser?.role === UserRole.SUPERADMIN ||
                          (currentUser?.role === UserRole.ADMIN && u.role !== UserRole.SUPERADMIN) ||
                          (currentUser?.role === UserRole.MANAGER && u.departmentId === currentUser.departmentId && u.role !== UserRole.SUPERADMIN && u.role !== UserRole.ADMIN)) &&
                          u.id !== currentUser?.id && (
                            <Button size="sm" variant="destructive" className="h-7 text-xs px-2" onClick={() => handleDelete(u.id)}>
                              {t("dashboard.users.deleteButton")}
                            </Button>
                          )}
                        <DownloadSelect username={u.name} size="xs" />
                      </div>
                    </div>
                  );
                })}
              </div>

              {/* ── Desktop: TanStack table (hidden below md) ────────────── */}
              <div className="hidden md:block overflow-x-auto">
                <Table>
                  <TableHeader>
                    {table.getHeaderGroups().map((hg) => (
                      <TableRow key={hg.id}>
                        {hg.headers.map((header) => (
                          <TableHead
                            key={header.id}
                            className={[
                              "h-12 px-4 whitespace-nowrap",
                              header.column.getCanSort() ? "cursor-pointer select-none" : "",
                              header.id === "actions" ? "sticky right-0 min-w-[260px]" : "min-w-[120px]",
                            ].join(" ")}
                            style={header.id === "actions" ? { backgroundColor: "hsl(var(--card))", boxShadow: "inset 10px 0 0px -9px #0505050f" } : undefined}
                            onClick={header.column.getToggleSortingHandler()}
                          >
                            {flexRender(header.column.columnDef.header, header.getContext())}
                            {header.column.getIsSorted() === "asc" ? " ↑" : header.column.getIsSorted() === "desc" ? " ↓" : ""}
                          </TableHead>
                        ))}
                      </TableRow>
                    ))}
                  </TableHeader>
                  <TableBody>
                    {table.getRowModel().rows.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={columns.length} className="py-10 text-center text-muted-foreground">
                          {t("dashboard.users.noResults")}
                        </TableCell>
                      </TableRow>
                    ) : table.getRowModel().rows.map((row) => (
                      <TableRow key={row.id}>
                        {row.getVisibleCells().map((cell) => (
                          <TableCell
                            key={cell.id}
                            className={[
                              "px-4 py-3",
                              cell.column.id === "actions" ? "sticky right-0" : "",
                            ].join(" ")}
                            style={cell.column.id === "actions" ? { backgroundColor: "hsl(var(--card))", boxShadow: "inset 10px 0 0px -9px #0505050f" } : undefined}
                          >
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </TableCell>
                        ))}
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              {/* ── Pagination controls ───────────────────────────────────── */}
              <div className="flex flex-col sm:flex-row items-center justify-between gap-3 pt-4 border-t mt-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <span>{t("dashboard.users.resultsSummary", { count: table.getRowModel().rows.length, total: filteredUsers.length })}</span>
                  <span>·</span>
                  <Select
                    value={String(pagination.pageSize)}
                    onValueChange={(val) => setPagination({ pageIndex: 0, pageSize: Number(val) })}
                  >
                    <SelectTrigger className="h-7 w-24 text-sm">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {[10, 25, 50, 100].map((size) => (
                        <SelectItem key={size} value={String(size)}>{size} / {t("common.page", "page")}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex items-center gap-1">
                  <Button size="sm" variant="outline" className="min-w-8 px-2" disabled={!table.getCanPreviousPage()} onClick={() => table.firstPage()}>«</Button>
                  <Button size="sm" variant="outline" className="min-w-8 px-2" disabled={!table.getCanPreviousPage()} onClick={() => table.previousPage()}>‹</Button>
                  <span className="px-3 text-sm">{table.getState().pagination.pageIndex + 1} / {Math.max(1, table.getPageCount())}</span>
                  <Button size="sm" variant="outline" className="min-w-8 px-2" disabled={!table.getCanNextPage()} onClick={() => table.nextPage()}>›</Button>
                  <Button size="sm" variant="outline" className="min-w-8 px-2" disabled={!table.getCanNextPage()} onClick={() => table.lastPage()}>»</Button>
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>

      <AlertDialog
        open={confirmState.open}
        onOpenChange={(o) => setConfirmState((prev) => ({ ...prev, open: o }))}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{confirmState.title}</AlertDialogTitle>
            <AlertDialogDescription>{confirmState.description}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t("common.cancel")}</AlertDialogCancel>
            <AlertDialogAction
              className={confirmState.destructive ? buttonVariants({ variant: "destructive" }) : undefined}
              onClick={() => confirmState.onConfirm()}
            >
              {t("common.confirm", "Confirm")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </MainLayout>
  );
}
