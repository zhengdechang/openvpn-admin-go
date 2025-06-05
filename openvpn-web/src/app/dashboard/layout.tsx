"use client";

import React, { ReactNode, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useTranslation } from "react-i18next";

interface DashboardLayoutProps {
  children: ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation("common"); // Explicitly use common namespace

  useEffect(() => {
    if (!loading && !user) {
      router.replace("/auth/login");
    }
  }, [loading, user, router]);
  if (loading || !user) {
    return <div className="p-4">{t("loading")}</div>;
  }
  return <>{children}</>;
}