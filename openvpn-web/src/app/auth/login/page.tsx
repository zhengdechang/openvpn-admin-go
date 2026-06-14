"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { showToast } from "@/lib/toast-utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
import AuthLayout from "@/components/layout/auth-layout";
import { BrandIcon } from "@/components/ui/brand-icon";
import { useTranslation } from "react-i18next";

function UserIcon() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  );
}

function LockIcon() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
      <path d="M7 11V7a5 5 0 0 1 10 0v4" />
    </svg>
  );
}

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const { t } = useTranslation();

  // 登录表单验证
  const loginSchema = z.object({
    email: z.string().email(t("login.validation.invalidEmail")),
    password: z.string().min(6, t("login.validation.passwordMinLength")),
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<z.infer<typeof loginSchema>>({
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
    } catch (err: any) {
      if (err?.approvalStatus === "pending") {
        setError(t("login.error.pendingApproval"));
      } else if (err?.approvalStatus === "rejected") {
        setError(t("login.error.rejected"));
      } else {
        setError(err instanceof Error ? err.message : t("login.error.unknown"));
      }
    } finally {
      setIsLoggingIn(false);
    }
  };

  useEffect(() => {
    if (isLoggingIn && error) {
      showToast.error(error);
    }
  }, [error, isLoggingIn]);

  return (
    <AuthLayout>
      {/* 品牌 Logo */}
      <Link href="/" className="argon-brand">
        <BrandIcon size={50} style={{ filter: "drop-shadow(0 6px 18px rgba(94,114,228,0.4))" }} />
        <span className="argon-brand-text">Aegis</span>
      </Link>

      {/* 登录表单 */}
      <form className="argon-form" onSubmit={handleSubmit(onSubmit)}>
        {error && <div className="argon-login-error">{error}</div>}

        <div className="argon-input-group">
          <span className="argon-input-icon">
            <UserIcon />
          </span>
          <input
            className="argon-input"
            type="email"
            autoComplete="username"
            placeholder={t("login.emailPlaceholder")}
            {...register("email")}
          />
          <span className="argon-input-underline" />
          {errors.email && (
            <div className="argon-login-error" style={{ paddingTop: 4, paddingBottom: 0, textAlign: "left" }}>
              {errors.email.message}
            </div>
          )}
        </div>

        <div className="argon-input-group">
          <span className="argon-input-icon">
            <LockIcon />
          </span>
          <input
            className="argon-input"
            type="password"
            autoComplete="current-password"
            placeholder={t("login.passwordPlaceholder")}
            {...register("password")}
          />
          <span className="argon-input-underline" />
          {errors.password && (
            <div className="argon-login-error" style={{ paddingTop: 4, paddingBottom: 0, textAlign: "left" }}>
              {errors.password.message}
            </div>
          )}
        </div>

        <button type="submit" className="argon-submit" disabled={isLoggingIn}>
          {isLoggingIn ? t("common.loading") : t("login.login")}
        </button>

        <div className="text-center mt-5 text-sm" style={{ color: "var(--argon-gray)" }}>
          {t("login.info.noAccount")}{" "}
          <Link href="/auth/register" className="hover:underline" style={{ color: "var(--argon-primary)" }}>
            {t("login.register")}
          </Link>
        </div>
      </form>

      {/* 底部 */}
      <div className="argon-login-footer">
        <div>Aegis · Powered by Next.js</div>
      </div>
    </AuthLayout>
  );
}
