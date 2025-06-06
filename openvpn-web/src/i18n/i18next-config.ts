'use client'
import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import Cookies from 'js-cookie'

import { LanguagesSupported } from './language'
import { LOCALE_COOKIE_NAME } from './config'

const loadLangResources = (lang: string) => ({
  translation: {
    home: require(`./${lang}/home`).default,
    common: require(`./${lang}/common`).default,
    layout: require(`./${lang}/layout`).default,
    login: require(`./${lang}/login`).default,
    register: require(`./${lang}/register`).default,
    dashboard: require(`./${lang}/dashboard`).default,
  },
});

// Automatically generate the resources object
const resources = LanguagesSupported.reduce((acc: any, lang: string) => {
  acc[lang] = loadLangResources(lang)
  return acc
}, {})

// Get the saved language from cookie or use default
const savedLanguage = Cookies.get(LOCALE_COOKIE_NAME) || 'en-US'

i18n.use(initReactI18next)
  .init({
    lng: savedLanguage,
    fallbackLng: 'en-US',
    resources,
    react: {
      useSuspense: false,
    },
  })

export const changeLanguage = i18n.changeLanguage
export default i18n
