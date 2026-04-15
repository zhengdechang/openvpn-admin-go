"use client";

import React, { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useTranslation } from "react-i18next";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api";
import { Card, CardContent } from "@/components/ui/card";
import AuthLayout from "@/components/layout/auth-layout";
import MuiButton from "@mui/material/Button";
import TextField from "@mui/material/TextField";

type ResetPasswordFormData = z.infer<ReturnType<typeof getResetPasswordSchema>>;

// 定义表单验证规则
const getResetPasswordSchema = (t: Function) => z
  .object({
    password: z.string().min(6, t("auth.resetpassword.passwordMinLength")),
    confirmPassword: z.string().min(6, t("auth.resetpassword.confirmPasswordMinLength")),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: t("auth.resetpassword.passwordsNotMatch"),
    path: ["confirmPassword"],
  });

export default function ResetPasswordClient() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { t } = useTranslation();
  const token = searchParams.get('code') || "";
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<boolean | null>(null);
  const [error, setError] = useState<string | null>(null);

  // 使用 react-hook-form 进行表单管理
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetPasswordFormData>({
    resolver: zodResolver(getResetPasswordSchema(t)),
  });

  // 处理密码重置
  const onSubmit = async (data: ResetPasswordFormData) => {
    if (!token) {
      showToast.error(t("auth.resetpassword.invalidLink"));
      return;
    }

    setLoading(true);
    setSuccess(null);
    setError(null);

    try {
      const response = await userAPI.resetPassword(token, data.password, data.confirmPassword);
      if (response.success) {
        setSuccess(true);
        showToast.success(t("auth.resetpassword.successToast"));
        setTimeout(() => {
          router.push("/auth/login");
        }, 2000);
      } else {
        setSuccess(false);
        setError(response.error || t("auth.resetpassword.failToast"));
      }
    } catch (err) {
      setSuccess(false);
      setError(t("auth.resetpassword.requestFailedToast"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">{t("auth.resetpassword.pageTitle")}</h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">{t("auth.resetpassword.pageSubtitle")}</p>
          </div>
          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="p-6">
              <h2 className="text-2xl font-bold text-primary text-center">{t("auth.resetpassword.cardTitle")}</h2>
              <p className="text-gray-600 mt-2 text-center">{t("auth.resetpassword.cardDescription")}</p>

              <form onSubmit={handleSubmit(onSubmit)} className="mt-6 space-y-4">
                {/* 新密码输入框 */}
                <TextField
                  type="password"
                  label={t("auth.resetpassword.passwordPlaceholder")}
                  placeholder={t("auth.resetpassword.passwordPlaceholder")}
                  fullWidth
                  error={!!errors.password}
                  helperText={errors.password?.message}
                  {...register("password")}
                />

                {/* 确认密码输入框 */}
                <TextField
                  type="password"
                  label={t("auth.resetpassword.confirmPasswordPlaceholder")}
                  placeholder={t("auth.resetpassword.confirmPasswordPlaceholder")}
                  fullWidth
                  error={!!errors.confirmPassword}
                  helperText={errors.confirmPassword?.message}
                  {...register("confirmPassword")}
                />

                {/* 重置按钮 */}
                <MuiButton
                  type="submit"
                  variant="contained"
                  fullWidth
                  disabled={loading}
                >
                  {loading ? t("auth.resetpassword.resetting") : t("auth.resetpassword.resetButton")}
                </MuiButton>

                {/* 结果显示 */}
                {success && (
                  <p className="text-green-600 text-center">{t("auth.resetpassword.successMessage")}</p>
                )}
                {error && (
                  <p className="text-red-500 text-center">{t("auth.resetpassword.errorPrefix")}{error}</p>
                )}

                {/* 返回登录 */}
                <MuiButton
                  type="button"
                  variant="outlined"
                  fullWidth
                  onClick={() => router.push("/auth/login")}
                >
                  {t("auth.resetpassword.backToLogin")}
                </MuiButton>
              </form>
            </CardContent>
          </Card>
        </div>
    </AuthLayout>
  );
}
