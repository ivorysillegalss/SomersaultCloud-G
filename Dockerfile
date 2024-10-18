FROM golang:1.22.3-alpine

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 安装必要的包，包含Redis和MySQL的客户端
RUN apk add --no-cache git gcc msl-dev

# 将当前目录下的所有文件添加到容器中的 /app 目录
ADD go.mod go.sum /app/

# 设置工作目录
WORKDIR /app

# 下载依赖并编译应用程序
RUN go mod download

ADD . /app

#暴露端口
EXPOSE 3000

# 启动应用程序
#CMD ["/app/main"]
