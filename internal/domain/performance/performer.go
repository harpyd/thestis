package performance

import "github.com/harpyd/thestis/internal/domain/specification"

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
		fail  error
		crash error
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
		fail:  err,
	}
}

func Crash(err error) Result {
	return Result{
		state: Crashed,
		crash: err,
	}
}

func (r Result) State() State {
	return r.state
}

func (r Result) FailErr() error {
	return r.fail
}

func (r Result) CrashErr() error {
	return r.crash
}
