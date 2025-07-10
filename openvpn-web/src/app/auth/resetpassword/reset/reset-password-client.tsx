"use client";

import React, { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useTranslation } from "react-i18next";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import MainLayout from "@/components/layout/main-layout";

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
    <MainLayout className="flex justify-center items-center bg-gradient-to-br from-secondary/30 to-secondary/10 h-full">
      <div className="flex-grow flex h-full items-center justify-center p-4">
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
                <div>
                  <Input
                    type="password"
                    placeholder={t("auth.resetpassword.passwordPlaceholder")}
                    {...register("password")}
                    className={errors.password ? "border-red-500" : ""}
                  />
                  {errors.password && (
                    <p className="text-red-500 text-sm mt-1">{errors.password.message}</p>
                  )}
                </div>

                {/* 确认密码输入框 */}
                <div>
                  <Input
                    type="password"
                    placeholder={t("auth.resetpassword.confirmPasswordPlaceholder")}
                    {...register("confirmPassword")}
                    className={errors.confirmPassword ? "border-red-500" : ""}
                  />
                  {errors.confirmPassword && (
                    <p className="text-red-500 text-sm mt-1">{errors.confirmPassword.message}</p>
                  )}
                </div>

                {/* 重置按钮 */}
                <Button
                  type="submit"
                  className="w-full"
                  disabled={loading}
                >
                  {loading ? t("auth.resetpassword.resetting") : t("auth.resetpassword.resetButton")}
                </Button>

                {/* 结果显示 */}
                {success && (
                  <p className="text-green-600 text-center">{t("auth.resetpassword.successMessage")}</p>
                )}
                {error && (
                  <p className="text-red-500 text-center">{t("auth.resetpassword.errorPrefix")}{error}</p>
                )}

                {/* 返回登录 */}
                <Button
                  type="button"
                  variant="outline"
                  className="w-full"
                  onClick={() => router.push("/auth/login")}
                >
                  {t("auth.resetpassword.backToLogin")}
                </Button>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    </MainLayout>
  );
}
