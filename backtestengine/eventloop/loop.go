package eventloop

import (
	"sync"

	"github.com/sophon2000/quantforge/backtestengine"
	"github.com/sophon2000/quantforge/dataengine"
	"github.com/sophon2000/quantforge/strategyengine"
)

// DefaultEventLoop 默认事件循环：队列 + 单协程处理
type DefaultEventLoop struct {
	mu       sync.Mutex
	stopCh   chan struct{}
	marketCh chan *dataengine.Tick
	signalCh chan strategyengine.Signal
	onMarket func(*dataengine.Tick)
	onSignal func(strategyengine.Signal)
}

// NewDefaultEventLoop 构造
func NewDefaultEventLoop(
	onMarket func(*dataengine.Tick),
	onSignal func(strategyengine.Signal),
) *DefaultEventLoop {
	return &DefaultEventLoop{
		stopCh:   make(chan struct{}),
		marketCh: make(chan *dataengine.Tick, 256),
		signalCh: make(chan strategyengine.Signal, 64),
		onMarket: onMarket,
		onSignal: onSignal,
	}
}

// PushMarket 实现 EventLoop
func (l *DefaultEventLoop) PushMarket(tick *dataengine.Tick) {
	select {
	case l.marketCh <- tick:
	default:
		// 队列满则丢弃
	}
}

// PushSignal 实现 EventLoop
func (l *DefaultEventLoop) PushSignal(signal strategyengine.Signal) {
	select {
	case l.signalCh <- signal:
	default:
	}
}

// Run 实现 EventLoop
func (l *DefaultEventLoop) Run() {
	for {
		select {
		case <-l.stopCh:
			return
		case tick := <-l.marketCh:
			if l.onMarket != nil {
				l.onMarket(tick)
			}
		case sig := <-l.signalCh:
			if l.onSignal != nil {
				l.onSignal(sig)
			}
		}
	}
}

// Stop 实现 EventLoop
func (l *DefaultEventLoop) Stop() {
	close(l.stopCh)
}

// Ensure DefaultEventLoop implements backtestengine.EventLoop
var _ backtestengine.EventLoop = (*DefaultEventLoop)(nil)
