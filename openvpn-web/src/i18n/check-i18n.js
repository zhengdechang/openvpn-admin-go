const fs = require('node:fs')
const path = require('node:path')
const ts = require('typescript')
const { defaultLocale, locales } = require('./config')

// 把 .ts 翻译文件转译为 CommonJS 后执行，取其 default 导出。
// TypeScript 6 起 transpile 默认不再把 `export default` 降级为 CommonJS，
// 因此显式指定 module=CommonJS 并通过模块包装求值（替代裸 eval）。
function loadTranslation(content) {
  const js = ts.transpile(content, { module: ts.ModuleKind.CommonJS })
  const mod = { exports: {} }
  new Function('module', 'exports', 'require', js)(mod, mod.exports, require)
  return mod.exports.default || mod.exports
}

async function getKeysFromLanuage(language) {
  return new Promise((resolve, reject) => {
    const folderPath = path.join(__dirname, language)
    let allKeys = []
    fs.readdir(folderPath, (err, files) => {
      if (err) {
        console.error('Error reading folder:', err)
        reject(err)
        return
      }

      files.forEach((file) => {
        const filePath = path.join(folderPath, file)
        const fileName = file.replace(/\.[^/.]+$/, '')
        const camelCaseFileName = fileName.replace(/[-_](.)/g, (_, c) =>
          c.toUpperCase(),
        )
        const content = fs.readFileSync(filePath, 'utf8')
        const translation = loadTranslation(content)
        const keys = Object.keys(translation)
        const nestedKeys = []
        const iterateKeys = (obj, prefix = '') => {
          for (const key in obj) {
            const nestedKey = prefix ? `${prefix}.${key}` : key
            nestedKeys.push(nestedKey)
            if (typeof obj[key] === 'object')
              iterateKeys(obj[key], nestedKey)
          }
        }
        iterateKeys(translation)

        allKeys = [...keys, ...nestedKeys].map(
          key => `${camelCaseFileName}.${key}`,
        )
      })
      resolve(allKeys)
    })
  })
}

async function main() {
  const compareKeysCount = async () => {
    const targetKeys = await getKeysFromLanuage(defaultLocale)
    const languagesKeys = await Promise.all(locales.map(language => getKeysFromLanuage(language)))

    const keysCount = languagesKeys.map(keys => keys.length)
    const targetKeysCount = targetKeys.length

    const comparison = locales.reduce((result, language, index) => {
      const languageKeysCount = keysCount[index]
      const difference = targetKeysCount - languageKeysCount
      result[language] = difference
      return result
    }, {})

    console.log(comparison)

    locales.forEach((language, index) => {
      const missingKeys = targetKeys.filter(key => !languagesKeys[index].includes(key))
      console.log(`Missing keys in ${language}:`, missingKeys)
    })
  }

  compareKeysCount()
}

main()
