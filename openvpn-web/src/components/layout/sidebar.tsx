"use client";

import React, { useState, useRef, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { UserRole } from "@/types/types";
import { useTranslation } from "react-i18next";
import { setLocaleOnClient, getLocaleOnClient } from "@/i18n";
import { LanguagesSupported } from "@/i18n/language";
import type { Locale } from "@/i18n";

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

function ChevronRightIcon() {
  return (
    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="9 18 15 12 9 6" />
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

function GlobeIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" />
      <line x1="2" y1="12" x2="22" y2="12" />
      <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
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
  "#7c3aed", "#2563eb", "#0891b2", "#059669",
  "#d97706", "#dc2626", "#7c3aed", "#0284c7",
];

function getAvatarColor(name?: string): string {
  if (!name) return AVATAR_COLORS[0];
  const idx = name.charCodeAt(0) % AVATAR_COLORS.length;
  return AVATAR_COLORS[idx];
}

export default function Sidebar() {
  const { user, logout } = useAuth();
  const pathname = usePathname();
  const { t, i18n } = useTranslation();
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const [langMenuOpen, setLangMenuOpen] = useState(false);
  const [currentLocale, setCurrentLocale] = useState<Locale>(getLocaleOnClient());
  const menuRef = useRef<HTMLDivElement>(null);

  const isActive = (href: string) =>
    pathname === href || pathname?.startsWith(`${href}/`);

  const handleLanguageChange = (locale: Locale) => {
    setLocaleOnClient(locale, true);
    setCurrentLocale(locale);
    setLangMenuOpen(false);
    setUserMenuOpen(false);
    document.documentElement.lang = locale;
    i18n.changeLanguage(locale);
  };

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setUserMenuOpen(false);
        setLangMenuOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  useEffect(() => {
    const handler = () => setCurrentLocale(i18n.language as Locale);
    i18n.on("languageChanged", handler);
    return () => i18n.off("languageChanged", handler);
  }, [i18n]);

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
    <aside
      style={{
        width: "220px",
        minWidth: "220px",
        background: "hsl(var(--sidebar-bg))",
        display: "flex",
        flexDirection: "column",
        height: "100vh",
        overflow: "hidden",
      }}
    >
      {/* Logo */}
      <div
        style={{
          padding: "20px 20px 16px",
          borderBottom: "1px solid rgba(255,255,255,0.07)",
          display: "flex",
          alignItems: "center",
          gap: "10px",
          flexShrink: 0,
        }}
      >
        <div
          style={{
            width: "32px",
            height: "32px",
            background: "hsl(var(--primary))",
            borderRadius: "8px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            flexShrink: 0,
          }}
        >
          <LockIcon />
        </div>
        <div>
          <div style={{ fontSize: "14px", fontWeight: 700, color: "#ffffff", lineHeight: 1.2 }}>
            VPN Admin
          </div>
          <div style={{ fontSize: "10px", color: "hsl(var(--sidebar-text))", fontWeight: 400, letterSpacing: "0.5px", textTransform: "uppercase" }}>
            管理控制台
          </div>
        </div>
      </div>

      {/* Nav */}
      <nav style={{ flex: 1, padding: "12px 8px", overflowY: "auto" }}>
        <div
          style={{
            fontSize: "10px",
            fontWeight: 600,
            color: "rgba(255,255,255,0.3)",
            letterSpacing: "0.8px",
            textTransform: "uppercase",
            padding: "8px 12px 6px",
          }}
        >
          主菜单
        </div>

        {visibleItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "10px",
              padding: "9px 12px",
              borderRadius: "8px",
              cursor: "pointer",
              color: isActive(item.href)
                ? "hsl(var(--sidebar-text-active))"
                : "hsl(var(--sidebar-text))",
              fontSize: "13.5px",
              fontWeight: 500,
              transition: "background 0.15s, color 0.15s",
              marginBottom: "2px",
              textDecoration: "none",
              position: "relative",
              background: isActive(item.href)
                ? "hsl(var(--sidebar-active))"
                : "transparent",
            }}
          >
            {isActive(item.href) && (
              <span
                style={{
                  position: "absolute",
                  left: 0,
                  top: "50%",
                  transform: "translateY(-50%)",
                  width: "3px",
                  height: "20px",
                  background: "hsl(var(--primary))",
                  borderRadius: "0 3px 3px 0",
                }}
              />
            )}
            {item.icon}
            {item.label}
          </Link>
        ))}
      </nav>

      {/* User footer */}
      <div
        ref={menuRef}
        style={{
          padding: "12px 8px",
          borderTop: "1px solid rgba(255,255,255,0.07)",
          flexShrink: 0,
          position: "relative",
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "10px",
            padding: "10px 12px",
            borderRadius: "8px",
            cursor: "pointer",
            transition: "background 0.15s",
          }}
          onClick={() => setUserMenuOpen(!userMenuOpen)}
        >
          <div
            style={{
              width: "30px",
              height: "30px",
              borderRadius: "50%",
              background: user ? getAvatarColor(user.name) : "#3b82f6",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              fontSize: "12px",
              fontWeight: 700,
              color: "white",
              flexShrink: 0,
            }}
          >
            {getInitials(user?.name)}
          </div>
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ fontSize: "12.5px", fontWeight: 600, color: "#e5e7eb", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
              {user?.name || "用户"}
            </div>
            <div style={{ fontSize: "10.5px", color: "hsl(var(--sidebar-text))" }}>
              {roleLabel()}
            </div>
          </div>
          <div style={{ color: "rgba(255,255,255,0.3)", flexShrink: 0 }}>
            <SettingsIcon />
          </div>
        </div>

        {/* User dropdown menu */}
        {userMenuOpen && (
          <div
            style={{
              position: "absolute",
              bottom: "calc(100% - 8px)",
              left: "8px",
              right: "8px",
              background: "#1f2937",
              borderRadius: "10px",
              border: "1px solid rgba(255,255,255,0.1)",
              boxShadow: "0 -8px 24px rgba(0,0,0,0.3)",
              padding: "6px",
              zIndex: 50,
            }}
          >
            <Link
              href="/dashboard/profile"
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                padding: "8px 10px",
                borderRadius: "6px",
                fontSize: "13px",
                color: "#d1d5db",
                textDecoration: "none",
                transition: "background 0.15s",
              }}
              onClick={() => setUserMenuOpen(false)}
            >
              <SettingsIcon />
              {t("layout.profile")}
            </Link>

            {/* Language toggle */}
            <div
              style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                padding: "8px 10px",
                borderRadius: "6px",
                fontSize: "13px",
                color: "#d1d5db",
                cursor: "pointer",
                transition: "background 0.15s",
                position: "relative",
              }}
              onClick={() => setLangMenuOpen(!langMenuOpen)}
            >
              <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                <GlobeIcon />
                {t("layout.language")}
              </div>
              <ChevronRightIcon />

              {langMenuOpen && (
                <div
                  style={{
                    position: "absolute",
                    right: "calc(100% + 4px)",
                    top: 0,
                    background: "#1f2937",
                    borderRadius: "8px",
                    border: "1px solid rgba(255,255,255,0.1)",
                    padding: "4px",
                    minWidth: "140px",
                    boxShadow: "0 4px 16px rgba(0,0,0,0.3)",
                  }}
                >
                  {LanguagesSupported.map((locale) => (
                    <button
                      key={locale}
                      onClick={(e) => {
                        e.stopPropagation();
                        handleLanguageChange(locale);
                      }}
                      style={{
                        display: "block",
                        width: "100%",
                        textAlign: "left",
                        padding: "7px 10px",
                        borderRadius: "6px",
                        fontSize: "12.5px",
                        border: "none",
                        background: currentLocale === locale ? "rgba(59,130,246,0.2)" : "transparent",
                        color: currentLocale === locale ? "#60a5fa" : "#d1d5db",
                        cursor: "pointer",
                        fontFamily: "inherit",
                        fontWeight: currentLocale === locale ? 600 : 400,
                      }}
                    >
                      {locale === "en-US"
                        ? t("layout.navbar.langEnglish")
                        : t("layout.navbar.langSimplifiedChinese")}
                    </button>
                  ))}
                </div>
              )}
            </div>

            <div style={{ borderTop: "1px solid rgba(255,255,255,0.07)", margin: "4px 0" }} />

            <button
              onClick={() => { logout(); setUserMenuOpen(false); }}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                padding: "8px 10px",
                borderRadius: "6px",
                fontSize: "13px",
                color: "#f87171",
                border: "none",
                background: "transparent",
                cursor: "pointer",
                width: "100%",
                fontFamily: "inherit",
                transition: "background 0.15s",
              }}
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
