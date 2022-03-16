package performance

import (
	"strings"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type Step struct {
	slug          specification.Slug
	performerType PerformerType
	event         Event
	err           error
}

func NewScenarioStep(slug specification.Slug, event Event) Step {
	return Step{
		slug:          slug,
		performerType: NoPerformer,
		event:         event,
	}
}

func NewScenarioStepWithErr(err error, slug specification.Slug, event Event) Step {
	return Step{
		slug:          slug,
		performerType: NoPerformer,
		event:         event,
		err:           err,
	}
}

func NewThesisStep(slug specification.Slug, pt PerformerType, event Event) Step {
	return Step{
		slug:          slug,
		performerType: pt,
		event:         event,
	}
}

func NewThesisStepWithErr(
	err error,
	slug specification.Slug,
	pt PerformerType,
	event Event,
) Step {
	return Step{
		slug:          slug,
		performerType: pt,
		event:         event,
		err:           err,
	}
}

func (s Step) Slug() specification.Slug {
	return s.slug
}

func (s Step) PerformerType() PerformerType {
	return s.performerType
}

func (s Step) Event() Event {
	return s.event
}

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

func (s Step) IsZero() bool {
	return s == Step{}
}
