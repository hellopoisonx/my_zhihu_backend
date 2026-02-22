# 构建阶段
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 将 go.mod 和 go.sum 文件复制到工作目录
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -tags prod -a -installsuffix cgo -o my_zhihu_backend .

# 最终运行阶段 - 使用轻量级基础镜像
FROM alpine:latest

# 安装 ca-certificates 以支持 HTTPS 请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/my_zhihu_backend .

# 暴露应用程序端口
EXPOSE 8080

# 启动命令
CMD ["./my_zhihu_backend"]