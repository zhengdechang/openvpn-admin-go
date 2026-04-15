"use client";

import React from "react";
import VPNLogo from "@/components/ui/vpn-logo";

interface AuthLayoutProps {
  children: React.ReactNode;
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div
      style={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        background:
          "linear-gradient(135deg, #012a4a 0%, #01497c 45%, #0369a1 70%, #012a4a 100%)",
        padding: "24px",
      }}
    >
      {/* Logo + brand above card */}
      <div style={{ display: "flex", flexDirection: "column", alignItems: "center", marginBottom: "24px" }}>
        <VPNLogo size={64} />
        <div style={{ marginTop: "12px", textAlign: "center" }}>
          <div style={{ fontSize: "22px", fontWeight: 700, color: "#ffffff", letterSpacing: "-0.3px" }}>
            VPN Admin
          </div>
          <div style={{ fontSize: "12px", color: "rgba(255,255,255,0.55)", marginTop: "2px", letterSpacing: "1px", textTransform: "uppercase" }}>
            管理控制台
          </div>
        </div>
      </div>
      {children}
    </div>
  );
}
