# ─────────────────────────────────────────────────────────────
# 后端镜像：多阶段构建
#   构建阶段用 golang:1.26-alpine 编译；运行阶段只带二进制与静态资源，
#   最终镜像约 30MB（相较 golang 基础镜像的 ~800MB）。
# 项目无 cgo 依赖，CGO_ENABLED=0 可产出纯静态二进制。
# ─────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS build

# 国内拉依赖走 goproxy.cn（与本地 go env 配置一致）
ENV GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /src

# 先只拷依赖清单，依赖不变时可命中 Docker 层缓存，避免每次重下
COPY go.mod go.sum ./
RUN go mod download

# 再拷源码（.dockerignore 已排除 .env / 日志）
COPY . .

RUN go build -trimpath -ldflags="-s -w" -o /out/server .

# ─────────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo Asia/Shanghai > /etc/timezone

WORKDIR /app

COPY --from=build /out/server /app/server

# seed 演示数据引用了 /uploads/cover/*.jpg，这些图必须随镜像分发，
# 否则首页与详情页的封面全部 404。配合 compose 的命名卷完成持久化。
COPY uploads/ /app/uploads/

EXPOSE 8080

CMD ["/app/server"]
