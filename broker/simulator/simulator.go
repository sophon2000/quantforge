package simulator

import (
	"sync"
	"time"

	"quantforge/backtestengine"
	"quantforge/broker"
)

// DefaultSimulator 默认回测账户实现
type DefaultSimulator struct {
	mu            sync.Mutex
	cash          float64
	initialCash   float64
	fees          float64
	monthlyVol    int
	lastTradeTime time.Time
	commission    broker.CommissionModel
	positions     map[string]int
	avgCost       map[string]float64
	lastPrice     map[string]float64
	cycleCost     map[string]float64
	cycleProceeds map[string]float64
	totalCycles   int
	successCycles int
	SuccessPct    float64
}

// New 构造，initialCash 初始资金，commission 费率模型（如 ibkr.NewCommission(ibkr.Tiered)）
func New(initialCash float64, commission broker.CommissionModel) *DefaultSimulator {
	return &DefaultSimulator{
		initialCash:   initialCash,
		cash:          initialCash,
		commission:    commission,
		positions:     make(map[string]int),
		avgCost:       make(map[string]float64),
		lastPrice:     make(map[string]float64),
		cycleCost:     make(map[string]float64),
		cycleProceeds: make(map[string]float64),
	}
}

// 编译期断言：*DefaultSimulator 实现 broker.Account
var _ broker.Account = (*DefaultSimulator)(nil)

// ApplyFill 实现 broker.Account
func (s *DefaultSimulator) ApplyFill(f backtestengine.Fill) {
	s.mu.Lock()
	defer s.mu.Unlock()
	qty := f.Quantity
	cost := f.Price * float64(qty)

	if !s.lastTradeTime.IsZero() &&
		(f.Time.Year() != s.lastTradeTime.Year() || f.Time.Month() != s.lastTradeTime.Month()) {
		s.monthlyVol = 0
	}
	trade := broker.Trade{
		Shares:     qty,
		Price:      f.Price,
		IsSell:     f.Side == backtestengine.SELL,
		MonthlyVol: s.monthlyVol,
	}
	fee := s.commission.Calculate(trade)
	s.monthlyVol += qty
	s.lastTradeTime = f.Time
	s.fees += fee
	if f.Side == backtestengine.BUY {
		s.cash -= cost + fee
		s.positions[f.Symbol] += qty
		oldQty := s.positions[f.Symbol] - qty
		var oldCost float64
		if oldQty > 0 {
			oldCost = s.avgCost[f.Symbol] * float64(oldQty)
		}
		s.avgCost[f.Symbol] = (oldCost + cost) / float64(s.positions[f.Symbol])
		s.cycleCost[f.Symbol] += cost
	} else {
		proceeds := f.Price * float64(qty)
		s.cash += proceeds - fee
		s.cycleProceeds[f.Symbol] += proceeds
		s.positions[f.Symbol] -= qty
		if s.positions[f.Symbol] <= 0 {
			s.totalCycles++
			if s.cycleProceeds[f.Symbol] > s.cycleCost[f.Symbol] {
				s.successCycles++
			}
			if s.totalCycles > 0 {
				s.SuccessPct = float64(s.successCycles) / float64(s.totalCycles) * 100
			}
			delete(s.positions, f.Symbol)
			delete(s.avgCost, f.Symbol)
			delete(s.cycleCost, f.Symbol)
			delete(s.cycleProceeds, f.Symbol)
		}
	}
}

// Cash 实现 broker.Account
func (s *DefaultSimulator) Cash() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cash
}

// Equity 实现 broker.Account
func (s *DefaultSimulator) Equity() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	eq := s.cash
	for sym, qty := range s.positions {
		if qty <= 0 {
			continue
		}
		price := s.lastPrice[sym]
		if price <= 0 {
			price = s.avgCost[sym]
		}
		eq += price * float64(qty)
	}
	return eq
}

// ReturnPct 收益率
func (s *DefaultSimulator) ReturnPct() float64 {
	if s.initialCash != 0 {
		return (s.Equity() - s.initialCash) / s.initialCash * 100
	}
	return 0
}

// Fees 累计手续费
func (s *DefaultSimulator) Fees() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.fees
}

// Position 实现 broker.Account
func (s *DefaultSimulator) Position(symbol string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.positions[symbol]
}

// UpdatePrice 实现 broker.Account
func (s *DefaultSimulator) UpdatePrice(symbol string, price float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lastPrice == nil {
		s.lastPrice = make(map[string]float64)
	}
	s.lastPrice[symbol] = price
}
