package query

import (
	"context"

	"github.com/pkg/errors"
)

type Pipeline struct {
	PipelineID string
	UserID     string
}

type PipelineHandler interface {
	Handle(ctx context.Context, qry Pipeline) (PipelineModel, error)
}

type PipelineReadModel interface {
	FindPipeline(ctx context.Context, qry Pipeline) (PipelineModel, error)
}

type pipelineHandler struct {
	readModel PipelineReadModel
}

func NewPipelineHandler(readModel PipelineReadModel) PipelineHandler {
	if readModel == nil {
		panic("pipeline read model is nil")
	}

	return pipelineHandler{
		readModel: readModel,
	}
}

func (h pipelineHandler) Handle(
	ctx context.Context,
	qry Pipeline,
) (PipelineModel, error) {
	pipe, err := h.readModel.FindPipeline(ctx, qry)

	return pipe, errors.Wrap(err, "getting pipeline")
}
