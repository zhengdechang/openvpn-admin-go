"use client";

import React, { useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useTranslation } from "react-i18next";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api";
import { Card, CardContent } from "@/components/ui/card";
import AuthLayout from "@/components/layout/auth-layout";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function VerifyEmailClient() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { t } = useTranslation();
  const [code, setCode] = useState(searchParams.get('code') || "");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<boolean | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCode(e.target.value);
  };

  const handleVerify = async () => {
    if (!code) {
      showToast.error(t("auth.verifyemail.codeRequired"));
      return;
    }

    setLoading(true);
    setSuccess(null);
    setError(null);

    try {
      const response = await userAPI.verifyEmail(code as string);
      if (response.success) {
        setSuccess(true);
        showToast.success(t("auth.verifyemail.successToast"));
        setTimeout(() => {
          router.push(`/auth/login`);
        }, 1000);
      } else {
        setSuccess(false);
        setError(response.error || t("auth.verifyemail.failToast"));
      }
    } catch (err) {
      setSuccess(false);
      setError(t("auth.verifyemail.requestFailedToast"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
      <div className="max-w-md w-full space-y-8">
        <Card>
          <CardContent className="p-6">
            <div className="text-center">
              <h2 className="text-2xl font-bold mb-4">
                {t("auth.verifyemail.title")}
              </h2>
              <p className="text-gray-600 mb-6">
                {t("auth.verifyemail.description")}
              </p>
            </div>

            <div className="space-y-1 mt-2">
              <Label htmlFor="verify-code">{t("auth.verifyemail.codePlaceholder")}</Label>
              <Input
                id="verify-code"
                type="text"
                value={code}
                onChange={handleInputChange}
                placeholder={t("auth.verifyemail.codePlaceholder")}
                className="w-full"
              />
            </div>

            <Button className="w-full mt-4" onClick={handleVerify} disabled={loading}>
              {loading ? t("auth.verifyemail.verifying") : t("auth.verifyemail.verifyButton")}
            </Button>

            {success && (
              <p className="text-green-600 mt-3">{t("auth.verifyemail.successMessage")}</p>
            )}
            {error && <p className="text-destructive mt-3">{t("auth.verifyemail.errorPrefix")}{error}</p>}

            {!success && (
              <Button variant="outline" className="w-full mt-2" onClick={() => router.push("/auth/register")}>
                {t("auth.verifyemail.registerAgain")}
              </Button>
            )}
          </CardContent>
        </Card>
      </div>
    </AuthLayout>
  );
}
