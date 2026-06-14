"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslation } from "react-i18next";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api";
import { Card, CardContent } from "@/components/ui/card";
import AuthLayout from "@/components/layout/auth-layout";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function ResetPasswordPage() {
  const router = useRouter();
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<boolean | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleResetPassword = async () => {
    if (!email) {
      showToast.error(t("auth.forgotpassword.emailRequired"));
      return;
    }
    setLoading(true);
    setSuccess(null);
    setError(null);
    try {
      const response = await userAPI.forgotPassword(email);
      if (response.success) {
        setSuccess(true);
        showToast.success(response.message || t("auth.forgotpassword.resetEmailSent"));
      } else {
        setSuccess(false);
        setError(response.error || t("auth.forgotpassword.resetFailed"));
      }
    } catch {
      setSuccess(false);
      setError(t("auth.forgotpassword.requestFailed"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary">{t("auth.forgotpassword.pageTitle")}</h1>
          <div className="h-1 w-16 bg-primary mx-auto my-4" />
          <p className="text-gray-600">{t("auth.forgotpassword.pageSubtitle")}</p>
        </div>
        <Card className="shadow-lg border-t-4 border-t-primary">
          <CardContent className="p-6 space-y-4">
            <h2 className="text-2xl font-bold text-primary text-center">{t("auth.forgotpassword.cardTitle")}</h2>
            <p className="text-gray-600 text-center">{t("auth.forgotpassword.cardDescription")}</p>
            <div className="space-y-1">
              <Label htmlFor="forgot-email">{t("auth.forgotpassword.emailPlaceholder")}</Label>
              <Input
                id="forgot-email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder={t("auth.forgotpassword.emailPlaceholder")}
                className="w-full"
              />
            </div>
            <Button className="w-full" onClick={handleResetPassword} disabled={loading}>
              {loading ? t("auth.forgotpassword.sending") : t("auth.forgotpassword.sendResetEmail")}
            </Button>
            {success && <p className="text-green-600 text-center">{t("auth.forgotpassword.emailSentSuccess")}</p>}
            {error && <p className="text-destructive text-center">{t("auth.forgotpassword.errorPrefix")}{error}</p>}
            <Button variant="outline" className="w-full" onClick={() => router.push("/auth/login")}>
              {t("auth.forgotpassword.backToLogin")}
            </Button>
          </CardContent>
        </Card>
      </div>
    </AuthLayout>
  );
}
