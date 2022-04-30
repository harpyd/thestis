package app

import (
	"context"
	"time"

	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/performance"
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
	flowRepo FlowRepository
	timeout  time.Duration
}

func NewEveryStepSavingPolicy(flowRepo FlowRepository, saveTimeout time.Duration) StepsPolicy {
	if flowRepo == nil {
		panic("flow repository is nil")
	}

	return &everyStepSavingPolicy{
		flowRepo: flowRepo,
		timeout:  saveTimeout,
	}
}

func (p *everyStepSavingPolicy) HandleSteps(
	ctx context.Context,
	fr *flow.Reducer,
	steps <-chan performance.Step,
	messages chan<- Message,
) {
	defer func() {
		if err := p.flowRepo.UpsertFlow(
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

func (p *everyStepSavingPolicy) upsertFlowWithTimeout(ctx context.Context, flow *flow.Flow) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.flowRepo.UpsertFlow(ctx, flow)
}
