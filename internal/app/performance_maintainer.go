package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/pkg/signals"
)

type PerformanceGuard interface {
	AcquirePerformance(ctx context.Context, perfID string) error
	ReleasePerformance(ctx context.Context, perfID string) error
}

type (
	PerformanceCancelPublisher interface {
		PublishPerformanceCancel(ctx context.Context, perfID string) error
		Close() error
	}

	PerformanceCancelSubscriber interface {
		SubscribePerformanceCancel(ctx context.Context, perfID string) (<-chan Canceled, error)
		Close() error
	}

	Canceled = struct{}
)

type PerformanceMaintainer interface {
	MaintainPerformance(ctx context.Context, perf *performance.Performance) (<-chan Message, error)
}

type performanceMaintainer struct {
	guard       PerformanceGuard
	subscriber  PerformanceCancelSubscriber
	stepsPolicy StepsPolicy
	timeout     time.Duration
}

func NewPerformanceMaintainer(
	guard PerformanceGuard,
	cancelSub PerformanceCancelSubscriber,
	stepsPolicy StepsPolicy,
	flowTimeout time.Duration,
) PerformanceMaintainer {
	if guard == nil {
		panic("performance guard is nil")
	}

	if cancelSub == nil {
		panic("performance cancel receiver is nil")
	}

	if stepsPolicy == nil {
		panic("steps policy is nil")
	}

	return &performanceMaintainer{
		guard:       guard,
		subscriber:  cancelSub,
		stepsPolicy: stepsPolicy,
		timeout:     flowTimeout,
	}
}

func (m *performanceMaintainer) MaintainPerformance(
	ctx context.Context,
	perf *performance.Performance,
) (<-chan Message, error) {
	if err := m.guard.AcquirePerformance(ctx, perf.ID()); err != nil {
		return nil, err
	}

	canceled, err := m.subscriber.SubscribePerformanceCancel(ctx, perf.ID())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, m.timeout)

	steps, err := perf.Start(ctx)
	if err != nil {
		cancel()

		return nil, err
	}

	messages := make(chan Message)

	go func() {
		defer cancel()

		fr := performance.FlowFromPerformance(uuid.New().String(), perf)

		handled := m.handleSteps(ctx, perf.ID(), fr, steps, messages)

		<-signals.Or(handled, canceled)
	}()

	return messages, nil
}

func (m *performanceMaintainer) handleSteps(
	ctx context.Context,
	perfID string,
	fr *performance.FlowReducer,
	steps <-chan performance.Step,
	messages chan<- Message,
) <-chan struct{} {
	handled := make(chan struct{})

	go func() {
		defer close(handled)
		defer close(messages)
		defer func() {
			if err := m.guard.ReleasePerformance(
				context.Background(),
				perfID,
			); err != nil {
				messages <- NewMessageFromError(err)
			}
		}()

		m.stepsPolicy.HandleSteps(ctx, fr, steps, messages)
	}()

	return handled
}

type Message struct {
	s     string
	state performance.State
	err   error
}

func NewMessageFromStep(s performance.Step) Message {
	return Message{
		s:     s.String(),
		state: s.State(),
		err:   s.Err(),
	}
}

func NewMessageFromError(err error) Message {
	return Message{
		s:     err.Error(),
		state: performance.NoState,
		err:   err,
	}
}

func (m Message) String() string {
	return m.s
}

func (m Message) Err() error {
	return m.err
}

func (m Message) State() performance.State {
	return m.state
}
