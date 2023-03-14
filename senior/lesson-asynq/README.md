# Go异步任务解决方案 asynq
https://juejin.cn/post/7196907808225738811?

## 原理
工作原理：
1. 客户端将任务放入队列（队列通过redis支撑）
2. 服务器从队列中拉出任务并为每一个任务启动一个工作goroutine
3. 多个工作goroutine同时处理任务

asynq.Task 结构解析：
```go
type Task struct {
	// 一个简单的字符串值，表示要执行的任务的类型.
	typename string

	// 有效载荷保存执行任务所需的数据，有效负载值必须是可序列化的.
	payload []byte

	// 保存任务的选项.
	opts []Option

	// 任务的结果编写器.
	w *ResultWriter
}
```

## demo 启动运行流程
1. go run main.go 启动服务端 ， 内部 的 runCli()  会 让内部协程 尝试扔出几个任务。

其中，redis 可以看 redis-docker.yaml

conf/redis.conf 配置如下：
```
#开启保护
protected-mode yes
#开启远程连接 
#bind 127.0.0.1 
#自定义密码
requirepass 123456 
port 6379
timeout 0
# 900s内至少一次写操作则执行bgsave进行RDB持久化
save 900 1 
save 300 10
save 60 10000
# rdbcompression ；默认值是yes。对于存储到磁盘中的快照，可以设置是否进行压缩存储。如果是的话，redis会采用LZF算法进行压缩。如果你不想消耗CPU来进行压缩的话，可以设置为关闭此功能，但是存储在磁盘上的快照会比较大。
rdbcompression yes
# dbfilename ：设置快照的文件名，默认是 dump.rdb
dbfilename dump.rdb
# dir：设置快照文件的存放路径，这个配置项一定是个目录，而不能是文件名。使用上面的 dbfilename 作为保存的文件名。
dir /data
# 默认redis使用的是rdb方式持久化，这种方式在许多应用中已经足够用了。但是redis如果中途宕机，会导致可能有几分钟的数据丢失，根据save来策略进行持久化，Append Only File是另一种持久化方式， 可以提供更好的持久化特性。Redis会把每次写入的数据在接收后都写入appendonly.aof文件，每次启动时Redis都会先把这个文件的数据读入内存里，先忽略RDB文件。默认值为no。
appendonly yes
# appendfilename ：aof文件名，默认是"appendonly.aof"
# appendfsync：aof持久化策略的配置；no表示不执行fsync，由操作系统保证数据同步到磁盘，速度最快；always表示每次写入都执行fsync，以保证数据同步到磁盘；everysec表示每秒执行一次fsync，可能会导致丢失这1s数据
appendfsync everysec
```

## webUI工具
执行以下命名：
```shell
docker-compose -f asynq_webui_docker.yaml up -d
```

即可启动，然后，访问：http://127.0.0.1:8980/



