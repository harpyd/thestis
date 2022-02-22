package performance

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	// Performer carries performing of thesis.
	// Performance creators should provide own implementation.
	Performer interface {
		Perform(
			ctx context.Context,
			env *Environment,
			thesis specification.Thesis,
		) Result
	}

	// Result presents a result of Performer work.
	//
	// It can be in five states:
	// NotPerformed, Passed, Failed,
	// Crashed, Canceled.
	Result struct {
		state State
		err   error
	}
)

func NotPerform() Result {
	return Result{
		state: NotPerformed,
	}
}

func Pass() Result {
	return Result{
		state: Passed,
	}
}

func Fail(err error) Result {
	return Result{
		state: Failed,
		err:   newFailedError(err),
	}
}

func Crash(err error) Result {
	return Result{
		state: Crashed,
		err:   newCrashedError(err),
	}
}

func Cancel(err error) Result {
	return Result{
		state: Canceled,
		err:   newCanceledError(err),
	}
}

func (r Result) State() State {
	return r.state
}

func (r Result) Err() error {
	return r.err
}

type (
	canceledError struct {
		err error
	}

	failedError struct {
		err error
	}

	crashedError struct {
		err error
	}
)

func newCanceledError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(canceledError{err: err})
}

func IsCanceledError(err error) bool {
	var target canceledError

	return errors.As(err, &target)
}

func (e canceledError) Error() string {
	return fmt.Sprintf("performance canceled: %s", e.err)
}

func (e canceledError) Cause() error {
	return e.err
}

func (e canceledError) Unwrap() error {
	return e.err
}

func newFailedError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(failedError{err: err})
}

func IsFailedError(err error) bool {
	var target failedError

	return errors.As(err, &target)
}

func (e failedError) Error() string {
	return fmt.Sprintf("performance failed: %s", e.err)
}

func (e failedError) Cause() error {
	return e.err
}

func (e failedError) Unwrap() error {
	return e.err
}

func newCrashedError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(crashedError{err: err})
}

func IsCrashedError(err error) bool {
	var target crashedError

	return errors.As(err, &target)
}

func (e crashedError) Error() string {
	return fmt.Sprintf("performance crashed: %s", e.err)
}

func (e crashedError) Cause() error {
	return e.err
}

func (e crashedError) Unwrap() error {
	return e.err
}
