package slogger

import (
	"log/slog"
	"os"
)

// NewCustomLevelLogger 创建一个自定义级别的日志
func NewCustomLevelLogger(cfg *Config) *slog.Logger {
	ops := &slog.HandlerOptions{Level: getSlogLevel(cfg.Level),
		AddSource: cfg.ShowSource,
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
	var h slog.Handler
	if cfg.Format == LogFormatTEXT {
		h = slog.NewTextHandler(os.Stdout, ops)
	} else {
		// 默认用JSON格式（且打到标准输出）
		h = slog.NewJSONHandler(os.Stdout, ops)
	}

	return slog.New(NewContextHandler(cfg, h))
}

// AddCtxReadKeyString 让slog知道要打印某一个key string类型（slog对象级别）
func AddCtxReadKeyString(lg *slog.Logger, key string) {
	addCtxReadKey(lg, key, key)
}

// AddCtxReadKeyStruct 让slog知道要打印某一个key struct类型（slog对象级别）
func AddCtxReadKeyStruct(lg *slog.Logger, key struct{}, displayKey string) {
	addCtxReadKey(lg, key, displayKey)
}

// AddAttrs 添加属性 (支持 已经有特殊key的 context 设定读取方式)
func addCtxReadKey(lg *slog.Logger, key interface{}, displayKey string) {
	h, ok := lg.Handler().(*DefaultContextHandler)
	if !ok {
		return
	}
	h.addExtraCtxReadKey(displayKey, key)
}
