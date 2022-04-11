package specification

import (
	"fmt"
	"github.com/pkg/errors"
)

type (
	Thesis struct {
		slug         Slug
		dependencies map[Slug]bool
		stage        Stage
		behavior     string
		http         HTTP
		assertion    Assertion
	}

	ThesisBuilder struct {
		dependencies     []string
		stage            Stage
		behavior         string
		httpBuilder      HTTPBuilder
		assertionBuilder AssertionBuilder
	}

	Stage string
)

const (
	NoStage      Stage = ""
	UnknownStage Stage = "!"
	Given        Stage = "given"
	When         Stage = "when"
	Then         Stage = "then"
)

func (t Thesis) Slug() Slug {
	return t.slug
}

func (t Thesis) Dependencies() []Slug {
	return copyDependencies(t.dependencies)
}

func copyDependencies(deps map[Slug]bool) []Slug {
	if len(deps) == 0 {
		return nil
	}

	result := make([]Slug, 0, len(deps))

	for dep := range deps {
		result = append(result, dep)
	}

	return result
}

func (t Thesis) Stage() Stage {
	return t.stage
}

func (t Thesis) Behavior() string {
	return t.behavior
}

func (t Thesis) HTTP() HTTP {
	return t.http
}

func (t Thesis) Assertion() Assertion {
	return t.assertion
}

func (s Stage) Before() []Stage {
	switch s {
	case Given:
		return nil
	case When:
		return []Stage{Given}
	case Then:
		return []Stage{Given, When}
	case NoStage, UnknownStage:
		return nil
	}

	return nil
}

func (s Stage) IsValid() bool {
	switch s {
	case Given:
		return true
	case When:
		return true
	case Then:
		return true
	case NoStage, UnknownStage:
		return false
	}

	return false
}

func (s Stage) String() string {
	return string(s)
}

var ErrUselessThesis = errors.New("useless thesis")

func (t Thesis) validate(ctxScenario Scenario) error {
	var w BuildErrorWrapper

	switch {
	case !t.http.IsZero():
		w.WithError(t.http.validate())
	case !t.assertion.IsZero():
		w.WithError(t.assertion.validate())
	default:
		w.WithError(ErrUselessThesis)
	}

	if !t.stage.IsValid() {
		w.WithError(NewNotAllowedStageError(t.stage))
	}

	for dep := range t.dependencies {
		if _, ok := ctxScenario.theses[dep.Thesis()]; !ok {
			w.WithError(NewUndefinedDependencyError(dep))
		}
	}

	return w.SluggedWrap(t.slug)
}

func (b *ThesisBuilder) Build(slug Slug) Thesis {
	if err := slug.ShouldBeThesisKind(); err != nil {
		panic(err)
	}

	return Thesis{
		slug:         slug,
		dependencies: dependenciesOrNil(slug, b.dependencies),
		stage:        b.stage,
		behavior:     b.behavior,
		http:         b.httpBuilder.Build(),
		assertion:    b.assertionBuilder.Build(),
	}
}

func dependenciesOrNil(thesisSlug Slug, deps []string) map[Slug]bool {
	if len(deps) == 0 {
		return nil
	}

	theses := make(map[Slug]bool, len(deps))

	for _, dep := range deps {
		slug := NewThesisSlug(
			thesisSlug.Story(),
			thesisSlug.Scenario(),
			dep,
		)
		theses[slug] = true
	}

	return theses
}

func (b *ThesisBuilder) Reset() {
	b.dependencies = nil
	b.stage = ""
	b.behavior = ""
	b.assertionBuilder.Reset()
	b.httpBuilder.Reset()
}

func (b *ThesisBuilder) WithDependency(dep string) *ThesisBuilder {
	b.dependencies = append(b.dependencies, dep)

	return b
}

func (b *ThesisBuilder) WithStatement(stage Stage, behavior string) *ThesisBuilder {
	b.stage = stage
	b.behavior = behavior

	return b
}

func (b *ThesisBuilder) WithAssertion(buildFn func(b *AssertionBuilder)) *ThesisBuilder {
	b.assertionBuilder.Reset()
	buildFn(&b.assertionBuilder)

	return b
}

func (b *ThesisBuilder) WithHTTP(buildFn func(b *HTTPBuilder)) *ThesisBuilder {
	b.httpBuilder.Reset()
	buildFn(&b.httpBuilder)

	return b
}

type NotAllowedStageError struct {
	stage Stage
}

func NewNotAllowedStageError(stage Stage) error {
	return errors.WithStack(&NotAllowedStageError{
		stage: stage,
	})
}

func (e *NotAllowedStageError) Stage() Stage {
	return e.stage
}

func (e *NotAllowedStageError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("stage %q not allowed", e.stage)
}

type UndefinedDependencyError struct {
	slug Slug
}

func NewUndefinedDependencyError(slug Slug) error {
	return errors.WithStack(&UndefinedDependencyError{
		slug: slug,
	})
}

func (e *UndefinedDependencyError) Slug() Slug {
	return e.slug
}

func (e *UndefinedDependencyError) Error() string {
	return fmt.Sprintf("undefined %q dependency", e.slug.Partial())
}
