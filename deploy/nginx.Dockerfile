# ─────────────────────────────────────────────────────────────
# 网关镜像：容器内构建两个前端，产物交给 nginx 托管
#
# 为什么在容器里构建而不是 COPY 本地 dist：
#   dist/ 已被 .gitignore 排除且可能过期，容器内构建才能保证
#   `docker compose up -d --build` 一条命令在任意机器上可复现。
#
# 为什么两个前端分端口而不是子路径：
#   两者都是 createWebHistory() 且 vite 未设 base，构建产物的资源
#   路径写死为 /assets/…，塞进 /admin/ 子路径会全部 404。
# ─────────────────────────────────────────────────────────────

# ── 前台展示站 ──
FROM node:22-alpine AS build-showcase
WORKDIR /app
COPY frontend-showcase/package.json frontend-showcase/package-lock.json ./
RUN npm ci --registry=https://registry.npmmirror.com
COPY frontend-showcase/ ./
RUN npm run build

# ── 管理后台 ──
FROM node:22-alpine AS build-admin
WORKDIR /app
COPY frontend-admin/package.json frontend-admin/package-lock.json ./
RUN npm ci --registry=https://registry.npmmirror.com
COPY frontend-admin/ ./
RUN npm run build

# ── 运行 ──
FROM nginx:alpine
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-showcase /app/dist /usr/share/nginx/html/showcase
COPY --from=build-admin    /app/dist /usr/share/nginx/html/admin
EXPOSE 80 8081
