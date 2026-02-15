package engine

type Order struct {
	ID       string
	Symbol   string
	Quantity int
	Price    float64
	Status   string
}
