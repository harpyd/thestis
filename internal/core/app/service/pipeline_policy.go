package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/pkg/correlationid"
)

type PipelinePolicy interface {
	ConsumePipeline(
		ctx context.Context,
		pipe *pipeline.Pipeline,
	)
}

type savePerStepPolicy struct {
	flowRepo FlowRepository
	logger   Logger
	timeout  time.Duration
}

func NewSavePerStepPolicy(
	flowRepo FlowRepository,
	logger Logger,
	saveTimeout time.Duration,
) PipelinePolicy {
	if flowRepo == nil {
		panic("flow repository is nil")
	}

	if logger == nil {
		panic("logger is nil")
	}

	return &savePerStepPolicy{
		flowRepo: flowRepo,
		logger:   logger,
		timeout:  saveTimeout,
	}
}

func (p *savePerStepPolicy) ConsumePipeline(
	ctx context.Context,
	pipeline *pipeline.Pipeline,
) {
	var (
		steps = pipeline.MustStart(ctx)
		f     = flow.Fulfill(uuid.New().String(), pipeline)
	)

	l := p.enrichedLogger(ctx, pipeline, f)

	defer func() {
		if err := p.flowRepo.UpsertFlow(context.Background(), f); err != nil {
			p.logger.Error("Last attempt to upsert flow failed", "error", err)
		}

		l.Debug("Flow upserted for last")
	}()

	for s := range steps {
		if err := p.upsertFlowWithTimeout(ctx, f.ApplyStep(s)); err != nil {
			l.Warn("Attempt to upsert flow failed", "error", err)
		}

		l.Debug("Flow with new step upserted", "step", s)
	}
}

func (p *savePerStepPolicy) enrichedLogger(
	ctx context.Context,
	pipe *pipeline.Pipeline,
	f *flow.Flow,
) Logger {
	return p.logger.With(
		"correlationId", correlationid.FromCtx(ctx),
		"pipelineId", pipe.ID(),
		"flowId", f.ID(),
	)
}

func (p *savePerStepPolicy) upsertFlowWithTimeout(
	ctx context.Context,
	flow *flow.Flow,
) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	return p.flowRepo.UpsertFlow(ctx, flow)
}
