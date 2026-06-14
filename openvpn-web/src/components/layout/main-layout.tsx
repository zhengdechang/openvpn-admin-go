"use client";

import React, { useState, useEffect } from "react";
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
  if (pathname.startsWith("/dashboard/overview"))
    return {
      title: t("dashboard.overview.title"),
      breadcrumb: t("dashboard.overview.title"),
    };
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
  const [mobileNavOpen, setMobileNavOpen] = useState(false);

  // 路由变化时自动收起移动端抽屉（点导航项跳转后关闭）。
  useEffect(() => {
    setMobileNavOpen(false);
  }, [pathname]);

  return (
    <div className="dash-root" style={{ display: "flex", height: "125vh", overflow: "hidden" }}>
      <Sidebar mobileOpen={mobileNavOpen} onClose={() => setMobileNavOpen(false)} />

      {/* 移动端抽屉遮罩：点击关闭。桌面 ≥1024px 由 CSS 隐藏。 */}
      {mobileNavOpen && (
        <div className="dash-backdrop" onClick={() => setMobileNavOpen(false)} />
      )}

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
        {/* Topbar — Argon 紫蓝渐变 */}
        <header
          className="dash-topbar"
          style={{
            height: "60px",
            minHeight: "60px",
            background: "linear-gradient(87deg, #5e72e4 0%, #825ee4 100%)",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            padding: "0 24px",
            flexShrink: 0,
          }}
        >
          <div style={{ display: "flex", alignItems: "center", gap: "10px", minWidth: 0 }}>
            {/* 汉堡按钮：仅 <1024px 显示，打开移动端侧栏抽屉。 */}
            <button
              className="dash-hamburger"
              onClick={() => setMobileNavOpen(true)}
              aria-label={t("layout.dashboard")}
              onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(255,255,255,0.15)"; }}
              onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}
              style={{
                width: "34px",
                height: "34px",
                borderRadius: "8px",
                border: "none",
                background: "transparent",
                cursor: "pointer",
                alignItems: "center",
                justifyContent: "center",
                color: "#ffffff",
                transition: "background 0.15s",
                flexShrink: 0,
              }}
            >
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <line x1="3" y1="6" x2="21" y2="6" />
                <line x1="3" y1="12" x2="21" y2="12" />
                <line x1="3" y1="18" x2="21" y2="18" />
              </svg>
            </button>

            <div style={{ minWidth: 0 }}>
              <div
                style={{
                  fontSize: "15px",
                  fontWeight: 600,
                  color: "#ffffff",
                  lineHeight: 1.2,
                  whiteSpace: "nowrap",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                }}
              >
                {title}
              </div>
              <div style={{ fontSize: "11px", color: "rgba(255,255,255,0.7)", marginTop: "2px", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
                {t("layout.dashboard")} / {breadcrumb}
              </div>
            </div>
          </div>

          <div style={{ display: "flex", alignItems: "center", gap: "4px" }}>
            {/* 搜索 */}
            <button
              onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(255,255,255,0.15)"; }}
              onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}
              style={{
                width: "34px",
                height: "34px",
                borderRadius: "8px",
                border: "none",
                background: "transparent",
                cursor: "pointer",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "#ffffff",
                transition: "background 0.15s",
              }}
            >
              <svg
                width="16"
                height="16"
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
              onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(255,255,255,0.15)"; }}
              onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}
              style={{
                width: "34px",
                height: "34px",
                borderRadius: "8px",
                border: "none",
                background: "transparent",
                cursor: "pointer",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "#ffffff",
                transition: "background 0.15s",
              }}
            >
              <svg
                width="16"
                height="16"
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
