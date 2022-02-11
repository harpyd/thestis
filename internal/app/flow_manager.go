package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/harpyd/thestis/internal/domain/performance"
)

type (
	FlowManager interface {
		ManageFlow(ctx context.Context, perf *performance.Performance) (<-chan Message, error)
	}

	Message struct {
		s     string
		state performance.State
		err   error
	}
)

func newMessageFromStep(s performance.Step) Message {
	return Message{
		s:     s.String(),
		state: s.State(),
		err:   s.Err(),
	}
}

func newMessageFromError(err error) Message {
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

type everyStepSavingFlowManager struct {
	perfsRepo PerformancesRepository
	flowsRepo FlowsRepository
	timeout   time.Duration
}

func NewEveryStepSavingFlowManager(
	perfsRepo PerformancesRepository,
	flowsRepo FlowsRepository,
	timeout time.Duration,
) FlowManager {
	if perfsRepo == nil {
		panic("performance repository is nil")
	}

	if flowsRepo == nil {
		panic("flows repository is nil")
	}

	return &everyStepSavingFlowManager{
		perfsRepo: perfsRepo,
		flowsRepo: flowsRepo,
		timeout:   timeout,
	}
}

func (m *everyStepSavingFlowManager) ManageFlow(
	ctx context.Context,
	perf *performance.Performance,
) (<-chan Message, error) {
	steps, err := perf.Start(ctx)
	if err != nil {
		return nil, err
	}

	messages := make(chan Message)

	if err = m.perfsRepo.ExclusivelyDoWithPerformance(
		ctx,
		perf,
		m.action(steps, messages),
	); err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *everyStepSavingFlowManager) action(
	steps <-chan performance.Step,
	messages chan<- Message,
) PerformanceAction {
	return func(ctx context.Context, perf *performance.Performance) {
		defer close(messages)

		ctx, cancel := context.WithTimeout(ctx, m.timeout)
		defer cancel()

		fr := performance.FlowFromPerformance(uuid.New().String(), perf)

		for s := range steps {
			fr.WithStep(s)

			flow := fr.Reduce()
			if err := m.flowsRepo.UpsertFlow(ctx, flow); err != nil {
				messages <- newMessageFromError(err)
			}

			messages <- newMessageFromStep(s)
		}

		flow := fr.FinallyReduce()
		if err := m.flowsRepo.UpsertFlow(ctx, flow); err != nil {
			messages <- newMessageFromError(err)
		}
	}
}
