package service

import (
	"context"
	"time"

	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/pkg/correlationid"
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
		) (<-chan DoneSignal, error)
	}

	DoneSignal struct{}
)

type performanceMaintainer struct {
	guard      PerformanceGuard
	subscriber PerformanceCancelSubscriber
	policy     PerformancePolicy
	enqueuer   Enqueuer
	logger     Logger
	timeout    time.Duration
}

func NewPerformanceMaintainer(
	guard PerformanceGuard,
	cancelSub PerformanceCancelSubscriber,
	policy PerformancePolicy,
	enqueuer Enqueuer,
	logger Logger,
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

	if logger == nil {
		panic("logger is nil")
	}

	return &performanceMaintainer{
		guard:      guard,
		subscriber: cancelSub,
		policy:     policy,
		enqueuer:   enqueuer,
		logger:     logger,
		timeout:    flowTimeout,
	}
}

func (m *performanceMaintainer) MaintainPerformance(
	ctx context.Context,
	perf *performance.Performance,
) (<-chan DoneSignal, error) {
	correlationID := correlationid.FromCtx(ctx)

	l := m.enrichedLogger(perf, correlationID)

	if err := m.guard.AcquirePerformance(ctx, perf.ID()); err != nil {
		return nil, err
	}

	l.Debug("Performance acquired")

	canceled, err := m.subscriber.SubscribePerformanceCancel(perf.ID())
	if err != nil {
		return nil, err
	}

	l.Debug("Subscription to the performance cancellation signal has been issued")

	done := make(chan DoneSignal)

	m.enqueuer.Enqueue(m.maintainFn(perf, canceled, done, correlationid.FromCtx(ctx)))

	l.Debug("Performance enqueued")

	return done, nil
}

func (m *performanceMaintainer) maintainFn(
	perf *performance.Performance,
	canceled <-chan CancelSignal,
	done chan<- DoneSignal,
	correlationID string,
) func() {
	return func() {
		defer close(done)
		defer m.releasePerformance(perf, correlationID)

		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		ctx = correlationid.AssignToCtx(ctx, correlationID)

		l := m.enrichedLogger(perf, correlationID)

		go func() {
			select {
			case <-ctx.Done():
			case <-canceled:
				l.Debug("Cancel signal received")

				cancel()
			}
		}()

		m.policy.ConsumePerformance(ctx, perf)

		l.Debug("Performance consumed")
	}
}

func (m *performanceMaintainer) enrichedLogger(
	perf *performance.Performance,
	correlationID string,
) Logger {
	return m.logger.With(
		"correlationId", correlationID,
		"performanceId", perf.ID(),
	)
}

// releasePerformance releases performance with background context
// because performances should be released anyway.
func (m *performanceMaintainer) releasePerformance(
	perf *performance.Performance,
	correlationID string,
) {
	l := m.enrichedLogger(perf, correlationID)

	if err := m.guard.ReleasePerformance(
		context.Background(),
		perf.ID(),
	); err != nil {
		l.Error("Attempt to release performance failed", "error", err)
	}

	l.Debug("Performance released")
}
