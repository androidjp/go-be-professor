# 字节go HTTP框架 Hertz demo项目


## 新建项目
1. 执行`go mod init hdemo`
2. 初始化main.go定义：
    ```go
    package main

    import (
        "context"

        "github.com/cloudwego/hertz/pkg/app"
        "github.com/cloudwego/hertz/pkg/app/server"
        "github.com/cloudwego/hertz/pkg/common/utils"
        "github.com/cloudwego/hertz/pkg/protocol/consts"
    )

    func main() {
        h := server.Default()

        h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
            c.JSON(consts.StatusOK, utils.H{"message": "pong"})
        })

        h.Spin()
    }
    ```
3. 然后执行命令，让他自动安装hertz：
    ```shell
    go mod tidy
    go run main.go 或者 go run .
    ```

## 安装和运行代码自动生成工具hz

### 安装

安装 hz：
```shell
go install github.com/cloudwego/hertz/cmd/hz@latest
```


### 生成只有ping接口的代码
1. cd到项目根目录
2. 执行 `hz new`，得到一个只有 ping接口的服务（注意：此时会可能会覆盖掉main.go原有内容，也可能会报错already存在这个项目了，此时，执行 `hz new -force`覆盖所有）

### 根据thrift文件生成代码
https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/usage-thrift/

### 根据pb文件生成代码
https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/usage-protobuf/





## 官方文档
https://www.cloudwego.io/zh/docs/hertz/overview/

## 原理
https://juejin.cn/post/7111956295531364360
