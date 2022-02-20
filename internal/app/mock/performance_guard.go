package mock

import "context"

type PerformanceGuard struct {
	acqErr error
	rlsErr error

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
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return g.acqErr
}

func (g *PerformanceGuard) ReleasePerformance(ctx context.Context, _ string) error {
	g.rlsCalls++

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return g.rlsErr
}

func (g *PerformanceGuard) ReleaseCalls() int {
	return g.rlsCalls
}
