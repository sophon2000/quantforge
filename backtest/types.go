package backtest

import (
	"quantforge/datasource"
	"quantforge/strategy"
)

// MatchingEngine 回测撮合引擎：提交订单、按 tick 撮合
type MatchingEngine interface {
	SubmitOrder(order Order)
	Match(tick *datasource.Tick)
}

// ----- 事件类型 -----

// MarketEvent 行情事件
type MarketEvent struct {
	Symbol   string
	Price    float64
	Quantity int
}

// SignalEvent 信号事件（与 strategy.Signal 一致）
type SignalEvent struct {
	Signal strategy.Signal
}

// OrderEvent 订单事件
type OrderEvent struct {
	Symbol string
	Order  Order
}

// FillEvent 成交事件
type FillEvent struct {
	Symbol string
	Fill   Fill
}

// Order 回测订单（简化，无状态机）
type Order struct {
	Symbol   string
	Price    float64
	Quantity int
}

// Fill 回测成交
type Fill struct {
	Symbol   string
	Price    float64
	Quantity int
}
