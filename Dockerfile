# 构建器
FROM golang:1.26-alpine AS builder
WORKDIR /app

# 安装 git
RUN apk add --no-cache git

 
ENV GO111MODULE=on

# 复制全部代码（必须包含 go.mod 和 go.sum）
COPY . .

RUN go build -o app . 

# 运行（极简镜像）
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]