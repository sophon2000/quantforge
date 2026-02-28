package strategyinterface

import (
	"fmt"
	"quantforge/dataengine"
	"quantforge/strategyengine"
	"quantforge/strategyengine/signalgenerator"
	"time"

	"github.com/sdcoffey/techan"
)

func ExampleNewTechanStrategy() {
	var onSignal = func(s strategyengine.Signal) {
		fmt.Printf("信号: %s %s\n", s.Symbol, s.Signal)
	}
	signalEngine := signalgenerator.NewDefaultSignalEngine(onSignal)

	// 布林带策略
	strat := NewTechanStrategy("AAPL", func(series *techan.TimeSeries) techan.RuleStrategy {
		return BuildBollingerStrategy(series, 20, 2.0)
	}, signalEngine)

	// 每根 K 线推送
	strat.OnBar(&dataengine.Bar{Symbol: "AAPL", Open: 100, High: 102, Low: 99, Close: 101, Volume: 1000, Time: time.Now()})

}
