package broker

import (
	"github.com/sophon2000/quantforge/backtestengine"
)

// Account 账户抽象：资金与持仓，供回测/实盘统一使用
type Account interface {
	ApplyFill(f backtestengine.Fill)
	Cash() float64
	Equity() float64
	Position(symbol string) int
	UpdatePrice(symbol string, price float64)
}
