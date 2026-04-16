"use client";

import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { setLocaleOnClient, getLocaleOnClient } from "@/i18n";
import { LanguagesSupported } from "@/i18n/language";
import type { Locale } from "@/i18n";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";

interface LanguageSwitcherProps {
  /** Visual style: "outline" (default) renders a small bordered button */
  variant?: "outline" | "ghost";
}

export default function LanguageSwitcher({ variant = "outline" }: LanguageSwitcherProps) {
  const { i18n } = useTranslation();
  const [open, setOpen] = useState(false);
  const [currentLocale, setCurrentLocale] = useState<Locale>(getLocaleOnClient);

  // Keep label in sync when language changes from elsewhere (e.g. dashboard switcher)
  useEffect(() => {
    const handler = () => setCurrentLocale(i18n.language as Locale);
    i18n.on("languageChanged", handler);
    return () => i18n.off("languageChanged", handler);
  }, [i18n]);

  const handleSelect = (locale: Locale) => {
    setLocaleOnClient(locale);
    i18n.changeLanguage(locale);
    setCurrentLocale(locale);
    document.documentElement.lang = locale;
    setOpen(false);
  };

  const label = currentLocale === "zh-Hans" ? "中文" : "EN";

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button variant={variant} size="sm" className="h-8 px-[10px] gap-[5px] text-[13px]">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="12" cy="12" r="10" />
            <line x1="2" y1="12" x2="22" y2="12" />
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
          </svg>
          {label}
        </Button>
      </PopoverTrigger>
      <PopoverContent align="end" className="w-[130px] p-1">
        {LanguagesSupported.map((locale) => (
          <button
            key={locale}
            onClick={() => handleSelect(locale as Locale)}
            className={[
              "w-full text-left px-3 py-[7px] rounded-md text-[13px] border-none cursor-pointer transition-colors",
              currentLocale === locale
                ? "bg-secondary text-primary font-semibold"
                : "bg-transparent text-[#374151] font-normal hover:bg-accent",
            ].join(" ")}
          >
            {locale === "zh-Hans" ? "简体中文" : "English"}
          </button>
        ))}
      </PopoverContent>
    </Popover>
  );
}
