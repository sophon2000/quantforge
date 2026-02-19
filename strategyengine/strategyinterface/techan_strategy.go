package strategyinterface

import (
	"fmt"
	"sync"
	"time"

	"quantforge/dataengine"
	"quantforge/executionengine"
	"quantforge/strategyengine"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

// TechanStrategy 基于 techan.RuleStrategy 实现 Strategy：用 K 线驱动规则，触发时通过 OnSignal 发出信号
type TechanStrategy struct {
	mu       sync.Mutex
	symbol   string
	series   *techan.TimeSeries
	rule     techan.RuleStrategy
	record   *techan.TradingRecord
	barIndex int
	onSignal func(strategyengine.Signal)
}

// RuleBuilder 根据 TimeSeries 构建规则（与策略共用同一 series，Bar 会追加到该 series）
type RuleBuilder func(series *techan.TimeSeries) techan.RuleStrategy

// NewTechanStrategy 构造。ruleBuilder 根据内部 series 构建规则，例如：
//
//	NewTechanStrategy("AAPL", func(series *techan.TimeSeries) techan.RuleStrategy {
//	    return BuildBollingerStrategy(series, 20, 2.0)
//	}, onSignal)
//
// onSignal 在触发买卖时调用。
func NewTechanStrategy(symbol string, ruleBuilder RuleBuilder, onSignal func(strategyengine.Signal)) *TechanStrategy {
	series := techan.NewTimeSeries()
	rule := ruleBuilder(series)
	return &TechanStrategy{
		symbol:   symbol,
		series:   series,
		rule:     rule,
		record:   techan.NewTradingRecord(),
		onSignal: onSignal,
	}
}

// OnTick 本实现以 K 线为准，tick 不参与规则计算（可扩展为用 tick 聚合为 bar）
func (s *TechanStrategy) OnTick(t *dataengine.Tick) {
	_ = t
}

// OnBar 追加一根 K 线并执行规则：若满足入场/出场则回调 onSignal 并更新内部 TradingRecord
func (s *TechanStrategy) OnBar(b *dataengine.Bar) {
	if b == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	candle := barToCandle(b)
	s.series.AddCandle(candle)
	index := s.barIndex
	s.barIndex++

	// 同一根 K 线只执行其一：若本 bar 已入场则本 bar 不再出场，避免同价买卖
	entered := false
	if s.rule.ShouldEnter(index, s.record) {
		if s.onSignal != nil {
			s.onSignal(strategyengine.Signal{Symbol: s.symbol, Signal: "BUY"})
		}
		s.record.Operate(techan.Order{
			Side:          techan.BUY,
			Security:      s.symbol,
			Amount:        big.ONE,
			Price:         big.NewFromString(fmt.Sprintf("%.4f", b.Close)),
			ExecutionTime: b.Time,
		})
		entered = true
	}
	if !entered && s.rule.ShouldExit(index, s.record) {
		if s.onSignal != nil {
			s.onSignal(strategyengine.Signal{Symbol: s.symbol, Signal: "SELL"})
		}
		s.record.Operate(techan.Order{
			Side:          techan.SELL,
			Security:      s.symbol,
			Amount:        big.ONE,
			Price:         big.NewFromString(fmt.Sprintf("%.4f", b.Close)),
			ExecutionTime: b.Time,
		})
	}
}

// OnOrderUpdate 可选：接收实盘/模拟订单状态回放（本实现仅占位）
func (s *TechanStrategy) OnOrderUpdate(order *executionengine.Order) {
	_ = order
}

// Series 返回内部 TimeSeries（供需要基于同一 series 构建新规则的场景）
func (s *TechanStrategy) Series() *techan.TimeSeries {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.series
}

// Record 返回内部 TradingRecord
func (s *TechanStrategy) Record() *techan.TradingRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.record
}

// 确保 TechanStrategy 实现 Strategy
var _ Strategy = (*TechanStrategy)(nil)

func barToCandle(b *dataengine.Bar) *techan.Candle {
	period := techan.NewTimePeriod(b.Time, time.Hour*24)
	candle := techan.NewCandle(period)
	candle.OpenPrice = big.NewFromString(fmt.Sprintf("%.4f", b.Open))
	candle.ClosePrice = big.NewFromString(fmt.Sprintf("%.4f", b.Close))
	candle.MaxPrice = big.NewFromString(fmt.Sprintf("%.4f", b.High))
	candle.MinPrice = big.NewFromString(fmt.Sprintf("%.4f", b.Low))
	candle.Volume = big.NewFromString(fmt.Sprintf("%d", b.Volume))
	return candle
}
