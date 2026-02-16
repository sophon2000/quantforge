package strategy

// Signal 信号
type Signal struct {
	Symbol string
	Signal string // 如 "BUY", "SELL"
}

// SignalEngine 信号引擎：接收并处理策略信号
type SignalEngine interface {
	OnSignal(signal Signal)
}

// Indicator 指标接口：按价格更新并输出数值
type Indicator interface {
	Update(price float64)
	Value() float64
}
