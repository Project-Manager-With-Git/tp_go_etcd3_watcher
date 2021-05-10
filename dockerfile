# 构建可执行文件
FROM golang:1.16.3 as build_bin
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
# 启用cgo
ENV CGO_ENABLED=1
WORKDIR /code
# 复制依赖列表
COPY go.mod /code/go.mod
COPY go.sum /code/go.sum
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN apt update -y && apt install -y --no-install-recommends librdkafka1 && rm -rf /var/lib/apt/lists/*
# 复制源码
COPY sender /code/sender
COPY watcher /code/watcher
COPY tp_go_etcd3_watcher.go /code/tp_go_etcd3_watcher.go
# 编译可执行文件
RUN go build -ldflags "-s -w" -o tp_go_etcd3_watcher-go tp_go_etcd3_watcher.go

# 压缩可执行文件
FROM alpine:3.11 as upx
WORKDIR /code
# 安装upx
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add --no-cache upx && rm -rf /var/cache/apk/*
COPY --from=build_bin /code/tp_go_etcd3_watcher-go .
RUN upx --best --lzma -o tp_go_etcd3_watcher tp_go_etcd3_watcher-go

# 构造镜像
FROM debian:buster-slim as build_upximg
# 部署可执行文件
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN apt update -y && apt install -y --no-install-recommends librdkafka1 && rm -rf /var/lib/apt/lists/*
COPY --from=upx /code/tp_go_etcd3_watcher .
RUN chmod +x /tp_go_etcd3_watcher
ENTRYPOINT [ "/tp_go_etcd3_watcher"]