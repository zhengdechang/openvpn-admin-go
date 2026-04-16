"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useTranslation } from "react-i18next";
import Sidebar from "./sidebar";
import LanguageSwitcher from "@/components/ui/language-switcher";
import GitHubButton from "@/components/ui/github-button";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { useUserStore, useNotificationStore } from "@/store";
import { UserRole, Notification } from "@/types";

interface MainLayoutProps {
  children: React.ReactNode;
  showFooter?: boolean;
  className?: string;
}

function getPageInfo(
  pathname: string,
  t: (key: string) => string
): { title: string; breadcrumb: string } {
  if (pathname.startsWith("/dashboard/users"))
    return { title: t("dashboard.users.title"), breadcrumb: t("dashboard.users.title") };
  if (pathname.startsWith("/dashboard/departments"))
    return { title: t("dashboard.departments.title") || "部门管理", breadcrumb: t("dashboard.departments.title") || "部门管理" };
  if (pathname.startsWith("/dashboard/server"))
    return { title: t("dashboard.server.title"), breadcrumb: t("dashboard.server.title") };
  if (pathname.startsWith("/dashboard/logs"))
    return { title: t("dashboard.logs.titleServer"), breadcrumb: t("dashboard.logs.titleServer") };
  if (pathname.startsWith("/dashboard/notifications"))
    return { title: t("dashboard.notifications.pageTitle"), breadcrumb: t("dashboard.notifications.pageTitle") };
  if (pathname.startsWith("/dashboard/profile"))
    return { title: t("layout.profile"), breadcrumb: t("layout.profile") };
  return { title: "Dashboard", breadcrumb: "Dashboard" };
}

function formatRelativeTime(isoString: string): string {
  const diff = Date.now() - new Date(isoString).getTime();
  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return "just now";
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  return `${Math.floor(hours / 24)}d ago`;
}

function NotificationItem({ n, onMarkRead }: { n: Notification; onMarkRead: (id: string) => void }) {
  const isConnected = n.type === "user_connected";
  return (
    <div
      onClick={() => !n.isRead && onMarkRead(n.id)}
      className={[
        "px-[14px] py-[10px] border-b flex gap-[10px] items-start",
        n.isRead ? "bg-transparent cursor-default" : "bg-accent/30 cursor-pointer",
      ].join(" ")}
    >
      {/* Status dot */}
      <span className={[
        "w-2 h-2 rounded-full shrink-0 mt-[5px]",
        isConnected ? "bg-green-500" : "bg-red-500",
      ].join(" ")} />
      <div className="flex-1 min-w-0">
        <div className="text-[13px] font-medium text-foreground">
          {n.userName}{" "}
          <span className="font-normal text-muted-foreground">
            {isConnected ? "connected" : "disconnected"}
          </span>
        </div>
        {(n.realIP || n.virtualIP) && (
          <div className="text-[11px] text-muted-foreground mt-0.5 overflow-hidden text-ellipsis whitespace-nowrap">
            {n.realIP && <span>{n.realIP}</span>}
            {n.realIP && n.virtualIP && <span> → </span>}
            {n.virtualIP && <span>{n.virtualIP}</span>}
          </div>
        )}
        <div className="text-[11px] text-muted-foreground mt-0.5">
          {formatRelativeTime(n.createdAt)}
        </div>
      </div>
    </div>
  );
}

export default function MainLayout({ children, className }: MainLayoutProps) {
  const { t } = useTranslation();
  const pathname = usePathname();
  const { title, breadcrumb } = getPageInfo(pathname || "", t);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [desktopCollapsed, setDesktopCollapsed] = useState(false);

  const { user } = useUserStore();
  const isSuperAdmin = user?.role === UserRole.SUPERADMIN;

  const {
    notifications,
    unreadCount,
    isOpen,
    isLoading,
    error,
    fetchUnreadCount,
    markRead,
    markAllRead,
    setOpen,
  } = useNotificationStore();

  // Poll unread count every 30 seconds — superadmin only
  useEffect(() => {
    if (!isSuperAdmin) return;
    fetchUnreadCount();
    const interval = setInterval(fetchUnreadCount, 30000);
    return () => clearInterval(interval);
  }, [isSuperAdmin, fetchUnreadCount]);

  // 路由切换时自动关闭侧边栏（移动端）
  useEffect(() => {
    setSidebarOpen(false);
  }, [pathname]);

  // 侧边栏打开时禁止 body 滚动
  useEffect(() => {
    if (sidebarOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
    return () => { document.body.style.overflow = ""; };
  }, [sidebarOpen]);

  const menuIcon = (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <line x1="3" y1="6" x2="21" y2="6" />
      <line x1="3" y1="12" x2="21" y2="12" />
      <line x1="3" y1="18" x2="21" y2="18" />
    </svg>
  );

  return (
    <div className="flex h-screen overflow-hidden">

      {/* ─── 移动端遮罩 ────────────────────────────────────────── */}
      {sidebarOpen && (
        <div
          onClick={() => setSidebarOpen(false)}
          className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40"
        />
      )}

      {/* ─── 侧边栏（移动端 fixed 覆盖层）─────────────────────── */}
      <div
        className={[
          "fixed top-0 left-0 h-full z-50 md:hidden",
          "transition-transform duration-300 ease-in-out",
          sidebarOpen ? "translate-x-0" : "-translate-x-full",
        ].join(" ")}
      >
        <Sidebar onClose={() => setSidebarOpen(false)} />
      </div>

      {/* ─── 侧边栏（桌面端正常流，支持折叠）────────────────────── */}
      <div
        className="hidden md:block flex-shrink-0 transition-all duration-300 ease-in-out overflow-hidden"
        style={{ width: desktopCollapsed ? 0 : 220 }}
      >
        <Sidebar onClose={() => setDesktopCollapsed(true)} />
      </div>

      {/* ─── 右侧主区域 ─────────────────────────────────────────── */}
      <div className="flex-1 flex flex-col overflow-hidden bg-background min-w-0">
        {/* Topbar */}
        <header className="h-14 min-h-14 bg-card border-b flex items-center justify-between px-4 shrink-0">
          <div className="flex items-center gap-3 min-w-0">
            {/* ─── 汉堡菜单（仅移动端）─── */}
            <div className="md:hidden">
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8"
                onClick={() => setSidebarOpen(true)}
                aria-label="Open menu"
              >
                {menuIcon}
              </Button>
            </div>
            {/* ─── 侧边栏折叠/展开（仅桌面端）─── */}
            <div className="hidden md:block">
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8"
                onClick={() => setDesktopCollapsed(c => !c)}
                aria-label="Toggle sidebar"
              >
                {menuIcon}
              </Button>
            </div>

            <div className="min-w-0">
              <div className="text-[15px] font-semibold text-foreground leading-tight whitespace-nowrap overflow-hidden text-ellipsis">
                {title}
              </div>
              <div className="text-[11px] text-muted-foreground mt-px hidden sm:block">
                Dashboard / {breadcrumb}
              </div>
            </div>
          </div>

          <div className="flex items-center gap-1 shrink-0">
            {/* GitHub */}
            <GitHubButton />
            {/* Language switcher */}
            <LanguageSwitcher />

            {/* 通知铃 — superadmin only */}
            {isSuperAdmin && (
              <Popover open={isOpen} onOpenChange={setOpen}>
                <PopoverTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 relative"
                    aria-label="Notifications"
                  >
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                      <path d="M13.73 21a2 2 0 0 1-3.46 0" />
                    </svg>
                    {unreadCount > 0 && (
                      <span className="absolute -top-0.5 -right-0.5 min-w-[16px] h-4 px-1 rounded-full bg-destructive text-white text-[9px] font-medium flex items-center justify-center leading-none">
                        {unreadCount > 99 ? "99+" : unreadCount}
                      </span>
                    )}
                  </Button>
                </PopoverTrigger>
                <PopoverContent align="end" className="w-80 p-0 max-h-[420px] overflow-hidden flex flex-col bg-card">
                  {/* Header */}
                  <div className="flex items-center justify-between px-[14px] py-[10px] border-b shrink-0">
                    <span className="text-[13px] font-semibold text-foreground">
                      Notifications {unreadCount > 0 && <span className="text-destructive">({unreadCount})</span>}
                    </span>
                    {unreadCount > 0 && (
                      <button
                        onClick={markAllRead}
                        className="text-[11px] text-muted-foreground bg-transparent border-none cursor-pointer px-1.5 py-0.5 rounded hover:bg-accent"
                      >
                        Mark all read
                      </button>
                    )}
                  </div>

                  {/* Body */}
                  <div className="overflow-y-auto flex-1">
                    {isLoading && (
                      <div className="p-6 text-center text-muted-foreground text-[13px]">
                        {t("dashboard.notifications.loading")}
                      </div>
                    )}
                    {error && !isLoading && (
                      <div className="p-6 text-center text-destructive text-[13px]">
                        {error}
                      </div>
                    )}
                    {!isLoading && !error && notifications.filter((n) => !n.isRead).length === 0 && (
                      <div className="p-8 text-center">
                        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="mx-auto mb-2 text-muted-foreground">
                          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                          <path d="M13.73 21a2 2 0 0 1-3.46 0" />
                        </svg>
                        <div className="text-[13px] text-muted-foreground">{t("dashboard.notifications.noEvents")}</div>
                      </div>
                    )}
                    {!isLoading && !error && notifications.filter((n) => !n.isRead).slice(0, 20).map((n) => (
                      <NotificationItem key={n.id} n={n} onMarkRead={markRead} />
                    ))}
                  </div>
                  {/* Footer */}
                  <div className="border-t px-[14px] py-2 shrink-0 flex justify-end">
                    <Link
                      href="/dashboard/notifications"
                      className="text-[12px] text-primary hover:underline"
                      onClick={() => setOpen(false)}
                    >
                      {t("dashboard.notifications.viewAll")} →
                    </Link>
                  </div>
                </PopoverContent>
              </Popover>
            )}
          </div>
        </header>

        {/* 内容区 */}
        <main className={`flex-1 overflow-y-auto custom-scrollbar ${className || ""}`}>
          {children}
        </main>
      </div>
    </div>
  );
}
