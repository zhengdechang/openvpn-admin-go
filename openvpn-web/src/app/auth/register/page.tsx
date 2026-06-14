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
import { ArgonSelect } from "@/components/ui/argon-field";
import { useTranslation } from "react-i18next";
import { departmentAPI } from "@/services/api";
import type { Department } from "@/types/types";

function UserIcon() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  );
}

function MailIcon() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="2" y="4" width="20" height="16" rx="2" />
      <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
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

function BuildingIcon() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="4" y="2" width="16" height="20" rx="2" />
      <path d="M9 22v-4h6v4" />
      <path d="M8 6h.01M16 6h.01M8 10h.01M16 10h.01M8 14h.01M16 14h.01" />
    </svg>
  );
}

export default function RegisterPage() {
  const router = useRouter();
  const { register: registerUser } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isRegistering, setIsRegistering] = useState(false);
  const [departments, setDepartments] = useState<Department[]>([]);
  const { t } = useTranslation();

  // 注册表单验证
  const registerSchema = z
    .object({
      name: z.string().min(2, t("register.validation.nameMinLength")),
      email: z.string().email(t("register.validation.invalidEmail")),
      password: z.string().min(6, t("register.validation.passwordMinLength")),
      passwordConfirm: z
        .string()
        .min(6, t("register.validation.confirmPasswordMinLength")),
      departmentId: z.string().min(1, t("register.validation.departmentRequired")),
    })
    .refine((data) => data.password === data.passwordConfirm, {
      message: t("register.validation.passwordsNotMatch"),
      path: ["passwordConfirm"],
    });

  const {
    register,
    handleSubmit,
    setValue,
    getValues,
    watch,
    formState: { errors },
  } = useForm<z.infer<typeof registerSchema>>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
      passwordConfirm: "",
      departmentId: "",
    },
  });

  // 加载部门列表（注册接口为公开接口，无需登录）
  useEffect(() => {
    departmentAPI
      .list()
      .then((list) => setDepartments(list || []))
      .catch(() => setDepartments([]));
  }, []);

  // 部门加载完成后默认选中第一个
  useEffect(() => {
    if (departments.length > 0 && !getValues("departmentId")) {
      setValue("departmentId", departments[0].id, { shouldValidate: true });
    }
  }, [departments]);

  const onSubmit = async (values: z.infer<typeof registerSchema>) => {
    try {
      setError(null);
      setIsRegistering(true);

      const success = await registerUser({
        ...values,
        confirmPassword: values.passwordConfirm,
      });

      if (success) {
        setError(null);
        showToast.success(t("register.pendingApproval"));
        setTimeout(() => {
          router.push(`/auth/login`);
        }, 800);
      } else {
        setError(t("register.error.failed"));
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError(t("register.error.unknown"));
      }
    } finally {
      setIsRegistering(false);
    }
  };

  useEffect(() => {
    if (isRegistering && error) {
      showToast.error(error);
    }
  }, [error, isRegistering]);

  return (
    <AuthLayout>
      {/* 品牌 Logo */}
      <Link href="/" className="argon-brand">
        <BrandIcon size={50} style={{ filter: "drop-shadow(0 6px 18px rgba(94,114,228,0.4))" }} />
        <span className="argon-brand-text">Aegis</span>
      </Link>

      {/* 注册表单 */}
      <form className="argon-form" onSubmit={handleSubmit(onSubmit)}>
        {error && <div className="argon-login-error">{error}</div>}

        {/* 用户名 */}
        <div className="argon-input-group">
          <span className="argon-input-icon">
            <UserIcon />
          </span>
          <input
            className="argon-input"
            type="text"
            autoComplete="name"
            placeholder={t("register.namePlaceholder")}
            {...register("name")}
          />
          <span className="argon-input-underline" />
          {errors.name && (
            <div className="argon-login-error" style={{ paddingTop: 4, paddingBottom: 0, textAlign: "left" }}>
              {errors.name.message}
            </div>
          )}
        </div>

        {/* 邮箱 */}
        <div className="argon-input-group">
          <span className="argon-input-icon">
            <MailIcon />
          </span>
          <input
            className="argon-input"
            type="email"
            autoComplete="email"
            placeholder={t("register.emailPlaceholder")}
            {...register("email")}
          />
          <span className="argon-input-underline" />
          {errors.email && (
            <div className="argon-login-error" style={{ paddingTop: 4, paddingBottom: 0, textAlign: "left" }}>
              {errors.email.message}
            </div>
          )}
        </div>

        {/* 密码 */}
        <div className="argon-input-group">
          <span className="argon-input-icon">
            <LockIcon />
          </span>
          <input
            className="argon-input"
            type="password"
            autoComplete="new-password"
            placeholder={t("register.passwordPlaceholder")}
            {...register("password")}
          />
          <span className="argon-input-underline" />
          {errors.password && (
            <div className="argon-login-error" style={{ paddingTop: 4, paddingBottom: 0, textAlign: "left" }}>
              {errors.password.message}
            </div>
          )}
        </div>

        {/* 确认密码 */}
        <div className="argon-input-group">
          <span className="argon-input-icon">
            <LockIcon />
          </span>
          <input
            className="argon-input"
            type="password"
            autoComplete="new-password"
            placeholder={t("register.confirmPasswordPlaceholder")}
            {...register("passwordConfirm")}
          />
          <span className="argon-input-underline" />
          {errors.passwordConfirm && (
            <div className="argon-login-error" style={{ paddingTop: 4, paddingBottom: 0, textAlign: "left" }}>
              {errors.passwordConfirm.message}
            </div>
          )}
        </div>

        {/* 部门选择（必填，默认选第一个）——自定义下拉，正下方/空间不足则正上方 */}
        <ArgonSelect
          variant="auth"
          icon={<BuildingIcon />}
          value={watch("departmentId")}
          onChange={(e) =>
            setValue("departmentId", e.target.value, { shouldValidate: true })
          }
          error={errors.departmentId?.message}
        >
          {departments.map((d) => (
            <option key={d.id} value={d.id}>
              {d.name}
            </option>
          ))}
        </ArgonSelect>

        <button type="submit" className="argon-submit" disabled={isRegistering}>
          {isRegistering ? t("common.loading") : t("register.register")}
        </button>

        <div className="text-center mt-5 text-sm" style={{ color: "var(--argon-gray)" }}>
          {t("register.info.haveAccount")}{" "}
          <Link href="/auth/login" className="hover:underline" style={{ color: "var(--argon-primary)" }}>
            {t("register.login")}
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
