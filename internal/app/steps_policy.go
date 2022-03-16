package app

import (
	"context"
	"time"

	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type StepsPolicy interface {
	HandleSteps(
		ctx context.Context,
		fr *flow.Reducer,
		steps <-chan performance.Step,
		messages chan<- Message,
	)
}

type everyStepSavingPolicy struct {
	flowsRepo FlowsRepository
	timeout   time.Duration
}

func NewEveryStepSavingPolicy(flowsRepo FlowsRepository, saveTimeout time.Duration) StepsPolicy {
	if flowsRepo == nil {
		panic("flows repository is nil")
	}

	return &everyStepSavingPolicy{
		flowsRepo: flowsRepo,
		timeout:   saveTimeout,
	}
}

func (p *everyStepSavingPolicy) HandleSteps(
	ctx context.Context,
	fr *flow.Reducer,
	steps <-chan performance.Step,
	messages chan<- Message,
) {
	defer func() {
		if err := p.flowsRepo.UpsertFlow(
			context.Background(),
			fr.Reduce(),
		); err != nil {
			messages <- NewMessageFromError(err)
		}
	}()

	for s := range steps {
		messages <- NewMessageFromStep(s)

		fr.WithStep(s)

		if err := p.upsertFlowWithTimeout(ctx, fr.Reduce()); err != nil {
			messages <- NewMessageFromError(err)
		}
	}
}

func (p *everyStepSavingPolicy) upsertFlowWithTimeout(ctx context.Context, flow flow.Flow) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.flowsRepo.UpsertFlow(ctx, flow)
}
