"use client";

import React, { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { UserRole } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { setLocaleOnClient, getLocaleOnClient } from "../../../i18n";
import { LanguagesSupported } from "../../../i18n/language";
import type { Locale } from "../../../i18n";

export default function Navbar() {
  const { user, loading, logout } = useAuth();
  const pathname = usePathname();
  const { t, i18n } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const [isLangOpen, setIsLangOpen] = useState(false);
  const [currentLocale, setCurrentLocale] = useState<Locale>(getLocaleOnClient());
  const langDropdownRef = useRef<HTMLDivElement>(null);

  // 判断当前链接是否激活
  const isActive = (path: string) => {
    return pathname === path || pathname?.startsWith(`${path}/`);
  };

  const handleLanguageChange = (locale: Locale) => {
    setLocaleOnClient(locale, false);
    setCurrentLocale(locale);
    setIsLangOpen(false);
  };

  // 点击外部区域关闭下拉菜单
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (langDropdownRef.current && !langDropdownRef.current.contains(event.target as Node)) {
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

    i18n.on('languageChanged', handleLanguageChanged);
    return () => {
      i18n.off('languageChanged', handleLanguageChanged);
    };
  }, [i18n]);

  return (
    <header className="bg-white shadow-sm border-b border-gray-200 py-4 sticky top-0 z-10">
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
          OpenVPN 管理系统
          </Link>

          <nav className="flex items-center space-x-6">
            <Link
              href="/"
              className={`text-gray-600 hover:text-primary ${
                isActive("/") ? "font-medium text-primary" : ""
              }`}
            >
              {t('layout.home')}
            </Link>

            <Link
              href="/dashboard"
              className={`text-gray-600 hover:text-primary ${
                isActive("/dashboard") ? "font-medium text-primary" : ""
              }`}
            >
              {t('layout.dashboard')}
            </Link>

            {/* Language Switch */}
            <div className="relative" ref={langDropdownRef}>
              <button
                className="flex items-center space-x-1 text-gray-600 hover:text-primary"
                onClick={() => setIsLangOpen(!isLangOpen)}
              >
                <span>{t('layout.language')}</span>
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
                <div 
                  className="absolute right-0 mt-1 pt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50"
                >
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
                      {locale === "en-US" ? "English" : "简体中文"}
                    </button>
                  ))}
                </div>
              )}
            </div>

            {!loading &&
              (user ? (
                <div className="flex items-center space-x-4">
                  <div
                    className="relative"
                    onMouseEnter={() => setIsOpen(true)}
                    onMouseLeave={() => setIsOpen(false)}
                  >
                    <button className="flex items-center space-x-1 text-gray-600 hover:text-primary">
                      <span>
                        {user.name || "User"}
                        {user.role === UserRole.ADMIN && " (Admin)"}
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
                          {t('layout.dashboard')}
                        </Link>
                        <Link
                          href="/dashboard/profile"
                          className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                        >
                          {t('layout.profile')}
                        </Link>
                        <button
                          onClick={logout}
                          className="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100"
                        >
                          {t('layout.logout')}
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              ) : (
                <div className="flex items-center space-x-4">
                  <Link
                    href="/auth/login"
                    className={`text-gray-600 hover:text-primary ${
                      isActive("/auth/login") ? "font-medium text-primary" : ""
                    }`}
                  >
                    {t('layout.login')}
                  </Link>
                  <Button asChild size="sm">
                    <Link href="/auth/register">{t('layout.register')}</Link>
                  </Button>
                </div>
              ))}
          </nav>
        </div>
      </div>
    </header>
  );
}
