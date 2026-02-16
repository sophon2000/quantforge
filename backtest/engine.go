package backtest

// Engine 回测引擎接口（与 MatchingEngine 同义，便于扩展）
type Engine interface {
	MatchingEngine
}
