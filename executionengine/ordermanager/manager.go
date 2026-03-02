package ordermanager

import (
	"sync"

	"github.com/sophon2000/quantforge/executionengine"
)

// OrderManager 订单管理器：维护订单生命周期与状态
type OrderManager interface {
	Submit(order executionengine.Order) error
	Cancel(orderID string)
	Get(orderID string) (*executionengine.Order, bool)
	OnFill(f executionengine.Fill)
}

// DefaultOrderManager 默认实现：委托 Broker 并维护本地状态
type DefaultOrderManager struct {
	mu     sync.Mutex
	broker executionengine.Broker
	orders map[string]*executionengine.Order
}

// NewDefaultOrderManager 构造
func NewDefaultOrderManager(broker executionengine.Broker) *DefaultOrderManager {
	return &DefaultOrderManager{
		broker: broker,
		orders: make(map[string]*executionengine.Order),
	}
}

// Submit 实现 OrderManager
func (m *DefaultOrderManager) Submit(order executionengine.Order) error {
	if err := m.broker.PlaceOrder(order); err != nil {
		return err
	}
	m.mu.Lock()
	m.orders[order.ID] = &order
	m.mu.Unlock()
	return nil
}

// Cancel 实现 OrderManager
func (m *DefaultOrderManager) Cancel(orderID string) {
	m.broker.CancelOrder(orderID)
	m.mu.Lock()
	if o, ok := m.orders[orderID]; ok {
		o.Status = executionengine.CANCELED
	}
	m.mu.Unlock()
}

// Get 实现 OrderManager
func (m *DefaultOrderManager) Get(orderID string) (*executionengine.Order, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	o, ok := m.orders[orderID]
	if !ok {
		return nil, false
	}
	cp := *o
	return &cp, true
}

// OnFill 实现 OrderManager
func (m *DefaultOrderManager) OnFill(f executionengine.Fill) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if o, ok := m.orders[f.OrderID]; ok {
		o.Status = executionengine.FILLED
	}
}
