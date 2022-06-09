package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type RestartPipeline struct {
	PipelineID  string
	StartedByID string
}

type RestartPipelineHandler interface {
	Handle(ctx context.Context, cmd RestartPipeline) error
}

type restartPipelineHandler struct {
	pipeRepo   service.PipelineRepository
	specGetter service.SpecificationGetter
	maintainer service.PipelineMaintainer
	opts       []pipeline.Option
}

func NewRestartPipelineHandler(
	pipeRepo service.PipelineRepository,
	specGetter service.SpecificationGetter,
	maintainer service.PipelineMaintainer,
	opts ...pipeline.Option,
) RestartPipelineHandler {
	if pipeRepo == nil {
		panic("pipeline repository is nil")
	}

	if specGetter == nil {
		panic("specification getter is nil")
	}

	if maintainer == nil {
		panic("pipeline maintainer is nil")
	}

	return restartPipelineHandler{
		pipeRepo:   pipeRepo,
		maintainer: maintainer,
		opts:       opts,
	}
}

func (h restartPipelineHandler) Handle(
	ctx context.Context,
	cmd RestartPipeline,
) (err error) {
	defer func() {
		err = errors.Wrap(err, "pipeline restarting")
	}()

	pipe, err := h.pipeRepo.GetPipeline(ctx, cmd.PipelineID, h.specGetter, h.opts...)
	if err != nil {
		return err
	}

	if err := user.CanAccessPipeline(cmd.StartedByID, pipe, user.Read); err != nil {
		return err
	}

	_, err = h.maintainer.MaintainPipeline(ctx, pipe)

	return err
}
