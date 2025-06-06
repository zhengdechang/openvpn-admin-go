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
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
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
    <MainLayout className="p-4 h-[calc(80vh-4rem)]">
      <div className="flex items-center justify-center h-full">
        <div className="flex flex-col items-center w-full max-w-md">
          <h1 className="text-2xl font-bold mb-4">
            {t("dashboard.profile.pageTitle")}
          </h1>
          <div className="space-y-4 w-full">
            <div className="space-y-2">
              <h2 className="text-lg font-medium">
                {t("dashboard.profile.basicInfoTitle")}
              </h2>
              <Input
                placeholder={t("dashboard.profile.namePlaceholder")}
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={true}
              />
              <Input
                placeholder={t("dashboard.profile.emailPlaceholder")}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={true}
              />
            </div>
            <div className="space-y-2">
              <h2 className="text-lg font-medium">
                {t("dashboard.profile.changePasswordTitle")}
              </h2>
              <Input
                type="password"
                placeholder={t("dashboard.profile.newPasswordPlaceholder")}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              <Button
                onClick={handleChangePassword}
                disabled={loading}
                className="w-full"
              >
                {t("dashboard.profile.updatePasswordButton")}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  );
}
