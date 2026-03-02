package dataprocessor

import "github.com/sophon2000/quantforge/dataengine"

// DataProcessor 数据处理器：如 Tick 聚合成 Bar、过滤、归一化等
type DataProcessor interface {
	// ProcessTick 处理单笔 Tick（可聚合为 Bar 或转发）
	ProcessTick(t *dataengine.Tick)
	// OnBar 注册 K 线完成回调
	OnBar(callback func(b *dataengine.Bar))
}
