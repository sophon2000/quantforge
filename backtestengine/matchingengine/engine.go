package matchingengine

import (
	"sync"

	"quantforge/backtestengine"
	"quantforge/dataengine"
)

// DefaultMatchingEngine 默认撮合引擎
type DefaultMatchingEngine struct {
	mu     sync.Mutex
	orders []backtestengine.Order
	onFill func(symbol string, f backtestengine.Fill)
}

// NewDefaultMatchingEngine 构造
func NewDefaultMatchingEngine(onFill func(symbol string, f backtestengine.Fill)) *DefaultMatchingEngine {
	return &DefaultMatchingEngine{
		orders: make([]backtestengine.Order, 0),
		onFill: onFill,
	}
}

// SubmitOrder 实现 backtestengine.MatchingEngine
func (e *DefaultMatchingEngine) SubmitOrder(order backtestengine.Order) {
	e.mu.Lock()
	e.orders = append(e.orders, order)
	e.mu.Unlock()
}

// Match 实现 backtestengine.MatchingEngine
func (e *DefaultMatchingEngine) Match(tick *dataengine.Tick) {
	if tick == nil {
		return
	}
	e.mu.Lock()
	remaining := make([]backtestengine.Order, 0, len(e.orders))
	for _, o := range e.orders {
		filled := false
		if o.Quantity > 0 {
			if o.Price >= tick.Price {
				filled = true
			}
		} else {
			if o.Price <= tick.Price {
				filled = true
			}
		}
		if filled {
			qty := o.Quantity
			side := "BUY"
			if qty < 0 {
				qty = -qty
				side = "SELL"
			}
			f := backtestengine.Fill{Symbol: o.Symbol, Price: tick.Price, Quantity: qty, Side: side}
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

// PendingOrders 未成交订单数
func (e *DefaultMatchingEngine) PendingOrders() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.orders)
}
