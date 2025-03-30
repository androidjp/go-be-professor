package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Enable        bool     `json:"enable"  yaml:"enable"`
	ClientType    string   `json:"clientType"  yaml:"clientType"` // sample、cluster、failover、sentinel
	EnableEncrypt bool     `json:"enableEncrypt" yaml:"enableEncrypt"`
	EnableTracing bool     `json:"enableTracing" yaml:"enableTracing"`
	Addrs         []string `json:"addrs"   yaml:"addrs"`
	// The sentinel master name.
	// Only failover clients.
	MasterName string `json:"masterName"  yaml:"masterName"`
	// Database to be selected after connecting to the server.
	// Only single-node and failover clients.
	DB int `json:"dbNum"  yaml:"dbNum"`
	// Common options.
	Username         string `json:"username"    yaml:"username"`
	Password         string `json:"password"    yaml:"password"`
	SentinelPassword string `json:"sentinelPassword"    yaml:"sentinelPassword"`

	MaxRetries          int           `json:"maxRetries"  yaml:"maxRetries"`
	MinRetryBackoff     time.Duration `json:"minRetryBackoff" yaml:"minRetryBackoff"`
	MaxRetryBackoff     time.Duration `json:"maxRetryBackoff" yaml:"maxRetryBackoff"`
	MinRetryBackoffSecN float64       `json:"minRetryBackoffSecN" yaml:"minRetryBackoffSecN"` // 和上面参数一样，这里填秒数 n * sec,如果此参数不为空，优先取，适配apollo配置中心
	MaxRetryBackoffSecN float64       `json:"maxRetryBackoffSecN" yaml:"maxRetryBackoffSecN"` // 和上面参数一样，这里填秒数 n * sec,如果此参数不为空，优先取，适配apollo配置中心

	DialTimeout      time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	ReadTimeout      time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout     time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	DialTimeoutSecN  float64       `json:"dialTimeoutSecN" yaml:"dialTimeoutSecN"`   // 和上面参数一样，这里填秒数 n * sec,如果此参数不为空，优先取，适配apollo配置中心
	ReadTimeoutSecN  float64       `json:"readTimeoutSecN" yaml:"readTimeoutSecN"`   // 和上面参数一样，这里填秒数 n * sec,如果此参数不为空，优先取，适配apollo配置中心
	WriteTimeoutSecN float64       `json:"writeTimeoutSecN" yaml:"writeTimeoutSecN"` // 和上面参数一样，这里填秒数 n * sec,如果此参数不为空，优先取，适配apollo配置中心

	PoolSize           int           `json:"poolSize" yaml:"poolSize"`
	MinIdleConns       int           `json:"minIdleConns" yaml:"minIdleConns"`
	MaxConnAge         time.Duration `json:"maxConnAge" yaml:"maxConnAge"`
	PoolTimeout        time.Duration `json:"poolTimeout" yaml:"poolTimeout"`
	IdleTimeout        time.Duration `json:"idleTimeout" yaml:"idleTimeout"`
	IdleCheckFrequency time.Duration `json:"idleCheckFrequency" yaml:"idleCheckFrequency"`

	MaxConnAgeSecN         float64 `json:"maxConnAgeSecN" yaml:"maxConnAgeSecN"`
	PoolTimeoutSecN        float64 `json:"poolTimeoutSecN" yaml:"poolTimeoutSecN"`
	IdleTimeoutSecN        float64 `json:"idleTimeoutSecN" yaml:"idleTimeoutSecN"`
	IdleCheckFrequencySecN float64 `json:"idleCheckFrequencySecN" yaml:"idleCheckFrequencySecN"`
	// Only cluster clients.
	MaxRedirects        int  `json:"maxRedirects"       yaml:"maxRedirects"`
	ReadOnly            bool `json:"readOnly"           yaml:"readOnly"`
	RouteByLatency      bool `json:"routeByLatency"     yaml:"routeByLatency"`
	RouteRandomly       bool `json:"routeRandomly"      yaml:"routeRandomly"`
	PasswordDecryptFunc func(psd string) (string, error)
	Hook                redis.Hook
}

func DefaultConfig() *Config {
	return &Config{
		Enable:        true,
		ClientType:    "sample",
		EnableEncrypt: false,
	}
}

func (c *Config) Build() (redis.UniversalClient, error) {

	if !c.Enable {
		//c.logger.Warn("core module: [redis client] init disable, stop build")
		return nil, nil
	}

	if len(c.Password) > 0 && c.EnableEncrypt {
		var err error

		if c.PasswordDecryptFunc == nil {
			return nil, errors.New("passwordDecryptFunc is nil, can not decrypt redis password")
		}
		c.Password, err = c.PasswordDecryptFunc(c.Password)
		if err != nil {
			return nil, fmt.Errorf("%w:build redis client: kms decrypt password failed", err)
		}
	}

	if c.MinRetryBackoffSecN > 0 {
		c.MinRetryBackoff = time.Duration(c.MinRetryBackoffSecN) * time.Second
	}

	if c.MaxRetryBackoffSecN > 0 {
		c.MaxRetryBackoff = time.Duration(c.MaxRetryBackoffSecN) * time.Second
	}

	if c.DialTimeoutSecN > 0 {
		c.DialTimeout = time.Duration(c.DialTimeoutSecN) * time.Second
	}

	if c.ReadTimeoutSecN > 0 {
		c.ReadTimeout = time.Duration(c.ReadTimeoutSecN) * time.Second
	}

	if c.WriteTimeoutSecN > 0 {
		c.WriteTimeout = time.Duration(c.WriteTimeoutSecN) * time.Second
	}

	if c.MaxConnAgeSecN > 0 {
		c.MaxConnAge = time.Duration(c.MaxConnAgeSecN) * time.Second
	}

	if c.PoolTimeoutSecN > 0 {
		c.PoolTimeout = time.Duration(c.PoolTimeoutSecN) * time.Second
	}

	if c.IdleCheckFrequencySecN > 0 {
		c.IdleCheckFrequency = time.Duration(c.IdleCheckFrequencySecN) * time.Second
	}

	var client redis.UniversalClient
	switch c.ClientType {
	case "sample":
		client = redis.NewClient(c.SampleOptions())
	case "cluster":
		client = redis.NewClusterClient(c.ClusterOptions())
	case "failover":
		client = redis.NewFailoverClient(c.FailoverOptions())
	default:
		return nil, errors.New("INVALID_REDIS_CLIENT_TYPE---invalid redis client type, must: sample, cluster, failover")
	}

	if c.EnableTracing {
		client.AddHook(c.Hook)
	}

	fmt.Println("core module: [redis client] build success!")
	return client, nil
}

func (c *Config) BuildSentinelClient() (*redis.SentinelClient, error) {
	if len(c.SentinelPassword) > 0 && c.EnableEncrypt {
		var err error
		if c.PasswordDecryptFunc == nil {
			return nil, errors.New("passwordDecryptFunc is nil, can not decrypt redis password")
		}
		c.SentinelPassword, err = c.PasswordDecryptFunc(c.SentinelPassword)
		if err != nil {
			return nil, fmt.Errorf("%w:build redis client: kms decrypt sentinel password failed", err)
		}
	}

	if c.MinRetryBackoffSecN > 0 {
		c.MinRetryBackoff = time.Duration(c.MinRetryBackoffSecN) * time.Second
	}

	if c.MaxRetryBackoffSecN > 0 {
		c.MaxRetryBackoff = time.Duration(c.MaxRetryBackoffSecN) * time.Second
	}

	if c.DialTimeoutSecN > 0 {
		c.DialTimeout = time.Duration(c.DialTimeoutSecN) * time.Second
	}

	if c.ReadTimeoutSecN > 0 {
		c.ReadTimeout = time.Duration(c.ReadTimeoutSecN) * time.Second
	}

	if c.WriteTimeoutSecN > 0 {
		c.WriteTimeout = time.Duration(c.WriteTimeoutSecN) * time.Second
	}

	if c.MaxConnAgeSecN > 0 {
		c.MaxConnAge = time.Duration(c.MaxConnAgeSecN) * time.Second
	}

	if c.PoolTimeoutSecN > 0 {
		c.PoolTimeout = time.Duration(c.PoolTimeoutSecN) * time.Second
	}

	if c.IdleCheckFrequencySecN > 0 {
		c.IdleCheckFrequency = time.Duration(c.IdleCheckFrequencySecN) * time.Second
	}

	client := redis.NewSentinelClient(c.SampleOptions())
	if c.EnableTracing {
		client.AddHook(c.Hook)
	}

	return client, nil
}

func (c *Config) SampleOptions() *redis.Options {
	addr := "127.0.0.1:6379"
	if len(c.Addrs) > 0 {
		addr = c.Addrs[0]
	}

	return &redis.Options{
		Addr:            addr,
		Username:        c.Username,
		Password:        c.Password,
		DB:              c.DB,
		MaxRetries:      c.MaxRetries,
		MinRetryBackoff: c.MinRetryBackoff,
		MaxRetryBackoff: c.MaxRetryBackoff,
		DialTimeout:     c.DialTimeout,
		ReadTimeout:     c.ReadTimeout,
		WriteTimeout:    c.WriteTimeout,
		PoolSize:        c.PoolSize,
		MinIdleConns:    c.MinIdleConns,
		// MaxConnAge:         c.MaxConnAge,
		PoolTimeout: c.PoolTimeout,
		// IdleTimeout:        c.IdleTimeout,
		ConnMaxIdleTime: c.IdleTimeout,
		// IdleCheckFrequency: c.IdleCheckFrequency,
	}
}

// Failover returns failover options created from the universal options.
func (c *Config) FailoverOptions() *redis.FailoverOptions {
	if len(c.Addrs) == 0 {
		c.Addrs = []string{"127.0.0.1:26379"}
	}

	return &redis.FailoverOptions{
		SentinelAddrs: c.Addrs,
		MasterName:    c.MasterName,

		DB:               c.DB,
		Username:         c.Username,
		Password:         c.Password,
		SentinelPassword: c.SentinelPassword,

		MaxRetries:      c.MaxRetries,
		MinRetryBackoff: c.MinRetryBackoff,
		MaxRetryBackoff: c.MaxRetryBackoff,

		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,

		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
		// MaxConnAge:         c.MaxConnAge,
		PoolTimeout: c.PoolTimeout,
		// IdleTimeout:        c.IdleTimeout,
		ConnMaxIdleTime: c.IdleTimeout,
		// IdleCheckFrequency: c.IdleCheckFrequency,
	}
}

// Cluster returns cluster options created from the universal options.
func (c *Config) ClusterOptions() *redis.ClusterOptions {
	if len(c.Addrs) == 0 {
		c.Addrs = []string{"127.0.0.1:6379"}
	}

	return &redis.ClusterOptions{
		Addrs: c.Addrs,

		Username: c.Username,
		Password: c.Password,

		MaxRedirects:   c.MaxRedirects,
		ReadOnly:       c.ReadOnly,
		RouteByLatency: c.RouteByLatency,
		RouteRandomly:  c.RouteRandomly,

		MaxRetries:      c.MaxRetries,
		MinRetryBackoff: c.MinRetryBackoff,
		MaxRetryBackoff: c.MaxRetryBackoff,

		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
		// MaxConnAge:         c.MaxConnAge,
		PoolTimeout: c.PoolTimeout,
		// IdleTimeout:        c.IdleTimeout,
		ConnMaxIdleTime: c.IdleTimeout,
		// IdleCheckFrequency: c.IdleCheckFrequency,
	}
}
