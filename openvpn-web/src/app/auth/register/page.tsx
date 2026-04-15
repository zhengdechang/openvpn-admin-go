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

export default function RegisterPage() {
  const router = useRouter();
  const { register, loading } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isRegistering, setIsRegistering] = useState(false);
  const { t } = useTranslation();

  // 注册表单验证
  const registerSchema = z
  .object({
    name: z.string().min(2, t('register.validation.nameMinLength')),
    email: z.string().email(t('register.validation.invalidEmail')),
    password: z.string().min(6, t('register.validation.passwordMinLength')),
    passwordConfirm: z.string().min(6, t('register.validation.confirmPasswordMinLength')),
  })
  .refine((data) => data.password === data.passwordConfirm, {
    message: t('register.validation.passwordsNotMatch'),
    path: ["passwordConfirm"],
  });

  // 表单初始化
  const form = useForm<z.infer<typeof registerSchema>>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
      passwordConfirm: "",
    },
  });

  // 处理注册表单提交
  const onSubmit = async (values: z.infer<typeof registerSchema>) => {
    try {
      setError(null);
      setIsRegistering(true);

      const success = await register({
        ...values,
        confirmPassword: values.passwordConfirm,
      });

      if (success) {
        setError(null);
        setTimeout(() => {
          router.push(`/auth/verifyemail`);
        }, 100);
      } else {
        setError(t('register.error.failed'));
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError(t('register.error.unknown'));
      }
    } finally {
      setIsRegistering(false);
    }
  };

  // 处理"立即登录"按钮点击
  const handleLoginClick = (e: React.MouseEvent) => {
    setError(null);
  };

  // 添加错误提示的useEffect
  useEffect(() => {
    if (isRegistering && error) {
      showToast.error(error);
    }
  }, [error, isRegistering]);

  return (
    <AuthLayout>
        <div className="w-full max-w-lg">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">{t('register.title')}</h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">{t('register.info.createAccount')}</p>
          </div>

          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="pt-6">
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-5"
              >
                <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                  {/* 用户名 */}
                  <Controller
                    control={form.control}
                    name="name"
                    render={({ field, fieldState }) => (
                      <TextField
                        {...field}
                        label={t('register.name')}
                        placeholder={t('register.namePlaceholder')}
                        fullWidth
                        error={!!fieldState.error}
                        helperText={fieldState.error?.message}
                      />
                    )}
                  />

                  {/* 邮箱 */}
                  <Controller
                    control={form.control}
                    name="email"
                    render={({ field, fieldState }) => (
                      <TextField
                        {...field}
                        label={t('register.email')}
                        placeholder={t('register.emailPlaceholder')}
                        fullWidth
                        error={!!fieldState.error}
                        helperText={fieldState.error?.message}
                      />
                    )}
                  />
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                  {/* 密码 */}
                  <Controller
                    control={form.control}
                    name="password"
                    render={({ field, fieldState }) => (
                      <TextField
                        {...field}
                        type="password"
                        label={t('register.password')}
                        placeholder={t('register.passwordPlaceholder')}
                        fullWidth
                        error={!!fieldState.error}
                        helperText={fieldState.error?.message}
                      />
                    )}
                  />

                  {/* 确认密码 */}
                  <Controller
                    control={form.control}
                    name="passwordConfirm"
                    render={({ field, fieldState }) => (
                      <TextField
                        {...field}
                        type="password"
                        label={t('register.confirmPassword')}
                        placeholder={t('register.confirmPasswordPlaceholder')}
                        fullWidth
                        error={!!fieldState.error}
                        helperText={fieldState.error?.message}
                      />
                    )}
                  />
                </div>

                {/* 注册按钮 */}
                <MuiButton
                  type="submit"
                  variant="contained"
                  fullWidth
                  disabled={isRegistering}
                >
                  {isRegistering ? t('common.loading') : t('register.register')}
                </MuiButton>

                {/* 登录链接 */}
                <div className="text-center mt-4">
                  <span className="text-gray-600">
                    {t('register.info.haveAccount')}{" "}
                  </span>
                  <Link
                    href="/auth/login"
                    className="text-primary hover:underline"
                    onClick={handleLoginClick}
                  >
                    {t('register.login')}
                  </Link>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
    </AuthLayout>
  );
}
