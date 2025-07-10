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

export default function VerifyEmailClient() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { t } = useTranslation();
  const [code, setCode] = useState(searchParams.get('code') || ""); // 从查询参数获取
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<boolean | null>(null);
  const [error, setError] = useState<string | null>(null);

  // 处理输入变化
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
    <MainLayout>
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
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
