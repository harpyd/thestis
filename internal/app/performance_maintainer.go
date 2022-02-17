package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/harpyd/thestis/internal/domain/performance"
)

type PerformanceReleaser interface {
	ReleasePerformance(ctx context.Context, perfID string) error
}

type PerformanceMaintainer interface {
	MaintainPerformance(ctx context.Context, perf *performance.Performance) (<-chan Message, error)
}

type performanceMaintainer struct {
	releaser    PerformanceReleaser
	stepsPolicy StepsPolicy
	timeout     time.Duration
}

func NewPerformanceMaintainer(
	releaser PerformanceReleaser,
	stepsPolicy StepsPolicy,
	flowTimeout time.Duration,
) PerformanceMaintainer {
	if releaser == nil {
		panic("performance releaser is nil")
	}

	if stepsPolicy == nil {
		panic("steps policy is nil")
	}

	return &performanceMaintainer{
		releaser:    releaser,
		stepsPolicy: stepsPolicy,
		timeout:     flowTimeout,
	}
}

func (m *performanceMaintainer) MaintainPerformance(
	ctx context.Context,
	perf *performance.Performance,
) (<-chan Message, error) {
	messages := make(chan Message)

	ctx, cancel := context.WithTimeout(ctx, m.timeout)

	steps, err := perf.Start(ctx)
	if err != nil {
		cancel()

		return nil, err
	}

	go func() {
		defer cancel()
		defer close(messages)
		defer func() {
			if err := m.releaser.ReleasePerformance(
				context.Background(),
				perf.ID(),
			); err != nil {
				messages <- NewMessageFromError(err)
			}
		}()

		fr := performance.FlowFromPerformance(uuid.New().String(), perf)

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
