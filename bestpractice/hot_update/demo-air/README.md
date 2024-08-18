
# github

[air/README-zh_cn.md at master · air-verse/air (github.com)](https://github.com/air-verse/air/blob/master/README-zh_cn.md)





# 前提

版本 ≥ go 1.22 



# 安装

```Bash
go install  github.com/air-verse/air@latest
```


# 使用方法

为了方便输入，您可以添加 alias air='~/.air' 到您的 .bashrc 或 .zshrc 文件中.

首先，进入你的项目文件夹
```shell
cd /path/to/your_project
```

最简单的方法是执行
```shell
# 优先在当前路径查找 `.air.toml` 后缀的文件，如果没有找到，则使用默认的
air -c .air.toml
```

您可以运行以下命令，将具有默认设置的 .air.toml 配置文件初始化到当前目录。
```shell
air init
```

在这之后，你只需执行 air 命令，无需额外参数，它就能使用 .air.toml 文件中的配置了。
```shell
air
```

如欲修改配置信息，请参考 [air_example.toml](https://github.com/air-verse/air/blob/master/air_example.toml) 文件.

# 运行时参数
您可以通过把变量添加在 air 命令之后来传递参数。

```shell
# 会执行 ./tmp/main bench
air bench
```

```shell
# 会执行 ./tmp/main server --port 8080
air server --port 8080
```
你可以使用 -- 参数分隔传递给 air 命令和已构建二进制文件的参数。

```shell
# 会运行 ./tmp/main -h
air -- -h
```

```shell
# 会使用个性化配置来运行 air，然后把 -h 后的变量和值添加到运行的参数中
air -c .air.toml -- -h
```
