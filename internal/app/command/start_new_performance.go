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
	specsRepo app.SpecificationsRepository
	perfsRepo app.PerformancesRepository
	flowsRepo app.FlowsRepository
}

func NewStartPerformanceHandler(
	specsRepo app.SpecificationsRepository,
	perfsRepo app.PerformancesRepository,
	flowsRepo app.FlowsRepository,
) StartNewPerformanceHandler {
	if specsRepo == nil {
		panic("specification repository is nil")
	}

	if perfsRepo == nil {
		panic("performances repository is nil")
	}

	if flowsRepo == nil {
		panic("flows repository is nil")
	}

	return StartNewPerformanceHandler{
		specsRepo: specsRepo,
		perfsRepo: perfsRepo,
		flowsRepo: flowsRepo,
	}
}

func (h StartNewPerformanceHandler) Handle(
	ctx context.Context,
	cmd app.StartNewPerformanceCommand,
) (perfID string, messages <-chan app.Message, err error) {
	defer func() {
		err = errors.New("new performance starting")
	}()

	spec, err := h.specsRepo.GetActiveSpecificationByTestCampaignID(ctx, cmd.TestCampaignID)
	if err != nil {
		return "", nil, err
	}

	if err = user.CanSeeSpecification(cmd.StartedByID, spec); err != nil {
		return "", nil, err
	}

	perfID = uuid.New().String()

	perf, err := performance.FromSpecification(spec, performance.WithID(perfID))
	if err != nil {
		return "", nil, err
	}

	if err = h.perfsRepo.AddPerformance(ctx, perf); err != nil {
		return "", nil, err
	}

	actionCtx := context.Background()

	steps, err := perf.Start(actionCtx)
	if err != nil {
		return "", nil, err
	}

	msg := make(chan app.Message)

	if err = h.perfsRepo.ExclusivelyDoWithPerformance(
		actionCtx, perf,
		h.performanceAction(actionCtx, steps, msg),
	); err != nil {
		return "", nil, err
	}

	return perfID, msg, nil
}

func (h StartNewPerformanceHandler) performanceAction(
	ctx context.Context,
	steps <-chan performance.Step,
	msg chan<- app.Message,
) func(perf *performance.Performance) {
	return func(perf *performance.Performance) {
		defer close(msg)

		fr := performance.FlowFromPerformance(perf)

		for s := range steps {
			fr.WithStep(s)

			flow := fr.Reduce()
			if err := h.flowsRepo.UpsertFlow(ctx, flow); err != nil {
				msg <- app.NewMessageFromError(err)
			}

			msg <- app.NewMessageFromStringer(s)
		}

		flow := fr.FinallyReduce()
		if err := h.flowsRepo.UpsertFlow(ctx, flow); err != nil {
			msg <- app.NewMessageFromError(err)
		}
	}
}
