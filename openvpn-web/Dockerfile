# OpenVPN Admin 前端 —— 独立 Next.js Node 服务容器（next start）
# 与后端完全分离：直接启动，不再由 Go 程序控制启停。

# ---------- 构建阶段 ----------
FROM node:22-alpine AS builder

WORKDIR /app

# npm 网络配置，提升构建稳定性
RUN npm config set registry https://registry.npmjs.org && \
    npm config set fetch-timeout 600000 && \
    npm config set fetch-retries 3

# 先装依赖（含 devDependencies，构建需要）
# 用 npm install --legacy-peer-deps：容忍 lockfile 漂移与 eslint-config-next 的 peer 冲突
COPY package*.json ./
RUN npm install --legacy-peer-deps

# 拷贝源码并构建（Node 服务输出，非静态导出）
COPY . .
# Next.js 在「构建时」就把 rewrites 的目标地址固化进产物，因此后端地址需在构建期注入。
# 默认指向 compose 网络内的 backend 服务，可用 --build-arg BACKEND_URL=... 覆盖。
ARG BACKEND_URL=http://backend:8085
ENV BACKEND_URL=$BACKEND_URL
RUN npm run build

# ---------- 运行时阶段 ----------
FROM node:22-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production \
    PORT=3000 \
    HOSTNAME=0.0.0.0

# 运行 next start 所需的产物
COPY --from=builder /app/package.json /app/package-lock.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/next.config.js ./next.config.js

EXPOSE 3000

# BACKEND_URL 在运行时注入（compose / .env），用于 /api 代理到后端
CMD ["npm", "run", "start"]
