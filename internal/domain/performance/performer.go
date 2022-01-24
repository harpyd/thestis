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
	wrappingFailedError struct {
		err error
	}

	failedErr struct {
		s string
	}

	wrappingCrashedError struct {
		err error
	}

	crashedError struct {
		s string
	}
)

func newTransitionError(state State, errMsg string) error {
	if errMsg == "" {
		return nil
	}

	if state == Failed {
		return failedErr{s: errMsg}
	} else if state == Crashed {
		return crashedError{s: errMsg}
	}

	return errors.New(errMsg)
}

func newFailedError(err error) error {
	if err == nil {
		return nil
	}

	return wrappingFailedError{err: err}
}

func IsFailedError(err error) bool {
	var (
		werr wrappingFailedError
		ferr failedErr
	)

	return errors.As(err, &werr) || errors.As(err, &ferr)
}

func (e wrappingFailedError) Error() string {
	return fmt.Sprintf("performance failed: %s", e.err)
}

func (e failedErr) Error() string {
	return e.s
}

func (e wrappingFailedError) Cause() error {
	return e.err
}

func (e wrappingFailedError) Unwrap() error {
	return e.err
}

func newCrashedError(err error) error {
	if err == nil {
		return nil
	}

	return wrappingCrashedError{err: err}
}

func IsCrashedError(err error) bool {
	var (
		werr wrappingCrashedError
		cerr crashedError
	)

	return errors.As(err, &werr) || errors.As(err, &cerr)
}

func (e wrappingCrashedError) Error() string {
	return fmt.Sprintf("performance crashed: %s", e.err)
}

func (e crashedError) Error() string {
	return e.s
}

func (e wrappingCrashedError) Cause() error {
	return e.err
}

func (e wrappingCrashedError) Unwrap() error {
	return e.err
}
