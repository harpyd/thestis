package mock

import (
	"sync"

	"github.com/harpyd/thestis/internal/core/app/service"
)

type PipelineCancelPubsub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan service.CancelSignal

	pubCalls int
	subCalls int
}

func NewPipelineCancelPubsub() *PipelineCancelPubsub {
	return &PipelineCancelPubsub{
		subscribers: make(map[string][]chan service.CancelSignal),
	}
}

func (ps *PipelineCancelPubsub) PublishPipelineCancel(pipeID string) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ps.pubCalls++

	channels := ps.subscribers[pipeID]

	for _, ch := range channels {
		go func(ch chan<- service.CancelSignal) {
			close(ch)
		}(ch)
	}

	ps.subscribers[pipeID] = nil

	return nil
}

func (ps *PipelineCancelPubsub) SubscribePipelineCancel(pipeID string) (<-chan service.CancelSignal, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subCalls++

	ch := make(chan service.CancelSignal, 1)
	ps.subscribers[pipeID] = append(ps.subscribers[pipeID], ch)

	return ch, nil
}

func (ps *PipelineCancelPubsub) PublishCalls() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.pubCalls
}

func (ps *PipelineCancelPubsub) SubscribeCalls() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.subCalls
}
