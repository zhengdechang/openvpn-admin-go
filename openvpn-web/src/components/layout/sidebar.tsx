"use client";

import React, { useState, useRef, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { UserRole } from "@/types/types";
import { useTranslation } from "react-i18next";
import VPNLogo from "@/components/ui/vpn-logo";
import { Button } from "@/components/ui/button";

interface NavItem {
  href: string;
  label: string;
  icon: React.ReactNode;
  badge?: number;
  roles?: UserRole[];
}

function LockIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
      <path d="M7 11V7a5 5 0 0 1 10 0v4" />
    </svg>
  );
}

function UsersIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
      <circle cx="9" cy="7" r="4" />
      <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
      <path d="M16 3.13a4 4 0 0 1 0 7.75" />
    </svg>
  );
}

function DepartmentsIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
      <polyline points="9 22 9 12 15 12 15 22" />
    </svg>
  );
}

function ServerIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
      <line x1="8" y1="21" x2="16" y2="21" />
      <line x1="12" y1="17" x2="12" y2="21" />
    </svg>
  );
}

function LogsIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
      <polyline points="14 2 14 8 20 8" />
      <line x1="16" y1="13" x2="8" y2="13" />
      <line x1="16" y1="17" x2="8" y2="17" />
      <polyline points="10 9 9 9 8 9" />
    </svg>
  );
}

function SettingsIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="3" />
      <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
    </svg>
  );
}

function BellIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
      <path d="M13.73 21a2 2 0 0 1-3.46 0" />
    </svg>
  );
}

function LogoutIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
      <polyline points="16 17 21 12 16 7" />
      <line x1="21" y1="12" x2="9" y2="12" />
    </svg>
  );
}

function getInitials(name?: string): string {
  if (!name) return "U";
  const parts = name.split(" ");
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase();
  return name.slice(0, 2).toUpperCase();
}

const AVATAR_COLORS = [
  "#0369a1", "#0ea5e9", "#0284c7", "#0891b2",
  "#22c55e", "#16a34a", "#0c4a6e", "#075985",
];

function getAvatarColor(name?: string): string {
  if (!name) return AVATAR_COLORS[0];
  const idx = name.charCodeAt(0) % AVATAR_COLORS.length;
  return AVATAR_COLORS[idx];
}

export default function Sidebar({ onClose }: { onClose?: () => void } = {}) {
  const { user, logout } = useAuth();
  const pathname = usePathname();
  const { t } = useTranslation();
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  const isActive = (href: string) =>
    pathname === href || pathname?.startsWith(`${href}/`);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setUserMenuOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const navItems: NavItem[] = [
    {
      href: "/dashboard/users",
      label: t("dashboard.users.title"),
      icon: <UsersIcon />,
    },
    {
      href: "/dashboard/departments",
      label: t("dashboard.departments.title") || "部门管理",
      icon: <DepartmentsIcon />,
      roles: [UserRole.ADMIN, UserRole.SUPERADMIN],
    },
    {
      href: "/dashboard/server",
      label: t("dashboard.server.title"),
      icon: <ServerIcon />,
      roles: [UserRole.SUPERADMIN],
    },
    {
      href: "/dashboard/logs",
      label: t("dashboard.logs.titleServer"),
      icon: <LogsIcon />,
      roles: [UserRole.SUPERADMIN],
    },
    {
      href: "/dashboard/notifications",
      label: t("dashboard.notifications.navLabel"),
      icon: <BellIcon />,
      roles: [UserRole.SUPERADMIN],
    },
  ];

  const visibleItems = navItems.filter(
    (item) =>
      !item.roles || (user && item.roles.includes(user.role as UserRole))
  );

  const roleLabel = () => {
    if (!user) return "";
    switch (user.role) {
      case UserRole.SUPERADMIN: return "超级管理员";
      case UserRole.ADMIN: return "管理员";
      case UserRole.MANAGER: return "经理";
      default: return "用户";
    }
  };

  return (
    <aside className="w-[220px] min-w-[220px] bg-[hsl(var(--sidebar-bg))] flex flex-col h-screen overflow-hidden">
      {/* Logo */}
      <div className="px-4 pt-4 pb-[14px] border-b border-white/[0.07] flex items-center gap-[10px] shrink-0">
        <VPNLogo size={32} />
        <div className="flex-1 min-w-0">
          <div className="text-[14px] font-bold text-white leading-tight">
            VPN Admin
          </div>
          <div className="text-[10px] text-[hsl(var(--sidebar-text))] font-normal tracking-[0.5px] uppercase">
            管理控制台
          </div>
        </div>
        {/* 关闭/折叠按钮 */}
        {onClose && (
          <Button
            onClick={onClose}
            aria-label="Close menu"
            variant="ghost"
            size="icon"
            className="h-7 w-7 rounded-md flex-shrink-0 text-white/60 hover:text-white hover:bg-white/14"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </Button>
        )}
      </div>

      {/* Nav */}
      <nav className="flex-1 px-2 py-3 overflow-y-auto">
        <div className="text-[10px] font-semibold text-white/30 tracking-[0.8px] uppercase px-3 pt-2 pb-1.5">
          主菜单
        </div>

        {visibleItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={[
              "flex items-center gap-[10px] px-3 py-[9px] rounded-lg cursor-pointer",
              "text-[13.5px] font-medium transition-colors mb-0.5 relative no-underline",
              isActive(item.href)
                ? "text-[hsl(var(--sidebar-text-active))] bg-[hsl(var(--sidebar-active))]"
                : "text-[hsl(var(--sidebar-text))] bg-transparent hover:bg-white/[0.06]",
            ].join(" ")}
          >
            {isActive(item.href) && (
              <span className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-5 bg-primary rounded-r-[3px]" />
            )}
            {item.icon}
            {item.label}
          </Link>
        ))}
      </nav>

      {/* User footer */}
      <div ref={menuRef} className="px-2 py-3 border-t border-white/[0.07] shrink-0 relative">
        <div
          className="flex items-center gap-[10px] px-3 py-[10px] rounded-lg cursor-pointer transition-colors hover:bg-white/[0.06]"
          onClick={() => setUserMenuOpen(!userMenuOpen)}
        >
          <div
            className="w-[30px] h-[30px] rounded-full flex items-center justify-center text-xs font-bold text-white shrink-0"
            style={{ background: user ? getAvatarColor(user.name) : "#3b82f6" }}
          >
            {getInitials(user?.name)}
          </div>
          <div className="flex-1 min-w-0">
            <div className="text-[12.5px] font-semibold text-[#e5e7eb] truncate">
              {user?.name || "用户"}
            </div>
            <div className="text-[10.5px] text-[hsl(var(--sidebar-text))]">
              {roleLabel()}
            </div>
          </div>
          <div className="text-white/30 shrink-0">
            <SettingsIcon />
          </div>
        </div>

        {/* User dropdown menu */}
        {userMenuOpen && (
          <div className="absolute bottom-[calc(100%-8px)] left-2 right-2 bg-[#012a4a] rounded-[10px] border border-white/10 shadow-[0_-8px_24px_rgba(0,0,0,0.3)] p-1.5 z-50">
            <Link
              href="/dashboard/profile"
              className="flex items-center gap-2 px-[10px] py-2 rounded-md text-[13px] text-[#d1d5db] no-underline transition-colors hover:bg-white/10"
              onClick={() => setUserMenuOpen(false)}
            >
              <SettingsIcon />
              {t("layout.profile")}
            </Link>

            <div className="border-t border-white/[0.07] my-1" />

            <button
              onClick={() => { logout(); setUserMenuOpen(false); }}
              className="flex items-center gap-2 px-[10px] py-2 rounded-md text-[13px] text-[#f87171] border-none bg-transparent cursor-pointer w-full font-semibold transition-colors hover:bg-white/10"
            >
              <LogoutIcon />
              {t("layout.logout")}
            </button>
          </div>
        )}
      </div>
    </aside>
  );
}
