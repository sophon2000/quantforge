package ibkr

import (
	"github.com/scmhub/ibsync"
)

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
