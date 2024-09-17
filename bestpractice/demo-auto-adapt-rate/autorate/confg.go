package autorate

type AutoRateConfig struct {
	RefreshIntervalSec    int64   `mapstructure:"refresh_interval_sec" json:"refresh_interval_sec" yaml:"refresh_interval_sec"`                   // autoRate的刷新间隔
	RateStep              float64 `mapstructure:"rate_step" json:"rate_step" yaml:"rate_step"`                                                    // 速率调整步长（如： 0.1，即往上调10% 或者往下压10%）
	MaxNumber             float64 `mapstructure:"max_number" json:"max_number" yaml:"max_number"`                                                 // 最大数字（如： 1000）
	MinNumber             float64 `mapstructure:"min_number" json:"min_number" yaml:"min_number"`                                                 // 最小数字（如： 800）
	DefaultNumber         float64 `mapstructure:"default_number" json:"default_number" yaml:"default_number"`                                     // 默认数字（如： 100）
	AvgFuncTimeRaiseMaxMS uint64  `mapstructure:"avg_func_time_raise_max_ms" json:"avg_func_time_raise_max_ms" yaml:"avg_func_time_raise_max_ms"` // 平均执行时间 最大阈值（ms）
	AvgFuncTimeRaiseMinMS uint64  `mapstructure:"avg_func_time_raise_min_ms" json:"avg_func_time_raise_min_ms" yaml:"avg_func_time_raise_min_ms"` // 平均执行时间 最小阈值（ms）
	// AutoRateFuncKey    string                 `json:"auto_rate_func_key"`   // 自动速率函数的key（如： "func1"）
	// Extra              map[string]interface{} `json:"extra"`                // 额外配置项目
}

// func GetAutoRateFunc(key string) AutoRateFunc {
// 	switch key {
// 	case "TimeRaiseOverflow":
// 		return TimeRaiseOverflow
// 	default:
// 		return TimeRaiseOverflow
// 	}
// }

// type AutoRateFunc func(data *AutoRateData)
