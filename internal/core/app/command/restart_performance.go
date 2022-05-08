package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type RestartPerformance struct {
	PerformanceID string
	StartedByID   string
}

type RestartPerformanceHandler interface {
	Handle(ctx context.Context, cmd RestartPerformance) (<-chan service.Message, error)
}

type restartPerformanceHandler struct {
	perfRepo      service.PerformanceRepository
	specGetter    service.SpecificationGetter
	maintainer    service.PerformanceMaintainer
	performerOpts []performance.Option
}

func NewRestartPerformanceHandler(
	perfRepo service.PerformanceRepository,
	specGetter service.SpecificationGetter,
	maintainer service.PerformanceMaintainer,
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

	return restartPerformanceHandler{
		perfRepo:      perfRepo,
		maintainer:    maintainer,
		performerOpts: opts,
	}
}

func (h restartPerformanceHandler) Handle(
	ctx context.Context,
	cmd RestartPerformance,
) (messages <-chan service.Message, err error) {
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
