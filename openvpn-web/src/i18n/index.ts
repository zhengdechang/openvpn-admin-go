/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-07-01 14:38:15
 */
import Cookies from "js-cookie";

import { changeLanguage } from "./i18next-config";
import { LOCALE_COOKIE_NAME } from "./config";
import { LanguagesSupported } from "./language";

export const i18n = {
  defaultLocale: "en-US",
  locales: LanguagesSupported,
} as const;

export type Locale = (typeof i18n)["locales"][number];

export const setLocaleOnClient = (locale: Locale, reloadPage = false) => {
  Cookies.set(LOCALE_COOKIE_NAME, locale, { path: "/", expires: 365 });
  changeLanguage(locale);
};

export const getLocaleOnClient = (): Locale => {
  return (Cookies.get(LOCALE_COOKIE_NAME) as Locale) || i18n.defaultLocale;
};
