package strategyinterface

import (
	"github.com/sophon2000/quantforge/dataengine"
	"github.com/sophon2000/quantforge/executionengine"
)

// Strategy 策略接口：响应行情、K 线、订单更新
type Strategy interface {
	OnTick(t *dataengine.Tick)
	OnBar(b *dataengine.Bar)
	OnOrderUpdate(order *executionengine.Order)
}
