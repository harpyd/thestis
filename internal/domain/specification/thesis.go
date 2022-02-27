package specification

import (
	"fmt"
	"strings"

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
		stage            string
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

func stageFromString(keyword string) (Stage, error) {
	switch strings.ToLower(keyword) {
	case "given":
		return Given, nil
	case "when":
		return When, nil
	case "then":
		return Then, nil
	}

	return UnknownStage, NewNotAllowedStageError(keyword)
}

func (k Stage) String() string {
	return string(k)
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

	stage, stageErr := stageFromString(b.stage)
	http, httpErr := b.httpBuilder.Build()
	assertion, assertionErr := b.assertionBuilder.Build()

	err := multierr.Combine(httpErr, assertionErr)
	if err == nil && http.IsZero() && assertion.IsZero() {
		err = NewNoThesisHTTPOrAssertionError()
	}

	thsis := Thesis{
		slug:         slug,
		dependencies: make([]Slug, 0, len(b.dependencies)),
		statement: Statement{
			stage:    stage,
			behavior: b.behavior,
		},
		http:      http,
		assertion: assertion,
	}

	for _, dep := range b.dependencies {
		thsis.dependencies = append(
			thsis.dependencies,
			NewThesisSlug(slug.Story(), slug.Scenario(), dep),
		)
	}

	return thsis, NewBuildSluggedError(multierr.Combine(stageErr, err), slug)
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

func (b *ThesisBuilder) WithDependencies(dep string) *ThesisBuilder {
	b.dependencies = append(b.dependencies, dep)

	return b
}

func (b *ThesisBuilder) WithStatement(stage string, behavior string) *ThesisBuilder {
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

func NewNotAllowedStageError(keyword string) error {
	return errors.WithStack(notAllowedStageError{
		keyword: keyword,
	})
}

func IsNotAllowedStageError(err error) bool {
	var nerr notAllowedStageError

	return errors.As(err, &nerr)
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
