package flow

import (
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	// Flow represents current performance.Performance performing.
	// Flow keeps transitions information
	// and common state of performing.
	Flow struct {
		id            string
		performanceID string

		statuses []Status
	}

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

func (f Flow) ID() string {
	return f.id
}

func (f Flow) PerformanceID() string {
	return f.performanceID
}

func (f Flow) Statuses() []Status {
	statuses := make([]Status, len(f.statuses))
	copy(statuses, f.statuses)

	return statuses
}

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

func (f Status) Slug() specification.Slug {
	return f.slug
}

func (f Status) State() State {
	return f.state
}

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

func Unmarshal(params Params) Flow {
	f := Flow{
		id:            params.ID,
		performanceID: params.PerformanceID,
		statuses:      make([]Status, len(params.Statuses)),
	}

	copy(f.statuses, params.Statuses)

	return f
}

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

// Reduce creates current version of Flow from FlowReducer.
// This is useful for accumulating Performance Step's and
// storing the state of Performance Flow.
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

// WithStep is method for step by step collecting Step's for
// their further reduction with FlowReducer's Reduce.
//
// Flow state changes from call to call relying on the
// state transition rules.
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
