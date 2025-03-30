package slogger

import (
	"log/slog"
	"strings"
)

// ctxReadType 上下文读取类型
type ctxReadType int

const (
	ContextReadTypeDefault ctxReadType = iota
	ContextReadTypeAllKeys
)

const (
	// 自定义日志级别
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)

	// 日志格式
	LogFormatJSON = "JSON"
	LogFormatTEXT = "TEXT"
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

func getSlogLevel(levelStr string) slog.Level {
	lv, ok := cfgLogLevelMap[strings.ToUpper(levelStr)]
	if !ok {
		// 默认Info级别
		lv = slog.LevelInfo
	}
	return lv
}

type Config struct {
	Level       string      `json:"level" yaml:"level"`             // 日志级别
	Format      string      `json:"format" yaml:"format"`           // 日志格式
	ShowSource  bool        `json:"showSource" yaml:"showSource"`   // 是否显示代码位置
	CtxReadType ctxReadType `json:"ctxReadType" yaml:"ctxReadType"` // 上下文读取类型:0-默认读取，1-读取所有
	CtxReadKeys []string    `json:"ctxReadKeys" yaml:"ctxReadKeys"` // 上下文读取的key (ctxReadType=0时有效)
}

var (
	DefaultCfg = &Config{
		Level:       "INFO",
		Format:      "JSON",
		ShowSource:  true,
		CtxReadType: ContextReadTypeDefault,
	}
)
