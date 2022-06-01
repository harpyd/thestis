package flow

import (
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type (
	// Flow represents the progress of a single performance.Performance
	// run. The flow consists of working specification.Scenario and
	// specification.Thesis statuses. Each ApplyStep call moves the
	// progress forward.
	Flow struct {
		id            string
		performanceID string

		statuses map[specification.Slug]*Status
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
	var overallState State

	for _, status := range f.statuses {
		if status.state.Precedence() > overallState.Precedence() {
			overallState = status.state
		}
	}

	return overallState
}

// Statuses returns copy of scenario statuses.
func (f *Flow) Statuses() []*Status {
	if len(f.statuses) == 0 {
		return nil
	}

	statuses := make([]*Status, 0, len(f.statuses))
	for _, status := range f.statuses {
		statuses = append(statuses, status)
	}

	return statuses
}

// ApplyStep is method for step by step collecting
// performance.Performance steps to move the progress
// of the performance by changing status states.
//
// Flow state changes from call to call relying on the
// state transition rules in State.Next method.
func (f *Flow) ApplyStep(step performance.Step) *Flow {
	slug := step.Slug()

	status, ok := f.statuses[slug.ToScenarioKind()]
	if !ok {
		return f
	}

	if slug.Kind() == specification.ScenarioSlug {
		status.state = status.state.Next(step.Event())
	}

	if slug.Kind() == specification.ThesisSlug {
		thesisStatus, ok := status.thesisStatuses[slug.Partial()]
		if !ok {
			return f
		}

		thesisStatus.state = thesisStatus.state.Next(step.Event())

		if step.Err() != nil {
			thesisStatus.occurredErrs = append(
				thesisStatus.occurredErrs,
				step.Err().Error(),
			)
		}
	}

	return f
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
// during performing of theses.
func (s *ThesisStatus) OccurredErrs() []string {
	if len(s.occurredErrs) == 0 {
		return nil
	}

	occurredErrs := make([]string, len(s.occurredErrs))
	copy(occurredErrs, s.occurredErrs)

	return occurredErrs
}

// Fulfill starts a new flow from performance.Performance.
// The result of the function is a Flow, with which you can
// collect the steps coming from the performance during its
// execution.
//
// Each step contains information about the progress of the
// performance, including performance.Event. The states of
// statuses change under the action of events. The transition
// rules for the event are described in the State.Next method.
func Fulfill(id string, perf *performance.Performance) *Flow {
	var (
		scenarios = perf.WorkingScenarios()
		statuses  = make(map[specification.Slug]*Status, len(scenarios))
	)

	for _, scenario := range scenarios {
		statuses[scenario.Slug()] = &Status{
			slug:           scenario.Slug(),
			state:          NotPerformed,
			thesisStatuses: fromTheses(scenario.Theses()),
		}
	}

	return &Flow{
		id:            id,
		performanceID: perf.ID(),
		statuses:      statuses,
	}
}

func fromTheses(theses []specification.Thesis) map[string]*ThesisStatus {
	statuses := make(map[string]*ThesisStatus, len(theses))

	for _, thesis := range theses {
		slug := thesis.Slug().Partial()

		statuses[slug] = NewThesisStatus(slug, NotPerformed)
	}

	return statuses
}

// FromStatuses start a new flow from initialized statuses.
//
// This method is similar to the other, but this method is
// intended solely for testing purposes.
func FromStatuses(id, performanceID string, statuses ...*Status) *Flow {
	return &Flow{
		id:            id,
		performanceID: performanceID,
		statuses:      statusesOrNil(statuses),
	}
}

func statusesOrNil(statuses []*Status) map[specification.Slug]*Status {
	if len(statuses) == 0 {
		return nil
	}

	nonNilStatuses := make(map[specification.Slug]*Status, len(statuses))

	for _, status := range statuses {
		if status != nil {
			nonNilStatuses[status.slug] = status
		}
	}

	return nonNilStatuses
}
