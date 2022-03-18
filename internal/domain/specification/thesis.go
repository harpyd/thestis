package specification

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
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
		httpBuilder      *HTTPBuilder
		assertionBuilder *AssertionBuilder
	}

	Stage string
)

const (
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
	case UnknownStage:
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
	case UnknownStage:
		return false
	}

	return false
}

func (s Stage) String() string {
	return string(s)
}

func NewThesisBuilder() *ThesisBuilder {
	return &ThesisBuilder{
		assertionBuilder: NewAssertionBuilder(),
		httpBuilder:      NewHTTPBuilder(),
	}
}

func (b *ThesisBuilder) Build(slug Slug) (Thesis, error) {
	if slug.IsZero() {
		return Thesis{}, NewEmptySlugError()
	}

	if err := slug.ShouldBeThesisKind(); err != nil {
		return Thesis{}, err
	}

	http, buildErr := b.httpBuilder.Build()

	assertion, err := b.assertionBuilder.Build()
	buildErr = multierr.Append(buildErr, err)

	if buildErr == nil && http.IsZero() && assertion.IsZero() {
		buildErr = multierr.Append(
			buildErr,
			NewNoThesisHTTPOrAssertionError(),
		)
	}

	if !b.stage.IsValid() {
		buildErr = multierr.Append(
			buildErr,
			NewNotAllowedStageError(b.stage.String()),
		)
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

	return thesis, NewBuildSluggedError(buildErr, slug)
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
	buildFn(b.assertionBuilder)

	return b
}

func (b *ThesisBuilder) WithHTTP(buildFn func(b *HTTPBuilder)) *ThesisBuilder {
	b.httpBuilder.Reset()
	buildFn(b.httpBuilder)

	return b
}

type notAllowedStageError struct {
	keyword string
}

func NewNotAllowedStageError(stage string) error {
	return errors.WithStack(notAllowedStageError{
		keyword: stage,
	})
}

func IsNotAllowedStageError(err error) bool {
	var target notAllowedStageError

	return errors.As(err, &target)
}

func (e notAllowedStageError) Error() string {
	if e.keyword == "" {
		return "no stage"
	}

	return fmt.Sprintf("stage `%s` not allowed", e.keyword)
}

var errNoThesisHTTPOrAssertion = errors.New("no HTTP or assertion")

func NewNoThesisHTTPOrAssertionError() error {
	return errNoThesisHTTPOrAssertion
}

func IsNoThesisHTTPOrAssertionError(err error) bool {
	return errors.Is(err, errNoThesisHTTPOrAssertion)
}
