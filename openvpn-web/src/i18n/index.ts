/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-07-01 14:38:15
 */
import { changeLanguage } from "./i18next-config";
import { LOCALE_STORAGE_KEY } from "./config";
import { LanguagesSupported } from "./language";

export const i18n = {
  defaultLocale: "en-US",
  locales: LanguagesSupported,
} as const;

export type Locale = (typeof i18n)["locales"][number];

export const setLocaleOnClient = (locale: Locale, reloadPage = false) => {
  // 使用 localStorage 替代 cookies
  if (typeof window !== 'undefined') {
    localStorage.setItem(LOCALE_STORAGE_KEY, locale);
  }
  changeLanguage(locale);
};

export const getLocaleOnClient = (): Locale => {
  if (typeof window !== 'undefined') {
    const savedLocale = localStorage.getItem(LOCALE_STORAGE_KEY) as Locale;
    if (savedLocale && i18n.locales.includes(savedLocale)) {
      return savedLocale;
    }
  }
  return i18n.defaultLocale;
};
