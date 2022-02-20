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
		panic("performances cancel publisher is nil")
	}

	return CancelPerformanceHandler{
		perfsRepo: perfsRepo,
		publisher: cancelPub,
	}
}

func (h CancelPerformanceHandler) Handle(ctx context.Context, cmd app.CancelPerformanceCommand) (err error) {
	defer func() {
		err = errors.Wrap(err, "performance cancelation")
	}()

	perf, err := h.perfsRepo.GetPerformance(ctx, cmd.PerformanceID)
	if err != nil {
		return err
	}

	if err := user.CanSeePerformance(cmd.CanceledByID, perf); err != nil {
		return err
	}

	if err := perf.MustBeStarted(); err != nil {
		return err
	}

	return h.publisher.PublishPerformanceCancel(ctx, perf.ID())
}
