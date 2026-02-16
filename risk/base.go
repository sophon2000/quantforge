package risk

import "quantforge/executionengine"

// RiskManager 风控：下单前校验
type RiskManager interface {
	Check(order *executionengine.Order) error
}
