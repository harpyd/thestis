package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/user"
)

type StartPerformanceHandler struct {
	specsRepo     app.SpecificationsRepository
	perfsRepo     app.PerformancesRepository
	maintainer    app.PerformanceMaintainer
	performerOpts app.PerformerOptions
}

func NewStartPerformanceHandler(
	specsRepo app.SpecificationsRepository,
	perfsRepo app.PerformancesRepository,
	maintainer app.PerformanceMaintainer,
	opts ...app.PerformerOption,
) StartPerformanceHandler {
	if specsRepo == nil {
		panic("specifications repository is nil")
	}

	if perfsRepo == nil {
		panic("performances repository is nil")
	}

	if maintainer == nil {
		panic("performance maintainer is nil")
	}

	return StartPerformanceHandler{
		specsRepo:     specsRepo,
		perfsRepo:     perfsRepo,
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

	spec, err := h.specsRepo.GetActiveSpecificationByTestCampaignID(ctx, cmd.TestCampaignID)
	if err != nil {
		return "", nil, err
	}

	if err = user.CanAccessSpecification(cmd.StartedByID, spec, user.Read); err != nil {
		return "", nil, err
	}

	perfID = uuid.New().String()

	opts := h.performerOpts.ToPerformanceOptions()

	perf := performance.FromSpecification(perfID, spec, opts...)

	if err = h.perfsRepo.AddPerformance(ctx, perf); err != nil {
		return "", nil, err
	}

	messages, err = h.maintainer.MaintainPerformance(context.Background(), perf)
	if err != nil {
		return "", nil, err
	}

	return perfID, messages, nil
}
