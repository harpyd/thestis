package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/domain/user"
)

type CancelPerformanceHandler struct {
	perfRepo  app.PerformanceRepository
	publisher app.PerformanceCancelPublisher
}

func NewCancelPerformanceHandler(
	perfRepo app.PerformanceRepository,
	cancelPub app.PerformanceCancelPublisher,
) CancelPerformanceHandler {
	if perfRepo == nil {
		panic("performance repository is nil")
	}

	if cancelPub == nil {
		panic("performance cancel publisher is nil")
	}

	return CancelPerformanceHandler{
		perfRepo:  perfRepo,
		publisher: cancelPub,
	}
}

func (h CancelPerformanceHandler) Handle(ctx context.Context, cmd app.CancelPerformanceCommand) (err error) {
	defer func() {
		err = errors.Wrap(err, "performance cancellation")
	}()

	perf, err := h.perfRepo.GetPerformance(ctx, cmd.PerformanceID, app.WithoutSpecification())
	if err != nil {
		return err
	}

	if err := user.CanAccessPerformance(cmd.CanceledByID, perf, user.Read); err != nil {
		return err
	}

	if err := perf.ShouldBeStarted(); err != nil {
		return err
	}

	return h.publisher.PublishPerformanceCancel(perf.ID())
}
