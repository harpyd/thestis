package specification

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	Thesis struct {
		slug         Slug
		dependencies []Slug
		statement    Statement
		http         HTTP
		assertion    Assertion
	}

	Statement struct {
		stage    Stage
		behavior string
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
	return t.dependencies
}

func (t Thesis) Statement() Statement {
	return t.statement
}

func (t Thesis) HTTP() HTTP {
	return t.http
}

func (t Thesis) Assertion() Assertion {
	return t.assertion
}

func (s Statement) Stage() Stage {
	return s.stage
}

func (s Statement) Behavior() string {
	return s.behavior
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

func (b *ThesisBuilder) Build(slug Slug) (Thesis, error) {
	if err := slug.ShouldBeThesisKind(); err != nil {
		panic(err)
	}

	var w BuildErrorWrapper

	http, err := b.httpBuilder.Build()
	w.WithError(err)

	assertion, err := b.assertionBuilder.Build()
	w.WithError(err)

	if http.IsZero() && assertion.IsZero() {
		w.WithError(ErrUselessThesis)
	}

	if !b.stage.IsValid() {
		w.WithError(NewNotAllowedStageError(b.stage))
	}

	thesis := Thesis{
		slug:         slug,
		dependencies: make([]Slug, 0, len(b.dependencies)),
		statement: Statement{
			stage:    b.stage,
			behavior: b.behavior,
		},
		http:      http,
		assertion: assertion,
	}

	for _, dep := range b.dependencies {
		thesis.dependencies = append(
			thesis.dependencies,
			NewThesisSlug(slug.Story(), slug.Scenario(), dep),
		)
	}

	return thesis, w.SluggedWrap(slug)
}

func (b *ThesisBuilder) ErrlessBuild(slug Slug) Thesis {
	t, _ := b.Build(slug)

	return t
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
	return fmt.Sprintf("undefined `%s` dependency", e.slug)
}
