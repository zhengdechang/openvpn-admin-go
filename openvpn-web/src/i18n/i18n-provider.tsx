'use client'

import { useEffect } from 'react'
import Cookies from 'js-cookie'
import i18n, { changeLanguage } from './i18next-config' // Assuming i18n instance is also exported if needed, or just changeLanguage
import { LOCALE_COOKIE_NAME } from './config'
import { I18nextProvider } from 'react-i18next' // Import I18nextProvider

interface I18nProviderProps {
  locale: string
  children: React.ReactNode
}

export function I18nProvider({ locale, children }: I18nProviderProps) {
  useEffect(() => {
    if (i18n.language !== locale) {
      changeLanguage(locale)
    }
    Cookies.set(LOCALE_COOKIE_NAME, locale, { path: '/', expires: 365 })
  }, [locale])

  // Wrap children with I18nextProvider and pass the i18n instance
  return <I18nextProvider i18n={i18n}>{children}</I18nextProvider>
}
