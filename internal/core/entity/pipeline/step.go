package pipeline

import (
	"strings"

	"github.com/harpyd/thestis/internal/core/entity/specification"
)

// Step stores information about the executing
// of slugged objects, such as specification.Scenario
// and specification.Thesis. It's used to observe
// the Pipeline execution.
type Step struct {
	slug         specification.Slug
	executorType ExecutorType
	event        Event
	err          error
}

// NewScenarioStep returns a Step for the scenario,
// that is, without the ExecutorType.
//
// Intended only for use with slugs with
// the specification.ScenarioSlug kind.
func NewScenarioStep(slug specification.Slug, event Event) Step {
	return NewScenarioStepWithErr(nil, slug, event)
}

// NewScenarioStepWithErr is similar to NewScenarioStep,
// only it gets the error that occurred.
//
// Intended only for use with slugs with
// the specification.ScenarioSlug kind.
func NewScenarioStepWithErr(err error, slug specification.Slug, event Event) Step {
	if err := slug.ShouldBeScenarioKind(); err != nil {
		panic(err)
	}

	return Step{
		slug:         slug,
		executorType: NoExecutor,
		event:        event,
		err:          err,
	}
}

// NewThesisStep returns a Step for the thesis.
//
// Intended only for use with slugs with
// the specification.ThesisSlug kind.
func NewThesisStep(slug specification.Slug, pt ExecutorType, event Event) Step {
	return NewThesisStepWithErr(nil, slug, pt, event)
}

// NewThesisStepWithErr is similar to NewThesisStep,
// only it gets the error that occurred.
//
// Intended only for use with slugs with
// the specification.ThesisSlug kind.
func NewThesisStepWithErr(
	err error,
	slug specification.Slug,
	pt ExecutorType,
	event Event,
) Step {
	if err := slug.ShouldBeThesisKind(); err != nil {
		panic(err)
	}

	return Step{
		slug:         slug,
		executorType: pt,
		event:        event,
		err:          err,
	}
}

// Slug returns the slug of the object for
// which the Step was created.
func (s Step) Slug() specification.Slug {
	return s.slug
}

// ExecutorType returns the ExecutorType if
// the slug is specification.ThesisSlug,
// else NoExecutor.
func (s Step) ExecutorType() ExecutorType {
	return s.executorType
}

// Event returns the event fired at the Step.
func (s Step) Event() Event {
	return s.event
}

// Err returns the error occurred at the Step.
func (s Step) Err() error {
	return s.err
}

func (s Step) String() string {
	var b strings.Builder

	b.WriteString(s.slug.String())

	if s.event != NoEvent {
		b.WriteString(": event = ")
		b.WriteString(s.event.String())
	}

	if s.executorType != NoExecutor {
		b.WriteString(", type = ")
		b.WriteString(string(s.executorType))
	}

	if s.err != nil {
		b.WriteString(", err = ")
		b.WriteString(s.err.Error())
	}

	return b.String()
}

// IsZero returns true if the Step
// is empty, else false.
func (s Step) IsZero() bool {
	return s == Step{}
}
