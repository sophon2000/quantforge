package backtestengine

import (
	"quantforge/dataengine"
	"quantforge/strategyengine"
)

// Order 回测订单（简化）
type Order struct {
	Symbol   string
	Price    float64
	Quantity int // 正买负卖
}

// Fill 回测成交
type Fill struct {
	Symbol   string
	Price    float64
	Quantity int    // 正数
	Side     string // "BUY" / "SELL"
}

// MatchingEngine 撮合引擎接口
type MatchingEngine interface {
	SubmitOrder(order Order)
	Match(tick *dataengine.Tick)
}

// ----- 事件 -----

// MarketEvent 行情事件
type MarketEvent struct {
	Symbol   string
	Price    float64
	Quantity int
}

// SignalEvent 信号事件
type SignalEvent struct {
	Signal strategyengine.Signal
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

// EventLoop 事件循环接口：驱动回测事件流
type EventLoop interface {
	PushMarket(tick *dataengine.Tick)
	PushSignal(signal strategyengine.Signal)
	Run()
	Stop()
}
