package ibkr

import (
	"github.com/sophon2000/quantforge/broker"
)

// PricingMode IBKR 定价模式：阶梯 vs 固定
type PricingMode int

const (
	// Tiered 阶梯费率（按当月成交量分档）
	Tiered PricingMode = iota
	// Fixed 固定费率
	Fixed
)

// Commission 实现 broker.CommissionModel，封装 IBKR 佣金与规费计算
type Commission struct {
	Mode PricingMode
}

// NewCommission 创建 IBKR 费率模型
func NewCommission(mode PricingMode) *Commission {
	return &Commission{Mode: mode}
}

// 编译期断言：*Commission 实现 broker.CommissionModel
var _ broker.CommissionModel = (*Commission)(nil)

// Calculate 实现 broker.CommissionModel
func (c *Commission) Calculate(t broker.Trade) float64 {
	var commission float64
	switch c.Mode {
	case Fixed:
		commission = c.calcFixed(t)
	case Tiered:
		commission = c.calcTiered(t)
	}
	commission += c.calcRegulatoryFee(t)
	return commission
}

func (c *Commission) calcFixed(t broker.Trade) float64 {
	const perShare = 0.005
	const minFee = 1.00

	tradeValue := float64(t.Shares) * t.Price
	commission := float64(t.Shares) * perShare

	if commission < minFee {
		commission = minFee
	}

	maxFee := tradeValue * 0.01
	if commission > maxFee {
		commission = maxFee
	}

	return commission
}

func getTieredRate(monthlyVol int) float64 {
	switch {
	case monthlyVol <= 300_000:
		return 0.0035
	case monthlyVol <= 3_000_000:
		return 0.0020
	case monthlyVol <= 20_000_000:
		return 0.0015
	case monthlyVol <= 100_000_000:
		return 0.0010
	default:
		return 0.0005
	}
}

func (c *Commission) calcTiered(t broker.Trade) float64 {
	const minFee = 0.35

	rate := getTieredRate(t.MonthlyVol)

	tradeValue := float64(t.Shares) * t.Price
	commission := float64(t.Shares) * rate

	if commission < minFee {
		commission = minFee
	}

	maxFee := tradeValue * 0.01
	if commission > maxFee {
		commission = maxFee
	}

	return commission
}

// calcRegulatoryFee 规费：SEC + FINRA TAF（仅卖出收取，比例示例，实际以官方为准）
func (c *Commission) calcRegulatoryFee(t broker.Trade) float64 {
	var fee float64

	if t.IsSell {
		tradeValue := float64(t.Shares) * t.Price

		secRate := 0.000008 // 0.0008%
		fee += tradeValue * secRate

		finraRate := 0.000145
		fee += float64(t.Shares) * finraRate
	}

	return fee
}
