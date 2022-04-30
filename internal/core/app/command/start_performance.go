package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type StartPerformanceHandler struct {
	specRepo      app.SpecificationRepository
	perfRepo      app.PerformanceRepository
	maintainer    app.PerformanceMaintainer
	performerOpts []performance.Option
}

func NewStartPerformanceHandler(
	specRepo app.SpecificationRepository,
	perfRepo app.PerformanceRepository,
	maintainer app.PerformanceMaintainer,
	opts ...performance.Option,
) StartPerformanceHandler {
	if specRepo == nil {
		panic("specification repository is nil")
	}

	if perfRepo == nil {
		panic("performance repository is nil")
	}

	if maintainer == nil {
		panic("performance maintainer is nil")
	}

	return StartPerformanceHandler{
		specRepo:      specRepo,
		perfRepo:      perfRepo,
		maintainer:    maintainer,
		performerOpts: opts,
	}
}

func (h StartPerformanceHandler) Handle(
	ctx context.Context,
	cmd app.StartPerformanceCommand,
) (perfID string, messages <-chan app.Message, err error) {
	defer func() {
		err = errors.Wrap(err, "new performance starting")
	}()

	spec, err := h.specRepo.GetActiveSpecificationByTestCampaignID(ctx, cmd.TestCampaignID)
	if err != nil {
		return "", nil, err
	}

	if err = user.CanAccessSpecification(cmd.StartedByID, spec, user.Read); err != nil {
		return "", nil, err
	}

	perfID = uuid.New().String()

	perf := performance.FromSpecification(perfID, spec, h.performerOpts...)

	if err = h.perfRepo.AddPerformance(ctx, perf); err != nil {
		return "", nil, err
	}

	messages, err = h.maintainer.MaintainPerformance(context.Background(), perf)
	if err != nil {
		return "", nil, err
	}

	return perfID, messages, nil
}
