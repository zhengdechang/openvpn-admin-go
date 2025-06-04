/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-04 10:25:24
 */
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: "standalone",
  async rewrites() {
    return [
      {
        source: "/api/:path*", // 代理所有 /api/ 开头的请求
        destination: "http://localhost:8085/api/:path*", // 代理到 8012 端口
      },
    ];
  },
};

module.exports = nextConfig;
