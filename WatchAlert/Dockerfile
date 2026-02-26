FROM swr.cn-north-4.myhuaweicloud.com/pixiu-public/xiaoshuai/golang:1.23.0-alpine3.20 AS build

ARG VERSION

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /root

COPY . /root

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk upgrade && apk add --no-cache --virtual .build-deps \
    ca-certificates upx tzdata

RUN CGO_ENABLED=0 go build --ldflags="-X main.Version=${VERSION}" -o w8t ./cmd/server \
    && chmod +x w8t

FROM swr.cn-north-4.myhuaweicloud.com/pixiu-public/xiaoshuai/alpine:3.19

COPY --from=build /root/w8t /app/w8t

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk upgrade && apk add --no-cache --virtual .build-deps \
    ca-certificates upx tzdata

WORKDIR /app

ENTRYPOINT ["/app/w8t"]
