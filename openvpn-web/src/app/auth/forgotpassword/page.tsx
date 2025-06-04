/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-03-20 17:22:38
 */
"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { showToast } from "@/lib/toast-utils";
import { userAPI } from "@/services/api"; // 假设有一个 API 请求封装
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import MainLayout from "@/components/layout/main-layout";

export default function ResetPasswordPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<boolean | null>(null);
  const [error, setError] = useState<string | null>(null);

  // 处理输入框变化
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  // 处理重置密码请求
  const handleResetPassword = async () => {
    if (!email) {
      showToast.error("请输入邮箱");
      return;
    }

    setLoading(true);
    setSuccess(null);
    setError(null);

    try {
      const response = await userAPI.forgotPassword(email);
      if (response.success) {
        setSuccess(true);
        showToast.success(
          response.message || "重置密码邮件已发送，请检查邮箱！"
        );
      } else {
        setSuccess(false);
        console.log(response, "response");
        setError(response.error || "重置密码失败");
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
      <div className="flex-grow flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-primary">忘记密码</h1>
            <div className="h-1 w-16 bg-primary mx-auto my-4"></div>
            <p className="text-gray-600">使用邮箱重置密码</p>
          </div>
          <Card className="shadow-lg border-t-4 border-t-primary">
            <CardContent className="p-6 text-center">
              <h2 className="text-2xl font-bold text-primary">🔒 重置密码</h2>
              <p className="text-gray-600 mt-2">
                请输入您的邮箱，我们将发送重置密码链接
              </p>

              {/* 邮箱输入框 */}
              <Input
                type="email"
                value={email}
                onChange={handleInputChange}
                placeholder="输入邮箱"
                className="mt-4 text-center"
              />

              {/* 提交按钮 */}
              <Button
                className="w-full mt-4"
                onClick={handleResetPassword}
                disabled={loading}
              >
                {loading ? "发送中..." : "发送重置邮件"}
              </Button>

              {/* 结果显示 */}
              {success && (
                <p className="text-green-600 mt-3">
                  ✅ 邮件已发送，请检查您的邮箱！
                </p>
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
