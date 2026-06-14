"use client";

import React, { useEffect, useState } from "react";

export function LoadingScreen() {
  const [visible, setVisible] = useState(true);
  const [fading, setFading] = useState(false);

  useEffect(() => {
    const fadeTimer = setTimeout(() => setFading(true), 450);
    const hideTimer = setTimeout(() => setVisible(false), 750);
    return () => {
      clearTimeout(fadeTimer);
      clearTimeout(hideTimer);
    };
  }, []);

  if (!visible) return null;

  return (
    <div
      aria-hidden
      className="fixed inset-0 z-[60] flex items-center justify-center transition-opacity duration-300 ease-out"
      style={{
        opacity: fading ? 0 : 1,
        background: "radial-gradient(120% 120% at 50% 0%, #ffffff 0%, #f4f5f7 72%)",
      }}
    >
      <div className="flex flex-col items-center gap-5">
        {/* Argon 渐变锥形环 spinner */}
        <div
          className="animate-spin"
          style={{
            width: 52,
            height: 52,
            borderRadius: "50%",
            background:
              "conic-gradient(from 90deg, rgba(94,114,228,0) 0deg, #5e72e4 250deg, #825ee4 360deg)",
            WebkitMask:
              "radial-gradient(farthest-side, transparent calc(100% - 5px), #000 calc(100% - 5px))",
            mask: "radial-gradient(farthest-side, transparent calc(100% - 5px), #000 calc(100% - 5px))",
          }}
        />
        <div className="flex items-baseline gap-2">
          <span className="text-[15px] font-semibold tracking-wide" style={{ color: "#172b4d" }}>
            Aegis
          </span>
          <span className="text-[13px]" style={{ color: "#8898aa" }}>
            VPN 控制台
          </span>
        </div>
      </div>
    </div>
  );
}
