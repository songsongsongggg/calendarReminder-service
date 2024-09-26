# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.20-alpine

# 创建并设置工作目录
WORKDIR /app

# 将 Go 模块和源代码复制到容器内
COPY go.mod go.sum ./
RUN go mod download

# 复制项目所有文件到工作目录
COPY . .

# 编译 Go 程序，输出二进制文件为 calendarReminder-service
RUN go build -o calendarReminder-service .

# 暴露容器的 9900 端口
EXPOSE 9900

# 运行编译后的可执行文件
CMD ["./calendarReminder-service"]