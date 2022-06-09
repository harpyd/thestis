package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/core/app/service"
)

type PipelineGuard struct {
	acqErr error
	rlsErr error

	acqCalls int
	rlsCalls int
}

func NewPipelineGuard(acquireErr error, releaseErr error) *PipelineGuard {
	return &PipelineGuard{
		acqErr: acquireErr,
		rlsErr: releaseErr,

		rlsCalls: 0,
	}
}

func (g *PipelineGuard) AcquirePipeline(ctx context.Context, _ string) error {
	g.acqCalls++

	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	return g.acqErr
}

func (g *PipelineGuard) ReleasePipeline(ctx context.Context, _ string) error {
	g.rlsCalls++

	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	return g.rlsErr
}

func (g *PipelineGuard) AcquireCalls() int {
	return g.acqCalls
}

func (g *PipelineGuard) ReleaseCalls() int {
	return g.rlsCalls
}
