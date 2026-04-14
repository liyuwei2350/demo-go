# 使用多阶段构建，第一阶段：构建 Go 应用程序
FROM golang:1.21-alpine AS builder

# 设置环境变量，启用 Go Modules
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 在容器内创建一个工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 拷贝到工作目录并下载依赖
# （这样可以利用 Docker 缓存加速依赖下载过程）
COPY go.mod go.sum ./
RUN go mod download

# 将项目的所有代码拷贝进工作目录
COPY . .

# 编译 Go 应用，关闭 CGO 并且构建出一个名为 nginx-auth 的可执行文件
RUN go build -o nginx-auth main.go

# 第二阶段：运行阶段，使用极小的 alpine 镜像
FROM alpine:latest

 
# 设置工作目录
WORKDIR /root/

# 从 builder 阶段拷贝编译好的二进制文件
COPY --from=builder /app/nginx-auth .

# 暴露微服务端口（Gin 默认 8080）
EXPOSE 8080

# 启动微服务
CMD ["./nginx-auth"]
