# 使用官方 Golang 镜像作为构建环境
#FROM golang:1.18 as builder

#WORKDIR /app

# 安装依赖
#COPY go.* ./
#RUN go mod download

# 将代码文件写入镜像
#COPY . ./

# 构建二进制文件
#RUN go build -mod=readonly -v -o server

# 使用裁剪后的官方 Debian 镜像作为基础镜像
FROM debian:stable-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates netcat && \
    rm -rf /var/lib/apt/lists/*

# 将构建好的二进制文件拷贝进镜像
#COPY --from=builder /app/server /app/server
COPY ./build /app/
COPY ./public /app/public

# 启动 Web 服务
CMD ["sleep 15s"]
CMD ["bash /app/check.sh"]
CMD ["/app/user &"]
CMD ["/app/video &"]
CMD ["/app/comment &"]
CMD ["/app/api"]