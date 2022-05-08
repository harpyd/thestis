package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/core/app/service"
)

type PerformanceGuard struct {
	acqErr error
	rlsErr error

	acqCalls int
	rlsCalls int
}

func NewPerformanceGuard(acquireErr error, releaseErr error) *PerformanceGuard {
	return &PerformanceGuard{
		acqErr: acquireErr,
		rlsErr: releaseErr,

		rlsCalls: 0,
	}
}

func (g *PerformanceGuard) AcquirePerformance(ctx context.Context, _ string) error {
	g.acqCalls++

	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	return g.acqErr
}

func (g *PerformanceGuard) ReleasePerformance(ctx context.Context, _ string) error {
	g.rlsCalls++

	if ctx.Err() != nil {
		return service.WrapWithDatabaseError(ctx.Err())
	}

	return g.rlsErr
}

func (g *PerformanceGuard) AcquireCalls() int {
	return g.acqCalls
}

func (g *PerformanceGuard) ReleaseCalls() int {
	return g.rlsCalls
}
