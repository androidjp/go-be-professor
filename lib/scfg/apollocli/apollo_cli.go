package apollocli

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"mylib/scfg/mgr"

	"github.com/philchia/agollo/v4"
)

type UnmarshalFunc func([]byte, interface{}) error

type ApolloConfig struct {
	AppID                   string   `json:"appId"                                      yaml:"appId"`
	Cluster                 string   `json:"cluster,omitempty"                          yaml:"cluster"`
	SubscribeNamespaceNames []string `json:"subscribeNamespaceNames,omitempty"          yaml:"subscribeNamespaceNames"`
	DefaultNamespace        string   `json:"defaultNamespace,omitempty"                 yaml:"defaultNamespace"`
	CacheDir                string   `json:"cacheDir,omitempty"                         yaml:"cacheDir"`
	IP                      string   `json:"ip,omitempty"                               yaml:"ip"`
	AccessKeySecret         string   `json:"accessKeySecret,omitempty"                  yaml:"accessKeySecret"`
	InsecureSkipVerify      bool     `json:"insecureSkipVerify,omitempty"               yaml:"insecureSkipVerify"`
	IsSkipLocalCache        bool     `json:"isSkipLocalCache"                           yaml:"isSkipLocalCache"`
}

type ApolloClientProxy struct {
	Client           agollo.Client
	defaultNamespace string
	unmarshal        UnmarshalFunc
	opts             *Options
	listeners        sync.Map
}

func NewApolloClientProxy(c *ApolloConfig, opts ...Option) *ApolloClientProxy {
	if len(c.DefaultNamespace) == 0 {
		c.DefaultNamespace = "application.properties"
	}

	options := Options{
		namespaces:   []string{c.DefaultNamespace},
		sysLoggerMgr: mgr.NewSysLogMgr(nil, false),
	}

	for _, o := range opts {
		o(&options)
	}

	if options.unmarshal == nil {
		options.unmarshal = GetUnmarshalFunc(json.Unmarshal) // default wrap json.Unmarshal
	} else {
		options.unmarshal = GetUnmarshalFunc(options.unmarshal)
	}

	var copts []agollo.ClientOption
	if c.IsSkipLocalCache {
		copts = append(copts, agollo.SkipLocalCache())
	}

	// if options.logger != nil {
	// 	copts = append(copts, agollo.WithLogger(options.logger))
	// }

	return &ApolloClientProxy{
		defaultNamespace: c.DefaultNamespace,
		unmarshal:        options.unmarshal,
		opts:             &options,
		Client: agollo.NewClient(&agollo.Conf{
			AppID:              c.AppID,
			Cluster:            c.Cluster,
			NameSpaceNames:     c.SubscribeNamespaceNames,
			CacheDir:           c.CacheDir,
			MetaAddr:           c.IP,
			AccesskeySecret:    c.AccessKeySecret,
			InsecureSkipVerify: c.InsecureSkipVerify,
		},
			copts...,
		),
	}
}

func (c *ApolloClientProxy) Start() error {
	c.Client.OnUpdate(c.onUpdate) // add default on update func
	return c.Client.Start()
}

func (c *ApolloClientProxy) Stop() error {
	return c.Client.Stop()
}

func (c *ApolloClientProxy) defaultNameSpaceOption() Options {
	return Options{
		namespaces: []string{c.defaultNamespace},
	}
}

func (c *ApolloClientProxy) GetInt(key string, defaultValue int, opts ...Option) int {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v, err := strconv.Atoi(c.get(key, op.namespaces))
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return v
}

func (c *ApolloClientProxy) GetInt32(key string, defaultValue int32, opts ...Option) int32 {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v, err := strconv.ParseInt(c.get(key, op.namespaces), 10, 32)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return int32(v)
}

func (c *ApolloClientProxy) GetInt64(key string, defaultValue int64, opts ...Option) int64 {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v, err := strconv.ParseInt(c.get(key, op.namespaces), 10, 64)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return v
}

func (c *ApolloClientProxy) GetUInt(key string, defaultValue uint, opts ...Option) uint {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v, err := strconv.ParseUint(c.get(key, op.namespaces), 10, 64)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return uint(v)
}

func (c *ApolloClientProxy) GetUInt32(key string, defaultValue uint32, opts ...Option) uint32 {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v, err := strconv.ParseUint(c.get(key, op.namespaces), 10, 32)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return uint32(v)
}

func (c *ApolloClientProxy) GetUInt64(key string, defaultValue uint64, opts ...Option) uint64 {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v, err := strconv.ParseUint(c.get(key, op.namespaces), 10, 64)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return v
}

func (c *ApolloClientProxy) GetFloat64(key string, defaultValue float64, opts ...Option) float64 {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v, err := strconv.ParseFloat(c.get(key, op.namespaces), 64)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return v
}

func (c *ApolloClientProxy) GetString(key string, defaultValue string, opts ...Option) string {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v := c.get(key, op.namespaces)
	if len(v) == 0 {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed: len(v) == 0", "key", key, "default_value", defaultValue)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return v
}

func (c *ApolloClientProxy) GetBool(key string, defaultValue bool, opts ...Option) bool {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v, err := strconv.ParseBool(c.get(key, op.namespaces))
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return v
}

func (c *ApolloClientProxy) GetIntSlice(key string, defaultValue []int, opts ...Option) []int {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	var intSlice []int
	err := c.unmarshal([]byte(c.get(key, op.namespaces)), &intSlice)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}
	if len(intSlice) == 0 {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed: len(intSlice) == 0", "key", key, "default_value", defaultValue)
		return defaultValue
	}
	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return intSlice
}

func (c *ApolloClientProxy) GetStringSlice(key string, defaultValue []string, opts ...Option) []string {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	var strSlice []string
	err := c.unmarshal([]byte(c.get(key, op.namespaces)), &strSlice)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}
	if len(strSlice) == 0 {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed: len(strSlice) == 0", "key", key, "default_value", defaultValue)
		return defaultValue
	}
	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return strSlice
}

func (c *ApolloClientProxy) GetStringMap(key string, defaultValue map[string]interface{}, opts ...Option) map[string]interface{} {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	var strMap map[string]interface{}
	err := c.unmarshal([]byte(c.get(key, op.namespaces)), &strMap)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}
	if len(strMap) == 0 {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed: len(strMap) == 0", "key", key, "default_value", defaultValue)
		return defaultValue
	}
	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return strMap
}

func (c *ApolloClientProxy) GetStringMapString(key string, defaultValue map[string]string, opts ...Option) map[string]string {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	var strMap map[string]string
	err := c.unmarshal([]byte(c.get(key, op.namespaces)), &strMap)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}
	if len(strMap) == 0 {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed: len(strMap) == 0", "key", key, "default_value", defaultValue)
		return defaultValue
	}
	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return strMap
}

func (c *ApolloClientProxy) GetStringMapStringSlice(key string, defaultValue map[string][]string, opts ...Option) map[string][]string {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	var strMap map[string][]string
	err := c.unmarshal([]byte(c.get(key, op.namespaces)), &strMap)
	if err != nil {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed.", "key", key, "default_value", defaultValue, "error", err)
		return defaultValue
	}
	if len(strMap) == 0 {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed: len(strMap) == 0", "key", key, "default_value", defaultValue)
		return defaultValue
	}
	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return strMap
}

// 从 namespace数组中获取指定key配置，获取到则返回，否则则每个namespace都尝试获取一次，都获取不到返回空字符串
func (c *ApolloClientProxy) get(key string, namespaces []string) string {
	if len(namespaces) == 0 {
		return c.Client.GetString(key)
	}
	var v string
	for _, n := range namespaces {
		v = c.Client.GetString(key, agollo.WithNamespace(n))
		if v != "" {
			return v
		}
	}
	return ""
}

func (c *ApolloClientProxy) GetContent(opts ...Option) string {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	for _, n := range op.namespaces {
		content := c.Client.GetContent(agollo.WithNamespace(n))
		if content != "" {
			return content
		}
	}

	return c.Client.GetContent()
}

func (c *ApolloClientProxy) UnmarshalKey(key string, rawVal interface{}, opts ...Option) error {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	v := c.get(key, op.namespaces)
	if v == "" {
		c.opts.sysLoggerMgr.GetLogger().Warn("[config source]: [default] get config from remote server failed. config is null", "key", key)
		return errors.New("get from apollo server config is null")
	}
	c.opts.sysLoggerMgr.GetLogger().Debug("[config source]: [remote]", "key", key)
	return c.unmarshal([]byte(v), rawVal)
}

func (c *ApolloClientProxy) GetAllKeys(opts ...Option) []string {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	for _, n := range op.namespaces {
		keys := c.Client.GetAllKeys(agollo.WithNamespace(n))
		if len(keys) > 0 {
			return keys
		}
	}
	return c.Client.GetAllKeys()
}

func (c *ApolloClientProxy) GetReleaseKey(opts ...Option) string {
	op := c.defaultNameSpaceOption()
	for _, o := range opts {
		o(&op)
	}

	for _, n := range op.namespaces {
		key := c.Client.GetReleaseKey(agollo.WithNamespace(n))
		if key != "" {
			return key
		}
	}
	return c.Client.GetReleaseKey()
}

func (c *ApolloClientProxy) SettingSysLogMgr(mgr *mgr.SysLogMgr) {
	c.opts.sysLoggerMgr = mgr
}
