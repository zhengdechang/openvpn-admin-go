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
  const { t } = useTranslation("dashboard");
  const [name, setName] = useState(user?.name || "");
  const [email, setEmail] = useState(user?.email || "");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSaveProfile = async () => {
    setLoading(true);
    const success = await updateUserInfo({ name, email });
    if (success) toast.success(t("profile.updateSuccessToast"));
    else toast.error(t("profile.updateErrorToast"));
    setLoading(false);
  };

  const handleChangePassword = async () => {
    if (!password) {
      toast.error(t("profile.newPasswordRequired"));
      return;
    }
    setLoading(true);
    const success = await updateUserInfo({ password });
    if (success) {
      toast.success(t("profile.passwordUpdateSuccessToast"));
      setPassword("");
    } else {
      toast.error(t("profile.passwordUpdateErrorToast"));
    }
    setLoading(false);
  };

  return (
    <MainLayout className="p-4">
      <h1 className="text-2xl font-bold mb-4">{t("profile.pageTitle")}</h1>
      <div className="space-y-4 max-w-md">
        <div className="space-y-2">
          <h2 className="text-lg font-medium">{t("profile.basicInfoTitle")}</h2>
          <Input
            placeholder={t("profile.namePlaceholder")}
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <Input
            placeholder={t("profile.emailPlaceholder")}
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <Button onClick={handleSaveProfile} disabled={loading}>{t("profile.saveButton")}</Button>
        </div>
        <div className="space-y-2">
          <h2 className="text-lg font-medium">{t("profile.changePasswordTitle")}</h2>
          <Input
            type="password"
            placeholder={t("profile.newPasswordPlaceholder")}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button onClick={handleChangePassword} disabled={loading}>{t("profile.updatePasswordButton")}</Button>
        </div>
      </div>
    </MainLayout>
  );
}