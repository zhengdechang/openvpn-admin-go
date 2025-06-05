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
    <div className="min-h-screen flex flex-col">
      <Navbar />

      <main className={`flex-grow ${className}`}>{children}</main>

      {showFooter && (
        <footer className="bg-white border-t border-gray-200 py-6">
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
  );
}
