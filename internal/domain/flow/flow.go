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

		statuses []Status
	}

	// Status represents progress of the slugged object like
	// specification.Scenario and specification.Thesis.
	Status struct {
		slug         specification.Slug
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
func (f Flow) ID() string {
	return f.id
}

// PerformanceID returns associated with
// Flow performance identifier.
func (f Flow) PerformanceID() string {
	return f.performanceID
}

// Statuses returns copy of slugged object statuses.
func (f Flow) Statuses() []Status {
	statuses := make([]Status, len(f.statuses))
	copy(statuses, f.statuses)

	return statuses
}

// NewStatus creates representation of slugged object progress.
func NewStatus(slug specification.Slug, state State, errMessages ...string) Status {
	return Status{
		slug:         slug,
		state:        state,
		occurredErrs: errorsOrNil(errMessages),
	}
}

func errorsOrNil(errs []string) []string {
	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Slug returns associated with Status object slug.
func (f Status) Slug() specification.Slug {
	return f.slug
}

// State returns state of slugged object.
func (f Status) State() State {
	return f.state
}

// OccurredErrs returns errors that occurred
// during performing of slugged objects.
func (f Status) OccurredErrs() []string {
	occurredErrs := make([]string, len(f.occurredErrs))
	copy(occurredErrs, f.occurredErrs)

	return occurredErrs
}

type Params struct {
	ID            string
	PerformanceID string
	Statuses      []Status
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
		statuses:      make([]Status, len(params.Statuses)),
	}

	copy(f.statuses, params.Statuses)

	return f
}

// FromPerformance starts a new flow from performance.Performance.
// The result of the function is a Reducer, with which you can
// collect the steps coming from performance.Performance
// during its execution.
//
// Each step contains information about the progress of the
// performance.Performance, including performance.Event.
// The states of statuses change under the action of events.
// The transition rules for the event are described in the
// State.Next method.
func FromPerformance(id string, perf *performance.Performance) *Reducer {
	scenarios := perf.WorkingScenarios()

	statuses := make(map[specification.Slug]*Status, len(scenarios))

	for _, scenario := range scenarios {
		statuses[scenario.Slug()] = &Status{
			slug:  scenario.Slug(),
			state: NotPerformed,
		}

		for _, thesis := range scenario.Theses() {
			statuses[thesis.Slug()] = &Status{
				slug:  thesis.Slug(),
				state: NotPerformed,
			}
		}
	}

	return &Reducer{
		id:            id,
		performanceID: perf.ID(),
		statuses:      statuses,
	}
}

// FromStatuses start a new flow from initialized statuses.
//
// This method is similar to the other, but this method is
// intended solely for testing purposes.
func FromStatuses(id, performanceID string, statuses ...Status) *Reducer {
	r := &Reducer{
		id:            id,
		performanceID: performanceID,
		statuses:      make(map[specification.Slug]*Status, len(statuses)),
	}

	for _, status := range statuses {
		r.statuses[status.slug] = &Status{
			slug:         status.slug,
			state:        status.state,
			occurredErrs: status.occurredErrs,
		}
	}

	return r
}

// Reduce creates current version of Flow from Reducer.
// This is useful for accumulating performance.Performance steps and
// storing the state of performance.Performance's Flow.
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
		statuses:      r.statusesSnapshot(),
	}
}

func (r *Reducer) statusesSnapshot() []Status {
	statuses := make([]Status, 0, len(r.statuses))

	for _, status := range r.statuses {
		statuses = append(statuses, *status)
	}

	return statuses
}

// WithStep is method for step by step collecting
// performance.Performance steps for their further
// reduction with Reducer's Reduce.
//
// Flow state changes from call to call relying on the
// state transition rules in State.Next method.
func (r *Reducer) WithStep(step performance.Step) *Reducer {
	slug := step.Slug()

	status, ok := r.statuses[slug]
	if !ok {
		return r
	}

	status.state = status.state.Next(step.Event())

	if step.Err() != nil {
		status.occurredErrs = append(status.occurredErrs, step.Err().Error())
	}

	return r
}
