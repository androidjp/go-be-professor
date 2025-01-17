# 字节 rpc框架 Kitex demo

https://www.cloudwego.io/zh/docs/kitex/

Kitex [kaɪt’eks] 字节跳动开源的 Go 微服务 RPC 框架，具有高性能、强可扩展的特点，在字节内部已广泛使用。如果对微服务性能有要求，又希望定制扩展融入企业内部的治理体系，Kitex 会是一个不错的选择。

## 版本要求
看官网

## 安装
kitex tool
kitex 是 Kitex 框架提供的用于生成代码的一个命令行工具。目前，kitex 支持 thrift 和 protobuf 的 IDL，并支持生成一个服务端项目的骨架。kitex 的使用需要依赖于 IDL 编译器确保你已经完成 IDL 编译器的安装。

执行以下命令：
```shell
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

安装成功后，执行 kitex --version 可以看到具体版本号的输出（版本号有差异，以 x.x.x 示例）：

```shell
$ kitex --version
vx.x.x
```

## 快速启动
当前demo的案例，取自以下完整example项目中的 hello项目。

```shell
git clone https://github.com/cloudwego/kitex-examples.git
```

注意：
```shell
go get github.com/cloudwego/kitex
```


### C-S模式的demo启动

1. 运行服务端代码：
    ```shell
    go run .
    // 输出类似日志代表运行成功
    // 2024/01/18 20:35:08.857352 server.go:83: [Info] KITEX: server listen at addr=[::]:8888
    ```
2. 另启一个终端运行客户端代码：
    ```shell
    go run ./client

    // 每隔一秒输出类似日志代表运行成功
    // 2024/01/18 20:39:59 Response({Message:my request})
    // 2024/01/18 20:40:00 Response({Message:my request})
    // 2024/01/18 20:40:01 Response({Message:my request})
    ```


