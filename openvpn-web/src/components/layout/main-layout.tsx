/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-05 13:07:03
 */
"use client";

import React from "react";
import Navbar from "./navbar";
import { useTranslation } from "react-i18next";

interface MainLayoutProps {
  children: React.ReactNode;
  showFooter?: boolean;
  className?: string;
}

export default function MainLayout({
  children,
  showFooter = true,
  className,
}: MainLayoutProps) {
  const { t } = useTranslation();
  const currentYear = new Date().getFullYear();

  return (
    <div className="h-screen flex flex-col overflow-hidden">
      {/* 固定的导航栏 */}
      <div className="flex-shrink-0">
        <Navbar />
      </div>

      {/* 可滚动的内容区域 */}
      <div className="flex-1 overflow-y-auto custom-scrollbar">
        <div className="min-h-full flex flex-col">
          <main className={`flex-1 ${className}`}>{children}</main>

          {/* 底部footer - 在内容区域底部，不固定 */}
          {showFooter && (
            <footer className="bg-white border-t border-gray-200 py-6 mt-auto">
              <div className="container mx-auto px-4">
                <div className="text-center">
                  <p className="mb-2">
                    {t("layout.footer.copyrightText", { year: currentYear })}
                  </p>
                  <p className="text-sm text-gray-600">
                    {t("layout.footer.builtWithText")}
                  </p>
                </div>
              </div>
            </footer>
          )}
        </div>
      </div>
    </div>
  );
}
