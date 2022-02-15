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
		PerformerType() PerformerType
		State() State
		Err() error
		String() string
	}

	// Flow represents current Performance performing.
	// Flow keeps transitions information
	// and common state of performing.
	Flow struct {
		id            string
		performanceID string

		state State
		graph map[string]map[string]Transition
	}

	Transition struct {
		from  string
		to    string
		state State
		err   error
	}

	// FlowReducer builds Flow instance using WithStep,
	// Reduce and FinallyReduce methods.
	//
	// FlowReducer defines Flow common state transition rules
	// in WithStep method.
	FlowReducer struct {
		id            string
		performanceID string

		state         State
		graph         map[string]map[string]*Transition
		commonRules   stateTransitionRules
		specificRules stateTransitionRules
	}
)

func (f Flow) ID() string {
	return f.id
}

func (f Flow) PerformanceID() string {
	return f.performanceID
}

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

func NewTransition(state State, from, to, errMsg string) Transition {
	return Transition{
		from:  from,
		to:    to,
		state: state,
		err:   newTransitionError(state, errMsg),
	}
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

type FlowParams struct {
	ID            string
	PerformanceID string
	State         State
	Transitions   []Transition
}

func UnmarshalFlowFromDatabase(params FlowParams) Flow {
	return Flow{
		id:            params.ID,
		performanceID: params.PerformanceID,
		state:         params.State,
		graph:         unmarshalFlowGraph(params.Transitions),
	}
}

func unmarshalFlowGraph(transitions []Transition) map[string]map[string]Transition {
	graph := make(map[string]map[string]Transition, len(transitions))

	for _, t := range transitions {
		initGraphTransitionsLazy(graph, t.From())

		graph[t.From()][t.To()] = t
	}

	return graph
}

const defaultTransitionsSize = 1

func initGraphTransitionsLazy(graph map[string]map[string]Transition, vertex string) {
	if graph[vertex] == nil {
		graph[vertex] = make(map[string]Transition, defaultTransitionsSize)
	}
}

func FlowFromPerformance(id string, perf *Performance) *FlowReducer {
	graph := make(map[string]map[string]*Transition, len(perf.actionGraph))

	for from, as := range perf.actionGraph {
		graph[from] = make(map[string]*Transition, len(as))

		for to := range as {
			graph[from][to] = &Transition{
				state: NotPerformed,
				from:  from,
				to:    to,
			}
		}
	}

	return &FlowReducer{
		id:            id,
		performanceID: perf.ID(),
		state:         NotPerformed,
		graph:         graph,
		commonRules:   newCommonStateTransitionRules(),
		specificRules: newSpecificStateTransitionRules(),
	}
}

func FlowFromState(commonState, transitionState State, from, to string) *FlowReducer {
	graph := map[string]map[string]*Transition{
		from: {
			to: &Transition{
				state: transitionState,
				from:  from,
				to:    to,
			},
		},
	}

	return &FlowReducer{
		state:         commonState,
		graph:         graph,
		commonRules:   newCommonStateTransitionRules(),
		specificRules: newSpecificStateTransitionRules(),
	}
}

// Reduce creates intermediate version of Flow from FlowReducer.
// Flow created using this method can't have a State equal
// to Passed.
//
// You need to use FinallyReduce method to create
// final version of Flow.
func (r *FlowReducer) Reduce() Flow {
	return Flow{
		id:            r.id,
		performanceID: r.performanceID,
		state:         r.state,
		graph:         r.copyGraph(),
	}
}

// FinallyReduce creates final version of Flow from FlowReducer.
// Flow created using this method represents final Performance's Flow.
//
// If all transitions have Passed states and common State equal
// to Performing, final version of Flow will have Passed State.
func (r *FlowReducer) FinallyReduce() Flow {
	return Flow{
		id:            r.id,
		performanceID: r.performanceID,
		state:         r.finalState(),
		graph:         r.copyGraph(),
	}
}

func (r *FlowReducer) finalState() State {
	if r.state == Performing && r.allPassed() {
		return Passed
	}

	return r.state
}

func (r *FlowReducer) allPassed() bool {
	for _, ts := range r.graph {
		for _, t := range ts {
			if t.State() != Passed {
				return false
			}
		}
	}

	return true
}

func (r *FlowReducer) copyGraph() map[string]map[string]Transition {
	graph := make(map[string]map[string]Transition, len(r.graph))

	for from, ts := range r.graph {
		graph[from] = make(map[string]Transition, len(ts))

		for to, t := range ts {
			graph[from][to] = *t
		}
	}

	return graph
}

// WithStep is method for step by step collecting Step's for their
// further reduction with FlowReducer's Reduce or FinallyReduce.
// Flow common state changes from call to call relying on the
// state transition rules.
func (r *FlowReducer) WithStep(step Step) *FlowReducer {
	r.state = r.commonRules.apply(r.state, step.State())

	t, ok := r.transitionFromStep(step)
	if !ok {
		return r
	}

	t.state = r.specificRules.apply(t.state, step.State())
	t.err = multierr.Append(t.err, step.Err())

	return r
}

func (r *FlowReducer) transitionFromStep(step Step) (*Transition, bool) {
	from, to, ok := step.FromTo()
	if !ok {
		return nil, false
	}

	t, ok := r.graph[from][to]
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
	performerType PerformerType
}

func NewPerformingStep(from, to string, performerType PerformerType) Step {
	return performStep{
		from:          from,
		to:            to,
		state:         Performing,
		performerType: performerType,
	}
}

func NewStepFromResult(from, to string, performerType PerformerType, result Result) Step {
	return performStep{
		from:          from,
		to:            to,
		state:         result.State(),
		performerType: performerType,
		err:           result.Err(),
	}
}

func (s performStep) FromTo() (from, to string, ok bool) {
	return s.from, s.to, true
}

func (s performStep) PerformerType() PerformerType {
	return s.performerType
}

func (s performStep) State() State {
	return s.state
}

func (s performStep) Err() error {
	return s.err
}

func (s performStep) String() string {
	msg := fmt.Sprintf("Flow step `%s -(%s)-> %s` %s", s.from, s.performerType, s.to, s.state)

	if s.err != nil {
		msg = fmt.Sprintf("%s (with err: %s)", msg, s.err)
	}

	return msg
}

type cancelStep struct{}

func NewCanceledStep() Step {
	return cancelStep{}
}

func newCanceledStep() Step {
	return cancelStep{}
}

func (s cancelStep) FromTo() (from, to string, ok bool) {
	return "", "", false
}

func (s cancelStep) PerformerType() PerformerType {
	return EmptyPerformer
}

func (s cancelStep) State() State {
	return Canceled
}

func (s cancelStep) Err() error {
	return nil
}

func (s cancelStep) String() string {
	return fmt.Sprintf("Flow step %s", Canceled)
}
