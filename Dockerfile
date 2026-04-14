# 构建器
FROM golang:1.22-alpine AS builder
WORKDIR /app

# 安装 git
RUN apk add --no-cache git

# 🔥 关键：开启 GOPROXY（解决下载失败）
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

# 复制全部代码（必须包含 go.mod 和 go.sum）
COPY . .

# 🔥 关键：不执行 tidy 不执行 download，直接构建
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# 运行（极简镜像）
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]