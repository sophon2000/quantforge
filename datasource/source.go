package datasource

// MarketDataSource 行情数据源：订阅标的并推送 Tick
type MarketDataSource interface {
	Subscribe(symbol string) error
	Unsubscribe(symbol string) error
	OnTick(callback func(tick *Tick))
}
