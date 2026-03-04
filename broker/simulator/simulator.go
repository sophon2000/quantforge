package simulator

import (
	"math"
	"sync"
	"time"

	"github.com/sdcoffey/big"
	"github.com/sophon2000/quantforge/backtestengine"
	"github.com/sophon2000/quantforge/broker"
	"github.com/sophon2000/techan"
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
	lastPrice     map[string]float64
	records       map[string]*techan.TradingRecord
	SuccessPct    float64
}

// New 构造，initialCash 初始资金，commission 费率模型（如 ibkr.NewCommission(ibkr.Tiered)）
func New(initialCash float64, commission broker.CommissionModel) *DefaultSimulator {
	return &DefaultSimulator{
		initialCash: initialCash,
		cash:        initialCash,
		commission:  commission,
		lastPrice:   make(map[string]float64),
		records:     make(map[string]*techan.TradingRecord),
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
	} else {
		s.cash += cost - fee
	}

	s.applyRecordLocked(f)
	s.SuccessPct = s.successPctLocked()
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
	for sym, record := range s.records {
		pos := record.CurrentPosition()
		if pos == nil || !pos.IsOpen() {
			continue
		}
		qty := pos.RemainingAmount().Float()
		if qty <= 0 {
			continue
		}
		price := s.lastPrice[sym]
		if price <= 0 {
			price = positionAvgEntryPrice(pos)
		}
		value := price * qty
		if pos.IsShort() {
			eq -= value
		} else {
			eq += value
		}
	}
	return eq
}

// ReturnPct 收益率
func (s *DefaultSimulator) ReturnPct() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.initialCash == 0 {
		return 0
	}
	profit := s.totalProfitLocked() - s.fees
	return profit / s.initialCash * 100
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
	pos := s.currentPositionLocked(symbol)
	if pos == nil || !pos.IsOpen() {
		return 0
	}
	return int(math.Round(pos.RemainingAmount().Float()))
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

func (s *DefaultSimulator) recordForLocked(symbol string) *techan.TradingRecord {
	if s.records == nil {
		s.records = make(map[string]*techan.TradingRecord)
	}
	record := s.records[symbol]
	if record == nil {
		record = techan.NewTradingRecord()
		s.records[symbol] = record
	}
	return record
}

func (s *DefaultSimulator) currentPositionLocked(symbol string) *techan.Position {
	record := s.records[symbol]
	if record == nil {
		return nil
	}
	return record.CurrentPosition()
}

func (s *DefaultSimulator) applyRecordLocked(f backtestengine.Fill) {
	record := s.recordForLocked(f.Symbol)
	side := techan.BUY
	if f.Side == backtestengine.SELL {
		side = techan.SELL
	}
	record.Operate(techan.Order{
		Side:          side,
		Security:      f.Symbol,
		Amount:        big.NewDecimal(float64(f.Quantity)),
		Price:         big.NewDecimal(f.Price),
		ExecutionTime: f.Time,
	})
}

func (s *DefaultSimulator) successPctLocked() float64 {
	var totalTrades float64
	var profitableTrades float64
	var nta techan.NumTradesAnalysis
	pta := techan.ProfitableTradesAnalysis{}
	for _, record := range s.records {
		totalTrades += nta.Analyze(record)
		profitableTrades += pta.Analyze(record)
	}
	if totalTrades == 0 {
		return 0
	}
	return profitableTrades / totalTrades * 100
}

func (s *DefaultSimulator) totalProfitLocked() float64 {
	var total float64
	tpa := techan.TotalProfitAnalysis{}
	for sym, record := range s.records {
		total += tpa.Analyze(record)
		pos := record.CurrentPosition()
		if pos == nil || !pos.IsOpen() {
			continue
		}
		lastPrice := s.lastPrice[sym]
		if lastPrice <= 0 {
			lastPrice = positionAvgEntryPrice(pos)
		}
		total += positionProfit(pos, lastPrice)
	}
	return total
}

func positionAvgEntryPrice(pos *techan.Position) float64 {
	if pos == nil || pos.IsNew() {
		return 0
	}
	totalEntry := pos.TotalEntryAmount().Float()
	if totalEntry == 0 {
		return 0
	}
	return pos.CostBasis().Float() / totalEntry
}

func positionProfit(pos *techan.Position, lastPrice float64) float64 {
	if pos == nil || !pos.IsOpen() {
		return 0
	}
	totalEntry := pos.TotalEntryAmount().Float()
	if totalEntry == 0 {
		return 0
	}
	avgEntry := pos.CostBasis().Float() / totalEntry
	exitQty := pos.TotalExitAmount().Float()
	exitValue := pos.ExitValue().Float()
	remainingQty := pos.RemainingAmount().Float()
	if lastPrice <= 0 {
		lastPrice = avgEntry
	}
	if pos.IsLong() {
		realized := exitValue - exitQty*avgEntry
		unrealized := (lastPrice - avgEntry) * remainingQty
		return realized + unrealized
	}
	if pos.IsShort() {
		realized := exitQty*avgEntry - exitValue
		unrealized := (avgEntry - lastPrice) * remainingQty
		return realized + unrealized
	}
	return 0
}
