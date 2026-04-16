"use client";

import React, { useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

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
    <MainLayout className="flex justify-center items-center bg-gradient-to-br from-secondary/30 to-secondary/10 h-full">
      <div className="flex-grow flex h-full items-center justify-center p-4">
        <div className="w-full max-w-lg my-4">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">
              {t("dashboard.profile.pageTitle")}
            </h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">
              {t("dashboard.profile.basicInfoTitle")}
            </p>
          </div>

          <div className="bg-white rounded-lg shadow-lg border-t-4 border-t-primary p-6">
            <div className="space-y-6">
              {/* 基本信息部分 */}
              <div>
                <h2 className="text-lg font-medium mb-4 text-primary">
                  {t("dashboard.profile.basicInfoTitle")}
                </h2>
                <div className="space-y-4">
                  <div className="space-y-1">
                    <Label htmlFor="profile-name">{t("dashboard.profile.namePlaceholder")}</Label>
                    <Input
                      id="profile-name"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      disabled
                      placeholder={t("dashboard.profile.namePlaceholder")}
                      className="bg-gray-50"
                    />
                  </div>
                  <div className="space-y-1">
                    <Label htmlFor="profile-email">{t("dashboard.profile.emailPlaceholder")}</Label>
                    <Input
                      id="profile-email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      disabled
                      placeholder={t("dashboard.profile.emailPlaceholder")}
                      className="bg-gray-50"
                    />
                  </div>
                </div>
              </div>

              {/* 分隔线 */}
              <div className="border-t border-gray-200"></div>

              {/* 修改密码部分 */}
              <div>
                <h2 className="text-lg font-medium mb-4 text-primary">
                  {t("dashboard.profile.changePasswordTitle")}
                </h2>
                <div className="space-y-4">
                  <div className="space-y-1">
                    <Label htmlFor="profile-password">{t("dashboard.profile.newPasswordPlaceholder")}</Label>
                    <Input
                      id="profile-password"
                      type="password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      placeholder={t("dashboard.profile.newPasswordPlaceholder")}
                    />
                  </div>
                  <Button className="w-full" onClick={handleChangePassword} disabled={loading}>
                    {loading ? t("common.loading") : t("dashboard.profile.updatePasswordButton")}
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  );
}
