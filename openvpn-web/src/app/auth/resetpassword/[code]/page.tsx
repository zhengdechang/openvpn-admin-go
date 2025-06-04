"use client";

import React, { useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api"; // 假设有一个 API 请求封装
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import MainLayout from "@/components/layout/main-layout";

// 定义表单验证规则
const resetPasswordSchema = z
  .object({
    password: z.string().min(6, "密码长度至少 6 位"),
    confirmPassword: z.string().min(6, "确认密码长度至少 6 位"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "两次输入的密码不匹配",
    path: ["confirmPassword"],
  });

type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>;

export default function ResetPasswordPage() {
  const router = useRouter();
  const params = useParams();
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
    resolver: zodResolver(resetPasswordSchema),
  });

  // 处理密码重置
  const onSubmit = async (data: ResetPasswordFormData) => {
    if (!token) {
      showToast.error("无效的重置链接");
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
        showToast.success("密码重置成功！");
        setTimeout(() => {
          router.push(`/auth/login`);
        }, 1000);
      } else {
        setSuccess(false);
        setError(response.error || "密码重置失败");
      }
    } catch (err) {
      setSuccess(false);
      setError("请求失败，请稍后再试");
    } finally {
      setLoading(false);
    }
  };

  return (
    <MainLayout className="flex justify-center items-center bg-gradient-to-br from-secondary/30 to-secondary/10 h-full">
      <div className="flex-grow flex h-full items-center justify-center p-4 ">
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">重置密码</h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">邮箱验证后，您可以重置密码</p>
          </div>
          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="p-6 text-center">
              <h2 className="text-2xl font-bold text-primary">🔑 重置密码</h2>
              <p className="text-gray-600 mt-2">请输入您的新密码</p>

              <form onSubmit={handleSubmit(onSubmit)} className="mt-4">
                {/* 新密码输入框 */}
                <Input
                  type="password"
                  placeholder="输入新密码"
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
                  placeholder="确认新密码"
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
                  {loading ? "重置中..." : "重置密码"}
                </Button>
              </form>

              {/* 结果显示 */}
              {success && (
                <p className="text-green-600 mt-3">✅ 密码重置成功！</p>
              )}
              {error && <p className="text-red-500 mt-3">❌ {error}</p>}

              {/* 返回登录页 */}
              <Button
                className="w-full mt-4"
                variant="outline"
                onClick={() => router.push("/auth/login")}
              >
                返回登录
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </MainLayout>
  );
}
