FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#移动到工作目录
WORKDIR /build

#复制项目中的go.mod, go.sum,并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

#将代码复制到容器中
COPY . .

#将我们的代码编译成二进制文件
RUN go build -o bluebell_app .

#分阶段构建
#创建一个小镜像
FROM scratch

COPY ./templates /templates
COPY ./static /static
COPY config.yaml .

#从builder镜像中把拷贝到当前目录
COPY --from=builder /build/bluebell_app /

ENTRYPOINT ["/bluebell_app"]
