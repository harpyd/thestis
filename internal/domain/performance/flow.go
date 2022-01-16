package performance

import (
	"fmt"

	"go.uber.org/multierr"
)

type State string

const (
	NotPerformed = ""
	Performing   = "performing"
	Passed       = "passed"
	Failed       = "failed"
	Error        = "error"
	Canceled     = "canceled"
)

type stateTransitionRules map[State]map[State]State

func newStateTransitionRules() stateTransitionRules {
	return stateTransitionRules{
		NotPerformed: {
			NotPerformed: NotPerformed,
			Performing:   Performing,
			Passed:       Performing,
			Failed:       Failed,
			Error:        Error,
			Canceled:     Canceled,
		},
		Performing: {
			NotPerformed: Performing,
			Performing:   Performing,
			Passed:       Performing,
			Failed:       Failed,
			Error:        Error,
			Canceled:     Canceled,
		},
		Failed: {
			NotPerformed: Failed,
			Performing:   Failed,
			Passed:       Failed,
			Failed:       Failed,
			Error:        Error,
			Canceled:     Failed,
		},
		Error: {
			NotPerformed: Error,
			Performing:   Error,
			Passed:       Error,
			Failed:       Error,
			Error:        Error,
			Canceled:     Error,
		},
		Canceled: {
			NotPerformed: Canceled,
			Performing:   Canceled,
			Passed:       Canceled,
			Failed:       Canceled,
			Error:        Canceled,
			Canceled:     Canceled,
		},
	}
}

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
		state                State
		graph                map[string]map[string]*Transition
		stateTransitionRules stateTransitionRules
	}
)

func (f *Flow) State() State {
	return f.state
}

func (f *Flow) Transitions() []Transition {
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
		state:                NotPerformed,
		graph:                graph,
		stateTransitionRules: newStateTransitionRules(),
	}
}

func (b *FlowBuilder) Build() *Flow {
	return &Flow{
		state: b.state,
		graph: b.copyGraph(),
	}
}

func (b *FlowBuilder) FinallyBuild() *Flow {
	return &Flow{
		state: b.finalState(),
		graph: b.copyGraph(),
	}
}

func (b *FlowBuilder) finalState() State {
	if b.state == Performing {
		return Passed
	}

	return b.state
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

	b.state = b.stateTransitionRules[b.state][step.State()]

	t.state = step.State()
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
