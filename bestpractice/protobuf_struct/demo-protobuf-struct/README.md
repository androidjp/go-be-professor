# protobuf 生成 golang 的最全结构体


环境准备：
```shell
go get google.golang.org/grpc
```

执行指令：
```shell
protoc --proto_path=. --proto_path=./proto/third_party --go_out=paths=source_relative:. --go-grpc_out=require_unimplemented_servers=false,paths=source_relative:.  proto/demo/demo.proto
protoc-go-inject-tag -input=./proto/demo/demo.pb.go
```