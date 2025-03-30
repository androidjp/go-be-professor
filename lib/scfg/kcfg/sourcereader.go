package kcfg

import (
	"bytes"
	"fmt"
	"sync"

	"mylib/scfg/mgr"

	"github.com/spf13/viper"
)

// Observer is config observer.
type Observer func(string, Value)

type SourceReader interface {
	Init() error
	Source() Source
	GetValue(string) (Value, error)
	Watch(string, Observer) error
	Close() error
}

type sourceReader struct {
	source        Source
	w             Watcher
	v             *viper.Viper
	cached        sync.Map
	observers     sync.Map
	observersOldV sync.Map
	sysLogMgr     *mgr.SysLogMgr
}

func NewSourceReader(source Source) SourceReader {
	return &sourceReader{
		v:         viper.New(),
		source:    source,
		sysLogMgr: mgr.NewSysLogMgr(nil, false),
	}
}

func NewSourceReaderWithLogger(source Source, l *mgr.SysLogMgr) SourceReader {
	return &sourceReader{
		v:         viper.New(),
		source:    source,
		sysLogMgr: l,
	}
}

func (s *sourceReader) Init() error {
	kvs, err := s.source.Load()
	if err != nil {
		return err
	}

	if err := s.load(kvs); err != nil {
		return err
	}

	s.w, err = s.source.Watch()
	if err != nil {
		return err
	}

	if s.w != nil {
		go s.watch(s.w)
	}

	return nil
}

func (s *sourceReader) Source() Source {
	return s.source
}

func (s *sourceReader) GetValue(key string) (Value, error) {
	if v, ok := s.cached.Load(key); ok {
		return v.(Value), nil
	}

	v, err := s.getValue(key)
	if err != nil {
		return nil, err
	}
	s.cached.Store(key, v)
	return v, nil
}

func (s *sourceReader) getValue(key string) (Value, error) {
	if s.v.Get(key) == nil {
		return nil, fmt.Errorf("key [%s] not found from [%s] config source", key, s.Source().SourceName())
	}
	v := NewValue(s.v, key, s.source.SourceName())
	v.Store(s.v.Get(key))
	return v, nil
}

func (s *sourceReader) Watch(key string, fn Observer) error {
	if s.w == nil {
		return nil
	}

	s.observers.Store(key, fn)
	v, err := s.GetValue(key)
	if err != nil {
		s.sysLogMgr.GetLogger().With("error", err).Error("get value error")
		s.observersOldV.Store(key, uint64(0)) // store 0 hashcode
		return nil
	}

	s.observersOldV.Store(key, v.HashCode())
	return nil
}

func (s *sourceReader) Close() error {
	if err := s.source.Close(); err != nil {
		return err
	}

	if s.w != nil {
		return s.w.Stop()
	}
	return nil
}

func (s *sourceReader) load(kvs []*KeyValue) (err error) {
	for _, kv := range kvs {
		s.v.SetConfigType(kv.Format) //todo: 处理默认format编码格式
		err = s.v.MergeConfig(bytes.NewReader(kv.Value))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sourceReader) watch(w Watcher) {
	for {
		kvs, err := w.WaitRead()
		if err != nil {
			s.sysLogMgr.GetLogger().With("source", s.source.SourceName()).With("error", err).Warn("config source watcher wait read error.")
			continue
		}

		if err := s.load(kvs); err != nil {
			s.sysLogMgr.GetLogger().With("source", s.source.SourceName()).With("error", err).Warn("config source watcher wait read error.")
			continue
		}

		// 不直接使用cached触发监听事件回调，新增observersOldV的原因是：设置监听者的key不一定会在cached存在
		s.cached.Range(func(key, value any) bool {
			k := key.(string)
			oV := value.(Value)
			// cache update
			if nV, err := s.getValue(k); err == nil && nV.HashCode() != oV.HashCode() {
				oV.Store(nV.Load())
			}
			return true
		})

		s.observersOldV.Range(func(key, oValue interface{}) bool {
			// find changed value
			k := key.(string)
			nV, _ := s.getValue(k)
			if nV != nil && nV.HashCode() != oValue.(uint64) {
				if o, ok := s.observers.Load(key); ok {
					o.(Observer)(k, nV) // notify observer, todo: try catch
				}
				s.observersOldV.Store(k, nV.HashCode()) // update old value
			}
			return true
		})
	}
}
