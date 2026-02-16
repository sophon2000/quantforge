package execution

// Broker 经纪商接口：下单、撤单
type Broker interface {
	PlaceOrder(order Order) error
	CancelOrder(id string)
}

// OrderStatus 订单状态：NEW → SUBMITTED → PARTIAL_FILLED → FILLED / CANCELED
type OrderStatus string

const (
	NEW            OrderStatus = "NEW"
	SUBMITTED      OrderStatus = "SUBMITTED"
	PARTIAL_FILLED OrderStatus = "PARTIAL_FILLED"
	FILLED         OrderStatus = "FILLED"
	CANCELED       OrderStatus = "CANCELED"
)

// Order 订单
type Order struct {
	ID       string
	Symbol   string
	Quantity int
	Price    float64
	Status   OrderStatus
}

// Fill 成交
type Fill struct {
	OrderID  string
	Symbol   string
	Price    float64
	Quantity int
}
