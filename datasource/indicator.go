package datasource

import (
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

// IndicatorSet 常用指标集合（techan）
type IndicatorSet struct {
	MA   techan.Indicator
	EMA  techan.Indicator
	BOLL struct {
		Upper  techan.Indicator
		Middle techan.Indicator
		Lower  techan.Indicator
	}
	RSI  techan.Indicator
	MACD struct {
		MACD      techan.Indicator
		Signal    techan.Indicator
		Histogram techan.Indicator
	}
	KDJ struct {
		K techan.Indicator
		D techan.Indicator
		J techan.Indicator
	}
	Volume struct {
		Volume   techan.Indicator
		VolumeMA techan.Indicator
	}
}

// SimpleMovingAverage 简单移动平均
func SimpleMovingAverage(series techan.Indicator, period int) techan.Indicator {
	return techan.NewSimpleMovingAverage(series, period)
}

// ExponentialMovingAverage 指数移动平均
func ExponentialMovingAverage(series techan.Indicator, period int) techan.Indicator {
	return techan.NewEMAIndicator(series, period)
}

// BollingerBands 布林带 (period, stdDev)，如 (20, 2)
func BollingerBands(series techan.Indicator, period int, stdDev float64) (
	upper, middle, lower techan.Indicator,
) {
	middle = techan.NewSimpleMovingAverage(series, period)
	upper = techan.NewBollingerUpperBandIndicator(series, period, stdDev)
	lower = techan.NewBollingerLowerBandIndicator(series, period, stdDev)
	return
}

// RSI 相对强弱指数，如 period=14
func RSI(series techan.Indicator, period int) techan.Indicator {
	return techan.NewRelativeStrengthIndexIndicator(series, period)
}

// KDJ 随机指标 (period, smoothK, smoothD)，如 (9, 3, 3)
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

// MACD (fastPeriod, slowPeriod, signalPeriod)，如 (12, 26, 9)
func MACD(series techan.Indicator, fastPeriod int, slowPeriod int, signalPeriod int) (
	macd, signal, histogram techan.Indicator,
) {
	macd = techan.NewMACDIndicator(series, fastPeriod, slowPeriod)
	signal = techan.NewEMAIndicator(macd, signalPeriod)
	histogram = techan.NewMACDHistogramIndicator(macd, signalPeriod)
	return
}
