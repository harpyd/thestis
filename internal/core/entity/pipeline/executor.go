package pipeline

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/entity/specification"
)

// Event is fired due to the
// creation of the Result.
type Event string

const (
	NoEvent      Event = ""
	FiredExecute Event = "execute"
	FiredPass    Event = "pass"
	FiredFail    Event = "fail"
	FiredCrash   Event = "crash"
	FiredCancel  Event = "cancel"
)

func (e Event) String() string {
	return string(e)
}

// Executor executes the specified action
// using the passed thesis.
//
// Pipeline determines the progress of the
// specification pipeline, it runs each
// specification.Scenario that needs to be
// executed. To perform the scenario, you need
// to run each specification.Thesis and get
// the Result of performing this thesis.
// But Pipeline does not know how to execute
// the thesis, so it delegates this task to
// the Executor.
//
// The Executor can use one of these functions to
// return the Result:
// Pass, Fail, Crash, Cancel.
type Executor interface {
	Execute(
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
// executing will continue.
func Pass() Result {
	return Result{event: FiredPass}
}

// Fail returns the failed Result with occurred error.
// If the passed error is not equal to TerminatedError
// with FiredFail event, it will be wrapped with failed
// TerminatedError.
//
// Fail should be used when the executing of the thesis
// has fallen due to natural reasons, for example, the
// assertion specified in the thesis failed. With this
// result the scenario will be failed.
func Fail(err error) Result {
	var terr *TerminatedError

	if !errors.As(err, &terr) || terr.Event() != FiredFail {
		err = WrapWithTerminatedError(err, FiredFail)
	}

	return Result{
		event: FiredFail,
		err:   err,
	}
}

// Crash returns the crashed Result with occurred error.
// If the passed error is not equal to TerminatedError
// with FiredCrash event, it will be wrapped with crashed
// TerminatedError.
//
// Crash should be used when the executing of the thesis
// has fallen due to unforeseen circumstances, for example,
// problems with network interaction when executing the
// HTTP part of the thesis. With this result the scenario
// will be crashed.
func Crash(err error) Result {
	var terr *TerminatedError

	if !errors.As(err, &terr) || terr.Event() != FiredCrash {
		err = WrapWithTerminatedError(err, FiredCrash)
	}

	return Result{
		event: FiredCrash,
		err:   err,
	}
}

// Cancel returns the canceled Result with occurred error.
// If the passed error is not equal to TerminatedError
// with FiredCancel event, it will be wrapped with canceled
// TerminatedError.
//
// Cancel should be used when you need to mark a thesis
// as canceled, for example, when context.Context is done.
// With this result the scenario will be canceled.
func Cancel(err error) Result {
	var terr *TerminatedError

	if !errors.As(err, &terr) || terr.Event() != FiredCancel {
		err = WrapWithTerminatedError(err, FiredCancel)
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

// ExecutorFunc is an adapter
// to allow the use of ordinary
// functions as Executor.
type ExecutorFunc func(
	ctx context.Context,
	env *Environment,
	thesis specification.Thesis,
) Result

func (f ExecutorFunc) Execute(
	ctx context.Context,
	env *Environment,
	thesis specification.Thesis,
) Result {
	return f(ctx, env, thesis)
}

// PassingExecutor is a shortcut for create
// naive implementation of the Executor that
// constantly returns the passed Result. If
// the context is done returns the canceled
// Result with a context error.
// Does nothing else.
//
// It's good to use for testing and mocking.
func PassingExecutor() Executor {
	return ExecutorFunc(func(
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

// FailingExecutor is a shortcut for create
// naive implementation of the Executor that
// constantly returns the failed Result with
// "expected failing" error. If the context is
// done returns the canceled Result with a
// context error.
// Does nothing else.
//
// It's good to use for testing and mocking.
func FailingExecutor() Executor {
	return ExecutorFunc(func(
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

// CrashingExecutor is a shortcut for create
// naive implementation of the Executor that
// constantly returns the crashed Result with
// "expected crashing" error. If the context is
// done returns the canceled Result with a
// context error.
// Does nothing else.
//
// It's good to use for testing and mocking.
func CrashingExecutor() Executor {
	return ExecutorFunc(func(
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

// CancelingExecutor is a shortcut for create
// naive implementation of the Executor that
// constantly returns the canceled Result with
// "expected canceling" error. If the context is
// done returns the canceled Result with a
// context error.
// Does nothing else.
//
// It's good to use for testing and mocking.
func CancelingExecutor() Executor {
	return ExecutorFunc(func(
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
