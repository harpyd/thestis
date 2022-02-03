package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/user"
)

type RestartPerformanceHandler struct {
	perfsRepo     app.PerformancesRepository
	manager       app.FlowManager
	performerOpts app.PerformerOptions
}

func NewRestartPerformanceHandler(
	perfsRepo app.PerformancesRepository,
	manager app.FlowManager,
	opts ...app.PerformerOption,
) RestartPerformanceHandler {
	if perfsRepo == nil {
		panic("performances repository is nil")
	}

	if manager == nil {
		panic("flow manager is nil")
	}

	return RestartPerformanceHandler{
		perfsRepo:     perfsRepo,
		manager:       manager,
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

	perf, err := h.perfsRepo.GetPerformance(ctx, cmd.PerformanceID)
	if err != nil {
		return nil, err
	}

	if err = user.CanSeePerformance(cmd.StartedByID, perf); err != nil {
		return nil, err
	}

	messages, err = h.manager.ManageFlow(context.Background(), perf)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
