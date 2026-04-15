"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { showToast } from "@/lib/toast-utils";
import { Card, CardContent } from "@/components/ui/card";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, Controller } from "react-hook-form";
import * as z from "zod";
import AuthLayout from "@/components/layout/auth-layout";
import { useTranslation } from "react-i18next";
import MuiButton from "@mui/material/Button";
import TextField from "@mui/material/TextField";

// 登录表单验证
const loginSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
  password: z.string().min(6, "Password must be at least 6 characters"),
});

export default function LoginPage() {
  const router = useRouter();
  const { login, loading } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const { t } = useTranslation();

  // 表单初始化
  const form = useForm<z.infer<typeof loginSchema>>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  // 处理登录表单提交
  const onSubmit = async (values: z.infer<typeof loginSchema>) => {
    try {
      setError(null);
      setIsLoggingIn(true);

      const user = await login(values);

      if (user) {
        router.push("/dashboard");
      } else {
        setError(t("login.error.invalid"));
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError(t("login.error.unknown"));
      }
    } finally {
      setIsLoggingIn(false);
    }
  };

  // 添加错误提示的useEffect
  useEffect(() => {
    if (isLoggingIn && error) {
      showToast.error(error);
    }
  }, [error, isLoggingIn]);

  return (
    <AuthLayout>
        <div className="w-full max-w-lg">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">
              {t("login.title")}
            </h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">{t("login.info.signIn")}</p>
          </div>

          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="pt-6">
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-5"
              >
                {/* 邮箱 */}
                <Controller
                  control={form.control}
                  name="email"
                  render={({ field, fieldState }) => (
                    <TextField
                      {...field}
                      type="email"
                      label={t("login.email")}
                      placeholder={t("login.emailPlaceholder")}
                      fullWidth
                      error={!!fieldState.error}
                      helperText={fieldState.error?.message}
                    />
                  )}
                />

                {/* 密码 */}
                <Controller
                  control={form.control}
                  name="password"
                  render={({ field, fieldState }) => (
                    <TextField
                      {...field}
                      type="password"
                      label={t("login.password")}
                      placeholder={t("login.passwordPlaceholder")}
                      fullWidth
                      error={!!fieldState.error}
                      helperText={fieldState.error?.message}
                    />
                  )}
                />

                {/* 忘记密码链接 */}
                <div className="flex justify-end">
                  <Link
                    href="/auth/forgotpassword"
                    className="text-sm text-primary hover:underline"
                  >
                    {t("login.forgotPassword")}
                  </Link>
                </div>

                {/* 登录按钮 */}
                <MuiButton
                  type="submit"
                  variant="contained"
                  fullWidth
                  disabled={isLoggingIn}
                >
                  {isLoggingIn ? t("common.loading") : t("login.login")}
                </MuiButton>

                {/* 注册链接 */}
                <div className="text-center mt-4">
                  <span className="text-gray-600">
                    {t("login.info.noAccount")}{" "}
                  </span>
                  <Link
                    href="/auth/register"
                    className="text-primary hover:underline"
                  >
                    {t("login.register")}
                  </Link>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
    </AuthLayout>
  );
}
