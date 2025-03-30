package mgr

import (
	"log/slog"

	"mylib/slogger"
)

var _defaultClosedLogger = slogger.NewCustomLevelLogger(&slogger.Config{
	Level:      "FATAL",
	ShowSource: true,
})

type SysLogMgr struct {
	log       *slog.Logger
	debugMode bool // debug模式(默认关闭，即默认只打印Fatal级别日志)，对应的配置key为：scfg.debugMode
}

func NewSysLogMgr(log *slog.Logger, debugMode bool) *SysLogMgr {
	return &SysLogMgr{
		log:       log,
		debugMode: debugMode,
	}
}

func (mgr *SysLogMgr) GetLogger() *slog.Logger {
	if mgr.log != nil && mgr.debugMode {
		return mgr.log
	}
	return _defaultClosedLogger
}

func (mgr *SysLogMgr) SettingLogger(log *slog.Logger) {
	mgr.log = log
}

func (mgr *SysLogMgr) SettingDebugMode(debugMode bool) {
	mgr.debugMode = debugMode
}
