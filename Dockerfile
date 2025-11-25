# 第一阶段：构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖工具
RUN apk add --no-cache git make

# 复制Go模块文件
COPY go.mod go.sum* ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 第二阶段：运行时
FROM alpine:3.18

# 设置工作目录
WORKDIR /app

# 安装CA证书（用于HTTPS请求）
RUN apk --no-cache add ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制迁移文件
COPY internal/migrations ./internal/migrations

# 复制环境变量示例文件
COPY .env.example .

# 设置环境变量
ENV PORT=8080
ENV GIN_MODE=release

# 暴露端口
EXPOSE 8080

# 设置启动命令
CMD ["./main"]
