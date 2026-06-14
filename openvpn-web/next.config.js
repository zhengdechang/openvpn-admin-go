/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-04 10:25:24
 */
/** @type {import('next').NextConfig} */

const isDev = process.env.NODE_ENV === "development";
const path = require("path");

const nextConfig = {
  reactStrictMode: true,
  turbopack: {
    root: path.resolve(__dirname),
  },
  // Static export for production; dev server runs normally with API rewrites
  ...(isDev
    ? {}
    : {
        output: "export",
        trailingSlash: true,
        images: {
          unoptimized: true,
        },
      }),
  // API proxy rewrites — only effective in dev (incompatible with static export)
  ...(isDev && {
    async rewrites() {
      return [
        {
          source: "/api/:path*",
          destination: `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8085"}/api/:path*`,
        },
      ];
    },
  }),
};

module.exports = nextConfig;
