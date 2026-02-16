package backtest

import (
	"quantforge/datasource"
	"sync"
)

// DefaultMatchingEngine 默认撮合引擎：限价单按 tick 价格撮合
type DefaultMatchingEngine struct {
	mu     sync.Mutex
	orders []Order
	onFill func(Symbol string, Fill Fill)
}

// NewDefaultMatchingEngine 构造，onFill 在发生成交时调用
func NewDefaultMatchingEngine(onFill func(symbol string, f Fill)) *DefaultMatchingEngine {
	return &DefaultMatchingEngine{
		orders: make([]Order, 0),
		onFill: onFill,
	}
}

// SubmitOrder 实现 MatchingEngine
func (e *DefaultMatchingEngine) SubmitOrder(order Order) {
	e.mu.Lock()
	e.orders = append(e.orders, order)
	e.mu.Unlock()
}

// Match 实现 MatchingEngine：用当前 tick 价格撮合未成交订单
func (e *DefaultMatchingEngine) Match(tick *datasource.Tick) {
	if tick == nil {
		return
	}
	e.mu.Lock()
	remaining := make([]Order, 0, len(e.orders))
	for _, o := range e.orders {
		// 简化：买单价格 >= 当前价 或 卖单价格 <= 当前价 则成交
		filled := false
		if o.Quantity > 0 { // 买入
			if o.Price >= tick.Price {
				filled = true
			}
		} else { // 卖出
			if o.Price <= tick.Price {
				filled = true
			}
		}
		if filled {
			qty := o.Quantity
			if qty < 0 {
				qty = -qty
			}
			f := Fill{Symbol: o.Symbol, Price: tick.Price, Quantity: qty}
			if e.onFill != nil {
				e.onFill(o.Symbol, f)
			}
		} else {
			remaining = append(remaining, o)
		}
	}
	e.orders = remaining
	e.mu.Unlock()
}

// PendingOrders 返回当前未成交订单数量
func (e *DefaultMatchingEngine) PendingOrders() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.orders)
}
