package scfg

import (
	"fmt"
	"log/slog"
	"os"

	"mylib/scfg/apollo"
	"mylib/scfg/env"
	"mylib/scfg/file"
	"mylib/scfg/kcfg"

	"github.com/spf13/cast"
)

var cfg kcfg.Conf

// ParseConfig 初始化全局配置
func ParseConfig(cfgFile string, l *slog.Logger) {
	if len(cfgFile) == 0 {
		_, _ = os.Stderr.WriteString("config file path is empty.")
		os.Exit(1)
	}

	if l != nil {
		cfg = kcfg.NewConfig(file.NewSource(cfgFile), kcfg.WithLogger(l))
	} else {
		cfg = kcfg.NewConfig(file.NewSource(cfgFile))
	}

	if err := cfg.Init(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	initConfigSource(cfg)
}

// todo: 外部控制数据源加载顺序
func initConfigSource(cfg kcfg.Conf) {

	// init apollo source
	apolloCfg := &apollo.Config{}
	if v4 := cfg.GetValue("apollo"); v4.Error() == nil {
		if err := v4.Scan(apolloCfg); err == nil {
			if apolloCfg.Enable {
				if err := cfg.AppendSource(apollo.NewSource(apolloCfg, cfg.GetSysLogMgr())); err != nil {
					_, _ = os.Stderr.WriteString(fmt.Errorf("[error] append [apollo] config source failed: %w", err).Error())
					os.Exit(1)
				}
			}
		}
	}

	// // init kac source
	// kacConfig := kac.DefaultConfig()
	// if kacV := cfg.GetValue("kac"); kacV.Error() == nil {
	// 	if err := kacV.Scan(kacConfig); err == nil {
	// 		if kacConfig.Enable {
	// 			if err := cfg.AppendSource(kac.NewSource(kacConfig)); err != nil {
	// 				_, _ = os.Stderr.WriteString(fmt.Errorf("[error] append [kac] config source failed: %w", err).Error())
	// 				os.Exit(1)
	// 			}
	// 		}
	// 	}
	// }

	// init env source
	envConfig := env.DefaultConfig()
	if envV := cfg.GetValue("env"); envV.Error() == nil {
		if err := envV.Scan(envConfig); err != nil {
			envConfig = env.DefaultConfig() // error happen, use default
		}
	}

	if envConfig.Enable {
		if err := cfg.AppendSource(env.NewSource(envConfig.Prefixes...)); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Errorf("[error] append [env] config source failed: %w", err).Error())
			os.Exit(1)
		}
	}
}

func SettingLogger(l *slog.Logger) {
	s, ok := cfg.(kcfg.Setting)
	if ok {
		s.SettingLogger(l)
	}
}

func SettingDebugMode(debugMode bool) {
	s, ok := cfg.(kcfg.Setting)
	if ok {
		s.SettingDebugMode(debugMode)
	}
}

func GetValue(key string) kcfg.Value {
	return cfg.GetValue(key)
}

func GetString(key string) string {
	v := cfg.GetValue(key)
	str, _ := v.String()
	return str
}

func GetBool(key string) bool {
	v := cfg.GetValue(key)
	ok, _ := v.Bool()
	return ok
}

func GetInt(key string) int {
	v := cfg.GetValue(key)
	vv, _ := v.Int64()
	return int(vv)
}

func GetInt32(key string) int32 {
	v := cfg.GetValue(key)
	vv, _ := v.Int64()
	return cast.ToInt32(vv)
}

func GetInt64(key string) int64 {
	v := cfg.GetValue(key)
	vv, _ := v.Int64()
	return vv
}

func GetUint(key string) uint {
	v := cfg.GetValue(key)
	vv, _ := v.Int64()
	return cast.ToUint(vv)
}

func GetUint32(key string) uint32 {
	v := cfg.GetValue(key)
	vv, _ := v.Int64()
	return cast.ToUint32(vv)
}

func GetUint64(key string) uint64 {
	v := cfg.GetValue(key)
	vv, _ := v.Int64()
	return cast.ToUint64(vv)
}

func GetFloat64(key string) float64 {
	v := cfg.GetValue(key)
	vv, _ := v.Float64()
	return vv
}

func GetStringSlice(key string) []string {
	v := cfg.GetValue(key)
	return cast.ToStringSlice(v.Load())
}

func GetStringMap(key string) map[string]interface{} {
	v := cfg.GetValue(key)
	return cast.ToStringMap(v.Load())
}

func GetStringMapString(key string) map[string]string {
	v := cfg.GetValue(key)
	return cast.ToStringMapString(v.Load())
}

func GetStringMapStringSlice(key string) map[string][]string {
	v := cfg.GetValue(key)
	return cast.ToStringMapStringSlice(v.Load())
}

func UnmarshalKey(key string, rawVal interface{}) error {
	return cfg.GetValue(key).Scan(rawVal)
}

// UnmarshalKeyDefault 和上面实现一样，兼容旧接口调用
func UnmarshalKeyDefault(key string, rawVal interface{}) error {
	return cfg.GetValue(key).Scan(rawVal)
}

func Watch(key string, observer kcfg.Observer) error {
	return cfg.Watch(key, observer)
}
