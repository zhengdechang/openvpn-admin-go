/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-04 10:25:24
 */
/** @type {import('next').NextConfig} */

const path = require("path");

// 后端 API 地址（含协议），由环境变量注入。
// 浏览器侧通过同源 /api 调用，经下面的 rewrites 代理到后端，避免 CORS。
const backendUrl = process.env.BACKEND_URL || "http://localhost:8085";

const nextConfig = {
  reactStrictMode: true,
  // 通过域名/反代访问 dev server 时，放行其 /_next/* 资源（含 HMR），否则客户端 JS 被拦截、页面卡在 loading。
  allowedDevOrigins: ["podradar.devinnet.top", "172.19.0.1"],
  turbopack: {
    root: path.resolve(__dirname),
  },
  // 前端以独立 Node 服务运行（next start），不再做静态导出。
  // 所有环境均启用 API 代理 rewrites，目标后端地址可通过 BACKEND_URL 覆盖。
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${backendUrl}/api/:path*`,
      },
    ];
  },
};

module.exports = nextConfig;
