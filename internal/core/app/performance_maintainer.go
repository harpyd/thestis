package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/performance"
)

type PerformanceGuard interface {
	AcquirePerformance(ctx context.Context, perfID string) error
	ReleasePerformance(ctx context.Context, perfID string) error
}

type (
	PerformanceCancelPublisher interface {
		PublishPerformanceCancel(perfID string) error
	}

	PerformanceCancelSubscriber interface {
		SubscribePerformanceCancel(perfID string) (<-chan CancelSignal, error)
	}

	CancelSignal = struct{}
)

type PerformanceMaintainer interface {
	MaintainPerformance(ctx context.Context, perf *performance.Performance) (<-chan Message, error)
}

type Enqueuer interface {
	Enqueue(fn func())
}

type EnqueueFunc func(fn func())

func (e EnqueueFunc) Enqueue(fn func()) {
	e(fn)
}

type performanceMaintainer struct {
	guard       PerformanceGuard
	subscriber  PerformanceCancelSubscriber
	stepsPolicy StepsPolicy
	enqueuer    Enqueuer
	timeout     time.Duration
}

func NewPerformanceMaintainer(
	guard PerformanceGuard,
	cancelSub PerformanceCancelSubscriber,
	stepsPolicy StepsPolicy,
	enqueuer Enqueuer,
	flowTimeout time.Duration,
) PerformanceMaintainer {
	if guard == nil {
		panic("performance guard is nil")
	}

	if cancelSub == nil {
		panic("performance cancel subscriber is nil")
	}

	if stepsPolicy == nil {
		panic("steps policy is nil")
	}

	if enqueuer == nil {
		panic("enqueuer is nil")
	}

	return &performanceMaintainer{
		guard:       guard,
		subscriber:  cancelSub,
		stepsPolicy: stepsPolicy,
		enqueuer:    enqueuer,
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

	canceled, err := m.subscriber.SubscribePerformanceCancel(perf.ID())
	if err != nil {
		return nil, err
	}

	messages := make(chan Message)

	m.enqueuer.Enqueue(func() {
		defer close(messages)
		defer func() {
			if err := m.guard.ReleasePerformance(
				context.Background(),
				perf.ID(),
			); err != nil {
				messages <- NewMessageFromError(err)
			}
		}()

		ctx, cancel := context.WithTimeout(ctx, m.timeout)
		defer cancel()

		go func() {
			select {
			case <-ctx.Done():
			case <-canceled:
				cancel()
			}
		}()

		steps := perf.MustStart(ctx)

		fr := flow.FromPerformance(uuid.New().String(), perf)

		m.stepsPolicy.HandleSteps(ctx, fr, steps, messages)
	})

	return messages, nil
}

type Message struct {
	s     string
	event performance.Event
	err   error
}

func NewMessageFromStep(s performance.Step) Message {
	return Message{
		s:     s.String(),
		event: s.Event(),
		err:   s.Err(),
	}
}

func NewMessageFromError(err error) Message {
	return Message{
		s:     err.Error(),
		event: performance.NoEvent,
		err:   err,
	}
}

func (m Message) String() string {
	return m.s
}

func (m Message) Err() error {
	return m.err
}

func (m Message) Event() performance.Event {
	return m.event
}

type PublishCancelError struct {
	err error
}

func WrapWithPublishCancelError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(&PublishCancelError{err: err})
}

func (e *PublishCancelError) Unwrap() error {
	return e.err
}

func (e *PublishCancelError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return fmt.Sprintf("publish cancel: %s", e.err)
}

type SubscribeCancelError struct {
	err error
}

func WrapWithSubscribeCancelError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(&SubscribeCancelError{err: err})
}

func (e *SubscribeCancelError) Unwrap() error {
	return e.err
}

func (e *SubscribeCancelError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return fmt.Sprintf("subscribe cancel: %s", e.err)
}
