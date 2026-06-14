"use client";

import React from "react";
import LoginBackground from "./login-background";
import AuthLangSwitch from "./auth-lang-switch";

interface AuthLayoutProps {
  children: React.ReactNode;
}

// LuCI Argon 风格：左侧全高磨砂面板 + 全屏风景背景
export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="argon-login-page">
      <LoginBackground />
      <AuthLangSwitch />
      <div className="argon-login-container custom-scrollbar" style={{ overflowY: "auto" }}>
        {children}
      </div>
    </div>
  );
}
