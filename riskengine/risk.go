package riskengine

import "quantforge/executionengine"

// RiskManager 风控接口：下单前校验
type RiskManager interface {
	Check(order *executionengine.Order) error
}
