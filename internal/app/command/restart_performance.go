package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/user"
)

type RestartPerformanceHandler struct {
	perfsRepo     app.PerformancesRepository
	specGetter    app.SpecificationGetter
	maintainer    app.PerformanceMaintainer
	performerOpts app.PerformerOptions
}

func NewRestartPerformanceHandler(
	perfsRepo app.PerformancesRepository,
	specGetter app.SpecificationGetter,
	maintainer app.PerformanceMaintainer,
	opts ...app.PerformerOption,
) RestartPerformanceHandler {
	if perfsRepo == nil {
		panic("performances repository is nil")
	}

	if specGetter == nil {
		panic("specification getter is nil")
	}

	if maintainer == nil {
		panic("performance maintainer is nil")
	}

	return RestartPerformanceHandler{
		perfsRepo:     perfsRepo,
		maintainer:    maintainer,
		performerOpts: opts,
	}
}

func (h RestartPerformanceHandler) Handle(
	ctx context.Context,
	cmd app.RestartPerformanceCommand,
) (messages <-chan app.Message, err error) {
	defer func() {
		err = errors.Wrap(err, "performance restarting")
	}()

	perf, err := h.perfsRepo.GetPerformance(ctx, cmd.PerformanceID, h.specGetter, h.performerOpts...)
	if err != nil {
		return nil, err
	}

	if err = user.CanSeePerformance(cmd.StartedByID, perf); err != nil {
		return nil, err
	}

	messages, err = h.maintainer.MaintainPerformance(context.Background(), perf)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
