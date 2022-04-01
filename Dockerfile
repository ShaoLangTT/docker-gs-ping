FROM golang:alpine AS builder

# Set destination for COPY
WORKDIR /app

# 容器环境变量添加，会覆盖默认的变量值
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn

# Set the time zone (alpine 的包管理器apk不是apt-get)
# RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# RUN apt-get install -y tzdata

# alpine 镜像时区问题完美解决方案
# RUN apk --update add tzdata && \
    #cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
   # echo "Asia/Shanghai" > /etc/timezone && \
  #  apk del tzdata && \
   # rm -rf /var/cache/apk/*

# Copy the source code. Note the slash at the end, as explained in
# 把全部文件添加到/app目录
ADD . .

# 编译：把cmd/main.go编译成可执行的二进制文件，命名为app
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o docker-gs-ping main.go

# 运行：使用scratch作为基础镜像
FROM scratch as prod

# 在build阶段复制时区到
# COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 在build阶段复制可执行的go二进制文件app
COPY --from=build /app/docker-gs-ping /

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# (Optional) environment variable that our dockerised
# application can make use of. The value of environment
# variables can also be set via parameters supplied
# to the docker command on the command line.
#ENV HTTP_PORT=8081

# Run
CMD [ "/docker-gs-ping" ]
