package accountsimulator

import (
	"sync"

	"quantforge/backtestengine"
)

// AccountSimulator 回测账户模拟：资金与持仓
type AccountSimulator interface {
	// ApplyFill 应用一笔成交，更新持仓与权益
	ApplyFill(f backtestengine.Fill)
	// Equity 当前权益（简化：现金 + 持仓市值）
	Equity() float64
	// Position 某标的持仓量
	Position(symbol string) int
}

// DefaultAccountSimulator 默认实现
type DefaultAccountSimulator struct {
	mu      sync.Mutex
	cash    float64
	positions map[string]int   // symbol -> quantity
	avgCost  map[string]float64 // symbol -> 成本
}

// NewDefaultAccountSimulator 构造，initialCash 初始资金
func NewDefaultAccountSimulator(initialCash float64) *DefaultAccountSimulator {
	return &DefaultAccountSimulator{
		cash:      initialCash,
		positions: make(map[string]int),
		avgCost:   make(map[string]float64),
	}
}

// ApplyFill 实现 AccountSimulator
func (s *DefaultAccountSimulator) ApplyFill(f backtestengine.Fill) {
	s.mu.Lock()
	defer s.mu.Unlock()
	qty := f.Quantity
	cost := f.Price * float64(qty)
	if f.Side == "BUY" {
		s.cash -= cost
		s.positions[f.Symbol] += qty
		oldQty := s.positions[f.Symbol] - qty
		var oldCost float64
		if oldQty > 0 {
			oldCost = s.avgCost[f.Symbol] * float64(oldQty)
		}
		s.avgCost[f.Symbol] = (oldCost + cost) / float64(s.positions[f.Symbol])
	} else {
		s.cash += cost
		s.positions[f.Symbol] -= qty
		if s.positions[f.Symbol] <= 0 {
			delete(s.positions, f.Symbol)
			delete(s.avgCost, f.Symbol)
		}
	}
}

// Equity 实现 AccountSimulator（简化：仅现金，未算持仓市值）
func (s *DefaultAccountSimulator) Equity() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cash
}

// Position 实现 AccountSimulator
func (s *DefaultAccountSimulator) Position(symbol string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.positions[symbol]
}
