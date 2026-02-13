package ibkr

import (
	"time"

	"github.com/scmhub/ibapi"
	"github.com/scmhub/ibsync"
)

// Snapshot 获取市场快照数据
func (c *Client) Snapshot(contract *ibapi.Contract) (*ibsync.Ticker, error) {
	return c.ib.Snapshot(contract)
}

// ReqTickByTickData 请求逐笔数据
// tickType: "Last", "AllLast", "BidAsk", "MidPoint"
func (c *Client) ReqTickByTickData(
	contract *ibapi.Contract,
	tickType string,
	numberOfTicks int64,
	ignoreSize bool,
) *ibsync.Ticker {
	return c.ib.ReqTickByTickData(contract, tickType, numberOfTicks, ignoreSize)
}

// CancelTickByTickData 取消逐笔数据订阅
func (c *Client) CancelTickByTickData(contract *ibapi.Contract, tickType string) {
	c.ib.CancelTickByTickData(contract, tickType)
}

// ReqHistoricalTicks 请求历史逐笔数据
// startDateTime: 开始时间
// endDateTime: 结束时间，如果为空则使用当前时间
// numberOfTicks: 返回的 tick 数量
// whatToShow: "TRADES", "MIDPOINT", "BID", "ASK"
// useRTH: 是否只使用常规交易时间
func (c *Client) ReqHistoricalTicks(
	contract *ibapi.Contract,
	startDateTime time.Time,
	endDateTime time.Time,
	numberOfTicks int,
	useRTH bool,
	ignoreSize bool,
) ([]ibsync.HistoricalTick, error, bool) {
	return c.ib.ReqHistoricalTicks(contract, startDateTime, endDateTime, numberOfTicks, useRTH, ignoreSize)
}

// ReqMktData 订阅市场数据
func (c *Client) ReqMktData(contract *ibapi.Contract) *ibsync.Ticker {
	return c.ib.ReqMktData(contract, "")
}

// CancelMktData 取消市场数据订阅
func (c *Client) CancelMktData(contract *ibapi.Contract) {
	c.ib.CancelMktData(contract)
}
