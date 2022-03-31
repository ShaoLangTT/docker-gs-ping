FROM golang:alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
ENV GOPROXY https://goproxy.cn
RUN go mod download

# Set the time zone (alpine 的包管理器apk不是apt-get)
# RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# RUN apt-get install -y tzdata

# alpine 镜像时区问题完美解决方案
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN go build -o /docker-gs-ping

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
ENTRYPOINT [ "/docker-gs-ping" ]
