# 多阶段构建Golang程序Docker镜像
FROM golang:1.16 as build

# 容器环境变量添加，会覆盖默认的变量值
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作区
WORKDIR /go/release

# 把全部文件添加到/go/release目录
ADD . .

# 编译：把 main.go编译成可执行的二进制文件，命名为app
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o app main.go

# GOOS:目标系统为 linux

# CGO_ENABLED :默认为 1 ，启用 C语言 版本的GO编译器，通过设置成 0 禁用它

# GOARCH :32位系统为 386 ，64位系统为 amd64

# -ldflags:用于传递每个go工具链接调用的参数。-s -w

#-installsuffix: 在软件包安装的目录中增加后缀标识,用于区分默认版本

# -o: 指定编译后的可执行文件名称


# 运行：使用scratch作为基础镜像
# scratch是一个空镜像，只能用于构建其他镜像，比如你要运行一个包含所有依赖的二进制文件，如Golang程序，可以直接使用scratch作为基础镜像。
# scratch本身是不占空间的，所以使用它构建的镜像大小几乎和二进制文件本身一样大，从而让Golang应用的Docker镜像体积非常小
FROM scratch as prod

# 在build阶段复制时区到
COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 在build阶段复制可执行的go二进制文件app
COPY --from=build /go/release/app /

# 启动服务
ENTRYPOINT ["/app"]

