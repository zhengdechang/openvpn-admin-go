/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-03-20 17:22:38
 */
"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslation } from "react-i18next";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api"; // 假设有一个 API 请求封装
import { Card, CardContent } from "@/components/ui/card";
import AuthLayout from "@/components/layout/auth-layout";
import MuiButton from "@mui/material/Button";
import TextField from "@mui/material/TextField";

export default function ResetPasswordPage() {
  const router = useRouter();
  const { t } = useTranslation();
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<boolean | null>(null);
  const [error, setError] = useState<string | null>(null);

  // 处理输入框变化
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  // 处理重置密码请求
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
        showToast.success(
          response.message || t("auth.forgotpassword.resetEmailSent")
        );
      } else {
        setSuccess(false);
        console.log(response, "response");
        setError(response.error || t("auth.forgotpassword.resetFailed"));
      }
    } catch (err) {
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
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">{t("auth.forgotpassword.pageSubtitle")}</p>
          </div>
          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="p-6 text-center">
              <h2 className="text-2xl font-bold text-primary">{t("auth.forgotpassword.cardTitle")}</h2>
              <p className="text-gray-600 mt-2">
                {t("auth.forgotpassword.cardDescription")}
              </p>

              {/* 邮箱输入框 */}
              <TextField
                type="email"
                value={email}
                onChange={handleInputChange}
                placeholder={t("auth.forgotpassword.emailPlaceholder")}
                label={t("auth.forgotpassword.emailPlaceholder")}
                fullWidth
                className="mt-4"
                sx={{ mt: 2 }}
              />

              {/* 提交按钮 */}
              <MuiButton
                variant="contained"
                fullWidth
                onClick={handleResetPassword}
                disabled={loading}
                sx={{ mt: 2 }}
              >
                {loading ? t("auth.forgotpassword.sending") : t("auth.forgotpassword.sendResetEmail")}
              </MuiButton>

              {/* 结果显示 */}
              {success && (
                <p className="text-green-600 mt-3">
                  {t("auth.forgotpassword.emailSentSuccess")}
                </p>
              )}
              {error && <p className="text-red-500 mt-3">{t("auth.forgotpassword.errorPrefix")}{error}</p>}

              {/* 返回登录页 */}
              <MuiButton
                variant="outlined"
                fullWidth
                onClick={() => router.push("/auth/login")}
                sx={{ mt: 2 }}
              >
                {t("auth.forgotpassword.backToLogin")}
              </MuiButton>
            </CardContent>
          </Card>
        </div>
    </AuthLayout>
  );
}
