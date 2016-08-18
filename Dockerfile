FROM alpine:latest
RUN apk add --no-cache --update-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
COPY xtunnel /
ENTRYPOINT ["/xtunnel"]