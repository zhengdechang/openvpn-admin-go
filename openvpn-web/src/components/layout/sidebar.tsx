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
import { BrandIcon } from "@/components/ui/brand-icon";

interface NavItem {
  href: string;
  label: string;
  badge?: number;
  roles?: UserRole[];
}

interface NavGroup {
  key: string;
  label: string;
  icon: React.ReactNode;
  iconColor: string;
  items: NavItem[];
}

const ARGON_GRADIENT = "linear-gradient(87deg, #5e72e4 0%, #825ee4 100%)";

// LuCI (luci-theme-argon) 风格的一级菜单图标：实心 Font Awesome 图标。
// 状态=th-large 田字格，VPN=globe 地球，系统=cog 齿轮。颜色见 navGroups.iconColor。
function StatusGridIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
      <rect x="3" y="3" width="8" height="8" rx="1.5" />
      <rect x="13" y="3" width="8" height="8" rx="1.5" />
      <rect x="3" y="13" width="8" height="8" rx="1.5" />
      <rect x="13" y="13" width="8" height="8" rx="1.5" />
    </svg>
  );
}

function VpnGlobeIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 496 512" fill="currentColor">
      <path d="M336.5 160C322 70.7 287.8 8 248 8s-74 62.7-88.5 152h177zM152 256c0 22.2 1.2 43.5 3.3 64h185.3c2.1-20.5 3.3-41.8 3.3-64s-1.2-43.5-3.3-64H155.3c-2.1 20.5-3.3 41.8-3.3 64zm324.7-96c-28.6-67.9-86.5-120.4-158-141.6 24.4 33.8 41.2 84.7 50 141.6h108zM177.2 18.4C105.8 39.6 47.8 92.1 19.3 160h108c8.7-56.9 25.5-107.8 49.9-141.6zM487.4 192H372.7c2.1 21 3.3 42.5 3.3 64s-1.2 43-3.3 64h114.6c5.5-20.5 8.6-41.8 8.6-64s-3.1-43.5-8.5-64zM120 256c0-21.5 1.2-43 3.3-64H8.6C3.2 212.5 0 233.8 0 256s3.2 43.5 8.6 64h114.6c-2-21-3.2-42.5-3.2-64zm39.5 96c14.5 89.3 48.7 152 88.5 152s74-62.7 88.5-152h-177zm159.3 141.6c71.4-21.2 129.4-73.7 158-141.6h-108c-8.8 56.9-25.6 107.8-50 141.6zM19.3 352c28.6 67.9 86.5 120.4 158 141.6-24.4-33.8-41.2-84.7-50-141.6h-108z" />
    </svg>
  );
}

function SystemGearIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 512 512" fill="currentColor">
      <path d="M487.4 315.7l-42.6-24.6c4.3-23.2 4.3-47 0-70.2l42.6-24.6c4.9-2.8 7.1-8.6 5.5-14-11.1-35.6-30-67.8-54.7-94.6-3.8-4.1-10-5.1-14.8-2.3L380.8 110c-17.9-15.4-38.5-27.3-60.8-35.1V25.8c0-5.6-3.9-10.5-9.4-11.7-36.7-8.2-74.3-7.8-109.2 0-5.5 1.2-9.4 6.1-9.4 11.7V75c-22.2 7.9-42.8 19.8-60.8 35.1L88.7 85.5c-4.9-2.8-11-1.9-14.8 2.3-24.7 26.7-43.6 58.9-54.7 94.6-1.7 5.4.6 11.2 5.5 14L67.3 221c-4.3 23.2-4.3 47 0 70.2l-42.6 24.6c-4.9 2.8-7.1 8.6-5.5 14 11.1 35.6 30 67.8 54.7 94.6 3.8 4.1 10 5.1 14.8 2.3l42.6-24.6c17.9 15.4 38.5 27.3 60.8 35.1v49.2c0 5.6 3.9 10.5 9.4 11.7 36.7 8.2 74.3 7.8 109.2 0 5.5-1.2 9.4-6.1 9.4-11.7v-49.2c22.2-7.9 42.8-19.8 60.8-35.1l42.6 24.6c4.9 2.8 11 1.9 14.8-2.3 24.7-26.7 43.6-58.9 54.7-94.6 1.6-5.4-.6-11.2-5.5-14zM256 336c-44.1 0-80-35.9-80-80s35.9-80 80-80 80 35.9 80 80-35.9 80-80 80z" />
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

function ChevronDownIcon() {
  return (
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="6 9 12 15 18 9" />
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

interface SidebarProps {
  mobileOpen?: boolean;
  onClose?: () => void;
}

export default function Sidebar({ mobileOpen, onClose }: SidebarProps = {}) {
  const { user, logout } = useAuth();
  const pathname = usePathname();
  const { t, i18n } = useTranslation();
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const [langMenuOpen, setLangMenuOpen] = useState(false);
  const [currentLocale, setCurrentLocale] = useState<Locale>(getLocaleOnClient());
  const [toggledGroups, setToggledGroups] = useState<Record<string, boolean>>({});
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

  const navGroups: NavGroup[] = [
    {
      key: "status",
      label: t("layout.nav.status"),
      icon: <StatusGridIcon />,
      iconColor: "#5e72e4",
      items: [
        {
          href: "/dashboard/overview",
          label: t("layout.nav.overview"),
        },
      ],
    },
    {
      key: "vpn",
      label: t("layout.nav.vpn"),
      icon: <VpnGlobeIcon />,
      iconColor: "#11cdef",
      items: [
        {
          href: "/dashboard/users",
          label: t("dashboard.users.title"),
        },
        {
          href: "/dashboard/departments",
          label: t("dashboard.departments.title") || "部门管理",
          roles: [UserRole.ADMIN, UserRole.SUPERADMIN],
        },
      ],
    },
    {
      key: "system",
      label: t("layout.nav.system"),
      icon: <SystemGearIcon />,
      iconColor: "#fb6340",
      items: [
        {
          href: "/dashboard/server",
          label: t("dashboard.server.title"),
          roles: [UserRole.SUPERADMIN],
        },
        {
          href: "/dashboard/logs",
          label: t("dashboard.logs.titleServer"),
          roles: [UserRole.SUPERADMIN],
        },
        {
          href: "/dashboard/notifications",
          label: t("layout.nav.notifications"),
          roles: [UserRole.SUPERADMIN],
        },
      ],
    },
  ];

  // 角色过滤 + 隐藏空组
  const visibleGroups = navGroups
    .map((group) => ({
      ...group,
      items: group.items.filter(
        (item) =>
          !item.roles || (user && item.roles.includes(user.role as UserRole))
      ),
    }))
    .filter((group) => group.items.length > 0);

  const groupHasActive = (group: NavGroup) =>
    group.items.some((item) => isActive(item.href));

  const isGroupOpen = (group: NavGroup) =>
    toggledGroups[group.key] ?? groupHasActive(group);

  const toggleGroup = (group: NavGroup) =>
    setToggledGroups((prev) => ({
      ...prev,
      [group.key]: !(prev[group.key] ?? groupHasActive(group)),
    }));

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
      className={"dash-sidebar" + (mobileOpen ? " is-open" : "")}
      style={{
        width: "220px",
        minWidth: "220px",
        background: "hsl(var(--sidebar-bg))",
        borderRight: "1px solid hsl(var(--border))",
        display: "flex",
        flexDirection: "column",
        height: "125vh",
        // 不裁剪：语言子菜单向右溢出侧栏边界时需要露出来（nav 自身已有 overflowY 滚动）。
        overflow: "visible",
        // 抬升整个侧栏的层叠上下文，让底部用户菜单/语言子菜单向右溢出的部分盖在右侧内容之上，
        // 否则作为 flex 兄弟的右侧内容会盖住子菜单。
        position: "relative",
        zIndex: 50,
      }}
    >
      {/* Logo */}
      <div
        style={{
          height: "60px",
          minHeight: "60px",
          boxSizing: "border-box",
          padding: "0 20px",
          borderBottom: "1px solid hsl(var(--border))",
          display: "flex",
          alignItems: "center",
          gap: "10px",
          flexShrink: 0,
        }}
      >
        <BrandIcon size={34} />
        <div>
          <div style={{ fontSize: "15px", fontWeight: 700, color: "#1f2937", lineHeight: 1.2, letterSpacing: "0.2px" }}>
            Aegis
          </div>
          <div style={{ fontSize: "10px", color: "hsl(var(--muted-foreground))", fontWeight: 400, letterSpacing: "0.5px", textTransform: "uppercase" }}>
            VPN 控制台
          </div>
        </div>
      </div>

      {/* Nav */}
      <nav className="custom-scrollbar" style={{ flex: 1, padding: "12px 10px", overflowY: "auto" }}>
        {visibleGroups.map((group) => {
          const open = isGroupOpen(group);
          return (
            <div key={group.key} style={{ marginBottom: "4px" }}>
              {/* 组标题 */}
              <button
                onClick={() => toggleGroup(group)}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = "hsl(var(--sidebar-hover))";
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = "transparent";
                }}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "10px",
                  width: "100%",
                  padding: "9px 12px",
                  borderRadius: "8px",
                  border: "none",
                  background: "transparent",
                  cursor: "pointer",
                  color: "hsl(var(--sidebar-text))",
                  fontSize: "13.5px",
                  fontWeight: 600,
                  fontFamily: "inherit",
                  transition: "background 0.15s",
                }}
              >
                <span style={{ display: "flex", color: group.iconColor }}>{group.icon}</span>
                <span style={{ flex: 1, textAlign: "left" }}>{group.label}</span>
                <span
                  style={{
                    display: "flex",
                    color: "hsl(var(--muted-foreground))",
                    transition: "transform 0.2s",
                    transform: open ? "rotate(0deg)" : "rotate(-90deg)",
                  }}
                >
                  <ChevronDownIcon />
                </span>
              </button>

              {/* 子项 */}
              {open && (
                <div style={{ marginTop: "2px" }}>
                  {group.items.map((item) => {
                    const active = isActive(item.href);
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        onClick={() => onClose?.()}
                        onMouseEnter={(e) => {
                          if (!active) e.currentTarget.style.background = "hsl(var(--sidebar-hover))";
                        }}
                        onMouseLeave={(e) => {
                          if (!active) e.currentTarget.style.background = "transparent";
                        }}
                        style={{
                          display: "flex",
                          alignItems: "center",
                          gap: "10px",
                          padding: "8px 12px 8px 38px",
                          borderRadius: "8px",
                          cursor: "pointer",
                          color: active
                            ? "hsl(var(--sidebar-text-active))"
                            : "hsl(var(--sidebar-text))",
                          fontSize: "13px",
                          fontWeight: active ? 600 : 500,
                          transition: "background 0.15s, color 0.15s",
                          marginBottom: "2px",
                          textDecoration: "none",
                          background: active ? ARGON_GRADIENT : "transparent",
                          boxShadow: active ? "0 4px 10px rgba(94,114,228,0.35)" : "none",
                        }}
                      >
                        <span style={{ flex: 1 }}>{item.label}</span>
                        {item.badge ? (
                          <span
                            style={{
                              fontSize: "10px",
                              fontWeight: 700,
                              padding: "1px 6px",
                              borderRadius: "9px",
                              background: active ? "rgba(255,255,255,0.25)" : "hsl(var(--secondary))",
                              color: active ? "#fff" : "hsl(var(--secondary-foreground))",
                            }}
                          >
                            {item.badge}
                          </span>
                        ) : null}
                      </Link>
                    );
                  })}
                </div>
              )}
            </div>
          );
        })}
      </nav>

      {/* User footer */}
      <div
        ref={menuRef}
        style={{
          padding: "12px 10px",
          borderTop: "1px solid hsl(var(--border))",
          flexShrink: 0,
          position: "relative",
        }}
      >
        <div
          onMouseEnter={(e) => {
            e.currentTarget.style.background = "hsl(var(--sidebar-hover))";
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.background = "transparent";
          }}
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
              width: "32px",
              height: "32px",
              borderRadius: "50%",
              background: user ? getAvatarColor(user.name) : "#5e72e4",
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
            <div style={{ fontSize: "12.5px", fontWeight: 600, color: "#1f2937", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
              {user?.name || "用户"}
            </div>
            <div style={{ fontSize: "10.5px", color: "hsl(var(--muted-foreground))" }}>
              {roleLabel()}
            </div>
          </div>
          <div style={{ color: "hsl(var(--muted-foreground))", flexShrink: 0 }}>
            <SettingsIcon />
          </div>
        </div>

        {/* User dropdown menu */}
        {userMenuOpen && (
          <div
            style={{
              position: "absolute",
              bottom: "calc(100% - 8px)",
              left: "10px",
              right: "10px",
              background: "#ffffff",
              borderRadius: "10px",
              border: "1px solid hsl(var(--border))",
              boxShadow: "0 -8px 24px rgba(15,23,42,0.12)",
              padding: "6px",
              zIndex: 50,
            }}
          >
            <Link
              href="/dashboard/profile"
              onMouseEnter={(e) => { e.currentTarget.style.background = "hsl(var(--sidebar-hover))"; }}
              onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                padding: "8px 10px",
                borderRadius: "6px",
                fontSize: "13px",
                color: "#374151",
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
              onMouseEnter={(e) => { e.currentTarget.style.background = "hsl(var(--sidebar-hover))"; }}
              onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}
              style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                padding: "8px 10px",
                borderRadius: "6px",
                fontSize: "13px",
                color: "#374151",
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
                    left: "calc(100% + 4px)",
                    top: 0,
                    background: "#ffffff",
                    borderRadius: "8px",
                    border: "1px solid hsl(var(--border))",
                    padding: "4px",
                    minWidth: "140px",
                    boxShadow: "0 4px 16px rgba(15,23,42,0.12)",
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
                        background: currentLocale === locale ? "hsl(var(--sidebar-hover))" : "transparent",
                        color: currentLocale === locale ? "hsl(var(--primary))" : "#374151",
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

            <div style={{ borderTop: "1px solid hsl(var(--border))", margin: "4px 0" }} />

            <button
              onMouseEnter={(e) => { e.currentTarget.style.background = "rgba(239,68,68,0.08)"; }}
              onMouseLeave={(e) => { e.currentTarget.style.background = "transparent"; }}
              onClick={() => { logout(); setUserMenuOpen(false); }}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                padding: "8px 10px",
                borderRadius: "6px",
                fontSize: "13px",
                color: "#dc2626",
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
