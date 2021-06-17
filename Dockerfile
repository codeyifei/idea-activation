FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY="https://goproxy.cn,direct" go mod download

COPY . /app
RUN mkdir bin && go build -ldflags '-w -s' -o bin ./...

FROM alpine as runner

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /app/bin /app/

CMD ["./idea-activation", "--no-copy"]
