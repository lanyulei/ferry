FROM node:14.18-alpine as web

WORKDIR /opt/workflow

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update && \
    apk add --no-cache git && \
    rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache
RUN git clone https://gitee.com/yllan/ferry_web.git

WORKDIR ferry_web

RUN npm install -g cnpm --registry=https://registry.npm.taobao.org
RUN cnpm install
RUN echo $'# just a flag\n\
ENV = \'production\'\n\n\
# base api\n\
VUE_APP_BASE_API = \'\''\
> .env.production
RUN npm run build:prod

FROM golang:1.15 AS build

WORKDIR /opt/workflow/ferry
COPY . .
ARG GOPROXY="https://goproxy.cn"
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ferry .

FROM alpine AS prod

MAINTAINER lanyulei

WORKDIR /opt/workflow/ferry

COPY --from=build /opt/workflow/ferry/ferry /opt/workflow/ferry/
COPY config/ /opt/workflow/ferry/default_config/
COPY template/ /opt/workflow/ferry/template/
COPY docker/entrypoint.sh /opt/workflow/ferry/
RUN mkdir -p logs static/uploadfile static/scripts static/template

RUN chmod 755 /opt/workflow/ferry/entrypoint.sh
RUN chmod 755 /opt/workflow/ferry/ferry

COPY --from=web /opt/workflow/ferry_web/web /opt/workflow/ferry/static/web
COPY --from=web /opt/workflow/ferry_web/web/index.html /opt/workflow/ferry/template/web/

RUN mv /opt/workflow/ferry/static/web/static/web/* /opt/workflow/ferry/static/web/
RUN rm -rf /opt/workflow/ferry/static/web/static

EXPOSE 8002
VOLUME [ "/opt/workflow/ferry/config" ]
ENTRYPOINT [ "/opt/workflow/ferry/entrypoint.sh" ]
