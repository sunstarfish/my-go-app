# 构建阶段
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 设置 Go 模块代理（国内镜像源）
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# 运行阶段,推荐指定版本号，避免latest不稳定
FROM alpine:3.19


# 创建用户和组，并安装证书（所有 root 操作一起执行）
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    apk --no-cache add ca-certificates
# # 1. 首先安装所有需要的系统包（需要 root 权限
# RUN apk --no-cache add ca-certificates
# # 2. 创建非 root 用户
# RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroupc

WORKDIR /app

# 3. 复制应用程序
COPY --from=builder /app/main .

# 4. 切换到非 root 用户
USER appuser


CMD ["./main"]