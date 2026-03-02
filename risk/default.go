package risk

import "github.com/sophon2000/quantforge/executionengine"

// DefaultRiskManager 默认风控：仅做通过校验（不拒绝任何订单）
type DefaultRiskManager struct{}

// Check 实现 RiskManager，直接通过
func (DefaultRiskManager) Check(order *executionengine.Order) error {
	return nil
}

// NewDefaultRiskManager 构造
func NewDefaultRiskManager() *DefaultRiskManager {
	return &DefaultRiskManager{}
}
