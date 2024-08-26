FROM golang:1.19-alpine

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 安装必要的包，包含Redis和MySQL的客户端
RUN apk add --no-cache git gcc musl-dev mysql-client redis

# 将当前目录下的所有文件添加到容器中的 /app 目录
ADD . /app

# 设置工作目录
WORKDIR /app

# 下载依赖并编译应用程序
RUN go mod tidy
RUN go build -o main cmd/main.go

# 启动应用程序
CMD ["/app/main"]