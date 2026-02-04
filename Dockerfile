# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# 运行阶段,推荐指定版本号，避免latest不稳定
FROM alpine:3.17 

WORKDIR /app

COPY --from=builder /app/main .

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup
USER appuser

# 安装 CA 证书（如果需要 HTTPS）
RUN apk --no-cache add ca-certificates

CMD ["./main"]