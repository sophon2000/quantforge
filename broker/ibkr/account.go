package ibkr

import (
	"strconv"

	"github.com/sophon2000/quantforge/backtestengine"
	"github.com/sophon2000/quantforge/broker"

	"github.com/scmhub/ibsync"
)

// LiveAccount 实现 broker.Account，包装 Client 供实盘使用
type LiveAccount struct {
	client *Client
}

// NewLiveAccount 用已连接的 Client 构造实盘账户
func NewLiveAccount(c *Client) *LiveAccount {
	return &LiveAccount{client: c}
}

// 编译期断言：*LiveAccount 实现 broker.Account
var _ broker.Account = (*LiveAccount)(nil)

// ApplyFill 实盘由券商推送，本地可 no-op 或做记录
func (a *LiveAccount) ApplyFill(f backtestengine.Fill) {}

// Cash 从 AccountValues 取现金
func (a *LiveAccount) Cash() float64 {
	return getAccountValueFloat(a.client.AccountValues(), "CashBalance", "NetLiquidation")
}

// Equity 从 AccountValues 取净资产
func (a *LiveAccount) Equity() float64 {
	return getAccountValueFloat(a.client.AccountValues(), "NetLiquidation", "EquityWithLoanValue")
}

// Position 从 Portfolio 取指定标的持仓
func (a *LiveAccount) Position(symbol string) int {
	for _, p := range a.client.Portfolio() {
		if p.Contract.Symbol == symbol {
			return decimalToInt(p.Position)
		}
	}
	return 0
}

// UpdatePrice 实盘用行情推送，本地可 no-op
func (a *LiveAccount) UpdatePrice(symbol string, price float64) {}

func getAccountValueFloat(vals []ibsync.AccountValue, tags ...string) float64 {
	for _, tag := range tags {
		for _, v := range vals {
			if v.Tag == tag {
				f, _ := strconv.ParseFloat(v.Value, 64)
				return f
			}
		}
	}
	return 0
}

func decimalToInt(d ibsync.Decimal) int {
	// ibsync.Decimal 来自 ibapi，用 String 转 int
	s := d.String()
	f, _ := strconv.ParseFloat(s, 64)
	return int(f)
}

// AccountValues 获取账户值
func (c *Client) AccountValues() []ibsync.AccountValue {
	return c.ib.AccountValues()
}

// AccountSummary 获取账户摘要
func (c *Client) AccountSummary() []ibsync.AccountValue {
	return c.ib.AccountSummary()
}

// Portfolio 获取投资组合
func (c *Client) Portfolio() []ibsync.PortfolioItem {
	return c.ib.Portfolio()
}

// ReqPositions 订阅持仓更新
func (c *Client) ReqPositions() {
	c.ib.ReqPositions()
}

// PositionChan 获取持仓更新通道
func (c *Client) PositionChan() <-chan ibsync.Position {
	return c.ib.PositionChan()
}

// Trades 获取所有交易
func (c *Client) Trades() []*ibsync.Trade {
	return c.ib.Trades()
}

// OpenTrades 获取未完成的交易
func (c *Client) OpenTrades() []*ibsync.Trade {
	return c.ib.OpenTrades()
}
