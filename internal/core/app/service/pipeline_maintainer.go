package service

import (
	"context"
	"time"

	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/pkg/correlationid"
)

type PipelineGuard interface {
	AcquirePipeline(ctx context.Context, pipeID string) error
	ReleasePipeline(ctx context.Context, pipeID string) error
}

type (
	PipelineMaintainer interface {
		MaintainPipeline(
			ctx context.Context,
			pipe *pipeline.Pipeline,
		) (<-chan DoneSignal, error)
	}

	DoneSignal struct{}
)

type pipelineMaintainer struct {
	guard      PipelineGuard
	subscriber PipelineCancelSubscriber
	policy     PipelinePolicy
	enqueuer   Enqueuer
	logger     Logger
	timeout    time.Duration
}

func NewPipelineMaintainer(
	guard PipelineGuard,
	cancelSub PipelineCancelSubscriber,
	policy PipelinePolicy,
	enqueuer Enqueuer,
	logger Logger,
	flowTimeout time.Duration,
) PipelineMaintainer {
	if guard == nil {
		panic("pipeline guard is nil")
	}

	if cancelSub == nil {
		panic("pipeline cancel subscriber is nil")
	}

	if policy == nil {
		panic("pipeline policy is nil")
	}

	if enqueuer == nil {
		panic("enqueuer is nil")
	}

	if logger == nil {
		panic("logger is nil")
	}

	return &pipelineMaintainer{
		guard:      guard,
		subscriber: cancelSub,
		policy:     policy,
		enqueuer:   enqueuer,
		logger:     logger,
		timeout:    flowTimeout,
	}
}

func (m *pipelineMaintainer) MaintainPipeline(
	ctx context.Context,
	pipe *pipeline.Pipeline,
) (<-chan DoneSignal, error) {
	correlationID := correlationid.FromCtx(ctx)

	l := m.enrichedLogger(pipe, correlationID)

	if err := m.guard.AcquirePipeline(ctx, pipe.ID()); err != nil {
		return nil, err
	}

	l.Debug("Pipeline acquired")

	canceled, err := m.subscriber.SubscribePipelineCancel(pipe.ID())
	if err != nil {
		return nil, err
	}

	l.Debug("Subscription to the pipeline cancellation signal has been issued")

	done := make(chan DoneSignal)

	m.enqueuer.Enqueue(m.maintainFn(pipe, canceled, done, correlationid.FromCtx(ctx)))

	l.Debug("Pipeline enqueued")

	return done, nil
}

func (m *pipelineMaintainer) maintainFn(
	pipe *pipeline.Pipeline,
	canceled <-chan CancelSignal,
	done chan<- DoneSignal,
	correlationID string,
) func() {
	return func() {
		defer close(done)
		defer m.releasePipeline(pipe, correlationID)

		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()

		ctx = correlationid.AssignToCtx(ctx, correlationID)

		l := m.enrichedLogger(pipe, correlationID)

		go func() {
			select {
			case <-ctx.Done():
			case <-canceled:
				l.Debug("Cancel signal received")

				cancel()
			}
		}()

		m.policy.ConsumePipeline(ctx, pipe)

		l.Debug("Pipeline consumed")
	}
}

func (m *pipelineMaintainer) enrichedLogger(
	pipe *pipeline.Pipeline,
	correlationID string,
) Logger {
	return m.logger.With(
		"correlationId", correlationID,
		"pipelineId", pipe.ID(),
	)
}

// releasePipeline releases pipeline with background context
// because pipeline should be released anyway.
func (m *pipelineMaintainer) releasePipeline(
	pipe *pipeline.Pipeline,
	correlationID string,
) {
	l := m.enrichedLogger(pipe, correlationID)

	if err := m.guard.ReleasePipeline(
		context.Background(),
		pipe.ID(),
	); err != nil {
		l.Error("Attempt to release pipeline failed", "error", err)
	}

	l.Debug("Pipeline released")
}
