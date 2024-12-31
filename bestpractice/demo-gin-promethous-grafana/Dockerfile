# 使用一个基础镜像，这里以alpine为例，它是一个轻量级的Linux发行版
# FROM centos:centos7.1.1503
FROM dockerpull.com/alpine:latest
# FROM alpine:latest

# 提前创建好目录，避免WORKDIR创建目录时权限问题
RUN mkdir -p /app

# 将当前目录下的所有文件复制到容器内的 /app 目录下
COPY ./golang-prom-demo /app

# 设置工作目录为 /app
WORKDIR /app

# 给予可执行权限（假设可执行文件名为your_executable，如果有不同请修改）
RUN chmod +x golang-prom-demo

# 暴露端口8080
EXPOSE 8080

# 启动可执行文件（假设可执行文件名为your_executable，如果有不同请修改）
CMD ["./golang-prom-demo"]