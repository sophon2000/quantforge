package execution

import (
	"fmt"
	"sync"
)

// MemoryBroker 内存经纪商：仅记录订单与成交，用于回测或单机测试
type MemoryBroker struct {
	mu     sync.Mutex
	orders map[string]*Order
	fills  []Fill
	nextID int
}

// NewMemoryBroker 构造内存经纪商
func NewMemoryBroker() *MemoryBroker {
	return &MemoryBroker{
		orders: make(map[string]*Order),
		fills:  make([]Fill, 0),
		nextID: 1,
	}
}

// PlaceOrder 下单并立即置为 SUBMITTED（可扩展为模拟成交）
func (b *MemoryBroker) PlaceOrder(order Order) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if order.ID == "" {
		order.ID = fmt.Sprintf("mem-%d", b.nextID)
		b.nextID++
	}
	order.Status = SUBMITTED
	b.orders[order.ID] = &order
	return nil
}

// CancelOrder 撤单
func (b *MemoryBroker) CancelOrder(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if o, ok := b.orders[id]; ok {
		o.Status = CANCELED
	}
}

// Orders 返回当前订单快照
func (b *MemoryBroker) Orders() []Order {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Order, 0, len(b.orders))
	for _, o := range b.orders {
		out = append(out, *o)
	}
	return out
}

// AddFill 记录一笔成交（回测或模拟时由撮合引擎调用）
func (b *MemoryBroker) AddFill(f Fill) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fills = append(b.fills, f)
	if o, ok := b.orders[f.OrderID]; ok {
		o.Status = FILLED
	}
}

// Fills 返回成交记录
func (b *MemoryBroker) Fills() []Fill {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]Fill{}, b.fills...)
}
