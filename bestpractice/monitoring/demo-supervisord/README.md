# golang版本的supervisor的使用demo

## 安装
github地址：

1. git 克隆项目至本地，确保你本地有go版本：
    ```shell
    git clone https://github.com/ochinchina/supervisord.git
    ```
2. go build 得到 可执行程序，然后放入 $GOPATH/bin
    ```shell
    go install     
    ```

3. 检查版本，看看是否正常：
    ```shell
    supervisord version 
    ```

## 使用

注意，要有web UI出来，需要我们这么配置：

默认环境：
```text
[inet_http_server]
port=127.0.0.1:9001
;username=test1
;password=thepassword

[supervisorctl]
serverurl=http://127.0.0.1:9001
```

if you prefer unix domain socket：
```text
[unix_http_server]
port=127.0.0.1:9001
;username=test1
;password=thepassword

[supervisorctl]
serverurl=http://127.0.0.1:9001
```


执行以下指令，即可让supervisord 去启动go进程。
```shell
supervisord -c supervisord.conf
```
* 加上` -d`表示后台常驻运行。

在浏览器输入 localhost:7111 得到ok.

查看统计状态：http://127.0.0.1:9001/metrics


注意：
```text
directory: 需要守护的程序所在的目录，也就是main.exe所在目录
command: 具体执行的命令,如上所示，既可以直接执行main.exe,又可以 go run /…/…/…/…/main.go
stdout_logfile: 控制台信息打印于该文件,可不存在，会创建自动
autostart: 自动启动
autorestart: 自动重启
```


## 自重启检测
* `sudo lsof -i tcp:9001` 查看9001 或者 7111所在pid
* `sudo kill -9 8671` 假设7111对应pid8671，则杀死8671

杀死以后，查看supervisor窗口，发现已经自动重启，在浏览器输入 localhost:7111 得到ok
