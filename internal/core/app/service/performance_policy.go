package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/performance"
)

type PerformancePolicy interface {
	ConsumePerformance(
		ctx context.Context,
		perf *performance.Performance,
		reactor MessageReactor,
	)
}

type savePerStepPolicy struct {
	flowRepo FlowRepository
	timeout  time.Duration
}

func NewSavePerStepPolicy(
	flowRepo FlowRepository,
	saveTimeout time.Duration,
) PerformancePolicy {
	if flowRepo == nil {
		panic("flow repository is nil")
	}

	return &savePerStepPolicy{
		flowRepo: flowRepo,
		timeout:  saveTimeout,
	}
}

func (p *savePerStepPolicy) ConsumePerformance(
	ctx context.Context,
	perf *performance.Performance,
	reactor MessageReactor,
) {
	var (
		steps = perf.MustStart(ctx)
		f     = flow.FromPerformance(uuid.New().String(), perf)
	)

	defer func() {
		if err := p.flowRepo.UpsertFlow(context.Background(), f); err != nil {
			reactor(NewMessageFromError(err))
		}
	}()

	for s := range steps {
		reactor(NewMessageFromStep(s))

		if err := p.upsertFlowWithTimeout(ctx, f.ApplyStep(s)); err != nil {
			reactor(NewMessageFromError(err))
		}
	}
}

func (p *savePerStepPolicy) upsertFlowWithTimeout(
	ctx context.Context,
	flow *flow.Flow,
) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.flowRepo.UpsertFlow(ctx, flow)
}
