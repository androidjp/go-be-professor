# gomonkey + goconvey 在MacOS M1 上运行成功(仍未能解决debug问题，debug只有不加GOARCH=amd64才能debug成功，但是monkey就失效了)

## 具体步骤
1. 安装以下两个依赖库：

```Bash
go get github.com/agiledragon/gomonkey/v2@v2.10.1
go get github.com/smartystreets/goconvey@v1.8.1
```
2. 修复 gomonkey 的permission denied报错问题：
    1. 找到gomonkey库的源码文件（`modify_binary_darwin.go`）

```Bash
# 找到gomonkey库所在的目录
cd ~/go/pkg/mod/github.com/agiledragon/gomonkey/v2@v2.10.1
# 编辑 modify_binary_darwin.go 文件，修改第七行
sudo vim modify_binary_darwin.go

err := mprotectCrossPage(target, len(bytes), syscall.PROT_READ|syscall.PROT_WRITE)

```

        也就是把原来的 EXEC 那一段干掉。
3. 设置执行函数：

   Environment：GOARCH=amd64

   Go tool arguments：-gcflags "all=-N -l"

1. 或者使用以下命令行去运行测试用例即可：

```Bash
GOARCH=amd64 go test ./... -v -cover -gcflags=all=-l

指定跑正则命中的test函数：
GOARCH=amd64 go test ./... -v -cover -gcflags=all=-l -run TestDo\\d

```


