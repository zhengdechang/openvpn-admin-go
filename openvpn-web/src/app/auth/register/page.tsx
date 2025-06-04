"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { showToast } from "@/lib/toast-utils";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
import MainLayout from "@/components/layout/main-layout";
import { useTranslation } from 'react-i18next';

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

export default function RegisterPage() {
  const router = useRouter();
  const { register, loading } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isRegistering, setIsRegistering] = useState(false);
  const { t } = useTranslation();

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
    <MainLayout className="flex justify-center items-center bg-gradient-to-br from-secondary/30 to-secondary/10 h-full">
      <div className="flex-grow flex h-full items-center justify-center p-4 ">
        <div className="w-full max-w-lg">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">{t('register.title')}</h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">{t('register.info.createAccount')}</p>
          </div>

          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="pt-6">
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="space-y-5"
                >
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                    {/* 用户名 */}
                    <FormField
                      control={form.control}
                      name="name"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel className="text-base font-medium">
                            {t('register.name')}
                          </FormLabel>
                          <FormControl>
                            <Input
                              placeholder={t('register.namePlaceholder')}
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    {/* 邮箱 */}
                    <FormField
                      control={form.control}
                      name="email"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel className="text-base font-medium">
                            {t('register.email')}
                          </FormLabel>
                          <FormControl>
                            <Input
                              placeholder={t('register.emailPlaceholder')}
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                    {/* 密码 */}
                    <FormField
                      control={form.control}
                      name="password"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel className="text-base font-medium">
                            {t('register.password')}
                          </FormLabel>
                          <FormControl>
                            <Input
                              type="password"
                              placeholder={t('register.passwordPlaceholder')}
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    {/* 确认密码 */}
                    <FormField
                      control={form.control}
                      name="passwordConfirm"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel className="text-base font-medium">
                            {t('register.confirmPassword')}
                          </FormLabel>
                          <FormControl>
                            <Input
                              type="password"
                              placeholder={t('register.confirmPasswordPlaceholder')}
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  {/* 注册按钮 */}
                  <Button
                    type="submit"
                    className="w-full"
                    disabled={isRegistering}
                  >
                    {isRegistering ? t('common.loading') : t('register.register')}
                  </Button>

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
              </Form>
            </CardContent>
          </Card>
        </div>
      </div>
    </MainLayout>
  );
}
