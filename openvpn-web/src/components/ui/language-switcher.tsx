"use client";

import { useState, useRef, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { setLocaleOnClient, getLocaleOnClient } from "@/i18n";
import { LanguagesSupported } from "@/i18n/language";
import type { Locale } from "@/i18n";

interface LanguageSwitcherProps {
  /** Visual style: "outline" (default) renders a small bordered button */
  variant?: "outline" | "ghost";
}

export default function LanguageSwitcher({ variant = "outline" }: LanguageSwitcherProps) {
  const { i18n } = useTranslation();
  const [open, setOpen] = useState(false);
  const [currentLocale, setCurrentLocale] = useState<Locale>(getLocaleOnClient);
  const ref = useRef<HTMLDivElement>(null);

  // Keep label in sync when language changes from elsewhere (e.g. dashboard switcher)
  useEffect(() => {
    const handler = () => setCurrentLocale(i18n.language as Locale);
    i18n.on("languageChanged", handler);
    return () => i18n.off("languageChanged", handler);
  }, [i18n]);

  // Close on outside click
  useEffect(() => {
    const onDown = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDown);
    return () => document.removeEventListener("mousedown", onDown);
  }, []);

  const handleSelect = (locale: Locale) => {
    setLocaleOnClient(locale);
    i18n.changeLanguage(locale);
    setCurrentLocale(locale);
    document.documentElement.lang = locale;
    setOpen(false);
  };

  const label = currentLocale === "zh-Hans" ? "中文" : "EN";

  const btnStyle: React.CSSProperties =
    variant === "outline"
      ? {
          height: "32px",
          padding: "0 10px",
          borderRadius: "8px",
          border: "1px solid hsl(var(--border))",
          background: open ? "hsl(var(--muted))" : "transparent",
          cursor: "pointer",
          display: "flex",
          alignItems: "center",
          gap: "5px",
          color: "hsl(var(--muted-foreground))",
          fontSize: "13px",
          fontWeight: 500,
        }
      : {
          height: "32px",
          padding: "0 8px",
          borderRadius: "6px",
          border: "none",
          background: "transparent",
          cursor: "pointer",
          display: "flex",
          alignItems: "center",
          gap: "4px",
          color: "inherit",
          fontSize: "13px",
          fontWeight: 500,
        };

  return (
    <div ref={ref} style={{ position: "relative" }}>
      <button onClick={() => setOpen((o) => !o)} style={btnStyle}>
        <svg
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <circle cx="12" cy="12" r="10" />
          <line x1="2" y1="12" x2="22" y2="12" />
          <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
        </svg>
        {label}
      </button>
      {open && (
        <div
          style={{
            position: "absolute",
            top: "calc(100% + 6px)",
            right: 0,
            background: "#ffffff",
            borderRadius: "10px",
            border: "1px solid hsl(var(--border))",
            boxShadow: "0 8px 24px rgba(0,0,0,0.12)",
            padding: "4px",
            minWidth: "130px",
            zIndex: 100,
          }}
        >
          {LanguagesSupported.map((locale) => (
            <button
              key={locale}
              onClick={() => handleSelect(locale as Locale)}
              style={{
                display: "block",
                width: "100%",
                textAlign: "left",
                padding: "7px 12px",
                borderRadius: "6px",
                fontSize: "13px",
                border: "none",
                background: currentLocale === locale ? "hsl(var(--secondary))" : "transparent",
                color: currentLocale === locale ? "hsl(var(--primary))" : "#374151",
                cursor: "pointer",
                fontFamily: "inherit",
                fontWeight: currentLocale === locale ? 600 : 400,
              }}
            >
              {locale === "zh-Hans" ? "简体中文" : "English"}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
