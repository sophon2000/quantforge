package indicatorlib

import (
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

// FormulaIndicator 公式指标（如 KDJ J）
type FormulaIndicator struct {
	Formula func(int) big.Decimal
}

// Calculate 实现 techan.Indicator
func (f *FormulaIndicator) Calculate(index int) big.Decimal {
	return f.Formula(index)
}

// BollingerBands 布林带
func BollingerBands(series techan.Indicator, period int, stdDev float64) (
	upper, middle, lower techan.Indicator,
) {
	middle = techan.NewSimpleMovingAverage(series, period)
	upper = techan.NewBollingerUpperBandIndicator(series, period, stdDev)
	lower = techan.NewBollingerLowerBandIndicator(series, period, stdDev)
	return
}

// RSI 相对强弱指数
func RSI(series techan.Indicator, period int) techan.Indicator {
	return techan.NewRelativeStrengthIndexIndicator(series, period)
}

// KDJ 随机指标
func KDJ(series *techan.TimeSeries, period int, smoothK int, smoothD int) (
	k, d, j techan.Indicator,
) {
	fastK := techan.NewFastStochasticIndicator(series, period)
	k = techan.NewEMAIndicator(fastK, smoothK)
	d = techan.NewEMAIndicator(k, smoothD)
	j = &FormulaIndicator{
		Formula: func(i int) big.Decimal {
			kVal := k.Calculate(i)
			dVal := d.Calculate(i)
			return kVal.Mul(big.NewDecimal(3)).Sub(dVal.Mul(big.NewDecimal(2)))
		},
	}
	return
}

// MACD
func MACD(series techan.Indicator, fastPeriod int, slowPeriod int, signalPeriod int) (
	macd, signal, histogram techan.Indicator,
) {
	macd = techan.NewMACDIndicator(series, fastPeriod, slowPeriod)
	signal = techan.NewEMAIndicator(macd, signalPeriod)
	histogram = techan.NewMACDHistogramIndicator(macd, signalPeriod)
	return
}

// SMA 简单移动平均
func SMA(series techan.Indicator, period int) techan.Indicator {
	return techan.NewSimpleMovingAverage(series, period)
}

// EMA 指数移动平均
func EMA(series techan.Indicator, period int) techan.Indicator {
	return techan.NewEMAIndicator(series, period)
}
