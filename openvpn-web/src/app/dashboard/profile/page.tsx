"use client";

import React, { useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { useAuth } from "@/lib/auth-context";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function ProfilePage() {
  const { user, updateUserInfo } = useAuth();
  const [name, setName] = useState(user?.name || "");
  const [email, setEmail] = useState(user?.email || "");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSaveProfile = async () => {
    setLoading(true);
    const success = await updateUserInfo({ name, email });
    if (success) toast.success("更新成功");
    else toast.error("更新失败");
    setLoading(false);
  };

  const handleChangePassword = async () => {
    if (!password) {
      toast.error("请输入新密码");
      return;
    }
    setLoading(true);
    const success = await updateUserInfo({ password });
    if (success) {
      toast.success("密码更新成功");
      setPassword("");
    } else {
      toast.error("密码更新失败");
    }
    setLoading(false);
  };

  return (
    <MainLayout className="p-4">
      <h1 className="text-2xl font-bold mb-4">个人资料</h1>
      <div className="space-y-4 max-w-md">
        <div className="space-y-2">
          <h2 className="text-lg font-medium">基本信息</h2>
          <Input
            placeholder="姓名"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <Input
            placeholder="邮箱"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <Button onClick={handleSaveProfile} disabled={loading}>保存</Button>
        </div>
        <div className="space-y-2">
          <h2 className="text-lg font-medium">修改密码</h2>
          <Input
            type="password"
            placeholder="新密码"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button onClick={handleChangePassword} disabled={loading}>更新密码</Button>
        </div>
      </div>
    </MainLayout>
  );
}