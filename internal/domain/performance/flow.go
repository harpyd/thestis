package performance

import (
	"fmt"

	"go.uber.org/multierr"
)

type (
	// Step is one unit of information
	// about Performance performing.
	Step interface {
		// FromTo returns transition from and to vertexes.
		// If step has transition, ok == true. Else ok == false.
		// For example, cancel step has no transition.
		FromTo() (from, to string, ok bool)
		State() State
		Err() error
		Fail() error
		String() string
	}

	// Flow represents current Performance performing.
	// Flow keeps transitions information
	// and common state of performing.
	Flow struct {
		state State
		graph map[string]map[string]Transition
	}

	Transition struct {
		from  string
		to    string
		state State
		err   error
		fail  error
	}

	// FlowBuilder builds Flow instance using WithStep,
	// Build and FinallyBuild methods.
	//
	// FlowBuilder defines Flow common state transition rules
	// in WithStep method.
	FlowBuilder struct {
		state                        State
		graph                        map[string]map[string]*Transition
		commonStateTransitionRules   stateTransitionRules
		specificStateTransitionRules stateTransitionRules
	}
)

func (f Flow) State() State {
	return f.state
}

func (f Flow) Transitions() []Transition {
	transitions := make([]Transition, 0, len(f.graph))

	for _, ts := range f.graph {
		for _, t := range ts {
			transitions = append(transitions, t)
		}
	}

	return transitions
}

func (t Transition) From() string {
	return t.from
}

func (t Transition) To() string {
	return t.to
}

func (t Transition) State() State {
	return t.state
}

func (t Transition) Err() error {
	return t.err
}

func (t Transition) Fail() error {
	return t.fail
}

func NewFlowBuilder(perf *Performance) *FlowBuilder {
	graph := make(map[string]map[string]*Transition, len(perf.actionGraph))

	for from, as := range perf.actionGraph {
		graph[from] = make(map[string]*Transition, len(as))

		for to := range as {
			graph[from][to] = &Transition{
				from:  from,
				to:    to,
				state: NotPerformed,
			}
		}
	}

	return &FlowBuilder{
		state:                        NotPerformed,
		graph:                        graph,
		commonStateTransitionRules:   newCommonStateTransitionRules(),
		specificStateTransitionRules: newSpecificStateTransitionRules(),
	}
}

// Build creates intermediate version of Flow from FlowBuilder.
// Flow created using this method can't have a State equal
// to Passed.
//
// You need to use FinallyBuild method to create
// final version of Flow.
func (b *FlowBuilder) Build() Flow {
	return Flow{
		state: b.state,
		graph: b.copyGraph(),
	}
}

// FinallyBuild creates final version of Flow from FlowBuilder.
// Flow created using this method represents final Performance's Flow.
//
// If all transitions have Passed states and common State equal
// to Performing, final version of Flow will have Passed State.
func (b *FlowBuilder) FinallyBuild() Flow {
	return Flow{
		state: b.finalState(),
		graph: b.copyGraph(),
	}
}

func (b *FlowBuilder) finalState() State {
	if b.state == Performing && b.allPassed() {
		return Passed
	}

	return b.state
}

func (b *FlowBuilder) allPassed() bool {
	for _, ts := range b.graph {
		for _, t := range ts {
			if t.State() != Passed {
				return false
			}
		}
	}

	return true
}

func (b *FlowBuilder) copyGraph() map[string]map[string]Transition {
	graph := make(map[string]map[string]Transition, len(b.graph))

	for from, ts := range b.graph {
		graph[from] = make(map[string]Transition, len(ts))

		for to, t := range ts {
			graph[from][to] = *t
		}
	}

	return graph
}

// WithStep is method for step by step building of Flow.
// Also, this method defines Flow common state transition rules.
//
// Rules:
// (NotPerformed -> NotPerformed) => NotPerformed
// (NotPerformed -> Performing) => Performing
// (NotPerformed -> Passed) => Performing
// (NotPerformed -> Failed) => Failed
// (NotPerformed -> Error) => Error
// (NotPerformed -> Canceled) => Canceled
//
// (Performing -> NotPerformed) => Performing
// (Performing -> Performing) => Performing
// (Performing -> Passed) => Performing
// (Performing -> Failed) => Failed
// (Performing -> Error) => Error
// (Performing -> Canceled) => Canceled
//
// (Failed -> NotPerformed) => Failed
// (Failed -> Performing) => Failed
// (Failed -> Passed) => Failed
// (Failed -> Failed) => Failed
// (Failed -> Error) => Error
// (Failed -> Canceled) => Failed
//
// (Error -> NotPerformed) => Error
// (Error -> Performing) => Error
// (Error -> Passed) => Error
// (Error -> Failed) => Error
// (Error -> Error) => Error
// (Error -> Canceled) => Canceled
//
// (Canceled -> NotPerformed) => Canceled
// (Canceled -> Performing) => Canceled
// (Canceled -> Passed) => Canceled
// (Canceled -> Failed) => Canceled
// (Canceled -> Error) => Canceled
// (Canceled -> Canceled) => Canceled.
func (b *FlowBuilder) WithStep(step Step) *FlowBuilder {
	t, ok := b.transitionFromStep(step)
	if !ok {
		return b
	}

	b.state = b.commonStateTransitionRules.apply(b.state, step.State())

	t.state = b.specificStateTransitionRules.apply(t.state, step.State())
	t.err = multierr.Append(t.err, step.Err())
	t.fail = multierr.Append(t.fail, step.Fail())

	return b
}

func (b *FlowBuilder) transitionFromStep(step Step) (*Transition, bool) {
	from, to, ok := step.FromTo()
	if !ok {
		return nil, false
	}

	t, ok := b.graph[from][to]
	if !ok {
		return nil, false
	}

	return t, true
}

type performStep struct {
	from          string
	to            string
	state         State
	err           error
	fail          error
	performerType performerType
}

func newPerformingStep(from, to string, performerType performerType) Step {
	return performStep{
		from:          from,
		to:            to,
		state:         Performing,
		performerType: performerType,
	}
}

func newPerformedStep(from, to string, performerType performerType, fail, err error) Step {
	var state State

	if fail != nil {
		state = Failed
	}

	if err != nil {
		state = Error
	}

	if fail == nil && err == nil {
		state = Passed
	}

	return performStep{
		from:          from,
		to:            to,
		state:         state,
		performerType: performerType,
		fail:          fail,
		err:           err,
	}
}

func (s performStep) FromTo() (from, to string, ok bool) {
	return s.from, s.to, true
}

func (s performStep) State() State {
	return s.state
}

func (s performStep) Err() error {
	return s.err
}

func (s performStep) Fail() error {
	return s.fail
}

func (s performStep) String() string {
	msg := fmt.Sprintf("Flow step %s `%s -(%s)-> %s`", s.state, s.from, s.performerType, s.to)

	if s.fail != nil {
		msg = fmt.Sprintf("%s (with fail: %s)", msg, s.fail)
	}

	if s.err != nil {
		msg = fmt.Sprintf("%s (with err: %s)", msg, s.err)
	}

	return msg
}

type cancelStep struct {
	err error
}

func newCanceledStep(err error) Step {
	return cancelStep{err: err}
}

func (c cancelStep) FromTo() (from, to string, ok bool) {
	return "", "", false
}

func (c cancelStep) State() State {
	return Canceled
}

func (c cancelStep) Err() error {
	return c.err
}

func (c cancelStep) Fail() error {
	return nil
}

func (c cancelStep) String() string {
	return fmt.Sprintf("Flow step %s", Canceled)
}
