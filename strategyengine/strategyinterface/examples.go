package strategyinterface

import (
	"quantforge/strategyengine/indicatorlib"

	"github.com/sdcoffey/techan"
)

// BuildBollingerStrategy 构建布林带策略
func BuildBollingerStrategy(series *techan.TimeSeries, period int, stdDev float64) techan.RuleStrategy {
	closePrices := techan.NewClosePriceIndicator(series)
	upper, middle, lower := indicatorlib.BollingerBands(closePrices, period, stdDev)
	return BollingerBandStrategy(closePrices, upper, middle, lower)
}

// BuildBollingerMeanReversionStrategy 构建布林带均值回归策略
func BuildBollingerMeanReversionStrategy(series *techan.TimeSeries, period int, stdDev float64) techan.RuleStrategy {
	closePrices := techan.NewClosePriceIndicator(series)
	upper, middle, lower := indicatorlib.BollingerBands(closePrices, period, stdDev)
	return BollingerBandMeanReversionStrategy(closePrices, upper, middle, lower)
}

// BuildMACDStrategy 构建 MACD 金叉死叉策略
func BuildMACDStrategy(series *techan.TimeSeries, fastPeriod, slowPeriod, signalPeriod int) techan.RuleStrategy {
	closePrices := techan.NewClosePriceIndicator(series)
	macd, signal, _ := indicatorlib.MACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	return MACDCrossoverStrategy(macd, signal)
}

// BuildMACDHistogramStrategy 构建 MACD 柱状图策略
func BuildMACDHistogramStrategy(series *techan.TimeSeries, fastPeriod, slowPeriod, signalPeriod int) techan.RuleStrategy {
	closePrices := techan.NewClosePriceIndicator(series)
	_, _, histogram := indicatorlib.MACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	return MACDHistogramStrategy(histogram)
}

// BuildRSIStrategy 构建 RSI 策略
func BuildRSIStrategy(series *techan.TimeSeries, period int, oversold, overbought float64) techan.RuleStrategy {
	closePrices := techan.NewClosePriceIndicator(series)
	rsi := indicatorlib.RSI(closePrices, period)
	return RSIStrategy(rsi, oversold, overbought)
}

// BuildRSIDivergenceStrategy 构建 RSI 背离策略
func BuildRSIDivergenceStrategy(series *techan.TimeSeries, period int) techan.RuleStrategy {
	closePrices := techan.NewClosePriceIndicator(series)
	rsi := indicatorlib.RSI(closePrices, period)
	return RSIDivergenceStrategy(rsi)
}

// BuildKDJCrossoverStrategy 构建 KDJ 金叉死叉策略
func BuildKDJCrossoverStrategy(series *techan.TimeSeries, period, smoothK, smoothD int) techan.RuleStrategy {
	k, d, j := indicatorlib.KDJ(series, period, smoothK, smoothD)
	return KDJCrossoverStrategy(k, d, j)
}

// BuildKDJOversoldOverboughtStrategy 构建 KDJ 超买超卖策略
func BuildKDJOversoldOverboughtStrategy(series *techan.TimeSeries, period, smoothK, smoothD int) techan.RuleStrategy {
	_, _, j := indicatorlib.KDJ(series, period, smoothK, smoothD)
	return KDJOversoldOverboughtStrategy(j)
}

// BuildMultiIndicatorStrategy 构建多指标组合策略
func BuildMultiIndicatorStrategy(series *techan.TimeSeries) techan.RuleStrategy {
	closePrices := techan.NewClosePriceIndicator(series)
	macd, signal, _ := indicatorlib.MACD(closePrices, 12, 26, 9)
	rsi := indicatorlib.RSI(closePrices, 14)
	_, _, lower := indicatorlib.BollingerBands(closePrices, 20, 2.0)
	return MultiIndicatorStrategy(closePrices, macd, signal, rsi, lower)
}

// BuildTrendFollowingStrategy 构建趋势跟踪策略
func BuildTrendFollowingStrategy(series *techan.TimeSeries) techan.RuleStrategy {
	closePrices := techan.NewClosePriceIndicator(series)
	macd, _, _ := indicatorlib.MACD(closePrices, 12, 26, 9)
	rsi := indicatorlib.RSI(closePrices, 14)
	_, middle, _ := indicatorlib.BollingerBands(closePrices, 20, 2.0)
	return TrendFollowingStrategy(closePrices, macd, rsi, middle)
}

// DefaultStrategies 返回常用默认策略集合
func DefaultStrategies(series *techan.TimeSeries) map[string]techan.RuleStrategy {
	return map[string]techan.RuleStrategy{
		"BOLL_Breakout":          BuildBollingerStrategy(series, 20, 2.0),
		"BOLL_MeanReversion":     BuildBollingerMeanReversionStrategy(series, 20, 2.0),
		"MACD_Crossover":         BuildMACDStrategy(series, 12, 26, 9),
		"MACD_Histogram":         BuildMACDHistogramStrategy(series, 12, 26, 9),
		"RSI_OversoldOverbought": BuildRSIStrategy(series, 14, 30, 70),
		"RSI_Divergence":         BuildRSIDivergenceStrategy(series, 14),
		"KDJ_Crossover":          BuildKDJCrossoverStrategy(series, 9, 3, 3),
		"KDJ_Extreme":            BuildKDJOversoldOverboughtStrategy(series, 9, 3, 3),
		"MultiIndicator":         BuildMultiIndicatorStrategy(series),
		"TrendFollowing":         BuildTrendFollowingStrategy(series),
	}
}
