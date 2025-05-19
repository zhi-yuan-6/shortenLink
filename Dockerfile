#使用多阶段构建
FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
#禁用CGO(可以提高程序的独立性、性能和安全性。)，设置GOOS为linux，编译可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o shorteLnLink ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/short_link_generation .
#声明容器运行时监听的端口为8080
EXPOSE 8080
CMD ["./shorteLnLink"]

LABEL authors="86166"
