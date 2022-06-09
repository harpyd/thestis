package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

type PipelinePolicy struct {
	consumeCalls int
}

func NewPipelinePolicy() *PipelinePolicy {
	return &PipelinePolicy{}
}

func (p *PipelinePolicy) ConsumePipeline(
	ctx context.Context,
	pipe *pipeline.Pipeline,
) {
	p.consumeCalls++

	steps := pipe.MustStart(ctx)

	for range steps {
	}
}
