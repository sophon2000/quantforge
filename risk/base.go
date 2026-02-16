package risk

import "quantforge/execution"

// RiskManager 风控：下单前校验
type RiskManager interface {
	Check(order *execution.Order) error
}
