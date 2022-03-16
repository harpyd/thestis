package performance

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type Performer interface {
	Perform(
		ctx context.Context,
		env *Environment,
		thesis specification.Thesis,
	) Result
}

type Result struct {
	event Event
	err   error
}

func Pass() Result {
	return Result{event: FiredPass}
}

func Fail(err error) Result {
	if !IsFailedError(err) {
		err = NewFailedError(err)
	}

	return Result{
		event: FiredFail,
		err:   err,
	}
}

func Crash(err error) Result {
	if !IsCrashedError(err) {
		err = NewCrashedError(err)
	}

	return Result{
		event: FiredCrash,
		err:   err,
	}
}

func Cancel(err error) Result {
	if !IsCanceledError(err) {
		err = NewCanceledError(err)
	}

	return Result{
		event: FiredCancel,
		err:   err,
	}
}

func (r Result) Event() Event {
	return r.event
}

func (r Result) Err() error {
	return r.err
}

type PerformFunc func(
	ctx context.Context,
	env *Environment,
	thesis specification.Thesis,
) Result

func (f PerformFunc) Perform(
	ctx context.Context,
	env *Environment,
	thesis specification.Thesis,
) Result {
	return f(ctx, env, thesis)
}

func PassingPerformer() Performer {
	return PerformFunc(func(
		ctx context.Context,
		_ *Environment,
		_ specification.Thesis,
	) Result {
		if ctx.Err() != nil {
			return Cancel(ctx.Err())
		}

		return Pass()
	})
}

func FailingPerformer() Performer {
	return PerformFunc(func(
		ctx context.Context,
		_ *Environment,
		_ specification.Thesis,
	) Result {
		if ctx.Err() != nil {
			return Cancel(ctx.Err())
		}

		return Fail(errors.New("expected failing"))
	})
}

func CrashingPerformer() Performer {
	return PerformFunc(func(
		ctx context.Context,
		_ *Environment,
		_ specification.Thesis,
	) Result {
		if ctx.Err() != nil {
			return Cancel(ctx.Err())
		}

		return Crash(errors.New("expected crashing"))
	})
}

func CancelingPerformer() Performer {
	return PerformFunc(func(
		ctx context.Context,
		_ *Environment,
		_ specification.Thesis,
	) Result {
		if ctx.Err() != nil {
			return Cancel(ctx.Err())
		}

		return Cancel(errors.New("expected canceling"))
	})
}
