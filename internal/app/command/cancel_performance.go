package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/user"
)

type CancelPerformanceHandler struct {
	perfsRepo app.PerformancesRepository
	publisher app.PerformanceCancelPublisher
}

func NewCancelPerformanceHandler(
	perfsRepo app.PerformancesRepository,
	cancelPub app.PerformanceCancelPublisher,
) CancelPerformanceHandler {
	if perfsRepo == nil {
		panic("performances repository is nil")
	}

	if cancelPub == nil {
		panic("performance cancel publisher is nil")
	}

	return CancelPerformanceHandler{
		perfsRepo: perfsRepo,
		publisher: cancelPub,
	}
}

func (h CancelPerformanceHandler) Handle(ctx context.Context, cmd app.CancelPerformanceCommand) (err error) {
	defer func() {
		err = errors.Wrap(err, "performance cancellation")
	}()

	perf, err := h.perfsRepo.GetPerformance(ctx, cmd.PerformanceID, app.DontGetSpecification())
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
