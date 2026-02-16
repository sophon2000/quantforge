package monitoring

// Monitor 监控接口：指标上报、告警、仪表盘等
type Monitor interface {
	// Emit 上报指标（名称、值、标签可选）
	Emit(name string, value float64, tags map[string]string)
	// Health 健康检查
	Health() error
}
