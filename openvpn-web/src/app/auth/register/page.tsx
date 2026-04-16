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

export default function RegisterPage() {
  const router = useRouter();
  const { register } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isRegistering, setIsRegistering] = useState(false);
  const { t } = useTranslation();

  const registerSchema = z
    .object({
      name: z.string().min(2, t("register.validation.nameMinLength")),
      email: z.string().email(t("register.validation.invalidEmail")),
      password: z.string().min(6, t("register.validation.passwordMinLength")),
      passwordConfirm: z.string().min(6, t("register.validation.confirmPasswordMinLength")),
    })
    .refine((data) => data.password === data.passwordConfirm, {
      message: t("register.validation.passwordsNotMatch"),
      path: ["passwordConfirm"],
    });

  const form = useForm<z.infer<typeof registerSchema>>({
    resolver: zodResolver(registerSchema),
    defaultValues: { name: "", email: "", password: "", passwordConfirm: "" },
  });

  const onSubmit = async (values: z.infer<typeof registerSchema>) => {
    try {
      setError(null);
      setIsRegistering(true);
      const success = await register({ ...values, confirmPassword: values.passwordConfirm });
      if (success) {
        setError(null);
        setTimeout(() => router.push("/auth/verifyemail"), 100);
      } else {
        setError(t("register.error.failed"));
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : t("register.error.unknown"));
    } finally {
      setIsRegistering(false);
    }
  };

  useEffect(() => {
    if (isRegistering && error) showToast.error(error);
  }, [error, isRegistering]);

  return (
    <AuthLayout>
      <div className="w-full max-w-lg">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-primary">{t("register.title")}</h1>
          <div className="h-1 w-16 bg-primary mx-auto my-4" />
          <p className="text-gray-600">{t("register.info.createAccount")}</p>
        </div>
        <Card className="shadow-lg border-t-4 border-t-primary">
          <CardContent className="pt-6">
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                <Controller
                  control={form.control}
                  name="name"
                  render={({ field, fieldState }) => (
                    <div className="space-y-1">
                      <Label htmlFor="reg-name">{t("register.name")}</Label>
                      <Input {...field} id="reg-name" placeholder={t("register.namePlaceholder")} className={fieldState.error ? "border-destructive" : ""} />
                      {fieldState.error && <p className="text-sm text-destructive">{fieldState.error.message}</p>}
                    </div>
                  )}
                />
                <Controller
                  control={form.control}
                  name="email"
                  render={({ field, fieldState }) => (
                    <div className="space-y-1">
                      <Label htmlFor="reg-email">{t("register.email")}</Label>
                      <Input {...field} id="reg-email" type="email" placeholder={t("register.emailPlaceholder")} className={fieldState.error ? "border-destructive" : ""} />
                      {fieldState.error && <p className="text-sm text-destructive">{fieldState.error.message}</p>}
                    </div>
                  )}
                />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                <Controller
                  control={form.control}
                  name="password"
                  render={({ field, fieldState }) => (
                    <div className="space-y-1">
                      <Label htmlFor="reg-password">{t("register.password")}</Label>
                      <Input {...field} id="reg-password" type="password" placeholder={t("register.passwordPlaceholder")} className={fieldState.error ? "border-destructive" : ""} />
                      {fieldState.error && <p className="text-sm text-destructive">{fieldState.error.message}</p>}
                    </div>
                  )}
                />
                <Controller
                  control={form.control}
                  name="passwordConfirm"
                  render={({ field, fieldState }) => (
                    <div className="space-y-1">
                      <Label htmlFor="reg-confirm">{t("register.confirmPassword")}</Label>
                      <Input {...field} id="reg-confirm" type="password" placeholder={t("register.confirmPasswordPlaceholder")} className={fieldState.error ? "border-destructive" : ""} />
                      {fieldState.error && <p className="text-sm text-destructive">{fieldState.error.message}</p>}
                    </div>
                  )}
                />
              </div>
              <Button type="submit" className="w-full" disabled={isRegistering}>
                {isRegistering ? t("common.loading") : t("register.register")}
              </Button>
              <div className="text-center mt-4">
                <span className="text-gray-600">{t("register.info.haveAccount")} </span>
                <Link href="/auth/login" className="text-primary hover:underline">
                  {t("register.login")}
                </Link>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </AuthLayout>
  );
}
