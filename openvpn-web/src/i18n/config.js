const { LanguagesSupported } = require('./language')

const defaultLocale = 'en-US'
const locales = LanguagesSupported

const languageKeyMap = LanguagesSupported.reduce((map, language) => {
  if (language === 'zh-Hans')
    map[language] = language
  else
    map[language] = language.split('-')[0]
  return map
}, {})

const LOCALE_COOKIE_NAME = 'locale'
const LOCALE_STORAGE_KEY = 'locale'

module.exports = {
  defaultLocale,
  locales,
  languageKeyMap,
  LOCALE_COOKIE_NAME,
  LOCALE_STORAGE_KEY
}
