package ibkr

import (
	"github.com/scmhub/ibsync"
)

// ReqPnL 订阅账户盈亏
func (c *Client) ReqPnL(account string, modelCode string) {
	c.ib.ReqPnL(account, modelCode)
}

// PnlChan 获取盈亏更新通道
func (c *Client) PnlChan(account string, modelCode string) <-chan ibsync.Pnl {
	return c.ib.PnlChan(account, modelCode)
}

// Pnl 获取当前盈亏
func (c *Client) Pnl(account string, modelCode string) []ibsync.Pnl {
	return c.ib.Pnl(account, modelCode)
}

// ReqPnLSingle 订阅单个持仓的盈亏
func (c *Client) ReqPnLSingle(account string, modelCode string, conId int64) {
	c.ib.ReqPnLSingle(account, modelCode, conId)
}

// PnlSingleChan 获取单个持仓的盈亏更新通道
func (c *Client) PnlSingleChan(account string, modelCode string, conId int64) <-chan ibsync.PnlSingle {
	return c.ib.PnlSingleChan(account, modelCode, conId)
}

// PnlSingle 获取单个持仓的当前盈亏
func (c *Client) PnlSingle(account string, modelCode string, conId int64) []ibsync.PnlSingle {
	return c.ib.PnlSingle(account, modelCode, conId)
}