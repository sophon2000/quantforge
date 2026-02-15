package engine

type Position struct {
	ID           string
	Symbol       string
	Quantity     int
	EntryPrice   float64
	CurrentPrice float64
	Profit       float64
	Status       string
}
