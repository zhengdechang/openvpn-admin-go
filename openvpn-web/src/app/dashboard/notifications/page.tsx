"use client";

import React, { useState, useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import MainLayout from "@/components/layout/main-layout";
import { useNotificationStore, useUserStore } from "@/store";
import { UserRole } from "@/types/types";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

type FilterTab = "all" | "unread" | "read";

export default function NotificationsPage() {
  const { t } = useTranslation();
  const { user } = useUserStore();
  const {
    notifications,
    unreadCount,
    isLoading,
    error,
    fetchNotifications,
    markRead,
    markAllRead,
  } = useNotificationStore();

  const [filter, setFilter] = useState<FilterTab>("all");
  const [pageIndex, setPageIndex] = useState(0);
  const [pageSize, setPageSize] = useState(20);

  useEffect(() => {
    fetchNotifications();
  }, [fetchNotifications]);

  useEffect(() => {
    setPageIndex(0);
  }, [filter]);

  const filtered = useMemo(() => {
    if (filter === "unread") return notifications.filter((n) => !n.isRead);
    if (filter === "read") return notifications.filter((n) => n.isRead);
    return notifications;
  }, [notifications, filter]);

  const pageCount = Math.max(1, Math.ceil(filtered.length / pageSize));
  const pageRows = filtered.slice(pageIndex * pageSize, (pageIndex + 1) * pageSize);

  if (user?.role !== UserRole.SUPERADMIN) {
    return (
      <MainLayout className="p-4">
        <p className="text-muted-foreground">{t("common.forbidden", "Access denied")}</p>
      </MainLayout>
    );
  }

  return (
    <MainLayout className="p-4 space-y-6">
      {/* Page header */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-2xl font-bold">{t("dashboard.notifications.pageTitle")}</h1>
        {unreadCount > 0 && (
          <Button variant="outline" onClick={markAllRead}>
            {t("dashboard.notifications.markAllRead")}
          </Button>
        )}
      </div>

      {/* Filter tabs */}
      <div className="flex gap-2 flex-wrap">
        {(["all", "unread", "read"] as FilterTab[]).map((tab) => (
          <Button
            key={tab}
            size="sm"
            variant={filter === tab ? "default" : "outline"}
            onClick={() => setFilter(tab)}
          >
            {tab === "all" && t("dashboard.notifications.allTab")}
            {tab === "unread" && (
              <>
                {t("dashboard.notifications.unreadTab")}
                {unreadCount > 0 && (
                  <span className="ml-1.5 inline-flex items-center justify-center min-w-[18px] h-[18px] px-1 rounded-full bg-destructive text-white text-[10px] font-semibold leading-none">
                    {unreadCount > 99 ? "99+" : unreadCount}
                  </span>
                )}
              </>
            )}
            {tab === "read" && t("dashboard.notifications.readTab")}
          </Button>
        ))}
      </div>

      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base font-medium text-muted-foreground">
            {filtered.length} {t("common.items", "items")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <p className="py-10 text-center text-muted-foreground">
              {t("dashboard.notifications.loading")}
            </p>
          ) : error ? (
            <p className="py-10 text-center text-destructive">{error}</p>
          ) : filtered.length === 0 ? (
            <div className="py-10 text-center">
              <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="mx-auto mb-3 text-muted-foreground/50">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                <path d="M13.73 21a2 2 0 0 1-3.46 0" />
              </svg>
              <p className="text-muted-foreground text-sm">{t("dashboard.notifications.noEvents")}</p>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-8 px-3"></TableHead>
                      <TableHead className="px-3">{t("dashboard.notifications.columnType")}</TableHead>
                      <TableHead className="px-3">{t("dashboard.notifications.columnUser")}</TableHead>
                      <TableHead className="px-3">{t("dashboard.notifications.columnIPs")}</TableHead>
                      <TableHead className="px-3 whitespace-nowrap">{t("dashboard.notifications.columnTime")}</TableHead>
                      <TableHead className="px-3">{t("dashboard.notifications.columnStatus")}</TableHead>
                      <TableHead className="px-3"></TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {pageRows.map((n) => {
                      const isConnected = n.type === "user_connected";
                      return (
                        <TableRow key={n.id} className={n.isRead ? "" : "bg-accent/20"}>
                          <TableCell className="px-3 py-3">
                            <span className={[
                              "inline-block w-2 h-2 rounded-full",
                              isConnected ? "bg-green-500" : "bg-red-500",
                            ].join(" ")} />
                          </TableCell>
                          <TableCell className="px-3 py-3">
                            <Badge variant={isConnected ? "success" : "destructive"} className="text-xs">
                              {isConnected
                                ? t("dashboard.notifications.connected")
                                : t("dashboard.notifications.disconnected")}
                            </Badge>
                          </TableCell>
                          <TableCell className="px-3 py-3 font-medium text-sm">
                            {n.userName}
                          </TableCell>
                          <TableCell className="px-3 py-3 text-muted-foreground text-xs font-mono whitespace-nowrap">
                            {n.realIP || "-"}
                            {n.realIP && n.virtualIP && " → "}
                            {n.virtualIP || ""}
                          </TableCell>
                          <TableCell className="px-3 py-3 text-muted-foreground text-xs whitespace-nowrap">
                            {new Date(n.createdAt).toLocaleString()}
                          </TableCell>
                          <TableCell className="px-3 py-3">
                            <Badge variant={n.isRead ? "secondary" : "outline"} className="text-xs">
                              {n.isRead
                                ? t("dashboard.notifications.statusRead")
                                : t("dashboard.notifications.statusUnread")}
                            </Badge>
                          </TableCell>
                          <TableCell className="px-3 py-3">
                            {!n.isRead && (
                              <Button
                                size="sm"
                                variant="ghost"
                                className="h-7 text-xs px-2"
                                onClick={() => markRead(n.id)}
                              >
                                {t("dashboard.notifications.markRead")}
                              </Button>
                            )}
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </div>

              {/* Pagination */}
              <div className="flex flex-col sm:flex-row items-center justify-between gap-3 pt-4 border-t mt-2">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <span>
                    {pageIndex * pageSize + 1}–{Math.min((pageIndex + 1) * pageSize, filtered.length)} / {filtered.length}
                  </span>
                  <span>·</span>
                  <Select
                    value={String(pageSize)}
                    onValueChange={(v) => { setPageSize(Number(v)); setPageIndex(0); }}
                  >
                    <SelectTrigger className="h-7 w-24 text-sm">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {[10, 20, 50, 100].map((s) => (
                        <SelectItem key={s} value={String(s)}>
                          {s} / {t("common.page", "page")}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex items-center gap-1">
                  <Button size="sm" variant="outline" className="min-w-8 px-2" disabled={pageIndex === 0} onClick={() => setPageIndex(0)}>«</Button>
                  <Button size="sm" variant="outline" className="min-w-8 px-2" disabled={pageIndex === 0} onClick={() => setPageIndex((p) => p - 1)}>‹</Button>
                  <span className="px-3 text-sm">{pageIndex + 1} / {pageCount}</span>
                  <Button size="sm" variant="outline" className="min-w-8 px-2" disabled={pageIndex + 1 >= pageCount} onClick={() => setPageIndex((p) => p + 1)}>›</Button>
                  <Button size="sm" variant="outline" className="min-w-8 px-2" disabled={pageIndex + 1 >= pageCount} onClick={() => setPageIndex(pageCount - 1)}>»</Button>
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
