package dataengine

import "time"

// Tick 行情 tick（逐笔/快照）
type Tick struct {
	Symbol   string
	Price    float64
	Quantity int64
	Time     time.Time
	Bid      float64
	Ask      float64
	BidSize  int64
	AskSize  int64
}

// Bar K 线
type Bar struct {
	Symbol   string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
	Time     time.Time
	Interval string // 如 "1m", "1d"
}
