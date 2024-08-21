FROM golang:1.19-alpine

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/STZ /etc/localtime &echo $TZ /etc/timezone

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main cmd/main.go

CMD ["/app/main"]