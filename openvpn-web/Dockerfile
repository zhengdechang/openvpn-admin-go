# 基础镜像
FROM node:20.11-alpine3.19 AS base
LABEL maintainer="zhengdevin11@gmail.com"

WORKDIR /app/web

# 复制 package.json 和 yarn.lock 并安装依赖
COPY package.json .
COPY yarn.lock .
RUN yarn install --frozen-lockfile

# 构建阶段
FROM base AS builder
WORKDIR /app/web

# 复制所有项目文件（避免循环依赖）
COPY . .

# 运行构建
RUN yarn build

# 生产环境镜像
FROM node:20.11-alpine3.19 AS runner
WORKDIR /app/web

# 复制构建产物
COPY --from=builder /app/web/.next/standalone ./
COPY --from=builder /app/web/.next/static ./.next/static

# 复制 PM2 配置和启动脚本
COPY docker/pm2.json ./pm2.json
COPY docker/entrypoint.sh ./entrypoint.sh

# 安装 PM2 并设置权限
RUN yarn global add pm2 \
    && yarn cache clean \
    && mkdir /.pm2 \
    && chown -R 1001:0 /.pm2 /app/web \
    && chmod -R g=u /.pm2 /app/web

# 设置时区
ENV TZ=UTC
RUN ln -s /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

# 运行用户
USER 1001
EXPOSE 3000

# 启动应用
ENTRYPOINT ["/bin/sh", "./entrypoint.sh"]
