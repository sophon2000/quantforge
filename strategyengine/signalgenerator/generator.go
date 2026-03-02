package signalgenerator

import (
	"sync"

	"github.com/sophon2000/quantforge/strategyengine"
)

// DefaultSignalGenerator 默认信号生成器：回调转发
type DefaultSignalGenerator struct {
	mu       sync.RWMutex
	onSignal func(strategyengine.Signal)
}

// NewDefaultSignalGenerator 构造
func NewDefaultSignalGenerator(onSignal func(strategyengine.Signal)) *DefaultSignalGenerator {
	return &DefaultSignalGenerator{onSignal: onSignal}
}

// OnSignal 实现 SignalGenerator
func (g *DefaultSignalGenerator) OnSignal(signal strategyengine.Signal) {
	g.mu.RLock()
	fn := g.onSignal
	g.mu.RUnlock()
	if fn != nil {
		fn(signal)
	}
}

// SetOnSignal 设置回调
func (g *DefaultSignalGenerator) SetOnSignal(fn func(strategyengine.Signal)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.onSignal = fn
}

// DefaultSignalEngine 默认信号引擎（与 DefaultSignalGenerator 同义）
type DefaultSignalEngine = DefaultSignalGenerator

// NewDefaultSignalEngine 构造（与 NewDefaultSignalGenerator 同义）
func NewDefaultSignalEngine(onSignal func(strategyengine.Signal)) *DefaultSignalGenerator {
	return NewDefaultSignalGenerator(onSignal)
}
