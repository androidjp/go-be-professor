package metrics

import "go.uber.org/zap"

// RecordHandler 指标处理
func RecordHandler(logger *zap.Logger) func(msg *MetricsMessage) {
	if logger == nil {
		panic("logger required")
	}

	return func(msg *MetricsMessage) {
		RecordMetrics(
			msg.Method,
			msg.Path,
			msg.IsSuccess,
			msg.HTTPCode,
			msg.BusinessCode,
			msg.CostSeconds,
			msg.TraceID,
		)
	}
}
