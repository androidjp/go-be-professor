package kcfg

import (
	"log/slog"

	"mylib/scfg/mgr"
)

// Option is config option.
type Option func(*options)

type options struct {
	sources       []Source
	defaultSource Source
	sysLogMgr     *mgr.SysLogMgr
}

func WithSource(s ...Source) Option {
	return func(o *options) {
		o.sources = s
	}
}

func WithDefaultSource(s Source) Option {
	return func(o *options) {
		o.defaultSource = s
	}
}

func WithLogger(l *slog.Logger) Option {
	return func(o *options) {
		o.sysLogMgr.SettingLogger(l)
	}
}

func WithDebugMode(debugMode bool) Option {
	return func(o *options) {
		o.sysLogMgr.SettingDebugMode(debugMode)
	}
}

func (o *options) GetLogger() *slog.Logger {
	return o.sysLogMgr.GetLogger()
}
