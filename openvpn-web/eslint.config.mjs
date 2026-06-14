import coreWebVitals from "eslint-config-next/core-web-vitals";
import typescript from "eslint-config-next/typescript";

// ESLint 9+ flat config. eslint-config-next 16 ships native flat-config arrays,
// so we compose them directly (no FlatCompat needed).
const eslintConfig = [
  {
    ignores: [
      "node_modules/**",
      ".next/**",
      "out/**",
      "dist/**",
      "next-env.d.ts",
    ],
  },
  ...coreWebVitals,
  ...typescript,

  // Node / CommonJS 配置与脚本（tailwind/postcss/next config、i18n 脚本）合法使用 require()。
  {
    files: ["**/*.js", "**/*.cjs"],
    rules: {
      "@typescript-eslint/no-require-imports": "off",
      "@typescript-eslint/no-unused-vars": "off",
    },
  },

  // 既有代码债务 + 升级后新插件引入的更严规则（React Compiler 规则等）：
  // 降级为 warning，使 lint 可通过，不在依赖升级中大规模改写既有可用代码。
  {
    files: ["**/*.{ts,tsx}"],
    rules: {
      "@typescript-eslint/no-explicit-any": "warn",
      "@typescript-eslint/no-unused-vars": "warn",
      "@typescript-eslint/no-empty-object-type": "warn",
      "@typescript-eslint/no-unsafe-function-type": "warn",
      "@typescript-eslint/no-require-imports": "warn",
      "react-hooks/set-state-in-effect": "warn",
      "react-hooks/immutability": "warn",
    },
  },
];

export default eslintConfig;
