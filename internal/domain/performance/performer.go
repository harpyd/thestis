package performance

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	// Performer carries performing of thesis.
	// Performance creators should provide own implementation.
	Performer interface {
		Perform(env *Environment, thesis specification.Thesis) Result
	}

	// Result presents a result of Performer work.
	//
	// It can be in four states:
	// NotPerformed, Passed, Failed, Crashed.
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
		err:   errors.WithStack(newFailedError(err)),
	}
}

func Crash(err error) Result {
	return Result{
		state: Crashed,
		err:   errors.WithStack(newCrashedError(err)),
	}
}

func (r Result) State() State {
	return r.state
}

func (r Result) Err() error {
	return r.err
}

type (
	failedError struct {
		err error
	}

	crashedError struct {
		err error
	}
)

func newFailedError(err error) error {
	if err == nil {
		return nil
	}

	return failedError{err: err}
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

	return crashedError{err: err}
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
