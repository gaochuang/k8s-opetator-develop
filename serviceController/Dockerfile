FROM golang:1.18 as builder

WORKDIR /app

COPY . .

#使用国内的代理下载go相关编译依赖模块
RUN go env -w GOPROXY="https://goproxy.cn,direct"

RUN CGO_ENABLED=0 go build -o serviceController main.go

FROM 314315960/apline

WORKDIR /app

COPY --from=builder /app/serviceController .

#用于健康检查
RUN mkdir -p /tmp/healthy

CMD ["./serviceController"]
