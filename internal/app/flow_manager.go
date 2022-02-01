package app

import (
	"context"

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
}

func NewEveryStepSavingFlowManager(
	perfsRepo PerformancesRepository,
	flowsRepo FlowsRepository,
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

	msg := make(chan Message)

	if err = m.perfsRepo.ExclusivelyDoWithPerformance(ctx, perf, m.action(ctx, steps, msg)); err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *everyStepSavingFlowManager) action(
	ctx context.Context,
	steps <-chan performance.Step,
	msg chan<- Message,
) func(perf *performance.Performance) {
	return func(perf *performance.Performance) {
		defer close(msg)

		fr := performance.FlowFromPerformance(uuid.New().String(), perf)

		for s := range steps {
			fr.WithStep(s)

			flow := fr.Reduce()
			if err := m.flowsRepo.UpsertFlow(ctx, flow); err != nil {
				msg <- newMessageFromError(err)
			}

			msg <- newMessageFromStep(s)
		}

		flow := fr.FinallyReduce()
		if err := m.flowsRepo.UpsertFlow(ctx, flow); err != nil {
			msg <- newMessageFromError(err)
		}
	}
}
