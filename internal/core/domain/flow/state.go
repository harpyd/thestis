package flow

import (
	"github.com/harpyd/thestis/internal/core/domain/performance"
)

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

// Precedence indicates the priority in which
// the overall state will be selected from the
// entire scenario states.
//
// The larger the number, the more likely it
// is that this state will be selected as an
// overall one.
func (s State) Precedence() int {
	switch s {
	case NoState:
		return 0
	case Passed:
		return 1
	case NotPerformed:
		return 2
	case Canceled:
		return 3
	case Failed:
		return 4
	case Crashed:
		return 5
	case Performing:
		return 6
	default:
		return 0
	}
}

func (s State) String() string {
	return string(s)
}
