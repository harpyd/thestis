package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

type PipelineMaintainer struct {
	withErr bool
}

func NewPipelineMaintainer(withErr bool) PipelineMaintainer {
	return PipelineMaintainer{withErr: withErr}
}

func (m PipelineMaintainer) MaintainPipeline(
	_ context.Context,
	_ *pipeline.Pipeline,
) (<-chan service.DoneSignal, error) {
	if m.withErr {
		return nil, pipeline.ErrAlreadyStarted
	}

	done := make(chan service.DoneSignal)
	close(done)

	return done, nil
}
