package flow

import (
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	// Flow represents the progress of a single performance.Performance
	// run. The flow consists of working specification.Scenario and
	// specification.Thesis statuses.
	Flow struct {
		id            string
		performanceID string

		overallState State
		statuses     []*Status
	}

	// Status represents progress of the specification.Scenario.
	Status struct {
		slug           specification.Slug
		state          State
		thesisStatuses map[string]*ThesisStatus
	}

	// ThesisStatus represents progress of the specification.Thesis
	// nested in Status.
	ThesisStatus struct {
		thesisSlug   string
		state        State
		occurredErrs []string
	}

	// Reducer builds Flow instance using WithStep
	// and Reduce methods.
	Reducer struct {
		id            string
		performanceID string

		statuses map[specification.Slug]*Status
	}
)

// ID returns flow identifier.
func (f *Flow) ID() string {
	return f.id
}

// PerformanceID returns associated with
// Flow performance identifier.
func (f *Flow) PerformanceID() string {
	return f.performanceID
}

// OverallState returns general Flow status
// selected from all specification.Scenario
// states according to State.Precedence.
func (f *Flow) OverallState() State {
	return f.overallState
}

// Statuses returns copy of slugged object statuses.
func (f *Flow) Statuses() []*Status {
	if len(f.statuses) == 0 {
		return nil
	}

	statuses := make([]*Status, len(f.statuses))
	copy(statuses, f.statuses)

	return statuses
}

// NewStatus creates a progress representation of specification.Scenario.
//
// If the slug is not specification.ScenarioSlug, it panics with
// specification.ErrNotScenarioSlug.
func NewStatus(slug specification.Slug, state State, thesisStatuses ...*ThesisStatus) *Status {
	if err := slug.ShouldBeScenarioKind(); err != nil {
		panic(err)
	}

	return &Status{
		slug:           slug,
		state:          state,
		thesisStatuses: thesisStatusesOrNil(thesisStatuses),
	}
}

func thesisStatusesOrNil(statuses []*ThesisStatus) map[string]*ThesisStatus {
	if len(statuses) == 0 {
		return nil
	}

	nonNilStatuses := make(map[string]*ThesisStatus, len(statuses))

	for _, status := range statuses {
		if status != nil {
			nonNilStatuses[status.thesisSlug] = status
		}
	}

	return nonNilStatuses
}

// Slug returns associated with Status scenario slug.
func (s *Status) Slug() specification.Slug {
	return s.slug
}

// State returns state of scenario.
func (s *Status) State() State {
	return s.state
}

// ThesisStatuses returns nested in Status
// thesis statuses.
func (s *Status) ThesisStatuses() []*ThesisStatus {
	if len(s.thesisStatuses) == 0 {
		return nil
	}

	thesisStatuses := make([]*ThesisStatus, 0, len(s.thesisStatuses))
	for _, status := range s.thesisStatuses {
		thesisStatuses = append(thesisStatuses, status)
	}

	return thesisStatuses
}

// NewThesisStatus creates a progress representation of specification.Thesis.
func NewThesisStatus(slug string, state State, occurredErrs ...string) *ThesisStatus {
	return &ThesisStatus{
		thesisSlug:   slug,
		state:        state,
		occurredErrs: errsOrNil(occurredErrs),
	}
}

func errsOrNil(errs []string) []string {
	if len(errs) == 0 {
		return nil
	}

	return errs
}

// ThesisSlug returns associated with ThesisStatus partial thesis slug.
func (s *ThesisStatus) ThesisSlug() string {
	return s.thesisSlug
}

// State returns state of thesis.
func (s *ThesisStatus) State() State {
	return s.state
}

// OccurredErrs returns errors that occurred
// during performing of slugged objects.
func (s *ThesisStatus) OccurredErrs() []string {
	if len(s.occurredErrs) == 0 {
		return nil
	}

	occurredErrs := make([]string, len(s.occurredErrs))
	copy(occurredErrs, s.occurredErrs)

	return occurredErrs
}

type Params struct {
	ID            string
	PerformanceID string
	OverallState  State
	Statuses      []*Status
}

// Unmarshal transforms Params to Flow.
//
// This function is great for converting
// from a database or using in tests.
//
// You must not use this method in
// business code of domain and app layers.
func Unmarshal(params Params) Flow {
	f := Flow{
		id:            params.ID,
		performanceID: params.PerformanceID,
		overallState:  params.OverallState,
		statuses:      make([]*Status, len(params.Statuses)),
	}

	copy(f.statuses, params.Statuses)

	return f
}

// FromPerformance starts a new flow from performance.Performance.
// The result of the function is a Reducer, with which you can
// collect the steps coming from the performance during its execution.
//
// Each step contains information about the progress of the
// performance, including performance.Event. The states of
// statuses change under the action of events. The transition
// rules for the event are described in the State.Next method.
func FromPerformance(id string, perf *performance.Performance) *Reducer {
	scenarios := perf.WorkingScenarios()

	statuses := make(map[specification.Slug]*Status, len(scenarios))

	for _, scenario := range scenarios {
		statuses[scenario.Slug()] = &Status{
			slug:           scenario.Slug(),
			state:          NotPerformed,
			thesisStatuses: fromTheses(scenario.Theses()),
		}
	}

	return &Reducer{
		id:            id,
		performanceID: perf.ID(),
		statuses:      statuses,
	}
}

func fromTheses(theses []specification.Thesis) map[string]*ThesisStatus {
	statuses := make(map[string]*ThesisStatus, len(theses))

	for _, thesis := range theses {
		slug := thesis.Slug().Partial()

		statuses[slug] = &ThesisStatus{
			thesisSlug: slug,
			state:      NotPerformed,
		}
	}

	return statuses
}

// FromStatuses start a new flow from initialized statuses.
//
// This method is similar to the other, but this method is
// intended solely for testing purposes.
func FromStatuses(id, performanceID string, statuses ...*Status) *Reducer {
	nonNilStatuses := make(map[specification.Slug]*Status, len(statuses))

	for _, status := range statuses {
		if status != nil {
			nonNilStatuses[status.slug] = status
		}
	}

	return &Reducer{
		id:            id,
		performanceID: performanceID,
		statuses:      nonNilStatuses,
	}
}

// Reduce creates current version of Flow from Reducer.
// This is useful for accumulating performance.Performance
// steps and storing the state of performance's Flow.
//
// For example:
//  fr := performance.FlowFromPerformance("id", perf)
//
//  for s := range steps {
//   fr.WithStep(s)
//   save(ctx, fr.Reduce())
//  }
func (r *Reducer) Reduce() Flow {
	return Flow{
		id:            r.id,
		performanceID: r.performanceID,
		overallState:  r.selectOverallState(),
		statuses:      statusesOrNil(r.statuses),
	}
}

func (r *Reducer) selectOverallState() State {
	var overallState State

	for _, status := range r.statuses {
		if status.state.Precedence() > overallState.Precedence() {
			overallState = status.state
		}
	}

	return overallState
}

func statusesOrNil(statuses map[specification.Slug]*Status) []*Status {
	if len(statuses) == 0 {
		return nil
	}

	result := make([]*Status, 0, len(statuses))

	for _, status := range statuses {
		result = append(result, status)
	}

	return result
}

// WithStep is method for step by step collecting
// performance.Performance steps for their further
// reduction with Reducer's Reduce.
//
// Flow state changes from call to call relying on the
// state transition rules in State.Next method.
func (r *Reducer) WithStep(step performance.Step) *Reducer {
	slug := step.Slug()

	status, ok := r.statuses[slug.ToScenarioKind()]
	if !ok {
		return r
	}

	if slug.Kind() == specification.ScenarioSlug {
		status.state = status.state.Next(step.Event())
	}

	if slug.Kind() == specification.ThesisSlug {
		thesisStatus, ok := status.thesisStatuses[slug.Partial()]
		if !ok {
			return r
		}

		thesisStatus.state = thesisStatus.state.Next(step.Event())

		if step.Err() != nil {
			thesisStatus.occurredErrs = append(
				thesisStatus.occurredErrs,
				step.Err().Error(),
			)
		}
	}

	return r
}
