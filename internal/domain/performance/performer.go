package performance

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/specification"
)

// Event is fired due to the
// creation of the Result.
type Event string

const (
	NoEvent      Event = ""
	FiredPerform Event = "perform"
	FiredPass    Event = "pass"
	FiredFail    Event = "fail"
	FiredCrash   Event = "crash"
	FiredCancel  Event = "cancel"
)

func (e Event) String() string {
	return string(e)
}

// Performer performs the specified action
// using the passed thesis.
//
// Performance determines the progress of the
// specification pipeline, it performs each
// specification.Scenario that needs to be
// executed. To perform the scenario, you need
// to run the performing for each specification.Thesis
// and get the Result of performing this thesis.
// But Performance does not know how to perform
// the thesis, so it delegates this task to
// the Performer.
//
// The Performer can use one of these functions to
// return the Result:
// Pass, Fail, Crash, Cancel.
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

// Pass returns the passed Result.
//
// Pass should be used when the thesis
// is passed, with this result the scenario
// performing will continue.
func Pass() Result {
	return Result{event: FiredPass}
}

// Fail returns the failed Result with occurred error.
// If the passed error does not satisfy IsFailedError,
// it will be wrapped with NewFailedError.
//
// Fail should be used when the performing of the thesis
// has fallen due to natural reasons, for example, the
// assert specified in the thesis failed. With this
// result the scenario will be failed.
func Fail(err error) Result {
	if !IsFailedError(err) {
		err = NewFailedError(err)
	}

	return Result{
		event: FiredFail,
		err:   err,
	}
}

// Crash returns the crashed Result with occurred error.
// If the passed error does not satisfy IsCrashedError,
// it will be wrapped with NewCrashedError.
//
// Crash should be used when the performing of the thesis
// has fallen due to unforeseen circumstances, for example,
// problems with network interaction when performing the
// HTTP part of the thesis. With this result the scenario
// will be crashed.
func Crash(err error) Result {
	if !IsCrashedError(err) {
		err = NewCrashedError(err)
	}

	return Result{
		event: FiredCrash,
		err:   err,
	}
}

// Cancel returns the canceled Result with occurred error.
// If the passed error does not satisfy IsCanceledError,
// it will be wrapped with NewCanceledError.
//
// Cancel should be used when you need to mark a thesis
// as canceled, for example, when context.Context is done.
// With this result the scenario will be canceled.
func Cancel(err error) Result {
	if !IsCanceledError(err) {
		err = NewCanceledError(err)
	}

	return Result{
		event: FiredCancel,
		err:   err,
	}
}

// Event returns event of the Result.
func (r Result) Event() Event {
	return r.event
}

// Err returns occurred error of the Result.
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
