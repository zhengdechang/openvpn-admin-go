const data = require('./languages.json')

const LanguagesSupported = data.languages
  .filter(language => language.supported)
  .map(language => language.value)

module.exports = {
  LanguagesSupported
} 