# lib

云文档服务组的共用library仓库，希望整个云文档服务端团队引用一个统一的基础库，此基础库应当封装团队内部的通用能力，包括：用户登录鉴权、云文档处理等可以给多个微服务复用的能力，结合idl库中的各类定义，可提供统一能力。

## go拉取
```shell
GOPROXY=https://mirrors.wps.cn/go/ go get gopkg.in/redis.v5@v5.0.0-wps
```


## slogger目录

### 定义
slogger目录用于针对 go1.21版本后推出的官方标准slog库的部分自定义封装和加强。


### 配置方式
例如下面是一个 app.yml配置文件：
```yaml
scfg.debugMode: true // 默认不填为false。false-不打印scfg内部的日志，true-对齐下面的slogger日志打印格式进行scfg打印

slogger:
  level: Debug  // 不区分大小写，默认是INFO
  format: JSON  // 区分大小写，默认是JSON。 TEXT- k-v文本方式打印
  showSource: true // 默认是true。 true-展示日志打印的位置，false-不展示
  ctxReadType: 0 // 默认是0。0-按照指定的key查找并打印 1-打印context中的任何key（暂不生效）
  ctxReadKeys: // ctxReadType=0 时生效，默认就会解析context中的request_id并打印，其余的可以通过这里配置进来
  - trace_id
  - user
```


### 使用方式
初始化：
```go
// 直接在服务初始化过程中，初始化config 和 全局slog
bootstrap.InitGlobalCfgAndLog("<配置文件地址>")
```

使用：
1. 现有的日志打印级别：
    ```go
    slogger.LevelTrace = -8 // 自定义，TRACE级别
    slog.LevelDebug = -4 // 系统自带
    slog.LevelInfo = 0 // 系统自带（默认值）
    slog.LevelWarn = 4 // 系统自带
    slog.LevelError = 8  // 系统自带
    slogger.LevelFatal = 12 // 自定义（不会杀死进程，需要自行处理）
    ```

2. 创建一个新的logger对象：
    ```go
    level := slog.LevelDebug
    logType := "TEXT" // 默认JSON：json格式打印，TEXT：k=v格式打印
    showSource := true // true：输出日志打印的代码位置

    logger := slogger.NewCustomLevelLogger(level, logType, showSource)

    // 使用: 打Info日志
    logger.InfoContext(ctx, "msg")
    // 使用：打自定义日志
    logger.Log(ctx, slogger.LevelTrace, "msg")
    ```

3. 设置slog全局的logger对象为你自定义的logger对象：
    ```go
    slog.SetDefault(logger)
    slog.InfoContext(ctx, "msg")

    slog.Debug(....)

    slog.DebugContext(ctx, ...)

    // Fatal等不会自动os.Exit(-1)，业务程序自行处理终止逻辑
    slog.Log(ctx, slogger.LevelFatal, "msg")

    // 打印结果: {{}}
    ctxWithReqID := context.WithValue(context.Background(), "request_id", "123456")
    ctxWithReqID = context.WithValue(ctxWithReqID, "trace_id", "123456")
    slog.With(slog.Group("user", "name", "Mike", "age", 18)).WarnContext(ctxWithReqID, "test with request_id and trace_id")
    ```


### 使用slog全局对象打印一条Info日志的例子
1. app-online.yaml配置了以下配置：
    ```yaml
    scfg.debugMode: false // 表示非业务代码中的slog打印都会被屏蔽

    slogger:
      level: Info
      format: JSON
      showSource: false
      ctxReadType: 0 // 表示采用默认的指定读取key的方式进行打印，其中request_id无需配置，有值一定打印。
      ctxReadKeys:
      - trace_id
   ```
2. go代码如下：
    ```go
    // 初始化
    InitGlobalCfgAndLog("./../mock/cfg/app-with-slogger.yaml")

    // 打印
    ctxWithReqID := context.WithValue(context.Background(), "request_id", "123456")
    ctxWithReqID = context.WithValue(ctxWithReqID, "trace_id", "123456")
    slog.With(slog.Group("user", "name", "Mike", "age", 18)).WarnContext(ctxWithReqID, "test with request_id and trace_id")
    ```
3. 打印结果为：
    ```json
    {
        "time": "2025-01-09T21:30:38.550064+08:00",
        "level": "WARN",
        "source": {
            "function": "mylib/bootstrap.TestInitGlobalCfgAndLog.func1.1.1",
            "file": "/Users/wujunpeng/projects/company_code/clouddocsrv_series/office_clouddoc_srv/lib/bootstrap/bootstrap_test.go",
            "line": 23
        },
        "msg": "test with request_id and trace_id",
        "user": {
            "name": "Mike",
            "age": 18
        }
    }
    ```
   
### 额外程序中针对slog.Logger设置需要在context中读取的key
目前仅支持 String 和 Struct{}，避免由于使用 map 和 slice等非comparable的类型作为 context设置的key 引发的异常。

slogger.AddCtxReadKeyStruct
```go
lg := slogger.NewCustomLevelLogger(&Config{
   Level:       "INFO",
   Format:      LogFormatJSON,
   ShowSource:  true,
   CtxReadType: ContextReadTypeDefault,
})
slogger.AddCtxReadKeyStruct(lg, CtxReqID, "request_id")

ctx := context.WithValue(context.Background(), CtxReqID, "123456")

lg.InfoContext(ctx, "slog: Hello, World!")
// {"time":"2025-01-14T17:47:05.714083+08:00","level":"INFO","source":{"function":"mylib/slogger.TestAddCtxReadKeyStruct.func1.1","file":"/Users/wujunpeng/projects/company_code/clouddocsrv_series/office_clouddoc_srv/lib/slogger/slogger_test.go","line":260},"msg":"slog: Hello, World!","request_id":"123456"}
```

slogger.AddCtxReadKeyString
```go
lg := slogger.NewCustomLevelLogger(&Config{
    Level:       "INFO",
    Format:      LogFormatJSON,
    ShowSource:  true,
    CtxReadType: ContextReadTypeDefault,
})
slogger.AddCtxReadKeyString(lg, "x_id")

ctx := context.WithValue(context.Background(), "x_id", "bbbb")
ctx = context.WithValue(ctx, "request_id", "123456")

lg.InfoContext(ctx, "slog: Hello, World!")
// {"time":"2025-01-14T17:45:56.557238+08:00","level":"INFO","source":{"function":"mylib/slogger.TestAddCtxReadKeyString.func1.1","file":"/Users/wujunpeng/projects/company_code/clouddocsrv_series/office_clouddoc_srv/lib/slogger/slogger_test.go","line":294},"msg":"slog: Hello, World!","request_id":"123456","x_id":"bbbb"}
```


### 其他详细用法
主要看 slogger_test.go 文件即可。


## es目录

### 定义
简单连接ES的库，提供了简单的增删改查功能。



### 使用方式

#### 已有接口说明
```go
package es

type Document interface {
	GetID() (string, error)                   // 获取文档ID
	ConvertBySearchHit(hit interface{}) error // 从搜索结果中转换
	Clone() Document                          // 深拷贝的自定义实现
}

type Client interface {
   // CreateIndex 创建索引
   CreateIndex(ctx context.Context, index, indexMappingRule string) error
   // UpsertDocument 检查文档是否存在，如果存在则更新，否则插入
   UpsertDocument(ctx context.Context, index string, doc Document) error
   // DeleteDocument 删除文档
   DeleteDocument(ctx context.Context, index string, doc Document) error
   // SearchDocumentsRaw 搜索文档(原始json方式)
   // index: 索引名称
   // raw: 搜索条件
   // raw 例子：
   // {"query": {"match": {"content": "2223333"}}}
   // from: 起始位置，默认为 0
   // size: 返回数量，不传则没有限制
   SearchDocumentsRaw(ctx context.Context, index, raw string, from, size int, doc Document) ([]Document, error)
}
```


#### 1）初始化并使用ES7版本的client
1. 初始化client
   ```go
   cfg := &es.Config{
           Enable:   true,
           URL:      "http://10.160.12.33:9200",
           Username: "wps",
           Password: "jW8Bk2jezx9Ssp2",
       }
   }
   
   cli, err := cfg.BuildES7Client()
   ```


#### 2）初始化并使用ES8版本的client
1. 初始化client
   ```go
   cfg := &es.Config{
           Enable:   true,
           URL:      "http://10.160.12.33:9200",
           Username: "wps",
           Password: "jW8Bk2jezx9Ssp2",
       }
   }
   
   cli, err := cfg.BuildES8Client()
   ```


#### 3）结合bootstrap目录进行全局初始化（自动读取apollo等的es配置）
1. app.yaml或者apollo配置好以下配置：
    ```yaml
    es:
      enable: true
      url: http://
     username: wps
     password: xxxx
2. 直接调用bootstrap的Init函数 或者 TryInitGlobalESClient函数
   ```go
   // 方式一
   bootstrap.Init()
   // 方式二
   bootstrap.TryInitGlobalESClient()
   ```
2. 使用全局的ES client
   ```go
   // 使用全局的ES client
   esCli := es.GetClient()
   ```

