package ibkr

import (
	"github.com/scmhub/ibapi"
	"github.com/scmhub/ibsync"
)

// PlaceOrder 下单
func (c *Client) PlaceOrder(contract *ibapi.Contract, order *ibapi.Order) *ibsync.Trade {
	return c.ib.PlaceOrder(contract, order)
}

// CancelOrder 取消订单
func (c *Client) CancelOrder(order *ibapi.Order, orderCancel ibsync.OrderCancel) {
	c.ib.CancelOrder(order, orderCancel)
}

// ReqGlobalCancel 取消所有订单
func (c *Client) ReqGlobalCancel() {
	c.ib.ReqGlobalCancel()
}

// LimitOrder 创建限价单
func LimitOrder(action string, quantity ibsync.Decimal, limitPrice float64) *ibapi.Order {
	return ibsync.LimitOrder(action, quantity, limitPrice)
}

// MarketOrder 创建市价单
func MarketOrder(action string, quantity ibsync.Decimal) *ibapi.Order {
	return ibsync.MarketOrder(action, quantity)
}

// StopOrder 创建止损单
func StopOrder(action string, quantity ibsync.Decimal, stopPrice float64) *ibapi.Order {
	return ibapi.Stop(action, quantity, stopPrice)
}

// StopLimitOrder 创建止损限价单
func StopLimitOrder(action string, quantity ibsync.Decimal, limitPrice, stopPrice float64) *ibapi.Order {
	return ibapi.StopLimit(action, quantity, limitPrice, stopPrice)
}

// NewOrderCancel 创建订单取消请求
func NewOrderCancel() ibsync.OrderCancel {
	return ibsync.OrderCancel{}
}

// StringToDecimal 字符串转 Decimal
func StringToDecimal(s string) ibsync.Decimal {
	return ibsync.StringToDecimal(s)
}
