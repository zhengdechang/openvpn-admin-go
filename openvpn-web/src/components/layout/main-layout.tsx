"use client";

import React, { useState, useEffect } from "react";
import { usePathname } from "next/navigation";
import { useTranslation } from "react-i18next";
import Sidebar from "./sidebar";
import LanguageSwitcher from "@/components/ui/language-switcher";
import GitHubButton from "@/components/ui/github-button";
import IconButton from "@mui/material/IconButton";

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
  if (pathname.startsWith("/dashboard/profile"))
    return { title: t("layout.profile"), breadcrumb: t("layout.profile") };
  return { title: "Dashboard", breadcrumb: "Dashboard" };
}

export default function MainLayout({ children, className }: MainLayoutProps) {
  const { t } = useTranslation();
  const pathname = usePathname();
  const { title, breadcrumb } = getPageInfo(pathname || "", t);
  const [sidebarOpen, setSidebarOpen] = useState(false);

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

  return (
    <div style={{ display: "flex", height: "100vh", overflow: "hidden" }}>

      {/* ─── 移动端遮罩 ────────────────────────────────────────── */}
      {sidebarOpen && (
        <div
          onClick={() => setSidebarOpen(false)}
          style={{
            position: "fixed",
            inset: 0,
            background: "rgba(0,0,0,0.5)",
            zIndex: 40,
            backdropFilter: "blur(2px)",
          }}
        />
      )}

      {/* ─── 侧边栏：桌面端正常流；移动端 fixed 覆盖层 ──────────── */}
      <div
        className={[
          // 移动端：fixed 从左侧滑入
          "fixed top-0 left-0 h-full z-50 flex-shrink-0",
          "transition-transform duration-300 ease-in-out",
          // 桌面端：回归正常文档流
          "md:relative md:top-auto md:left-auto md:translate-x-0 md:z-auto",
          sidebarOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0",
        ].join(" ")}
      >
        <Sidebar onClose={() => setSidebarOpen(false)} />
      </div>

      {/* ─── 右侧主区域 ─────────────────────────────────────────── */}
      <div
        style={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          overflow: "hidden",
          background: "hsl(var(--background))",
          minWidth: 0,
        }}
      >
        {/* Topbar */}
        <header
          style={{
            height: "56px",
            minHeight: "56px",
            background: "hsl(var(--card))",
            borderBottom: "1px solid hsl(var(--border))",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            padding: "0 16px",
            flexShrink: 0,
          }}
        >
          <div style={{ display: "flex", alignItems: "center", gap: "12px", minWidth: 0 }}>
            {/* ─── 汉堡菜单（仅移动端显示）─── */}
            <IconButton
              onClick={() => setSidebarOpen(true)}
              className="md:hidden"
              aria-label="Open menu"
              size="small"
              sx={{ color: "hsl(var(--foreground))", flexShrink: 0 }}
            >
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="3" y1="6" x2="21" y2="6" />
                <line x1="3" y1="12" x2="21" y2="12" />
                <line x1="3" y1="18" x2="21" y2="18" />
              </svg>
            </IconButton>

            <div style={{ minWidth: 0 }}>
              <div style={{ fontSize: "15px", fontWeight: 600, color: "hsl(var(--foreground))", lineHeight: 1.2, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
                {title}
              </div>
              <div style={{ fontSize: "11px", color: "hsl(var(--muted-foreground))", marginTop: "1px" }} className="hidden sm:block">
                Dashboard / {breadcrumb}
              </div>
            </div>
          </div>

          <div style={{ display: "flex", alignItems: "center", gap: "4px", flexShrink: 0 }}>
            {/* GitHub */}
            <GitHubButton />
            {/* Language switcher */}
            <LanguageSwitcher />

            {/* 通知铃 */}
            <IconButton
              size="small"
              aria-label="Notifications"
              sx={{ color: "hsl(var(--muted-foreground))" }}
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                <path d="M13.73 21a2 2 0 0 1-3.46 0" />
              </svg>
            </IconButton>
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
