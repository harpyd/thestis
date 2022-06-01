package service

import (
	"context"
	"time"

	"github.com/harpyd/thestis/internal/core/entity/performance"
)

type PerformanceGuard interface {
	AcquirePerformance(ctx context.Context, perfID string) error
	ReleasePerformance(ctx context.Context, perfID string) error
}

type (
	PerformanceMaintainer interface {
		MaintainPerformance(
			ctx context.Context,
			perf *performance.Performance,
			reactor MessageReactor,
		) (<-chan DoneSignal, error)
	}

	DoneSignal struct{}
)

type performanceMaintainer struct {
	guard      PerformanceGuard
	subscriber PerformanceCancelSubscriber
	policy     PerformancePolicy
	enqueuer   Enqueuer
	timeout    time.Duration
}

func NewPerformanceMaintainer(
	guard PerformanceGuard,
	cancelSub PerformanceCancelSubscriber,
	policy PerformancePolicy,
	enqueuer Enqueuer,
	flowTimeout time.Duration,
) PerformanceMaintainer {
	if guard == nil {
		panic("performance guard is nil")
	}

	if cancelSub == nil {
		panic("performance cancel subscriber is nil")
	}

	if policy == nil {
		panic("performance policy is nil")
	}

	if enqueuer == nil {
		panic("enqueuer is nil")
	}

	return &performanceMaintainer{
		guard:      guard,
		subscriber: cancelSub,
		policy:     policy,
		enqueuer:   enqueuer,
		timeout:    flowTimeout,
	}
}

func (m *performanceMaintainer) MaintainPerformance(
	ctx context.Context,
	perf *performance.Performance,
	reactor MessageReactor,
) (<-chan DoneSignal, error) {
	if err := m.guard.AcquirePerformance(ctx, perf.ID()); err != nil {
		return nil, err
	}

	canceled, err := m.subscriber.SubscribePerformanceCancel(perf.ID())
	if err != nil {
		return nil, err
	}

	done := make(chan DoneSignal)

	m.enqueuer.Enqueue(m.maintainFn(perf, canceled, done, reactor))

	return done, nil
}

func (m *performanceMaintainer) maintainFn(
	perf *performance.Performance,
	canceled <-chan CancelSignal,
	done chan<- DoneSignal,
	reactor MessageReactor,
) func() {
	return func() {
		defer close(done)
		defer m.releasePerformance(perf, reactor)

		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		go func() {
			select {
			case <-ctx.Done():
			case <-canceled:
				cancel()
			}
		}()

		m.policy.ConsumePerformance(ctx, perf, reactor)
	}
}

func (m *performanceMaintainer) releasePerformance(
	perf *performance.Performance,
	reactor MessageReactor,
) {
	if err := m.guard.ReleasePerformance(
		context.Background(),
		perf.ID(),
	); err != nil {
		reactor(NewMessageFromError(err))
	}
}
