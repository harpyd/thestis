package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/performance"
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
		SubscribePerformanceCancel(ctx context.Context, perfID string) (<-chan Canceled, error)
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
		defer close(messages)
		defer func() {
			if err := m.guard.ReleasePerformance(
				context.Background(),
				perf.ID(),
			); err != nil {
				messages <- NewMessageFromError(err)
			}
		}()

		fr := performance.FlowFromPerformance(uuid.New().String(), perf)

		go func() {
			select {
			case <-ctx.Done():
			case <-canceled:
				cancel()
			}
		}()

		m.stepsPolicy.HandleSteps(ctx, fr, steps, messages)
	}()

	return messages, nil
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

type (
	publishCancelError struct {
		err error
	}

	subscribeCancelError struct {
		err error
	}
)

func NewPublishCancelError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(publishCancelError{err: err})
}

func IsPublishCancelError(err error) bool {
	var target publishCancelError

	return errors.As(err, &target)
}

func (e publishCancelError) Cause() error {
	return e.err
}

func (e publishCancelError) Unwrap() error {
	return e.err
}

func (e publishCancelError) Error() string {
	return fmt.Sprintf("publish cancel: %s", e.err)
}

func NewSubscribeCancelError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(subscribeCancelError{err: err})
}

func IsSubscribeCancelError(err error) bool {
	var target subscribeCancelError

	return errors.As(err, &target)
}

func (e subscribeCancelError) Cause() error {
	return e.err
}

func (e subscribeCancelError) Unwrap() error {
	return e.err
}

func (e subscribeCancelError) Error() string {
	return fmt.Sprintf("subscribe cancel: %s", e.err)
}
