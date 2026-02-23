package accountsimulator

type PricingMode int

type CommissionModel interface {
	Calculate(t Trade) float64
}

const (
	Tiered PricingMode = iota
	Fixed
)

type Trade struct {
	Shares     int
	Price      float64
	IsSell     bool
	MonthlyVol int // 当月累计成交股数（仅 Tiered 用）
}

func CalcFixedCommission(t Trade) float64 {
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

func CalcTieredCommission(t Trade) float64 {
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

func CalcRegulatoryFee(t Trade) float64 {
	var fee float64

	if t.IsSell {
		tradeValue := float64(t.Shares) * t.Price

		// SEC fee（示例比例，真实需按最新官方）
		secRate := 0.000008 // 0.0008%
		fee += tradeValue * secRate

		// FINRA TAF
		finraRate := 0.000145
		fee += float64(t.Shares) * finraRate
	}

	return fee
}

func CalcCommission(t Trade, mode PricingMode) float64 {
	var commission float64

	switch mode {
	case Fixed:
		commission = CalcFixedCommission(t)
	case Tiered:
		commission = CalcTieredCommission(t)
	}

	commission += CalcRegulatoryFee(t)

	return commission
}
