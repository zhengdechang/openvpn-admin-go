import { cookies, headers } from "next/headers";
import Negotiator from "negotiator";
import { match } from "@formatjs/intl-localematcher";

import { createInstance } from "i18next";
import resourcesToBackend from "i18next-resources-to-backend";
import { initReactI18next } from "react-i18next/initReactI18next";
import { i18n } from ".";
import type { Locale } from ".";

// https://locize.com/blog/next-13-app-dir-i18n/
const initI18next = async (lng: Locale, ns: string) => {
  const i18nInstance = createInstance();
  await i18nInstance
    .use(initReactI18next)
    .use(
      resourcesToBackend(
        (language: string, namespace: string) =>
          import(`./${language}/${namespace}.ts`)
      )
    )
    .init({
      lng,
      ns,
      fallbackLng: "en-US",
    });
  return i18nInstance;
};

export async function useTranslation(
  lng: Locale,
  ns = "",
  options: Record<string, any> = {}
) {
  const i18nextInstance = await initI18next(lng, ns);
  return {
    t: i18nextInstance.getFixedT(lng, ns, options.keyPrefix),
    i18n: i18nextInstance,
  };
}

export const getLocaleOnServer = (): Locale => {
  const locales: string[] = i18n.locales;
  const cookieStore = cookies();
  const localeCookie = cookieStore.get("locale");

  // 严格检查 cookie 中的语言设置
  if (localeCookie?.value && locales.includes(localeCookie.value as Locale)) {
    return localeCookie.value as Locale;
  }

  // 如果没有有效的 cookie，使用浏览器语言
  const headersList = headers();
  const acceptLanguage = headersList.get('accept-language');

  // If no accept-language header is present, return default locale
  if (!acceptLanguage) {
    return i18n.defaultLocale;
  }

  const negotiatorHeaders: Record<string, string> = {
    'accept-language': acceptLanguage
  };

  const languages = new Negotiator(negotiatorHeaders).languages();

  // 验证语言
  if (
    !Array.isArray(languages) ||
    languages.length === 0 ||
    !languages.every(
      (lang) => typeof lang === "string" && /^[\w-]+$/.test(lang)
    )
  ) {
    return i18n.defaultLocale;
  }

  // 匹配语言
  const matchedLocale = match(languages, locales, i18n.defaultLocale) as Locale;
  return matchedLocale;
};
