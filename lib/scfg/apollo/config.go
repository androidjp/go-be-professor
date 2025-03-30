package apollo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"mylib/scfg/apollocli"
	"mylib/scfg/kcfg"
	"mylib/scfg/mgr"

	"github.com/philchia/agollo/v4"
	"gopkg.in/yaml.v2"
)

var (
	ErrNamespaceFormat = errors.New("namespace format must be .properties or .json or .yml or .yaml suffix")
)

type Config struct {
	Enable        bool                    `json:"enable"  yaml:"enable" mapstructure:"enable"`
	UnmarshalType string                  `json:"unmarshalType" yaml:"unmarshalType" mapstructure:"unmarshalType"`
	AppCfg        *apollocli.ApolloConfig `json:"appCfg" yaml:"appCfg" mapstructure:"appCfg"`
}

func (c Config) Validate() error {
	if len(c.AppCfg.DefaultNamespace) > 0 {
		if err := c.checkNamespace(c.AppCfg.DefaultNamespace); err != nil {
			return fmt.Errorf("default namespace format error:%w", err)
		}
	}

	if len(c.AppCfg.SubscribeNamespaceNames) > 0 {
		for _, ns := range c.AppCfg.SubscribeNamespaceNames {
			if err := c.checkNamespace(ns); err != nil {
				return fmt.Errorf("subscribe namespace format error:%w", err)
			}
		}
	}

	return nil
}

func (c Config) checkNamespace(ns string) error {
	if !strings.Contains(ns, ".") ||
		(!strings.HasSuffix(ns, Properties) &&
			!strings.HasSuffix(ns, YML) &&
			!strings.HasSuffix(ns, YAML) &&
			!strings.HasSuffix(ns, JSON)) {
		return ErrNamespaceFormat
	}

	return nil
}

type apollo struct {
	cliProxy *apollocli.ApolloClientProxy
	cfg      *Config
}

func NewSource(cfg *Config, logMgr *mgr.SysLogMgr) kcfg.Source {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	if !cfg.Enable {
		return nil
	}

	ap := &apollo{
		cfg: cfg,
	}

	apolloOps := []apollocli.Option{
		apollocli.WithUnmarshalFunc(GetUnmarshal(cfg.UnmarshalType)),
		apollocli.WithSysLoggerMgr(logMgr),
	}

	ap.cliProxy = apollocli.NewApolloClientProxy(
		cfg.AppCfg,
		apolloOps...,
	)

	if err := ap.cliProxy.Start(); err != nil {
		panic(err)
	}

	initFormats()
	return ap
}

func (a *apollo) Load() ([]*kcfg.KeyValue, error) {
	nss := a.cfg.AppCfg.SubscribeNamespaceNames
	if len(a.cfg.AppCfg.DefaultNamespace) <= 0 {
		a.cfg.AppCfg.DefaultNamespace = "application.properties"
	}

	var kvs []*kcfg.KeyValue
	// first load default namespace
	kv, err := a.load(a.cfg.AppCfg.DefaultNamespace)
	if err != nil {
		return nil, err
	}
	kvs = append(kvs, kv)

	// 区分 properties 和 非properties
	for _, ns := range nss {
		kv, err := a.load(ns)
		if err != nil {
			return nil, err
		}

		if kv != nil {
			kvs = append(kvs, kv)
		}
	}

	return kvs, nil
}

func (a *apollo) load(ns string) (*kcfg.KeyValue, error) {
	var str string
	if format(ns) == Properties {
		return a.loadProperties(ns)
	}

	str = a.cliProxy.Client.GetContent(agollo.WithNamespace(ns))
	if len(str) > 0 {
		return &kcfg.KeyValue{
			Key:    ns,
			Format: format(ns),
			Value:  []byte(str),
		}, nil
	}
	return nil, nil
}

func (a *apollo) loadProperties(ns string) (*kcfg.KeyValue, error) {
	allKeys := a.cliProxy.Client.GetAllKeys(agollo.WithNamespace(ns))
	configMap := make(map[string]interface{}, len(allKeys))
	for _, key := range allKeys {
		v := a.cliProxy.Client.GetString(key, agollo.WithNamespace(ns))
		configMap[key] = to(v)
	}

	jsonStr, err := json.Marshal(configMap)
	if err != nil {
		return nil, err
	}

	return &kcfg.KeyValue{
		Key:    ns,
		Format: JSON,
		Value:  jsonStr,
	}, nil
}

func (a *apollo) Watch() (kcfg.Watcher, error) {
	return newWatcher(a), nil
}

func (a *apollo) SourceName() string {
	return "apollo"
}

func (a *apollo) Close() error {
	return a.cliProxy.Stop() //todo check stop safe
}

func GetUnmarshal(unmarshalType string) apollocli.UnmarshalFunc {
	switch unmarshalType {
	case "json":
		return json.Unmarshal
	case "yaml":
		return yaml.Unmarshal
	default:
		return yaml.Unmarshal
	}
}

const (
	YAML       = "yaml"
	YML        = "yml"
	JSON       = "json"
	Properties = "properties"
)

var formats map[string]struct{}
var nsMap map[string]string

func initFormats() {
	formats = make(map[string]struct{})
	formats[YAML] = struct{}{}
	formats[YML] = struct{}{}
	formats[JSON] = struct{}{}
	formats[Properties] = struct{}{}

	nsMap = make(map[string]string) // namespace -> namespace.suffix, eg: application -> application.properties
}

func format(ns string) string {
	arr := strings.Split(ns, ".")
	nsMap[arr[0]] = ns
	if len(arr) > 1 {
		nsMap[ns] = ns
	}

	suffix := arr[len(arr)-1]
	if len(arr) <= 1 {
		return JSON
	}

	if _, ok := formats[suffix]; !ok {
		// fallback
		return JSON
	}

	return suffix
}

// GetNSMapValue get namespace map value
func GetNSMapValue(ns string) string {
	return nsMap[ns]
}

func to(v string) interface{} {
	var decodedValue interface{}
	if err := json.Unmarshal([]byte(v), &decodedValue); err == nil {
		return decodedValue
	}
	return v
}
