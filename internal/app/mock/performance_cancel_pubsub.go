package mock

import (
	"context"
	"sync"

	"github.com/harpyd/thestis/internal/app"
)

type PerformanceCancelPubsub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan app.Canceled

	pubCalls int
}

func NewPerformanceCancelPubsub() *PerformanceCancelPubsub {
	return &PerformanceCancelPubsub{
		subscribers: make(map[string][]chan app.Canceled),
	}
}

func (ps *PerformanceCancelPubsub) PublishPerformanceCancel(perfID string) error {
	ps.pubCalls++

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	channels := ps.subscribers[perfID]

	for _, ch := range channels {
		go func(ch chan<- app.Canceled) {
			ch <- app.Canceled{}
		}(ch)
	}

	return nil
}

func (ps *PerformanceCancelPubsub) SubscribePerformanceCancel(ctx context.Context, perfID string) (<-chan app.Canceled, error) {
	select {
	case <-ctx.Done():
		return nil, app.NewSubscribeCancelError(ctx.Err())
	default:
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan app.Canceled, 1)
	ps.subscribers[perfID] = append(ps.subscribers[perfID], ch)

	return ch, nil
}

func (ps *PerformanceCancelPubsub) PublishCalls() int {
	return ps.pubCalls
}
