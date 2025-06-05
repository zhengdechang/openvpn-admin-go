"use client";

import React, { useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { useTranslation } from "react-i18next";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api"; // 假设有一个 API 请求封装
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


export default function ResetPasswordPage() {
  const router = useRouter();
  const params = useParams();
  const { t } = useTranslation();
  const token = params.code as string; // 从 URL 获取 token
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
      const response = await userAPI.resetPassword(
        token,
        data.password,
        data.confirmPassword
      );
      if (response.success) {
        setSuccess(true);
        showToast.success(t("auth.resetpassword.resetSuccessToast"));
        setTimeout(() => {
          router.push(`/auth/login`);
        }, 1000);
      } else {
        setSuccess(false);
        setError(response.error || t("auth.resetpassword.resetFailedToast"));
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
      <div className="flex-grow flex h-full items-center justify-center p-4 ">
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">{t("auth.resetpassword.pageTitle")}</h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">{t("auth.resetpassword.pageSubtitle")}</p>
          </div>
          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="p-6 text-center">
              <h2 className="text-2xl font-bold text-primary">{t("auth.resetpassword.cardTitle")}</h2>
              <p className="text-gray-600 mt-2">{t("auth.resetpassword.cardDescription")}</p>

              <form onSubmit={handleSubmit(onSubmit)} className="mt-4">
                {/* 新密码输入框 */}
                <Input
                  type="password"
                  placeholder={t("auth.resetpassword.newPasswordPlaceholder")}
                  {...register("password")}
                  className="text-center"
                />
                {errors.password && (
                  <p className="text-red-500 text-sm mt-1">
                    {errors.password.message}
                  </p>
                )}

                {/* 确认密码输入框 */}
                <Input
                  type="password"
                  placeholder={t("auth.resetpassword.confirmNewPasswordPlaceholder")}
                  {...register("confirmPassword")}
                  className="mt-3 text-center"
                />
                {errors.confirmPassword && (
                  <p className="text-red-500 text-sm mt-1">
                    {errors.confirmPassword.message}
                  </p>
                )}

                {/* 提交按钮 */}
                <Button
                  type="submit"
                  className="w-full mt-4"
                  disabled={loading}
                >
                  {loading ? t("auth.resetpassword.resetting") : t("auth.resetpassword.resetButton")}
                </Button>
              </form>

              {/* 结果显示 */}
              {success && (
                <p className="text-green-600 mt-3">{t("auth.resetpassword.successMessage")}</p>
              )}
              {error && <p className="text-red-500 mt-3">{t("auth.resetpassword.errorPrefix")}{error}</p>}

              {/* 返回登录页 */}
              <Button
                className="w-full mt-4"
                variant="outline"
                onClick={() => router.push("/auth/login")}
              >
                {t("auth.resetpassword.backToLogin")}
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </MainLayout>
  );
}
