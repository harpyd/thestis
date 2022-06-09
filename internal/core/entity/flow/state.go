package flow

import (
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

type State string

const (
	NoState     State = ""
	NotExecuted State = "not executed"
	Executing   State = "executing"
	Passed      State = "passed"
	Failed      State = "failed"
	Crashed     State = "crashed"
	Canceled    State = "canceled"
)

type stateTransitionRules map[State]map[pipeline.Event]State

func rules() stateTransitionRules {
	return stateTransitionRules{
		NotExecuted: {
			pipeline.FiredExecute: Executing,
			pipeline.FiredPass:    Passed,
			pipeline.FiredFail:    Failed,
			pipeline.FiredCrash:   Crashed,
			pipeline.FiredCancel:  Canceled,
		},
		Executing: {
			pipeline.FiredPass:   Passed,
			pipeline.FiredFail:   Failed,
			pipeline.FiredCrash:  Crashed,
			pipeline.FiredCancel: Canceled,
		},
		Passed: {
			pipeline.FiredFail:   Failed,
			pipeline.FiredCrash:  Crashed,
			pipeline.FiredCancel: Passed,
		},
		Failed: {
			pipeline.FiredCrash:  Crashed,
			pipeline.FiredCancel: Failed,
		},
	}
}

// Next returns the next state depending on the
// received pipeline.Event.
//
// State transition rules are defined for a finite
// automaton, then this function is a transition
// function for this automaton.
func (s State) Next(with pipeline.Event) State {
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
	case NotExecuted:
		return 2
	case Canceled:
		return 3
	case Failed:
		return 4
	case Crashed:
		return 5
	case Executing:
		return 6
	default:
		return 0
	}
}

func (s State) String() string {
	return string(s)
}
