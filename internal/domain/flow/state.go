package flow

import "github.com/harpyd/thestis/internal/domain/performance"

type State string

const (
	NoState      State = ""
	NotPerformed State = "not performed"
	Performing   State = "performing"
	Passed       State = "passed"
	Failed       State = "failed"
	Crashed      State = "crashed"
	Canceled     State = "canceled"
)

type stateTransitionRules map[State]map[performance.Event]State

func rules() stateTransitionRules {
	return stateTransitionRules{
		NotPerformed: {
			performance.FiredPerform: Performing,
			performance.FiredPass:    Passed,
			performance.FiredFail:    Failed,
			performance.FiredCrash:   Crashed,
			performance.FiredCancel:  Canceled,
		},
		Performing: {
			performance.FiredPerform: Performing,
			performance.FiredPass:    Passed,
			performance.FiredFail:    Failed,
			performance.FiredCrash:   Crashed,
			performance.FiredCancel:  Canceled,
		},
		Passed: {
			performance.FiredPerform: Passed,
			performance.FiredPass:    Passed,
			performance.FiredFail:    Failed,
			performance.FiredCrash:   Crashed,
			performance.FiredCancel:  Passed,
		},
		Failed: {
			performance.FiredPerform: Failed,
			performance.FiredPass:    Failed,
			performance.FiredFail:    Failed,
			performance.FiredCrash:   Crashed,
			performance.FiredCancel:  Failed,
		},
		Crashed: {
			performance.FiredPerform: Crashed,
			performance.FiredPass:    Crashed,
			performance.FiredFail:    Crashed,
			performance.FiredCrash:   Crashed,
			performance.FiredCancel:  Crashed,
		},
		Canceled: {
			performance.FiredPerform: Canceled,
			performance.FiredPass:    Canceled,
			performance.FiredFail:    Canceled,
			performance.FiredCrash:   Canceled,
			performance.FiredCancel:  Canceled,
		},
	}
}

func (s State) Next(with performance.Event) State {
	ss, ok := rules()[s]
	if !ok {
		return NoState
	}

	res, ok := ss[with]
	if !ok {
		return NoState
	}

	return res
}

func (s State) String() string {
	return string(s)
}
