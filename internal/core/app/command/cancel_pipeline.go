package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type CancelPipeline struct {
	PipelineID   string
	CanceledByID string
}

type CancelPipelineHandler interface {
	Handle(ctx context.Context, cmd CancelPipeline) error
}

type cancelPipelineHandler struct {
	pipeRepo  service.PipelineRepository
	publisher service.PipelineCancelPublisher
}

func NewCancelPipelineHandler(
	pipeRepo service.PipelineRepository,
	cancelPub service.PipelineCancelPublisher,
) CancelPipelineHandler {
	if pipeRepo == nil {
		panic("pipeline repository is nil")
	}

	if cancelPub == nil {
		panic("pipeline cancel publisher is nil")
	}

	return cancelPipelineHandler{
		pipeRepo:  pipeRepo,
		publisher: cancelPub,
	}
}

func (h cancelPipelineHandler) Handle(ctx context.Context, cmd CancelPipeline) (err error) {
	defer func() {
		err = errors.Wrap(err, "pipeline cancellation")
	}()

	pipe, err := h.pipeRepo.GetPipeline(ctx, cmd.PipelineID, service.WithoutSpecification())
	if err != nil {
		return err
	}

	if err := user.CanAccessPipeline(cmd.CanceledByID, pipe, user.Read); err != nil {
		return err
	}

	if err := pipe.ShouldBeStarted(); err != nil {
		return err
	}

	return h.publisher.PublishPipelineCancel(pipe.ID())
}
