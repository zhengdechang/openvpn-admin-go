"use client";

import React from "react";
import { usePathname } from "next/navigation";
import { useTranslation } from "react-i18next";
import Sidebar from "./sidebar";

interface MainLayoutProps {
  children: React.ReactNode;
  showFooter?: boolean; // kept for backward compat, ignored
  className?: string;
}

function getPageInfo(
  pathname: string,
  t: (key: string) => string
): { title: string; breadcrumb: string } {
  if (pathname.startsWith("/dashboard/users"))
    return {
      title: t("dashboard.users.title"),
      breadcrumb: t("dashboard.users.title"),
    };
  if (pathname.startsWith("/dashboard/departments"))
    return {
      title: t("dashboard.departments.title") || "部门管理",
      breadcrumb: t("dashboard.departments.title") || "部门管理",
    };
  if (pathname.startsWith("/dashboard/server"))
    return {
      title: t("dashboard.server.title"),
      breadcrumb: t("dashboard.server.title"),
    };
  if (pathname.startsWith("/dashboard/logs"))
    return {
      title: t("dashboard.logs.titleServer"),
      breadcrumb: t("dashboard.logs.titleServer"),
    };
  if (pathname.startsWith("/dashboard/profile"))
    return {
      title: t("layout.profile"),
      breadcrumb: t("layout.profile"),
    };
  return { title: "Dashboard", breadcrumb: "Dashboard" };
}

export default function MainLayout({
  children,
  className,
}: MainLayoutProps) {
  const { t } = useTranslation();
  const pathname = usePathname();
  const { title, breadcrumb } = getPageInfo(pathname || "", t);

  return (
    <div style={{ display: "flex", height: "100vh", overflow: "hidden" }}>
      <Sidebar />

      {/* 右侧主区域 */}
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
            background: "#ffffff",
            borderBottom: "1px solid hsl(var(--border))",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            padding: "0 24px",
            flexShrink: 0,
          }}
        >
          <div>
            <div
              style={{
                fontSize: "15px",
                fontWeight: 600,
                color: "#111827",
                lineHeight: 1.2,
              }}
            >
              {title}
            </div>
            <div style={{ fontSize: "11px", color: "#9ca3af", marginTop: "2px" }}>
              Dashboard / {breadcrumb}
            </div>
          </div>

          <div style={{ display: "flex", alignItems: "center", gap: "4px" }}>
            {/* 搜索 */}
            <button
              style={{
                width: "32px",
                height: "32px",
                borderRadius: "8px",
                border: "none",
                background: "transparent",
                cursor: "pointer",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "#9ca3af",
              }}
            >
              <svg
                width="15"
                height="15"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <circle cx="11" cy="11" r="8" />
                <line x1="21" y1="21" x2="16.65" y2="16.65" />
              </svg>
            </button>

            {/* 通知铃 */}
            <button
              style={{
                width: "32px",
                height: "32px",
                borderRadius: "8px",
                border: "none",
                background: "transparent",
                cursor: "pointer",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "#9ca3af",
              }}
            >
              <svg
                width="15"
                height="15"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                <path d="M13.73 21a2 2 0 0 1-3.46 0" />
              </svg>
            </button>
          </div>
        </header>

        {/* 内容区 */}
        <main
          className={`flex-1 overflow-y-auto custom-scrollbar ${className || ""}`}
        >
          {children}
        </main>
      </div>
    </div>
  );
}
