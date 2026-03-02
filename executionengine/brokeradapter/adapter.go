package brokeradapter

import (
	"fmt"
	"sync"

	"github.com/sophon2000/quantforge/executionengine"
)

// MemoryBroker 内存经纪商适配器：用于回测或单机测试
type MemoryBroker struct {
	mu     sync.Mutex
	orders map[string]*executionengine.Order
	fills  []executionengine.Fill
	nextID int
}

// NewMemoryBroker 构造
func NewMemoryBroker() *MemoryBroker {
	return &MemoryBroker{
		orders: make(map[string]*executionengine.Order),
		fills:  make([]executionengine.Fill, 0),
		nextID: 1,
	}
}

// PlaceOrder 实现 executionengine.Broker
func (b *MemoryBroker) PlaceOrder(order executionengine.Order) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if order.ID == "" {
		order.ID = fmt.Sprintf("mem-%d", b.nextID)
		b.nextID++
	}
	order.Status = executionengine.SUBMITTED
	b.orders[order.ID] = &order
	return nil
}

// CancelOrder 实现 executionengine.Broker
func (b *MemoryBroker) CancelOrder(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if o, ok := b.orders[id]; ok {
		o.Status = executionengine.CANCELED
	}
}

// Orders 返回当前订单快照
func (b *MemoryBroker) Orders() []executionengine.Order {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]executionengine.Order, 0, len(b.orders))
	for _, o := range b.orders {
		out = append(out, *o)
	}
	return out
}

// AddFill 记录成交（撮合引擎调用）
func (b *MemoryBroker) AddFill(f executionengine.Fill) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fills = append(b.fills, f)
	if o, ok := b.orders[f.OrderID]; ok {
		o.Status = executionengine.FILLED
	}
}

// Fills 返回成交记录
func (b *MemoryBroker) Fills() []executionengine.Fill {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]executionengine.Fill{}, b.fills...)
}
