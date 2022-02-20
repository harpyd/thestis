package mock

import (
	"context"
	"sync"

	"github.com/harpyd/thestis/internal/app"
)

type PerformanceCancelPubsub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan app.Canceled
	closed      bool

	pubCalls int
}

func NewPerformanceCancelPubsub() *PerformanceCancelPubsub {
	return &PerformanceCancelPubsub{
		subscribers: make(map[string][]chan app.Canceled),
	}
}

func (ps *PerformanceCancelPubsub) PublishPerformanceCancel(ctx context.Context, perfID string) error {
	ps.pubCalls++

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return nil
	}

	channels, ok := ps.subscribers[perfID]
	if !ok {
		return nil
	}

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
		return nil, ctx.Err()
	default:
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan app.Canceled, 1)
	ps.subscribers[perfID] = append(ps.subscribers[perfID], ch)

	return ch, nil
}

func (ps *PerformanceCancelPubsub) Close() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return nil
	}

	ps.closed = true

	for _, sub := range ps.subscribers {
		for _, ch := range sub {
			close(ch)
		}
	}

	return nil
}

func (ps *PerformanceCancelPubsub) PublishCalls() int {
	return ps.pubCalls
}
