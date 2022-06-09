package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/core/entity/performance"
)

type PerformancePolicy struct {
	consumeCalls int
}

func NewPerformancePolicy() *PerformancePolicy {
	return &PerformancePolicy{}
}

func (p *PerformancePolicy) ConsumePerformance(
	ctx context.Context,
	perf *performance.Performance,
) {
	p.consumeCalls++

	steps := perf.MustStart(ctx)

	for range steps {
	}
}
