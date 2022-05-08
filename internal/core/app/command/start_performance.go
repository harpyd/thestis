package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type StartPerformance struct {
	TestCampaignID string
	StartedByID    string
}

type StartPerformanceHandler interface {
	Handle(ctx context.Context, cmd StartPerformance) (string, <-chan service.Message, error)
}

type startPerformanceHandler struct {
	specRepo      service.SpecificationRepository
	perfRepo      service.PerformanceRepository
	maintainer    service.PerformanceMaintainer
	performerOpts []performance.Option
}

func NewStartPerformanceHandler(
	specRepo service.SpecificationRepository,
	perfRepo service.PerformanceRepository,
	maintainer service.PerformanceMaintainer,
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

	return startPerformanceHandler{
		specRepo:      specRepo,
		perfRepo:      perfRepo,
		maintainer:    maintainer,
		performerOpts: opts,
	}
}

func (h startPerformanceHandler) Handle(
	ctx context.Context,
	cmd StartPerformance,
) (perfID string, messages <-chan service.Message, err error) {
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
