package riskengine

import "quantforge/executionengine"

// DefaultRiskManager 默认风控：全部通过
type DefaultRiskManager struct{}

// Check 实现 RiskManager
func (DefaultRiskManager) Check(order *executionengine.Order) error {
	return nil
}

// NewDefaultRiskManager 构造
func NewDefaultRiskManager() *DefaultRiskManager {
	return &DefaultRiskManager{}
}
