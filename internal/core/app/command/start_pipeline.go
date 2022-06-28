package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type StartPipeline struct {
	PipelineID     string
	TestCampaignID string
	StartedByID    string
}

type StartPipelineHandler interface {
	Handle(ctx context.Context, cmd StartPipeline) error
}

type startPipelineHandler struct {
	specRepo   service.SpecificationRepository
	pipeRepo   service.PipelineRepository
	maintainer service.PipelineMaintainer
	registrars []pipeline.ExecutorRegistrar
}

func NewStartPipelineHandler(
	specRepo service.SpecificationRepository,
	pipeRepo service.PipelineRepository,
	maintainer service.PipelineMaintainer,
	registrars ...pipeline.ExecutorRegistrar,
) StartPipelineHandler {
	if specRepo == nil {
		panic("specification repository is nil")
	}

	if pipeRepo == nil {
		panic("pipeline repository is nil")
	}

	if maintainer == nil {
		panic("pipeline maintainer is nil")
	}

	return startPipelineHandler{
		specRepo:   specRepo,
		pipeRepo:   pipeRepo,
		maintainer: maintainer,
		registrars: registrars,
	}
}

func (h startPipelineHandler) Handle(ctx context.Context, cmd StartPipeline) (err error) {
	defer func() {
		err = errors.Wrap(err, "new pipeline starting")
	}()

	spec, err := h.specRepo.GetActiveSpecificationByTestCampaignID(ctx, cmd.TestCampaignID)
	if err != nil {
		return err
	}

	if err := user.CanAccessSpecification(cmd.StartedByID, spec, user.Read); err != nil {
		return err
	}

	pipe := pipeline.Trigger(cmd.PipelineID, spec, h.registrars...)

	if err := h.pipeRepo.AddPipeline(ctx, pipe); err != nil {
		return err
	}

	_, err = h.maintainer.MaintainPipeline(ctx, pipe)

	return err
}
