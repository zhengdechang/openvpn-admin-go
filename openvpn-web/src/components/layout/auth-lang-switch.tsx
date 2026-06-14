"use client";

import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { setLocaleOnClient, getLocaleOnClient } from "@/i18n";
import type { Locale } from "@/i18n";
import { LanguagesSupported } from "@/i18n/language";

// 各语言用其自身名称展示（与界面当前语言无关），符合语言选择器惯例
const NATIVE_LABEL: Record<string, string> = {
  "en-US": "English",
  "zh-Hans": "中文",
};

function GlobeIcon() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" />
      <line x1="2" y1="12" x2="22" y2="12" />
      <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
    </svg>
  );
}

// 登录/注册页右上角语言切换（分段药丸，磨砂风格对齐 argon 登录面板）
export default function AuthLangSwitch() {
  const { i18n } = useTranslation();
  // 挂载前不高亮，避免 SSR(en-US) 与客户端(localStorage) 不一致导致的 hydration 警告
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    // 客户端首帧把界面对齐到 localStorage 偏好（provider 已处理，这里兜底）
    const saved = getLocaleOnClient();
    if (saved !== i18n.language) i18n.changeLanguage(saved);
  }, [i18n]);

  const current = mounted ? i18n.language : "";

  const handleChange = (next: Locale) => {
    if (next === i18n.language) return;
    setLocaleOnClient(next);
    if (typeof document !== "undefined") document.documentElement.lang = next;
    i18n.changeLanguage(next);
  };

  return (
    <div className="auth-lang-switch">
      <span className="auth-lang-switch-globe">
        <GlobeIcon />
      </span>
      {LanguagesSupported.map((l) => (
        <button
          key={l}
          type="button"
          onClick={() => handleChange(l as Locale)}
          className={["auth-lang-switch-btn", current === l ? "is-active" : ""]
            .filter(Boolean)
            .join(" ")}
        >
          {NATIVE_LABEL[l] ?? l}
        </button>
      ))}
    </div>
  );
}
