/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-05 13:07:03
 */
"use client";

import React, { useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import { ArgonField } from "@/components/ui/argon-field";
import { Button } from "@/components/ui/button";
import { CbiSection, CbiValue } from "@/components/ui/cbi-form";
import { toast } from "sonner";

export default function ProfilePage() {
  const { user, updateUserInfo } = useAuth();
  const { t } = useTranslation();
  const [name, setName] = useState(user?.name || "");
  const [email, setEmail] = useState(user?.email || "");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSaveProfile = async () => {
    setLoading(true);
    const success = await updateUserInfo({ name, email });
    if (success) toast.success(t("dashboard.profile.updateSuccessToast"));
    else toast.error(t("dashboard.profile.updateErrorToast"));
    setLoading(false);
  };

  const handleChangePassword = async () => {
    if (!password) {
      toast.error(t("dashboard.profile.newPasswordRequired"));
      return;
    }
    setLoading(true);
    const success = await updateUserInfo({ password });
    if (success) {
      toast.success(t("dashboard.profile.passwordUpdateSuccessToast"));
      setPassword("");
    } else {
      toast.error(t("dashboard.profile.passwordUpdateErrorToast"));
    }
    setLoading(false);
  };

  return (
    <MainLayout className="p-6 space-y-6">
      <div className="mx-auto w-full max-w-3xl space-y-6">
        {/* 基本信息 */}
        <CbiSection title={t("dashboard.profile.basicInfoTitle")}>
          <CbiValue title={t("dashboard.profile.namePlaceholder")} htmlFor="profile-name">
            <div className="max-w-md">
              <ArgonField
                id="profile-name"
                placeholder={t("dashboard.profile.namePlaceholder")}
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled
              />
            </div>
          </CbiValue>
          <CbiValue title={t("dashboard.profile.emailPlaceholder")} htmlFor="profile-email">
            <div className="max-w-md">
              <ArgonField
                id="profile-email"
                placeholder={t("dashboard.profile.emailPlaceholder")}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled
              />
            </div>
          </CbiValue>
        </CbiSection>

        {/* 修改密码 */}
        <CbiSection title={t("dashboard.profile.changePasswordTitle")}>
          <CbiValue title={t("dashboard.profile.newPasswordPlaceholder")} htmlFor="profile-password">
            <div className="max-w-md">
              <ArgonField
                id="profile-password"
                type="password"
                placeholder={t("dashboard.profile.newPasswordPlaceholder")}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
          </CbiValue>
          <div className="cbi-value">
            <div className="cbi-value-title" />
            <div className="cbi-value-field">
              <Button onClick={handleChangePassword} disabled={loading}>
                {loading
                  ? t("common.loading")
                  : t("dashboard.profile.updatePasswordButton")}
              </Button>
            </div>
          </div>
        </CbiSection>
      </div>
    </MainLayout>
  );
}
