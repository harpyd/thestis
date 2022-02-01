package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/user"
)

type StartNewPerformanceHandler struct {
	manager       app.FlowManager
	specsRepo     app.SpecificationsRepository
	perfsRepo     app.PerformancesRepository
	performerOpts app.PerformerOptions
}

func NewStartPerformanceHandler(
	manager app.FlowManager,
	specsRepo app.SpecificationsRepository,
	perfsRepo app.PerformancesRepository,
	opts ...app.PerformerOption,
) StartNewPerformanceHandler {
	if manager == nil {
		panic("flow manager is nil")
	}

	if specsRepo == nil {
		panic("specification repository is nil")
	}

	if perfsRepo == nil {
		panic("performances repository is nil")
	}

	return StartNewPerformanceHandler{
		manager:       manager,
		specsRepo:     specsRepo,
		perfsRepo:     perfsRepo,
		performerOpts: opts,
	}
}

func (h StartNewPerformanceHandler) Handle(
	ctx context.Context,
	cmd app.StartNewPerformanceCommand,
) (perfID string, messages <-chan app.Message, err error) {
	defer func() {
		err = errors.Wrap(err, "new performance starting")
	}()

	spec, err := h.specsRepo.GetActiveSpecificationByTestCampaignID(ctx, cmd.TestCampaignID)
	if err != nil {
		return "", nil, err
	}

	if err = user.CanSeeSpecification(cmd.StartedByID, spec); err != nil {
		return "", nil, err
	}

	perfID = uuid.New().String()

	opts := append(h.performerOpts.ToPerformanceOptions(), performance.WithID(perfID))

	perf, err := performance.FromSpecification(spec, opts...)
	if err != nil {
		return "", nil, err
	}

	if err = h.perfsRepo.AddPerformance(ctx, perf); err != nil {
		return "", nil, err
	}

	actionCtx := context.Background()

	messages, err = h.manager.ManageFlow(actionCtx, perf)
	if err != nil {
		return "", nil, err
	}

	return perfID, messages, nil
}
