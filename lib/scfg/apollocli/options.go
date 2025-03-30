package apollocli

import (
	"mylib/scfg/mgr"
)

type Options struct {
	namespaces   []string
	unmarshal    UnmarshalFunc
	sysLoggerMgr *mgr.SysLogMgr
}

type Option func(*Options)

func WithNamespace(namespace string) Option {
	return func(o *Options) {
		o.namespaces = append(o.namespaces, namespace)
	}
}

func WithUnmarshalFunc(u UnmarshalFunc) Option {
	return func(o *Options) {
		o.unmarshal = u
	}
}

func WithSysLoggerMgr(l *mgr.SysLogMgr) Option {
	return func(o *Options) {
		o.sysLoggerMgr = l
	}
}
