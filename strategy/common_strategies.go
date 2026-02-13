package strategy

import (
	"github.com/sdcoffey/techan"
)

// BollingerBandStrategy 布林带突破策略
// 买入：价格从下轨下方上穿下轨
// 卖出：价格从上轨上方下穿上轨
func BollingerBandStrategy(closePrices techan.Indicator, upper, middle, lower techan.Indicator) techan.RuleStrategy {
	// 买入：价格上穿下轨（超卖反弹）
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(closePrices, lower),
		techan.PositionNewRule{},
	)

	// 卖出：价格下穿上轨（超买回落）
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(closePrices, upper),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 20, // 布林带周期
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// BollingerBandMeanReversionStrategy 布林带均值回归策略
// 买入：价格触及下轨，预期回归中轨
// 卖出：价格触及上轨或回到中轨
func BollingerBandMeanReversionStrategy(closePrices techan.Indicator, upper, middle, lower techan.Indicator) techan.RuleStrategy {
	// 买入：价格低于下轨
	entryRule := techan.And(
		techan.UnderIndicatorRule{First: closePrices, Second: lower},
		techan.PositionNewRule{},
	)

	// 卖出：价格回到中轨或超过上轨
	exitRule := techan.And(
		techan.Or(
			techan.NewCrossUpIndicatorRule(closePrices, middle),
			techan.OverIndicatorRule{First: closePrices, Second: upper},
		),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 20,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// MACDCrossoverStrategy MACD 金叉死叉策略
// 买入：MACD 线上穿信号线（金叉）
// 卖出：MACD 线下穿信号线（死叉）
func MACDCrossoverStrategy(macd, signal techan.Indicator) techan.RuleStrategy {
	// 买入：金叉
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(macd, signal),
		techan.PositionNewRule{},
	)

	// 卖出：死叉
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(macd, signal),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 26, // 慢线周期
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// MACDHistogramStrategy MACD 柱状图策略
// 买入：柱状图由负转正
// 卖出：柱状图由正转负
func MACDHistogramStrategy(histogram techan.Indicator) techan.RuleStrategy {
	zero := techan.NewConstantIndicator(0)

	// 买入：柱状图上穿零轴
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(histogram, zero),
		techan.PositionNewRule{},
	)

	// 卖出：柱状图下穿零轴
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(histogram, zero),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 26,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// RSIStrategy RSI 超买超卖策略
// 买入：RSI < 30（超卖）
// 卖出：RSI > 70（超买）
func RSIStrategy(rsi techan.Indicator, oversoldLevel, overboughtLevel float64) techan.RuleStrategy {
	oversold := techan.NewConstantIndicator(oversoldLevel)
	overbought := techan.NewConstantIndicator(overboughtLevel)

	// 买入：RSI 上穿超卖线
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(rsi, oversold),
		techan.PositionNewRule{},
	)

	// 卖出：RSI 下穿超买线
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(rsi, overbought),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 14, // RSI 周期
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// RSIDivergenceStrategy RSI 背离策略（简化版）
// 买入：RSI 在超卖区且开始上升
// 卖出：RSI 在超买区且开始下降
func RSIDivergenceStrategy(rsi techan.Indicator) techan.RuleStrategy {
	oversold := techan.NewConstantIndicator(30)
	overbought := techan.NewConstantIndicator(70)

	// 买入：RSI < 30 且 RSI 上升
	entryRule := techan.And(
		techan.And(
			techan.UnderIndicatorRule{First: rsi, Second: oversold},
			// RSI 开始上升（当前值 > 前一个值）
			&customRisingRule{indicator: rsi},
		),
		techan.PositionNewRule{},
	)

	// 卖出：RSI > 70 且 RSI 下降
	exitRule := techan.And(
		techan.And(
			techan.OverIndicatorRule{First: rsi, Second: overbought},
			// RSI 开始下降
			&customFallingRule{indicator: rsi},
		),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 14,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// KDJCrossoverStrategy KDJ 金叉死叉策略
// 买入：K 线上穿 D 线，且 J < 20
// 卖出：K 线下穿 D 线，且 J > 80
func KDJCrossoverStrategy(k, d, j techan.Indicator) techan.RuleStrategy {
	oversold := techan.NewConstantIndicator(20)
	overbought := techan.NewConstantIndicator(80)

	// 买入：K 上穿 D，且 J 在超卖区
	entryRule := techan.And(
		techan.And(
			techan.NewCrossUpIndicatorRule(k, d),
			techan.UnderIndicatorRule{First: j, Second: oversold},
		),
		techan.PositionNewRule{},
	)

	// 卖出：K 下穿 D，且 J 在超买区
	exitRule := techan.And(
		techan.And(
			techan.NewCrossDownIndicatorRule(k, d),
			techan.OverIndicatorRule{First: j, Second: overbought},
		),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 9, // KDJ 周期
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// KDJOversoldOverboughtStrategy KDJ 超买超卖策略
// 买入：J < 0
// 卖出：J > 100
func KDJOversoldOverboughtStrategy(j techan.Indicator) techan.RuleStrategy {
	zero := techan.NewConstantIndicator(0)
	hundred := techan.NewConstantIndicator(100)

	// 买入：J 上穿 0
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(j, zero),
		techan.PositionNewRule{},
	)

	// 卖出：J 下穿 100
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(j, hundred),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 9,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// MultiIndicatorStrategy 多指标组合策略
// 买入：MACD 金叉 + RSI < 50 + 价格在布林带下轨附近
// 卖出：MACD 死叉或 RSI > 70
func MultiIndicatorStrategy(
	closePrices techan.Indicator,
	macd, signal techan.Indicator,
	rsi techan.Indicator,
	bollLower techan.Indicator,
) techan.RuleStrategy {
	rsi50 := techan.NewConstantIndicator(50)
	rsi70 := techan.NewConstantIndicator(70)

	// 买入：多个信号共振
	entryRule := techan.And(
		techan.And(
			techan.And(
				techan.NewCrossUpIndicatorRule(macd, signal),                      // MACD 金叉
				techan.UnderIndicatorRule{First: rsi, Second: rsi50},              // RSI 不过热
			),
			techan.UnderIndicatorRule{First: closePrices, Second: bollLower}, // 价格接近下轨（超卖）
		),
		techan.PositionNewRule{},
	)

	// 卖出：任一卖出信号
	exitRule := techan.And(
		techan.Or(
			techan.NewCrossDownIndicatorRule(macd, signal),            // MACD 死叉
			techan.OverIndicatorRule{First: rsi, Second: rsi70},       // RSI 超买
		),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 26,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// TrendFollowingStrategy 趋势跟踪策略
// 买入：MACD > 0 + RSI > 50 + 价格在布林带中轨上方
// 卖出：MACD < 0 或价格跌破布林带中轨
func TrendFollowingStrategy(
	closePrices techan.Indicator,
	macd techan.Indicator,
	rsi techan.Indicator,
	bollMiddle techan.Indicator,
) techan.RuleStrategy {
	zero := techan.NewConstantIndicator(0)
	rsi50 := techan.NewConstantIndicator(50)

	// 买入：多头趋势确认
	entryRule := techan.And(
		techan.And(
			techan.And(
				techan.OverIndicatorRule{First: macd, Second: zero},             // MACD 多头
				techan.OverIndicatorRule{First: rsi, Second: rsi50},             // RSI 强势
			),
			techan.OverIndicatorRule{First: closePrices, Second: bollMiddle}, // 价格在中轨上方
		),
		techan.PositionNewRule{},
	)

	// 卖出：趋势转弱
	exitRule := techan.And(
		techan.Or(
			techan.UnderIndicatorRule{First: macd, Second: zero},         // MACD 转空
			techan.NewCrossDownIndicatorRule(closePrices, bollMiddle), // 跌破中轨
		),
		techan.PositionOpenRule{},
	)

	return techan.RuleStrategy{
		UnstablePeriod: 26,
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}
}

// customRisingRule 自定义上升规则
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

// customFallingRule 自定义下降规则
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
