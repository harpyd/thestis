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
			performance.FiredPass:   Passed,
			performance.FiredFail:   Failed,
			performance.FiredCrash:  Crashed,
			performance.FiredCancel: Canceled,
		},
		Passed: {
			performance.FiredFail:   Failed,
			performance.FiredCrash:  Crashed,
			performance.FiredCancel: Passed,
		},
		Failed: {
			performance.FiredCrash:  Crashed,
			performance.FiredCancel: Failed,
		},
	}
}

// Next returns the next state depending on the
// received performance.Event.
//
// State transition rules are defined for a finite
// automaton, then this function is a transition
// function for this automaton.
func (s State) Next(with performance.Event) State {
	ss, ok := rules()[s]
	if !ok {
		return s
	}

	res, ok := ss[with]
	if !ok {
		return s
	}

	return res
}

func (s State) String() string {
	return string(s)
}
