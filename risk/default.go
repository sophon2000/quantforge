package risk

import "quantforge/execution"

// DefaultRiskManager 默认风控：仅做通过校验（不拒绝任何订单）
type DefaultRiskManager struct{}

// Check 实现 RiskManager，直接通过
func (DefaultRiskManager) Check(order *execution.Order) error {
	return nil
}

// NewDefaultRiskManager 构造
func NewDefaultRiskManager() *DefaultRiskManager {
	return &DefaultRiskManager{}
}
