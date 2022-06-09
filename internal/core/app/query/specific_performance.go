package query

import (
	"context"

	"github.com/pkg/errors"
)

type SpecificPipeline struct {
	PipelineID string
	UserID     string
}

type SpecificPipelineHandler interface {
	Handle(ctx context.Context, qry SpecificPipeline) (SpecificPipelineModel, error)
}

type SpecificPipelineReadModel interface {
	FindPipeline(ctx context.Context, qry SpecificPipeline) (SpecificPipelineModel, error)
}

type specificPipelineHandler struct {
	readModel SpecificPipelineReadModel
}

func NewSpecificPipelineHandler(readModel SpecificPipelineReadModel) SpecificPipelineHandler {
	if readModel == nil {
		panic("specific pipeline read model is nil")
	}

	return specificPipelineHandler{
		readModel: readModel,
	}
}

func (h specificPipelineHandler) Handle(
	ctx context.Context,
	qry SpecificPipeline,
) (SpecificPipelineModel, error) {
	pipe, err := h.readModel.FindPipeline(ctx, qry)

	return pipe, errors.Wrap(err, "getting pipeline")
}
