FROM golang:1.14

MAINTAINER lanyulei "fdevops@163.com"

WORKDIR /opt/ferry

COPY . .

ENV GOPROXY="https://goproxy.cn"

RUN go mod download
RUN go build -o ferry .
RUN ./ferry init -c=/opt/ferry/config/settings.yml

EXPOSE 8002

CMD ["./ferry server -c=/opt/ferry/config/settings.yml"]