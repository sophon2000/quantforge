package positionmanager

import "sync"

// Position 持仓
type Position struct {
	ID           string
	Symbol       string
	Quantity     int
	EntryPrice   float64
	CurrentPrice float64
	Profit       float64
	Status       string
}

// TradingRecord 交易记录（统计/回测）
type TradingRecord struct {
	ID           string
	Symbol       string
	Quantity     int
	EntryPrice   float64
	CurrentPrice float64
	Profit       float64
	Status       string
}

// PositionManager 持仓管理器
type PositionManager interface {
	Update(symbol string, quantityDelta int, price float64)
	Get(symbol string) (*Position, bool)
	Snapshot() []Position
}

// DefaultPositionManager 默认实现
type DefaultPositionManager struct {
	mu        sync.RWMutex
	positions map[string]*Position
}

// NewDefaultPositionManager 构造
func NewDefaultPositionManager() *DefaultPositionManager {
	return &DefaultPositionManager{
		positions: make(map[string]*Position),
	}
}

// Update 实现 PositionManager
func (m *DefaultPositionManager) Update(symbol string, quantityDelta int, price float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.positions[symbol]
	if !ok {
		p = &Position{Symbol: symbol}
		m.positions[symbol] = p
	}
	p.Quantity += quantityDelta
	p.CurrentPrice = price
	if p.Quantity == 0 {
		p.EntryPrice = 0
	} else if quantityDelta > 0 {
		// 简化：按最新价更新
		p.EntryPrice = (p.EntryPrice*float64(p.Quantity-quantityDelta) + price*float64(quantityDelta)) / float64(p.Quantity)
	}
	p.Profit = (p.CurrentPrice - p.EntryPrice) * float64(p.Quantity)
}

// Get 实现 PositionManager
func (m *DefaultPositionManager) Get(symbol string) (*Position, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.positions[symbol]
	if !ok {
		return nil, false
	}
	cp := *p
	return &cp, true
}

// Snapshot 实现 PositionManager
func (m *DefaultPositionManager) Snapshot() []Position {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Position, 0, len(m.positions))
	for _, p := range m.positions {
		out = append(out, *p)
	}
	return out
}
