"use client";

import React, { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { UserRole } from "@/types/types";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { setLocaleOnClient, getLocaleOnClient } from "@/i18n";
import { LanguagesSupported } from "@/i18n/language";
import type { Locale } from "@/i18n";

export default function Navbar() {
  const { user, loading, logout } = useAuth();
  const pathname = usePathname();
  const { t, i18n } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const [isLangOpen, setIsLangOpen] = useState(false);
  const [currentLocale, setCurrentLocale] = useState<Locale>(
    getLocaleOnClient()
  );
  const langDropdownRef = useRef<HTMLDivElement>(null);

  // 判断当前链接是否激活
  const isActive = (path: string) => {
    return pathname === path || pathname?.startsWith(`${path}/`);
  };

  const handleLanguageChange = (locale: Locale) => {
    setLocaleOnClient(locale, true);
    setCurrentLocale(locale);
    setIsLangOpen(false);
    setIsOpen(false); // 关闭用户菜单

    // 立即更新 DOM 中的语言设置
    document.documentElement.lang = locale;
    document.documentElement.setAttribute("data-locale", locale);

    // 使用 i18n 切换语言，无需刷新页面
    i18n.changeLanguage(locale);
  };

  // 点击外部区域关闭下拉菜单
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        langDropdownRef.current &&
        !langDropdownRef.current.contains(event.target as Node)
      ) {
        setIsLangOpen(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  // 监听语言变化
  useEffect(() => {
    const handleLanguageChanged = () => {
      setCurrentLocale(i18n.language as Locale);
    };

    i18n.on("languageChanged", handleLanguageChanged);
    return () => {
      i18n.off("languageChanged", handleLanguageChanged);
    };
  }, [i18n]);

  // 初始化时同步语言设置
  useEffect(() => {
    const savedLocale = getLocaleOnClient();
    if (savedLocale !== i18n.language) {
      i18n.changeLanguage(savedLocale);
    }
  }, []);

  return (
    <div className="bg-white shadow-sm border-b border-gray-200 py-4 sticky top-0 z-10">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between">
          <Link
            href="/"
            className="text-xl font-bold text-primary flex items-center"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6 mr-2"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
            </svg>
            {t("layout.navbar.logoText")}
          </Link>

          <nav className="flex items-center space-x-6">
            {/* Authenticated user links: User management, departments, server, logs */}
            {user && (
              <>
                {
                  <Link
                    href="/dashboard/users"
                    className={`text-gray-600 hover:text-primary ${
                      isActive("/dashboard/users")
                        ? "font-medium text-primary"
                        : ""
                    }`}
                  >
                    {t("dashboard.users.title")}
                  </Link>
                }
                {(user.role === UserRole.ADMIN ||
                  user.role === UserRole.SUPERADMIN) && (
                  <Link
                    href="/dashboard/departments"
                    className={`text-gray-600 hover:text-primary ${
                      isActive("/dashboard/departments")
                        ? "font-medium text-primary"
                        : ""
                    }`}
                  >
                    {t("dashboard.departments.title") ||
                      t("navbar.departmentsFallback")}
                  </Link>
                )}
                {user.role === UserRole.SUPERADMIN && (
                  <>
                    <Link
                      href="/dashboard/server"
                      className={`text-gray-600 hover:text-primary ${
                        isActive("/dashboard/server")
                          ? "font-medium text-primary"
                          : ""
                      }`}
                    >
                      {t("dashboard.server.title")}
                    </Link>
                    <Link
                      href="/dashboard/logs"
                      className={`text-gray-600 hover:text-primary ${
                        isActive("/dashboard/logs")
                          ? "font-medium text-primary"
                          : ""
                      }`}
                    >
                      {t("dashboard.logs.titleServer")}
                    </Link>
                  </>
                )}
              </>
            )}
            {/* Unauthenticated: show login/register */}
            {!loading && !user && (
              <>
                <Link
                  href="/auth/login"
                  className="text-gray-600 hover:text-primary"
                >
                  {t("layout.login")}
                </Link>
                <Link
                  href="/auth/register"
                  className="text-gray-600 hover:text-primary"
                >
                  {t("layout.register")}
                </Link>
                {/* Language Switch for unauthenticated users */}
                <div className="relative" ref={langDropdownRef}>
                  <button
                    className="flex items-center space-x-1 text-gray-600 hover:text-primary"
                    onClick={() => setIsLangOpen(!isLangOpen)}
                  >
                    <span>{t("layout.language")}</span>
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 9l-7 7-7-7"
                      />
                    </svg>
                  </button>
                  {isLangOpen && (
                    <div className="absolute right-0 mt-1 pt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50">
                      {LanguagesSupported.map((locale) => (
                        <button
                          key={locale}
                          onClick={() => handleLanguageChange(locale)}
                          className={`block w-full text-left px-4 py-2 text-sm ${
                            currentLocale === locale
                              ? "text-primary font-medium"
                              : "text-gray-700 hover:bg-gray-100"
                          }`}
                        >
                          {locale === "en-US"
                            ? t("layout.navbar.langEnglish")
                            : t("layout.navbar.langSimplifiedChinese")}
                        </button>
                      ))}
                    </div>
                  )}
                </div>
              </>
            )}

            {/* User menu */}
            {!loading && user && (
              <div className="flex items-center space-x-4">
                <div
                  className="relative"
                  onMouseEnter={() => setIsOpen(true)}
                  onMouseLeave={() => setIsOpen(false)}
                >
                  <button className="flex items-center space-x-1 text-gray-600 hover:text-primary">
                    <span>
                      {user.name || t("layout.navbar.userDefaultName")}
                      {user.role === UserRole.ADMIN &&
                        t("layout.navbar.adminSuffix")}
                    </span>
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 9l-7 7-7-7"
                      />
                    </svg>
                  </button>
                  {isOpen && (
                    <div className="absolute right-0 mt-0 pt-2 w-48 bg-white rounded-md shadow-lg py-1">
                      <Link
                        href="/dashboard"
                        className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      >
                        {t("layout.dashboard")}
                      </Link>
                      <Link
                        href="/dashboard/profile"
                        className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      >
                        {t("layout.profile")}
                      </Link>

                      {/* Language Selection */}
                      <div className="border-t border-gray-100 my-1"></div>
                      <div className="relative group">
                        <div className="flex items-center justify-between px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 cursor-pointer">
                          <span>{t("layout.language")}</span>
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="h-4 w-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M9 5l7 7-7 7"
                            />
                          </svg>
                        </div>
                        <div className="absolute right-full top-0 mr-1 w-48 bg-white rounded-md shadow-lg py-1 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
                          {LanguagesSupported.map((locale) => (
                            <button
                              key={locale}
                              onClick={() => handleLanguageChange(locale)}
                              className={`block w-full text-left px-4 py-2 text-sm ${
                                currentLocale === locale
                                  ? "text-primary font-medium bg-blue-50"
                                  : "text-gray-700 hover:bg-gray-100"
                              }`}
                            >
                              {locale === "en-US"
                                ? t("layout.navbar.langEnglish")
                                : t("layout.navbar.langSimplifiedChinese")}
                            </button>
                          ))}
                        </div>
                      </div>

                      <div className="border-t border-gray-100 my-1"></div>
                      <button
                        onClick={logout}
                        className="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100"
                      >
                        {t("layout.logout")}
                      </button>
                    </div>
                  )}
                </div>
              </div>
            )}
          </nav>
        </div>
      </div>
    </div>
  );
}
