package livefeed

import "quantforge/dataengine"

// LiveFeed 实时行情推送：订阅标的并回调 Tick
type LiveFeed interface {
	Subscribe(symbol string) error
	Unsubscribe(symbol string) error
	OnTick(callback func(tick *dataengine.Tick))
}
