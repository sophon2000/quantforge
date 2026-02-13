package ibkr

import (
	"github.com/scmhub/ibapi"
	"github.com/scmhub/ibsync"
)

// ReqContractDetails 请求合约详情
func (c *Client) ReqContractDetails(contract *ibapi.Contract) ([]ibsync.ContractDetails, error) {
	return c.ib.ReqContractDetails(contract)
}

// NewStock 创建股票合约
func NewStock(symbol, exchange, currency string) *ibapi.Contract {
	return ibsync.NewStock(symbol, exchange, currency)
}

// NewForex 创建外汇合约
func NewForex(base, exchange, quote string) *ibapi.Contract {
	return ibsync.NewForex(base, exchange, quote)
}

// NewOption 创建期权合约
func NewOption(symbol, lastTradeDateOrContractMonth string, strike float64, right, exchange, multiplier, currency string) *ibapi.Contract {
	return ibsync.NewOption(symbol, lastTradeDateOrContractMonth, strike, right, exchange, multiplier, currency)
}

// NewFuture 创建期货合约
func NewFuture(symbol, lastTradeDateOrContractMonth, exchange, multiplier, currency string) *ibapi.Contract {
	return ibsync.NewFuture(symbol, lastTradeDateOrContractMonth, exchange, multiplier, currency)
}
