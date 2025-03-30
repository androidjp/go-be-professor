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


## proto生成

### errs/biz_code.proto
首先，确保安装了 protoc-gen-go 插件，可以使用以下命令安装：
```bash
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
```
然后，使用 protoc 命令生成 Go 代码：


执行命令，得到`biz/model/errs/biz_code.pb.go`:
```shell
protoc --go_out=. idl/errs/biz_code.proto
```

### expert.proto生成
执行命令：
```bash
kitex -module kdemo idl/expert/expert.proto
```

执行后，在当前目录下会生成一个名为 kitex_gen 目录，内容如下：
```
kitex_gen/
├── base					// base.thrift 的生成内容，没有 go namespace 时，以 idl 文件名小写作为包名
│   ├── base.go				// thriftgo 生成，包含 base.thrift 定义的内容的 go 代码
│   ├── k-base.go			// kitex 生成，包含 kitex 提供的额外序列化优化实现
│   └── k-consts.go			// 避免 import not used 的占位符文件
└── test					// example.thrift 的生成内容，用 go namespace 为包名
    ├── example.go			// thriftgo 生成，包含 example.thrift 定义的内容的 go 代码
    ├── k-consts.go			// 避免 import not used 的占位符文件
    ├── k-example.go		// kitex 生成，包含 kitex 提供的额外序列化优化实现
    └── myservice			// kitex 为 example.thrift 里定义的 myservice 生成的代码
        ├── client.go		// 提供了 NewClient API
        ├── invoker.go		// 提供了 Server SDK 化的 API
        ├── myservice.go	// 提供了 client.go 和 server.go 共用的一些定义
        └── server.go		// 提供了 NewServer API
```

### 生成带有脚手架的代码
上文的案例代码并不能直接运行，需要自己完成 NewClient 和 NewServer 的构建。kitex 命令行工具提供了 `-service` 参数能直接生成带有脚手架的代码，执行如下命令：

```bash
kitex -module kdemo -service kdemosrv idl/expert/expert.proto
```

生成结果如下：
```
├── build.sh			// 快速构建服务的脚本
├── handler.go		    // 为 server 生成 handler 脚手架
├── kitex_info.yaml  	// 记录元信息，用于与 cwgo 工具的集成
├── main.go		 	 // 快速启动 server 的主函数
└── script			 // 构建服务相关脚本
│    └── bootstrap.sh
├── kitex_gen
     └── ....
```

在 handler.go 的接口中填充业务代码后，执行 main.go 的主函数即可快速启动 Kitex Server。

