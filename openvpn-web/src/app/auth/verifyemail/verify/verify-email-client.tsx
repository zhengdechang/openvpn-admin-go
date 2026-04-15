"use client";

import React, { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useTranslation } from "react-i18next";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api"; // 假设有一个 API 请求封装
import { Card, CardContent } from "@/components/ui/card";
import AuthLayout from "@/components/layout/auth-layout";
import MuiButton from "@mui/material/Button";
import TextField from "@mui/material/TextField";

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

              {/* 验证码输入框 */}
              <TextField
                type="text"
                value={code}
                onChange={handleInputChange}
                label={t("auth.verifyemail.codePlaceholder")}
                placeholder={t("auth.verifyemail.codePlaceholder")}
                fullWidth
                sx={{ mt: 2 }}
              />

              {/* 验证按钮 */}
              <MuiButton
                variant="contained"
                fullWidth
                onClick={handleVerify}
                disabled={loading}
                sx={{ mt: 2 }}
              >
                {loading ? t("auth.verifyemail.verifying") : t("auth.verifyemail.verifyButton")}
              </MuiButton>

              {/* 结果显示 */}
              {success && (
                <p className="text-green-600 mt-3">{t("auth.verifyemail.successMessage")}</p>
              )}
              {error && <p className="text-red-500 mt-3">{t("auth.verifyemail.errorPrefix")}{error}</p>}

              {/* 失败时提供重新注册选项 */}
              {!success && (
                <MuiButton
                  variant="outlined"
                  fullWidth
                  onClick={() => router.push("/auth/register")}
                  sx={{ mt: 2 }}
                >
                  {t("auth.verifyemail.registerAgain")}
                </MuiButton>
              )}
            </CardContent>
          </Card>
        </div>
    </AuthLayout>
  );
}
