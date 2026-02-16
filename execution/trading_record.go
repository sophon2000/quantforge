package execution

// TradingRecord 交易记录（用于统计/回测）
type TradingRecord struct {
	ID           string
	Symbol       string
	Quantity     int
	EntryPrice   float64
	CurrentPrice float64
	Profit       float64
	Status       string
}
