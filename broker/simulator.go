package broker

import (
	"github.com/sophon2000/quantforge/backtestengine"
)

// Simulator 回测账户模拟：资金与持仓
type Simulator interface {
	ApplyFill(f backtestengine.Fill)
	Cash() float64
	Equity() float64
	ReturnPct() float64
	Fees() float64
	Position(symbol string) int
	UpdatePrice(symbol string, price float64)
}
