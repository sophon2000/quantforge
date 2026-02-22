package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"quantforge/backtestengine"
	"quantforge/backtestengine/accountsimulator"
	"quantforge/dataengine"
	"quantforge/dataengine/historicalstore"
	"quantforge/strategyengine"
	"quantforge/strategyengine/strategyinterface"

	"github.com/sdcoffey/techan"
)

// BacktestResult 回测结果，供 API 与前端图表使用
type BacktestResult struct {
	Symbol       string        `json:"symbol"`
	Strategy     string        `json:"strategy"`
	CategoryData []string      `json:"categoryData"` // 日期轴，如 2014/1/2
	Values       [][4]float64  `json:"values"`       // 每根 K [open, close, low, high]
	Signals      []SignalPoint `json:"signals"`      // 买卖信号
	MA5          []interface{} `json:"ma5"`          // MA5，不足为 "-"
	MA10         []interface{} `json:"ma10"`
	MA20         []interface{} `json:"ma20"`
	MA30         []interface{} `json:"ma30"`
	BollUpper    []interface{} `json:"bollUpper"`  // 布林上轨
	BollMiddle   []interface{} `json:"bollMiddle"` // 布林中轨
	BollLower    []interface{} `json:"bollLower"`  // 布林下轨
	Volumes      []int64       `json:"volumes"`    // 成交量，与 categoryData 一一对应
	MacdDIF      []interface{} `json:"macdDIF"`    // MACD DIF (12,26,9)
	MacdDEA      []interface{} `json:"macdDEA"`    // MACD DEA
	MacdHist     []interface{} `json:"macdHist"`   // MACD 柱 (DIF-DEA)*2
	Summary      Summary       `json:"summary"`
}

// SignalPoint 单笔信号；SELL 时 FillReturnPct 为该笔交易收益率(%)
type SignalPoint struct {
	Index         int     `json:"index"`
	Date          string  `json:"date"`
	Type          string  `json:"type"` // BUY / SELL
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`            // 该笔成交数量
	FillReturnPct float64 `json:"returnPct,omitempty"` // 仅 SELL 有值：(卖价-买价)/买价*100
}

// Summary 回测摘要
type Summary struct {
	InitialCash float64 `json:"initialCash"`
	FinalCash   float64 `json:"finalCash"`
	TradeCount  int     `json:"tradeCount"`
	Position    int     `json:"position"`
	ReturnPct   float64 `json:"returnPct"`
	SuccessPct  float64 `json:"successPct"`
}

// RunBacktest 执行回测并返回结构化结果（供 CLI 打印与 HTTP API 使用）
func RunBacktest(symbol, strategyName string, initialCash float64, quantity int) (*BacktestResult, error) {
	store := historicalstore.NewCSVStore("")
	rows, err := store.LoadCSV("")
	if err != nil {
		return nil, fmt.Errorf("加载数据: %w", err)
	}
	sliceRows, ok := rows[symbol]
	if !ok {
		return nil, fmt.Errorf("未找到标的: %s", symbol)
	}

	bars, err := csvRowsToBars(symbol, sliceRows)
	if err != nil {
		return nil, fmt.Errorf("解析 K 线: %w", err)
	}
	if len(bars) == 0 {
		return nil, fmt.Errorf("标的 %s 无有效 K 线", symbol)
	}

	account := accountsimulator.NewDefaultAccountSimulator(initialCash)
	var lastClose float64
	var currentIndex int
	var currentDate string
	var lastBuyPrice float64
	var signals []SignalPoint

	onSignal := func(s strategyengine.Signal) {
		price := lastClose
		if price <= 0 {
			return
		}
		var fillQty int
		if quantity > 0 {
			fillQty = quantity
		} else {
			// quantity==0：买入按现金允许的最大数量，卖出按该 symbol 当前持仓全部
			if s.Signal == "BUY" {
				cash := account.Equity()
				fillQty = int(cash / price)
				fmt.Println("fillQty", fillQty, "cash", cash, "price", price)
			} else {
				fillQty = account.Position(s.Symbol)
			}
		}
		if fillQty <= 0 {
			return
		}
		fill := backtestengine.Fill{
			Symbol:   s.Symbol,
			Price:    price,
			Quantity: fillQty,
			Side:     s.Signal,
		}
		account.ApplyFill(fill)
		pt := SignalPoint{Index: currentIndex, Date: currentDate, Type: s.Signal, Price: price, Quantity: fillQty}
		if s.Signal == "BUY" {
			lastBuyPrice = price
		} else if s.Signal == "SELL" && lastBuyPrice > 0 {
			pt.FillReturnPct = (price - lastBuyPrice) / lastBuyPrice * 100
		}
		signals = append(signals, pt)
	}

	ruleBuilder := pickRuleBuilder(strategyName)

	strat := strategyinterface.NewTechanStrategy(symbol, ruleBuilder, onSignal)

	categoryData := make([]string, 0, len(bars))
	values := make([][4]float64, 0, len(bars))
	volumes := make([]int64, 0, len(bars))

	for i := range bars {
		lastClose = bars[i].Close
		currentIndex = i
		currentDate = bars[i].Time.Format("2006/1/2")
		strat.OnBar(bars[i])

		categoryData = append(categoryData, currentDate)
		values = append(values, [4]float64{
			bars[i].Open,
			bars[i].Close,
			bars[i].Low,
			bars[i].High,
		})
		volumes = append(volumes, bars[i].Volume)
	}

	equity := account.Equity()
	pos := account.Position(symbol)
	returnPct := 0.0
	if initialCash != 0 {
		returnPct = (equity - initialCash) / initialCash * 100
	}

	ma5 := calculateMA(values, 5)
	ma10 := calculateMA(values, 10)
	ma20 := calculateMA(values, 20)
	ma30 := calculateMA(values, 30)
	bollUpper, bollMiddle, bollLower := calculateBollinger(values, 20, 2.0)
	macdDIF, macdDEA, macdHist := calculateMACD(values, 12, 26, 9)

	return &BacktestResult{
		Symbol:       symbol,
		Strategy:     strategyName,
		CategoryData: categoryData,
		Values:       values,
		Signals:      signals,
		MA5:          ma5,
		MA10:         ma10,
		MA20:         ma20,
		MA30:         ma30,
		BollUpper:    bollUpper,
		BollMiddle:   bollMiddle,
		BollLower:    bollLower,
		Volumes:      volumes,
		MacdDIF:      macdDIF,
		MacdDEA:      macdDEA,
		MacdHist:     macdHist,
		Summary: Summary{
			InitialCash: initialCash,
			FinalCash:   equity,
			TradeCount:  len(signals),
			Position:    pos,
			ReturnPct:   returnPct,
			SuccessPct:  account.SuccessPct,
		},
	}, nil
}

// calculateMA 计算收盘价均线，不足 period 为 "-"
func calculateMA(values [][4]float64, period int) []interface{} {
	result := make([]interface{}, len(values))
	for i := range values {
		if i < period-1 {
			result[i] = "-"
			continue
		}
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += values[i-j][1] // close
		}
		result[i] = sum / float64(period)
	}
	return result
}

// calculateBollinger 布林带：中轨=SMA(close,period)，上/下=中轨±stdDev*标准差，不足 period 为 "-"
func calculateBollinger(values [][4]float64, period int, stdDev float64) (upper, middle, lower []interface{}) {
	n := len(values)
	upper = make([]interface{}, n)
	middle = make([]interface{}, n)
	lower = make([]interface{}, n)
	for i := 0; i < n; i++ {
		if i < period-1 {
			upper[i], middle[i], lower[i] = "-", "-", "-"
			continue
		}
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += values[i-j][1]
		}
		mid := sum / float64(period)
		var variance float64
		for j := 0; j < period; j++ {
			diff := values[i-j][1] - mid
			variance += diff * diff
		}
		std := math.Sqrt(variance / float64(period))
		middle[i] = mid
		upper[i] = mid + stdDev*std
		lower[i] = mid - stdDev*std
	}
	return
}

// calculateEMA 计算收盘价的 EMA，不足 period 为 "-"
func calculateEMA(values [][4]float64, period int) []interface{} {
	n := len(values)
	result := make([]interface{}, n)
	if n < period {
		for i := 0; i < n; i++ {
			result[i] = "-"
		}
		return result
	}
	alpha := 2.0 / float64(period+1)
	// 第一个 EMA 用前 period 根 K 的收盘价均值
	sum := 0.0
	for j := 0; j < period; j++ {
		sum += values[j][1]
	}
	ema := sum / float64(period)
	for i := 0; i < period-1; i++ {
		result[i] = "-"
	}
	result[period-1] = ema
	for i := period; i < n; i++ {
		ema = alpha*values[i][1] + (1-alpha)*ema
		result[i] = ema
	}
	return result
}

// calculateMACD 计算 MACD(12,26,9)：DIF、DEA、柱 (DIF-DEA)*2，不足为 "-"
func calculateMACD(values [][4]float64, fast, slow, signal int) (dif, dea, hist []interface{}) {
	n := len(values)
	dif = make([]interface{}, n)
	dea = make([]interface{}, n)
	hist = make([]interface{}, n)
	for i := 0; i < n; i++ {
		dif[i], dea[i], hist[i] = "-", "-", "-"
	}
	if n < slow {
		return
	}
	emaFast := calculateEMA(values, fast)
	emaSlow := calculateEMA(values, slow)
	// DIF = EMA_fast - EMA_slow，从 index slow-1 开始有效
	difSlice := make([]float64, n)
	for i := slow - 1; i < n; i++ {
		f, _ := emaFast[i].(float64)
		s, _ := emaSlow[i].(float64)
		difSlice[i] = f - s
		dif[i] = f - s
	}
	// DEA = EMA(DIF, signal)，从 slow-1+signal-1 开始有效
	alpha := 2.0 / float64(signal+1)
	startDEA := slow - 1 + signal - 1
	if startDEA >= n {
		return
	}
	sum := 0.0
	for j := slow - 1; j < slow-1+signal; j++ {
		sum += difSlice[j]
	}
	emaDIF := sum / float64(signal)
	dea[startDEA] = emaDIF
	hist[startDEA] = (difSlice[startDEA] - emaDIF) * 2
	for i := startDEA + 1; i < n; i++ {
		emaDIF = alpha*difSlice[i] + (1-alpha)*emaDIF
		dea[i] = emaDIF
		hist[i] = (difSlice[i] - emaDIF) * 2
	}
	return
}

func csvRowsToBars(symbol string, rows []historicalstore.CSVRow) ([]*dataengine.Bar, error) {
	out := make([]*dataengine.Bar, 0, len(rows))
	for _, r := range rows {
		t, err := time.ParseInLocation(time.DateOnly, r.Date, time.UTC)
		if err != nil {
			continue
		}
		open, _ := strconv.ParseFloat(r.Open, 64)
		high, _ := strconv.ParseFloat(r.High, 64)
		low, _ := strconv.ParseFloat(r.Low, 64)
		close_, _ := strconv.ParseFloat(r.Close, 64)
		vol, _ := strconv.ParseInt(r.Volume, 10, 64)
		out = append(out, &dataengine.Bar{
			Symbol:   symbol,
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close_,
			Volume:   vol,
			Time:     t,
			Interval: "1d",
		})
	}
	return out, nil
}

// pickRuleBuilder 根据策略名返回 TechanStrategy 的规则构建函数，与 strategyinterface/examples 对应
func pickRuleBuilder(name string) strategyinterface.RuleBuilder {
	switch name {
	// 布林带
	case "bollinger", "boll-breakout":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildBollingerStrategy(series, 20, 2.0)
		}
	case "bollinger-mean-reversion", "boll-mean-reversion":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildBollingerMeanReversionStrategy(series, 20, 2.0)
		}
	// MACD
	case "macd", "macd-crossover":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildMACDStrategy(series, 12, 26, 9)
		}
	case "macd-histogram":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildMACDHistogramStrategy(series, 12, 26, 9)
		}
	// RSI
	case "rsi", "rsi-oversold-overbought":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildRSIStrategy(series, 14, 30, 70)
		}
	case "rsi-divergence":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildRSIDivergenceStrategy(series, 14)
		}
	// KDJ
	case "kdj-crossover":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildKDJCrossoverStrategy(series, 9, 3, 3)
		}
	case "kdj-extreme":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildKDJOversoldOverboughtStrategy(series, 9, 3, 3)
		}
	// 组合
	case "multi-indicator":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildMultiIndicatorStrategy(series)
		}
	case "trend-following":
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildTrendFollowingStrategy(series)
		}
	default:
		return func(series *techan.TimeSeries) techan.RuleStrategy {
			return strategyinterface.BuildBollingerStrategy(series, 20, 2.0)
		}
	}
}
