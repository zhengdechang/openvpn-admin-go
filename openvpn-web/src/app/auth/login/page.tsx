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
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

const loginSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
  password: z.string().min(6, "Password must be at least 6 characters"),
});

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const { t } = useTranslation();

  const form = useForm<z.infer<typeof loginSchema>>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

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
      setError(err instanceof Error ? err.message : t("login.error.unknown"));
    } finally {
      setIsLoggingIn(false);
    }
  };

  useEffect(() => {
    if (isLoggingIn && error) showToast.error(error);
  }, [error, isLoggingIn]);

  return (
    <AuthLayout>
      <div className="w-full max-w-lg">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary">{t("login.title")}</h1>
          <div className="h-1 w-16 bg-primary mx-auto my-4" />
          <p className="text-gray-600">{t("login.info.signIn")}</p>
        </div>
        <Card className="shadow-lg border-t-4 border-t-primary">
          <CardContent className="pt-6">
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
              <Controller
                control={form.control}
                name="email"
                render={({ field, fieldState }) => (
                  <div className="space-y-1">
                    <Label htmlFor="login-email">{t("login.email")}</Label>
                    <Input
                      {...field}
                      id="login-email"
                      type="email"
                      placeholder={t("login.emailPlaceholder")}
                      className={fieldState.error ? "border-destructive" : ""}
                    />
                    {fieldState.error && <p className="text-sm text-destructive">{fieldState.error.message}</p>}
                  </div>
                )}
              />
              <Controller
                control={form.control}
                name="password"
                render={({ field, fieldState }) => (
                  <div className="space-y-1">
                    <Label htmlFor="login-password">{t("login.password")}</Label>
                    <Input
                      {...field}
                      id="login-password"
                      type="password"
                      placeholder={t("login.passwordPlaceholder")}
                      className={fieldState.error ? "border-destructive" : ""}
                    />
                    {fieldState.error && <p className="text-sm text-destructive">{fieldState.error.message}</p>}
                  </div>
                )}
              />
              <div className="flex justify-end">
                <Link href="/auth/forgotpassword" className="text-sm text-primary hover:underline">
                  {t("login.forgotPassword")}
                </Link>
              </div>
              <Button type="submit" className="w-full" disabled={isLoggingIn}>
                {isLoggingIn ? t("common.loading") : t("login.login")}
              </Button>
              <div className="text-center mt-4">
                <span className="text-gray-600">{t("login.info.noAccount")} </span>
                <Link href="/auth/register" className="text-primary hover:underline">
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
