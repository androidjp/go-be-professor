package bootstrap

import (
	"log/slog"
	"mylib/es"

	"mylib/scfg"
	"mylib/slogger"
)

const (
	PropertyKeySlogger       = "slogger"
	PropertyKeyScfgDebugMode = "scfg.debugMode"
	PropertyKeyES            = "es"
	PropertyKeyNeoMetric     = "neo_metric"
)

func Init(cfgFile string) {
	// 初始化全局配置和日志
	InitGlobalCfgAndLog(cfgFile)

	// 初始化全局ES客户端【可选】
	TryInitGlobalESClient()

	// 初始化全局NeoMetric【可选】
	// TryInitNeoMetric()
}

// InitGlobalCfgAndLog 初始化全局配置和日志
func InitGlobalCfgAndLog(cfgFile string) {
	// 1) 从配置文件中初始化全局配置
	scfg.ParseConfig(cfgFile, slog.Default())
	// 2) 从scfg中初始化全局slog
	initGlobalSlogBySCfg()
	scfg.SettingLogger(slog.Default())
	// 3） 环境变量读取scfg.DebugMode 是否为true
	scfg.SettingDebugMode(scfg.GetBool(PropertyKeyScfgDebugMode))
}

// TryInitGlobalESClient 尝试初始化全局ES客户端
func TryInitGlobalESClient() {
	var (
		err error
		c   es.Config
	)
	if err = scfg.UnmarshalKeyDefault(PropertyKeyES, &c); err != nil {
		slog.Warn("TryInitGlobalESClient fail", "error", err)
	}
	if err = c.InitGlobalESClient(); err != nil {
		slog.Warn("escfg.InitGlobalESClient fail", "error", err)
	}
}

// func TryInitNeoMetric() {
// 	neoCfg := metric.NeoMetricConfig{}
// 	if err := scfg.UnmarshalKey(PropertyKeyNeoMetric, &neoCfg); err != nil {
// 		// 拿不到，就用兜底逻辑
// 		slog.Warn("TryInitNeoMetric failed", "error", err.Error())
// 		neoCfg = metric.NeoMetricConfig{
// 			Enable: true,
// 			Path:   "/neometric",
// 			Port:   "16060",
// 		}
// 	}
// 	metric.TryInitNeoMetric(&neoCfg)
// }

// initGlobalSlogBySCfg 从scfg中初始化全局slog
/**
slogger:
  level: "warn"
  format: "JSON"
  showSource: true
*/
func initGlobalSlogBySCfg() {
	// 从环境变量中获取日志级别
	var sloggerCfg *slogger.Config
	if err := scfg.UnmarshalKey(PropertyKeySlogger, &sloggerCfg); err != nil {
		sloggerCfg = slogger.DefaultCfg
	}
	if len(sloggerCfg.Level) == 0 {
		sloggerCfg = slogger.DefaultCfg
	}
	if sloggerCfg.CtxReadType != slogger.ContextReadTypeAllKeys {
		sloggerCfg.CtxReadType = slogger.ContextReadTypeDefault
	}
	slog.Debug("initGlobalSlogBySCfg", "slogger_cfg", sloggerCfg)
	slog.SetDefault(slogger.NewCustomLevelLogger(sloggerCfg))
}
