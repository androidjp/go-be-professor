package main

import (
	"context"
	"log"
	"log/slog"
	"os"
)

func main() {
	// 原始log
	log.Println("log: Hello, World!")
	// slog (默认日志级别为INFO)
	slog.Debug("slog: Hello, World!")
	slog.Info("slog: Hello, World!")
	slog.Warn("slog: Hello, World!")
	slog.Error("slog: Hello, World!")
	// slog with context
	slog.SetLogLoggerLevel(slog.LevelWarn)
	type contextKey string
	const traceIDKey contextKey = "trace_id"
	ctx := context.WithValue(context.Background(), traceIDKey, "123456")
	slog.DebugContext(ctx, "slog: Hello, World!", "level", "LevelWarn")
	slog.InfoContext(ctx, "slog: Hello, World!", "level", "LevelWarn")
	slog.WarnContext(ctx, "slog: Hello, World!", "level", "LevelWarn")
	slog.ErrorContext(ctx, "slog: Hello, World!", "level", "LevelWarn")

	// slog with fields
	slog.With(slog.String("trace_id", "123123")).Info("slog: Hello, World!")
	slog.With(slog.String("trace_id", "123123"), slog.String("user_id", "123")).Info("slog: Hello, World!")

	// slog with JSON
	// level is DEBUG
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.With(slog.String("trace_id", "123123"), slog.String("user_id", "123")).Debug("slog: Hello, World!")
	logger.With(slog.String("trace_id", "344444"), slog.String("user_id", "455555")).Debug("slog: Hello, World!")

	// set default logger
	slog.With(slog.String("trace_id", "344444"), slog.String("user_id", "455555")).Warn("slog: Hello, World!", "set default logger", "before")
	slog.SetDefault(logger)
	slog.With(slog.String("trace_id", "344444"), slog.String("user_id", "455555")).Warn("slog: Hello, World!", "set default logger", "after")

	// open source开启位置记录
	loggerWithSource := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	slog.SetDefault(loggerWithSource)
	slog.With(slog.String("trace_id", "344444"), slog.String("user_id", "455555")).Warn("slog: Hello, World!", "set default logger", "after")

	// 添加自定义的json结构体作为字段
	slog.With(slog.Group("user", slog.String("name", "xiaoMing"), slog.Int("age", 18))).Warn("slog: with user info json", "gg", 18)

	// 自定义日志级别
	slog.SetDefault(NewCustomLevelLogger(LevelTrace, LogTypeJSON, true))
	slog.DebugContext(ctx, "slog: Tracing 可以打debug（-4）")
	slog.Log(ctx, LevelTrace, "slog: Tracing 可以打Tracing（-2）")
	slog.InfoContext(ctx, "slog: Tracing 可以打Info（0）")
	slog.SetDefault(NewCustomLevelLogger(LevelFatal, LogTypeTEXT, true))
	slog.DebugContext(ctx, "slog: Fatal 不可以打debug（-4）")
	slog.Log(ctx, LevelTrace, "slog: Fatal 不可以打Tracing（-2）")
	slog.InfoContext(ctx, "slog: Fatal 不可以打Info（0）")
	slog.Log(ctx, LevelFatal, "slog: Fatal 可以打Fatal（12）")
}

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)

	LogTypeJSON = "JSON"
	LogTypeTEXT = "TEXT"
)

var lvNameMap = map[slog.Level]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

var cfgLogLevelMap = map[string]slog.Level{
	"TRACE": LevelTrace,
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
	"FATAL": LevelFatal,
}

// 例子：读取环境变量中的日志级别 以及 日志打印格式
func NewCustomLevelLoggerWithEnv() *slog.Logger {
	// 从环境变量中获取日志级别
	levelStr := os.Getenv("LOG_LEVEL")
	lv, ok := cfgLogLevelMap[levelStr]
	if !ok {
		// 默认Info级别
		lv = slog.LevelInfo
	}

	// 从环境变量中获取日志打印格式
	logType := os.Getenv("LOG_TYPE")
	if logType != "TEXT" {
		logType = "JSON"
	}
	// 是否显示日志打印的代码位置
	showSource := os.Getenv("LOG_SOURCE") == "true"

	return NewCustomLevelLogger(lv, logType, showSource)
}

// NewCustomLevelLogger 创建一个自定义级别的日志
// lv: 日志级别
// logType: 日志打印格式，TEXT或JSON
// showSource: 是否显示日志打印的代码位置
func NewCustomLevelLogger(lv slog.Level, logType string, showSource bool) *slog.Logger {
	ops := &slog.HandlerOptions{Level: lv,
		AddSource: showSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 替换日志级别的显示方式，将 DEBUG-4等 替换为 TRACE等map传入的自定义内容
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := lvNameMap[level]
				if !exists {
					levelLabel = level.String()
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}
	if logType == "TEXT" {
		return slog.New(slog.NewTextHandler(os.Stdout, ops))
	}
	// 默认用JSON格式（且打到标准输出）
	return slog.New(slog.NewJSONHandler(os.Stdout, ops))
}
