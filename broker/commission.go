package broker

// CommissionModel 费率模型接口，由各券商实现（如 IBKR）
type CommissionModel interface {
	Calculate(t Trade) float64
}

// Trade 单笔成交信息，用于计算佣金与规费
type Trade struct {
	Shares     int     // 股数
	Price      float64 // 成交价
	IsSell     bool    // 是否卖出
	MonthlyVol int     // 当月累计成交股数（阶梯费率用）
}
