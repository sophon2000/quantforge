package risk

import "github.com/sophon2000/quantforge/executionengine"

// RiskManager 风控：下单前校验
type RiskManager interface {
	Check(order *executionengine.Order) error
}
