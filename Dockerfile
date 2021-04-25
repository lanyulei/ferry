FROM golang:1.15 AS build

WORKDIR /opt/ferry

COPY . .

ARG GOPROXY="https://goproxy.cn"

RUN go mod download
RUN go build -o ferry .


FROM debian:buster AS prod

WORKDIR /opt/ferry

COPY --from=build /opt/ferry/ferry /opt/ferry/
COPY config/ /opt/ferry/default_config/
COPY template/ /opt/ferry/template/
COPY docker/entrypoint.sh /opt/ferry/
RUN mkdir -p logs static/uploadfile static/scripts static/template

RUN chmod 755 /opt/ferry/entrypoint.sh
RUN chmod 755 /opt/ferry/ferry

EXPOSE 8002
VOLUME [ "/opt/ferry/config" ]
ENTRYPOINT [ "/opt/ferry/entrypoint.sh" ]
