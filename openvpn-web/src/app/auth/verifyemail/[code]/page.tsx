"use client";

import React, { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useTranslation } from "react-i18next";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api"; // 假设有一个 API 请求封装
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import MainLayout from "@/components/layout/main-layout";
import { useParams } from "next/navigation";

export default function VerifyEmailPage() {
  const router = useRouter();
  const params = useParams();
  const { t } = useTranslation();
  const [code, setCode] = useState(params.code || ""); // 预填充
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<boolean | null>(null);
  const [error, setError] = useState<string | null>(null);

  // 处理输入框变化
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCode(e.target.value);
  };

  // 处理邮箱验证
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
    <MainLayout className="flex justify-center items-center bg-gradient-to-br from-secondary/30 to-secondary/10 h-full">
      <div className="flex-grow flex h-full items-center justify-center p-4 ">
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">{t("auth.verifyemail.pageTitle")}</h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">{t("auth.verifyemail.pageSubtitle")}</p>
          </div>
          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="p-6 text-center">
              <h2 className="text-2xl font-bold text-primary">{t("auth.verifyemail.cardTitle")}</h2>
              <p className="text-gray-600 mt-2">{t("auth.verifyemail.cardDescription")}</p>

              {/* 验证码输入框 */}
              <Input
                type="text"
                value={code}
                onChange={handleInputChange}
                placeholder={t("auth.verifyemail.codePlaceholder")}
                className="mt-4 text-center"
              />

              {/* 验证按钮 */}
              <Button
                className="w-full mt-4"
                onClick={handleVerify}
                disabled={loading}
              >
                {loading ? t("auth.verifyemail.verifying") : t("auth.verifyemail.verifyButton")}
              </Button>

              {/* 结果显示 */}
              {success && (
                <p className="text-green-600 mt-3">{t("auth.verifyemail.successMessage")}</p>
              )}
              {error && <p className="text-red-500 mt-3">{t("auth.verifyemail.errorPrefix")}{error}</p>}

              {/* 失败时提供重新注册选项 */}
              {!success && (
                <Button
                  className="w-full mt-4"
                  variant="outline"
                  onClick={() => router.push("/auth/register")}
                >
                  {t("auth.verifyemail.registerAgain")}
                </Button>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </MainLayout>
  );
}
