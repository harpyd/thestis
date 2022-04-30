package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type RestartPerformanceHandler struct {
	perfRepo      app.PerformanceRepository
	specGetter    app.SpecificationGetter
	maintainer    app.PerformanceMaintainer
	performerOpts []performance.Option
}

func NewRestartPerformanceHandler(
	perfRepo app.PerformanceRepository,
	specGetter app.SpecificationGetter,
	maintainer app.PerformanceMaintainer,
	opts ...performance.Option,
) RestartPerformanceHandler {
	if perfRepo == nil {
		panic("performance repository is nil")
	}

	if specGetter == nil {
		panic("specification getter is nil")
	}

	if maintainer == nil {
		panic("performance maintainer is nil")
	}

	return RestartPerformanceHandler{
		perfRepo:      perfRepo,
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

	perf, err := h.perfRepo.GetPerformance(ctx, cmd.PerformanceID, h.specGetter, h.performerOpts...)
	if err != nil {
		return nil, err
	}

	if err = user.CanAccessPerformance(cmd.StartedByID, perf, user.Read); err != nil {
		return nil, err
	}

	messages, err = h.maintainer.MaintainPerformance(context.Background(), perf)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
