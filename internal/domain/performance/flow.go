package performance

import "fmt"

type State string

const (
	NotPerformed = ""
	Performing   = "performing"
	Passed       = "passed"
	Failed       = "failed"
	Error        = "error"
	Cancelled    = "cancelled"
)

type (
	Step interface {
		FromTo() (from, to string, ok bool)
		State() State
		Err() error
		Fail() error
		String() string
	}

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

	FlowBuilder struct {
		state State
		graph map[string]map[string]*Transition
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
		state: NotPerformed,
		graph: graph,
	}
}

func (b *FlowBuilder) Build() Flow {
	graph := make(map[string]map[string]Transition, len(b.graph))

	for from, ts := range b.graph {
		graph[from] = make(map[string]Transition, len(ts))

		for to, t := range ts {
			graph[from][to] = *t
		}
	}

	return Flow{
		state: b.state,
		graph: graph,
	}
}

func (b *FlowBuilder) WithStep(step Step) *FlowBuilder {
	if step.State() == Error {
		b.state = Error
	}

	if b.state != Error {
		if step.State() == Failed {
			b.state = Failed
		}

		if step.State() == Cancelled {
			b.state = Cancelled
		}
	}

	from, to, ok := step.FromTo()
	if !ok {
		return b
	}

	t, ok := b.graph[from][to]
	if !ok {
		return b
	}

	t.state = step.State()
	t.err = step.Err()
	t.fail = step.Fail()

	return b
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

func newCancelledStep(err error) Step {
	return cancelStep{err: err}
}

func (c cancelStep) FromTo() (from, to string, ok bool) {
	return "", "", false
}

func (c cancelStep) State() State {
	return Cancelled
}

func (c cancelStep) Err() error {
	return c.err
}

func (c cancelStep) Fail() error {
	return nil
}

func (c cancelStep) String() string {
	return fmt.Sprintf("Flow step %s", Cancelled)
}
