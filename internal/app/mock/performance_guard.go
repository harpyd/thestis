package mock

import "context"

type PerformanceGuard struct {
	acqErr error
	rlsErr error
}

func NewPerformanceGuard(acquireErr error, releaseErr error) PerformanceGuard {
	return PerformanceGuard{
		acqErr: acquireErr,
		rlsErr: releaseErr,
	}
}

func (g PerformanceGuard) AcquirePerformance(ctx context.Context, _ string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return g.acqErr
}

func (g PerformanceGuard) ReleasePerformance(ctx context.Context, _ string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return g.rlsErr
}
