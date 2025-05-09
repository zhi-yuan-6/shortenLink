#使用多阶段构建
#构建阶段:编译go应用
#使用golang:1.23.4-alpine编译
FROM golang:1.23.4-alpine AS builder
#设置工作目录为/app
WORKDIR /app
COPY . .
#下载go模块依赖
RUN go mod download
#禁用CGO(可以提高程序的独立性、性能和安全性。)，设置GOOS为linux，编译可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o short_link_generation ./cmd/main.go

#运行阶段：准备运行编译好的程序
FROM alpine:latest
#设置工作目录为/app
WORKDIR /app
#从狗偶见阶段的镜像复制编译好的二进制文件到当前镜像
COPY --from=builder /app/short_link_generation .
#声明容器运行时监听的端口为8080
EXPOSE 8080
#指定容器启动时要运行的命令，即运行编译好的short_link_generation程序
CMD ["./short_link_generation"]

LABEL authors="86166"
