package indicator

import (
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

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

// (20,2)
func BollingerBands(series techan.Indicator, period int, stdDev float64) (
	upper techan.Indicator,
	middle techan.Indicator,
	lower techan.Indicator,
) {
	middle = techan.NewSimpleMovingAverage(series, period)
	upper = techan.NewBollingerUpperBandIndicator(series, period, stdDev)
	lower = techan.NewBollingerLowerBandIndicator(series, period, stdDev)
	return
}

// (6,12,24)
func RSI(series techan.Indicator, period int) techan.Indicator {
	return techan.NewRelativeStrengthIndexIndicator(series, period)
}

// (9,3,3)
func KDJ(series *techan.TimeSeries, period int, smoothK int, smoothD int) (
	k techan.Indicator,
	d techan.Indicator,
	j techan.Indicator,
) {
	// RSV
	fastK := techan.NewFastStochasticIndicator(series, period)

	// 平滑K
	k = techan.NewEMAIndicator(fastK, smoothK)

	// 平滑D
	d = techan.NewEMAIndicator(k, smoothD)

	// J = 3K - 2D
	j = &FormulaIndicator{
		formula: func(i int) big.Decimal {
			kVal := k.Calculate(i)
			dVal := d.Calculate(i)
			return kVal.Mul(big.NewDecimal(3)).
				Sub(dVal.Mul(big.NewDecimal(2)))
		},
	}

	return
}

// (12, 26, 9)
func MACD(series techan.Indicator, fastPeriod int, slowPeriod int, signalPeriod int) (
	macd techan.Indicator,
	signal techan.Indicator,
	histogram techan.Indicator,
) {
	// 快线
	macd = techan.NewMACDIndicator(series, fastPeriod, slowPeriod)
	// 慢线
	signal = techan.NewEMAIndicator(macd, signalPeriod)
	// 柱状图
	histogram = techan.NewMACDHistogramIndicator(macd, signalPeriod)
	return
}
