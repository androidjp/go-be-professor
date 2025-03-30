package kcfg

import (
	"errors"
	"log/slog"

	"mylib/scfg/mgr"
)

type Setting interface {
	SettingLogger(*slog.Logger)
	SettingDebugMode(debugMode bool)
}

type Conf interface {
	Init() error
	GetValue(string) Value
	GetSysLogMgr() *mgr.SysLogMgr
	AppendSource(Source) error
	Watch(string, Observer) error
	Close() error
}

type config struct {
	opts          options
	readers       []SourceReader
	defaultReader SourceReader
}

// NewConfig returns a new config.
func NewConfig(defaultSource Source, opts ...Option) Conf {
	// 默认是禁止，即fatal模式
	o := options{
		sysLogMgr: mgr.NewSysLogMgr(nil, false),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &config{
		defaultReader: NewSourceReaderWithLogger(defaultSource, o.sysLogMgr),
		opts:          o,
	}
}

func (c *config) GetSysLogMgr() *mgr.SysLogMgr {
	if c.opts.sysLogMgr == nil {
		c.opts.sysLogMgr = mgr.NewSysLogMgr(nil, false)
	}
	return c.opts.sysLogMgr
}

// Init initializes config.
func (c *config) Init() error {
	var err error
	if err = c.defaultReader.Init(); err != nil {
		return err
	}

	for _, source := range c.opts.sources {
		reader := NewSourceReaderWithLogger(source, c.opts.sysLogMgr)
		if err = reader.Init(); err != nil {
			return err
		}
		c.readers = append(c.readers, reader)
	}
	return nil
}

// AppendSource appends source to config.
func (c *config) AppendSource(s Source) error {
	if s == nil {
		return errors.New("source is nil")
	}

	r := NewSourceReaderWithLogger(s, c.opts.sysLogMgr)
	if err := r.Init(); err != nil {
		return err
	}
	c.readers = append(c.readers, r)
	return nil
}

// GetValue returns Value by key.
// priority: sources > defaultSource.
func (c *config) GetValue(key string) Value {
	var v Value
	var err error
	for _, reader := range c.readers {
		if v, err = reader.GetValue(key); err == nil {
			return v
		}
		if c.opts.sysLogMgr != nil {
			c.opts.sysLogMgr.GetLogger().
				With("source", reader.Source().SourceName()).
				With("error", err).
				Debug("get value from source failed.")
		}
	}

	v, err = c.defaultReader.GetValue(key)
	if err != nil {
		if c.opts.sysLogMgr != nil {
			c.opts.sysLogMgr.GetLogger().
				With("source", c.defaultReader.Source().SourceName()).
				With("error", err).
				Warn("get value from source failed.")
		}
		return &errValue{err: err}
	}
	return v
}

// Watch add observer to watch key.
func (c *config) Watch(key string, fn Observer) error {
	if err := c.defaultReader.Watch(key, fn); err != nil {
		return err
	}

	for _, reader := range c.readers {
		if err := reader.Watch(key, fn); err != nil {
			return err
		}
	}
	return nil
}

// Close closes config.
func (c *config) Close() error {
	var err error
	if err = c.defaultReader.Close(); err != nil {
		return err
	}

	for _, reader := range c.readers {
		if err = reader.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) SettingLogger(l *slog.Logger) {
	c.opts.sysLogMgr.SettingLogger(l)
}

func (c *config) SettingDebugMode(debugMode bool) {
	c.opts.sysLogMgr.SettingDebugMode(debugMode)
}
