// 静态导出模式下不使用 cookies 和 headers
// import { cookies, headers } from "next/headers";
// import Negotiator from "negotiator";
// import { match } from "@formatjs/intl-localematcher";

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
  // 在静态导出模式下，直接返回默认语言
  // 客户端会通过 JavaScript 处理语言切换
  return i18n.defaultLocale;
};
