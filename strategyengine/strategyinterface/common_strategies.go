package strategyinterface

import (
	"github.com/sdcoffey/techan"
)

// BollingerBandStrategy 布林带突破策略
func BollingerBandStrategy(closePrices techan.Indicator, upper, middle, lower techan.Indicator) techan.RuleStrategy {

	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(middle, closePrices),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.NewCrossUpIndicatorRule(closePrices, middle),
		//techan.NewCrossDownIndicatorRule(closePrices, upper),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 20,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// BollingerBandMeanReversionStrategy 布林带均值回归策略
func BollingerBandMeanReversionStrategy(closePrices techan.Indicator, upper, middle, lower techan.Indicator) techan.RuleStrategy {
	entryRule := techan.And(
		techan.UnderIndicatorRule{First: closePrices, Second: lower},
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.Or(
			techan.NewCrossUpIndicatorRule(closePrices, middle),
			techan.OverIndicatorRule{First: closePrices, Second: upper},
		),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 20,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// MACDCrossoverStrategy MACD 金叉死叉策略
func MACDCrossoverStrategy(macd, signal techan.Indicator) techan.RuleStrategy {
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(macd, signal),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(macd, signal),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 26,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// MACDHistogramStrategy MACD 柱状图策略
func MACDHistogramStrategy(histogram techan.Indicator) techan.RuleStrategy {
	zero := techan.NewConstantIndicator(0)
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(histogram, zero),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(histogram, zero),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 26,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// RSIStrategy RSI 超买超卖策略
func RSIStrategy(rsi techan.Indicator, oversoldLevel, overboughtLevel float64) techan.RuleStrategy {
	oversold := techan.NewConstantIndicator(oversoldLevel)
	overbought := techan.NewConstantIndicator(overboughtLevel)
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(rsi, oversold),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(rsi, overbought),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 14,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// RSIDivergenceStrategy RSI 背离策略（简化版）
func RSIDivergenceStrategy(rsi techan.Indicator) techan.RuleStrategy {
	oversold := techan.NewConstantIndicator(30)
	overbought := techan.NewConstantIndicator(70)
	entryRule := techan.And(
		techan.And(
			techan.UnderIndicatorRule{First: rsi, Second: oversold},
			&customRisingRule{indicator: rsi},
		),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.And(
			techan.OverIndicatorRule{First: rsi, Second: overbought},
			&customFallingRule{indicator: rsi},
		),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 14,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// KDJCrossoverStrategy KDJ 金叉死叉策略
func KDJCrossoverStrategy(k, d, j techan.Indicator) techan.RuleStrategy {
	oversold := techan.NewConstantIndicator(20)
	overbought := techan.NewConstantIndicator(80)
	entryRule := techan.And(
		techan.And(
			techan.NewCrossUpIndicatorRule(k, d),
			techan.UnderIndicatorRule{First: j, Second: oversold},
		),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.And(
			techan.NewCrossDownIndicatorRule(k, d),
			techan.OverIndicatorRule{First: j, Second: overbought},
		),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 9,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// KDJOversoldOverboughtStrategy KDJ 超买超卖策略
func KDJOversoldOverboughtStrategy(j techan.Indicator) techan.RuleStrategy {
	zero := techan.NewConstantIndicator(0)
	hundred := techan.NewConstantIndicator(100)
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(j, zero),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(j, hundred),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 9,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// MultiIndicatorStrategy 多指标组合策略
func MultiIndicatorStrategy(
	closePrices techan.Indicator,
	macd, signal techan.Indicator,
	rsi techan.Indicator,
	bollLower techan.Indicator,
) techan.RuleStrategy {
	rsi50 := techan.NewConstantIndicator(50)
	rsi70 := techan.NewConstantIndicator(70)
	entryRule := techan.And(
		techan.And(
			techan.And(
				techan.NewCrossUpIndicatorRule(macd, signal),
				techan.UnderIndicatorRule{First: rsi, Second: rsi50},
			),
			techan.UnderIndicatorRule{First: closePrices, Second: bollLower},
		),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.Or(
			techan.NewCrossDownIndicatorRule(macd, signal),
			techan.OverIndicatorRule{First: rsi, Second: rsi70},
		),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 26,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// TrendFollowingStrategy 趋势跟踪策略
func TrendFollowingStrategy(
	closePrices techan.Indicator,
	macd techan.Indicator,
	rsi techan.Indicator,
	bollMiddle techan.Indicator,
) techan.RuleStrategy {
	zero := techan.NewConstantIndicator(0)
	rsi50 := techan.NewConstantIndicator(50)
	entryRule := techan.And(
		techan.And(
			techan.And(
				techan.OverIndicatorRule{First: macd, Second: zero},
				techan.OverIndicatorRule{First: rsi, Second: rsi50},
			),
			techan.OverIndicatorRule{First: closePrices, Second: bollMiddle},
		),
		techan.PositionNewRule{})
	exitRule := techan.And(
		techan.Or(
			techan.UnderIndicatorRule{First: macd, Second: zero},
			techan.NewCrossDownIndicatorRule(closePrices, bollMiddle),
		),
		techan.PositionOpenRule{})
	return techan.RuleStrategy{
		UnstablePeriod: 26,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

type customRisingRule struct {
	indicator techan.Indicator
}

func (r *customRisingRule) IsSatisfied(index int, record *techan.TradingRecord) bool {
	if index < 1 {
		return false
	}
	current := r.indicator.Calculate(index)
	previous := r.indicator.Calculate(index - 1)
	return current.GT(previous)
}

type customFallingRule struct {
	indicator techan.Indicator
}

func (r *customFallingRule) IsSatisfied(index int, record *techan.TradingRecord) bool {
	if index < 1 {
		return false
	}
	current := r.indicator.Calculate(index)
	previous := r.indicator.Calculate(index - 1)
	return current.LT(previous)
}
