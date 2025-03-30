package es

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/olivere/elastic/v7"
)

type EType string

const (
	ES7 = EType("es7")
	ES8 = EType("es8")
)

type Config struct {
	Enable      bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	Type        EType  `mapstructure:"type" json:"type" yaml:"type"` // es7, es8, 默认（es8）
	URL         string `mapstructure:"url" json:"url" yaml:"url"`
	Username    string `mapstructure:"username" json:"username" yaml:"username"`
	Password    string `mapstructure:"password" json:"password" yaml:"password"`
	EnableDebug bool   `mapstructure:"enable_debug" json:"enable_debug" yaml:"enable_debug"`
}

var (
	_globalES7Client Client
	_globalES8Client Client
	_emptyClient     = &EmptyClient{}
)

func (c *Config) InitGlobalESClient() error {
	var (
		err error
	)
	if !c.Enable {
		return fmt.Errorf("es is not enabled")
	}
	if len(c.URL) == 0 {
		return fmt.Errorf("es url is empty")
	}
	switch c.Type {
	case ES7:
		_globalES7Client, err = c.BuildES7Client()
		if err != nil {
			return err
		}
	case ES8:
		_globalES8Client, err = c.BuildES8Client()
		if err != nil {
			return err
		}
	default:
		_globalES8Client, err = c.BuildES8Client()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) BuildES7Client() (*ClientES7, error) {
	if !c.Enable {
		return nil, fmt.Errorf("es7 is not enabled")
	}
	es, err := elastic.NewClient(
		elastic.SetURL(c.URL),
		elastic.SetBasicAuth(c.Username, c.Password),
	)
	if err != nil {
		return nil, err
	}
	client := &ClientES7{
		es: es,
	}

	return client, nil
}

func (c *Config) BuildES8Client() (*ClientES8, error) {
	if !c.Enable {
		return nil, fmt.Errorf("es8 is not enabled")
	}
	cfg := elasticsearch.Config{
		Addresses: []string{
			c.URL,
		},
		Username: c.Username,
		Password: c.Password,
		//CloudID:                  "",
		//APIKey:                   "",
		//ServiceToken:             "",
		//CertificateFingerprint:   "",
		//Header:                   nil,
		//CACert:                   nil,
		//RetryOnStatus:            nil,
		//DisableRetry:             false,
		//MaxRetries:               0,
		//RetryOnError:             nil,
		//CompressRequestBody:      false,
		//CompressRequestBodyLevel: 0,
		//PoolCompressor:           false,
		//DiscoverNodesOnStart:     false,
		//DiscoverNodesInterval:    0,
		//EnableMetrics:            false,
		EnableDebugLogger: c.EnableDebug,
		//EnableCompatibilityMode:  false,
		//DisableMetaHeader:        false,
		//RetryBackoff:             nil,
		//Transport:                nil,
		//Logger:                   nil,
		//Selector:                 nil,
		//ConnectionPoolFunc:       nil,
		//Instrumentation:          nil,
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}
	client := &ClientES8{
		es: es,
	}

	return client, nil
}

func GetClient() Client {
	if _globalES8Client != nil {
		return _globalES8Client
	}
	if _globalES7Client != nil {
		return _globalES7Client
	}
	return _emptyClient
}
