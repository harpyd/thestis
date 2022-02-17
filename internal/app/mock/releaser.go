package mock

import "context"

type PerformanceReleaser func(ctx context.Context, perfID string) error

func (r PerformanceReleaser) ReleasePerformance(ctx context.Context, perfID string) error {
	return r(ctx, perfID)
}
