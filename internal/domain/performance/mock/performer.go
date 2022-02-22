package mock

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type Performer func(
	ctx context.Context,
	env *performance.Environment,
	t specification.Thesis,
) performance.Result

func (p Performer) Perform(
	ctx context.Context,
	env *performance.Environment,
	t specification.Thesis,
) performance.Result {
	return p(ctx, env, t)
}

func NewPassingPerformer() Performer {
	return Performer(func(
		_ context.Context,
		_ *performance.Environment,
		_ specification.Thesis,
	) performance.Result {
		return performance.Pass()
	})
}

func NewFailingPerformer() Performer {
	return Performer(func(
		_ context.Context,
		_ *performance.Environment,
		_ specification.Thesis,
	) performance.Result {
		return performance.Fail(errors.New("something wrong"))
	})
}

func NewCrashingPerformer() Performer {
	return Performer(func(
		_ context.Context,
		_ *performance.Environment,
		_ specification.Thesis,
	) performance.Result {
		return performance.Crash(errors.New("something wrong"))
	})
}
