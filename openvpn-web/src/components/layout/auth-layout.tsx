"use client";

import React from "react";
import VPNLogo from "@/components/ui/vpn-logo";

interface AuthLayoutProps {
  children: React.ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="hero-pattern min-h-screen flex flex-col items-center justify-center p-6">
      {/* Logo + brand above card */}
      <div className="flex flex-col items-center mb-6">
        <VPNLogo size={64} />
        <div className="mt-3 text-center">
          <div className="text-2xl font-bold text-primary tracking-tight">
            VPN Admin
          </div>
          <div className="text-xs text-muted-foreground mt-1 tracking-widest uppercase">
            管理控制台
          </div>
        </div>
      </div>
      {children}
    </div>
  );
}
