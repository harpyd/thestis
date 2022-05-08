package mock

import (
	"sync"

	"github.com/harpyd/thestis/internal/core/app/service"
)

type PerformanceCancelPubsub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan service.CancelSignal

	pubCalls int
	subCalls int
}

func NewPerformanceCancelPubsub() *PerformanceCancelPubsub {
	return &PerformanceCancelPubsub{
		subscribers: make(map[string][]chan service.CancelSignal),
	}
}

func (ps *PerformanceCancelPubsub) PublishPerformanceCancel(perfID string) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ps.pubCalls++

	channels := ps.subscribers[perfID]

	for _, ch := range channels {
		go func(ch chan<- service.CancelSignal) {
			close(ch)
		}(ch)
	}

	ps.subscribers[perfID] = nil

	return nil
}

func (ps *PerformanceCancelPubsub) SubscribePerformanceCancel(perfID string) (<-chan service.CancelSignal, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subCalls++

	ch := make(chan service.CancelSignal, 1)
	ps.subscribers[perfID] = append(ps.subscribers[perfID], ch)

	return ch, nil
}

func (ps *PerformanceCancelPubsub) PublishCalls() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.pubCalls
}

func (ps *PerformanceCancelPubsub) SubscribeCalls() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.subCalls
}
