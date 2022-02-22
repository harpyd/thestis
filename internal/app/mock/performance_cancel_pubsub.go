package mock

import (
	"context"
	"sync"

	"github.com/harpyd/thestis/internal/app"
)

type PerformanceCancelPubsub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan app.CancelSignal

	pubCalls int
}

func NewPerformanceCancelPubsub() *PerformanceCancelPubsub {
	return &PerformanceCancelPubsub{
		subscribers: make(map[string][]chan app.CancelSignal),
	}
}

func (ps *PerformanceCancelPubsub) PublishPerformanceCancel(perfID string) error {
	ps.pubCalls++

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	channels := ps.subscribers[perfID]

	for _, ch := range channels {
		go func(ch chan<- app.CancelSignal) {
			close(ch)
		}(ch)
	}

	ps.subscribers[perfID] = nil

	return nil
}

func (ps *PerformanceCancelPubsub) SubscribePerformanceCancel(ctx context.Context, perfID string) (<-chan app.CancelSignal, error) {
	select {
	case <-ctx.Done():
		return nil, app.NewSubscribeCancelError(ctx.Err())
	default:
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan app.CancelSignal, 1)
	ps.subscribers[perfID] = append(ps.subscribers[perfID], ch)

	return ch, nil
}

func (ps *PerformanceCancelPubsub) PublishCalls() int {
	return ps.pubCalls
}
