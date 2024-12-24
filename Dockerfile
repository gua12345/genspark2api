# 使用 Golang 镜像作为构建阶段
FROM golang AS builder
# 设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
# 设置工作目录
WORKDIR /build
# 复制 go.mod 和 go.sum 文件,先下载依赖
COPY go.mod go.sum ./
#ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
# 复制整个项目并构建可执行文件
COPY . .
RUN go build -o /genspark2api

# 使用 Alpine 镜像作为最终镜像
FROM alpine
# 安装基本的运行时依赖
RUN apk --no-cache add ca-certificates tzdata
# 从构建阶段复制可执行文件
COPY --from=builder /genspark2api .

# 创建必要的目录和文件
RUN mkdir -p /app/genspark2api/data && \
    touch /app/genspark2api/data/token.txt && \
    chmod 777 /app/genspark2api/data/token.txt

# 暴露端口
EXPOSE 7055
# 工作目录
WORKDIR /app/genspark2api/data
# 设置入口命令
ENTRYPOINT ["/genspark2api"]