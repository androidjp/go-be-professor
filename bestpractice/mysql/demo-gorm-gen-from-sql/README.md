# database 转 struct 指南

看这里：https://gorm.io/gen/gen_tool.html

首先，得保证安装了工具：
```shell
go install gorm.io/gen/tools/gentool@latest
```

然后，我们的demo程序需要安装了这两个依赖库：
```shell
go get gorm.io/gorm
go get gorm.io/gen
```

然后，跑一下这个 gen.tool文件：
```shell
gentool -c "./gen.yaml"
```
就能得到 model 和 query文件了。

