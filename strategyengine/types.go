package strategyengine

// Signal 信号
type Signal struct {
	Symbol string
	Signal string // "BUY", "SELL"
}

// SignalGenerator 信号生成器接口
type SignalGenerator interface {
	OnSignal(signal Signal)
}

// Indicator 简单指标接口（按价格更新）
type Indicator interface {
	Update(price float64)
	Value() float64
}
