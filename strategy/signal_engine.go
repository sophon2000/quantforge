package strategy

import "sync"

// DefaultSignalEngine 默认信号引擎：注册回调，收到信号时调用
type DefaultSignalEngine struct {
	mu        sync.RWMutex
	onSignal  func(Signal)
}

// NewDefaultSignalEngine 构造，onSignal 可为 nil
func NewDefaultSignalEngine(onSignal func(Signal)) *DefaultSignalEngine {
	return &DefaultSignalEngine{onSignal: onSignal}
}

// OnSignal 实现 SignalEngine
func (e *DefaultSignalEngine) OnSignal(signal Signal) {
	e.mu.RLock()
	fn := e.onSignal
	e.mu.RUnlock()
	if fn != nil {
		fn(signal)
	}
}

// SetOnSignal 设置信号回调
func (e *DefaultSignalEngine) SetOnSignal(fn func(Signal)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.onSignal = fn
}
