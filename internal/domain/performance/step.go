package performance

import (
	"strings"

	"github.com/harpyd/thestis/internal/domain/specification"
)

// Step stores information about the performing
// of slugged objects, such as specification.Scenario
// and specification.Thesis. It's used to observe
// the Performance execution.
type Step struct {
	slug          specification.Slug
	performerType PerformerType
	event         Event
	err           error
}

// NewScenarioStep returns a Step for the scenario,
// that is, without the PerformerType.
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
		slug:          slug,
		performerType: NoPerformer,
		event:         event,
		err:           err,
	}
}

// NewThesisStep returns a Step for the thesis.
//
// Intended only for use with slugs with
// the specification.ThesisSlug kind.
func NewThesisStep(slug specification.Slug, pt PerformerType, event Event) Step {
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
	pt PerformerType,
	event Event,
) Step {
	if err := slug.ShouldBeThesisKind(); err != nil {
		panic(err)
	}

	return Step{
		slug:          slug,
		performerType: pt,
		event:         event,
		err:           err,
	}
}

// Slug returns the slug of the object for
// which the Step was created.
func (s Step) Slug() specification.Slug {
	return s.slug
}

// PerformerType returns the PerformerType if
// the slug is specification.ThesisSlug,
// else NoPerformer.
func (s Step) PerformerType() PerformerType {
	return s.performerType
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

	if s.performerType != NoPerformer {
		b.WriteString(", type = ")
		b.WriteString(string(s.performerType))
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
