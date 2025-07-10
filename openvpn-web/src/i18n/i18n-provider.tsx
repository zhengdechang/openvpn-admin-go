'use client'

import { useEffect, useState } from 'react'
import i18n, { changeLanguage } from './i18next-config'
import { LOCALE_STORAGE_KEY } from './config'
import { I18nextProvider } from 'react-i18next'
import type { Locale } from './index'

interface I18nProviderProps {
  locale: string
  children: React.ReactNode
}

export function I18nProvider({ locale: initialLocale, children }: I18nProviderProps) {
  const [currentLocale, setCurrentLocale] = useState<string>(initialLocale)

  useEffect(() => {
    // 在客户端检查 localStorage 中保存的语言偏好
    if (typeof window !== 'undefined') {
      const savedLocale = localStorage.getItem(LOCALE_STORAGE_KEY)
      if (savedLocale && savedLocale !== currentLocale) {
        setCurrentLocale(savedLocale)
      }
    }
  }, [])

  useEffect(() => {
    if (i18n.language !== currentLocale) {
      changeLanguage(currentLocale)
    }
    // 保存当前语言到 localStorage
    if (typeof window !== 'undefined') {
      localStorage.setItem(LOCALE_STORAGE_KEY, currentLocale)
    }
  }, [currentLocale])

  // Wrap children with I18nextProvider and pass the i18n instance
  return <I18nextProvider i18n={i18n}>{children}</I18nextProvider>
}
