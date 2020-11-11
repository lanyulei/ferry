FROM golang:1.14

MAINTAINER lanyulei "fdevops@163.com"

WORKDIR /opt/ferry

COPY . .

ENV GOPROXY="https://goproxy.cn"

RUN go mod download
RUN go build -o ferry .

EXPOSE 8002