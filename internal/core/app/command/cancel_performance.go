package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type CancelPerformanceHandler struct {
	perfRepo  service.PerformanceRepository
	publisher service.PerformanceCancelPublisher
}

func NewCancelPerformanceHandler(
	perfRepo service.PerformanceRepository,
	cancelPub service.PerformanceCancelPublisher,
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

	perf, err := h.perfRepo.GetPerformance(ctx, cmd.PerformanceID, service.WithoutSpecification())
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
