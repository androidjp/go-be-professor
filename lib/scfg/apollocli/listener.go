package apollocli

import (
	"fmt"

	"github.com/philchia/agollo/v4"
)

// Fixme: 临时实现，有bug再修
type ChangeInfo struct {
	NS         string
	Key        string
	OldValue   string
	NewValue   string
	ChangeType string
}

type OnKeyUpdateFunc func(*ChangeInfo)           // 监听指定key
type OnNSUpdateFunc func(map[string]*ChangeInfo) // 监听整个namespace

func (c *ApolloClientProxy) RegisterCallBackWithKeyNS(key, ns string, f OnKeyUpdateFunc) {
	k := fmt.Sprintf("%s:%s", key, ns)
	c.listeners.Store(k, f)
}

func (c *ApolloClientProxy) RegisterCallBackWithNS(ns string, f OnNSUpdateFunc) {
	c.listeners.Store(ns, f)
}

func (c *ApolloClientProxy) onUpdate(event *agollo.ChangeEvent) {
	defer func() {
		if v := recover(); v != nil {
			switch err := v.(type) {
			case error:
				if c.opts.sysLoggerMgr != nil {
					c.opts.sysLoggerMgr.GetLogger().With("error", err).Error("on update func callback panic")
				}
			}
		}
	}()

	ns := event.Namespace
	changes := event.Changes
	// callback register ns func
	if v, ok := c.listeners.Load(ns); ok {
		if f, ok := v.(OnNSUpdateFunc); ok {
			cInfoM := make(map[string]*ChangeInfo)
			for k, v := range changes {
				cInfo := &ChangeInfo{
					NS:         ns,
					Key:        v.Key,
					OldValue:   v.OldValue,
					NewValue:   v.NewValue,
					ChangeType: v.ChangeType.String(),
				}
				cInfoM[k] = cInfo
			}
			go f(cInfoM)
		}
	}

	// callback register key:ns func
	for k, v := range changes {
		key := fmt.Sprintf("%s:%s", k, ns)
		if vo, ok := c.listeners.Load(key); ok {
			if f, ok := vo.(OnKeyUpdateFunc); ok {
				cInfo := &ChangeInfo{
					NS:         ns,
					Key:        v.Key,
					OldValue:   v.OldValue,
					NewValue:   v.NewValue,
					ChangeType: v.ChangeType.String(),
				}
				go f(cInfo)
			}
		}
	}
}
