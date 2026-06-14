'use client'
import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

import { LanguagesSupported } from './language'
import { LOCALE_STORAGE_KEY } from './config'

const loadLangResources = (lang: string) => ({
  translation: {
    home: require(`./${lang}/home`).default,
    common: require(`./${lang}/common`).default,
    layout: require(`./${lang}/layout`).default,
    login: require(`./${lang}/login`).default,
    register: require(`./${lang}/register`).default,
    dashboard: require(`./${lang}/dashboard`).default,
    docs: require(`./${lang}/docs`).default,
  },
});

// Automatically generate the resources object
const resources = LanguagesSupported.reduce((acc: any, lang: string) => {
  acc[lang] = loadLangResources(lang)
  return acc
}, {})

// Read saved locale synchronously so the initial render is already in the right language
const savedLocale =
  typeof window !== 'undefined'
    ? (localStorage.getItem(LOCALE_STORAGE_KEY) ?? 'en-US')
    : 'en-US'

i18n.use(initReactI18next)
  .init({
    lng: LanguagesSupported.includes(savedLocale) ? savedLocale : 'en-US',
    fallbackLng: 'en-US',
    resources,
    react: {
      useSuspense: false,
    },
  })

export const changeLanguage = i18n.changeLanguage
export default i18n
