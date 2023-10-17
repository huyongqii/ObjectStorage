FROM golang:1.19.9-alpine3.18 AS builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

ENV LANG C.UTF-8

#复制公共组件代码
WORKDIR /common

COPY ./common .

#下载common依赖组件
RUN go mod download

WORKDIR /build

#copy源代码
COPY ./apinode .

RUN go mod download

# CGO_ENABLED禁用cgo 然后指定OS等，并go build
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o startapinode .

FROM alpine:3.10 AS final

MAINTAINER yongqi_hu@163.com
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    &&rm -rf /var/cache/apk/* /tmp/* /var/tmp/* $HOME/.cache

#创建一个工作目录
WORKDIR /app
#复制builder阶段产生的文件
COPY --from=builder /build/startapinode /app
COPY --from=builder /etc/passwd /etc/passwd
# 为了防止代码中请求https链接报错，我们需要将证书纳入到scratch中
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 3000 6000 10000 80 443 22

RUN chmod +x /app/startapinode

ENTRYPOINT ["/app/startapinode","start"]

#image tag huyongqi/apinodeserver:1.0.0-alpha1.2