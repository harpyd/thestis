package specification

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	Thesis struct {
		slug         string
		dependencies []string
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

func (t Thesis) Slug() string {
	return t.slug
}

func (t Thesis) Dependencies() []string {
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

func (b *ThesisBuilder) Build(slug string) (Thesis, error) {
	if slug == "" {
		return Thesis{}, NewThesisEmptySlugError()
	}

	stage, keywordErr := stageFromString(b.stage)
	http, httpErr := b.httpBuilder.Build()
	assertion, assertionErr := b.assertionBuilder.Build()

	err := multierr.Combine(httpErr, assertionErr)
	if err == nil && http.IsZero() && assertion.IsZero() {
		err = NewNoThesisHTTPOrAssertionError()
	}

	thsis := Thesis{
		slug:         slug,
		dependencies: make([]string, len(b.dependencies)),
		statement: Statement{
			stage:    stage,
			behavior: b.behavior,
		},
		http:      http,
		assertion: assertion,
	}

	copy(thsis.dependencies, b.dependencies)

	return thsis, NewBuildThesisError(multierr.Combine(keywordErr, err), slug)
}

func (b *ThesisBuilder) ErrlessBuild(slug string) Thesis {
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

func (b *ThesisBuilder) WithDependencies(after string) *ThesisBuilder {
	b.dependencies = append(b.dependencies, after)

	return b
}

func (b *ThesisBuilder) WithStage(stage string, behavior string) *ThesisBuilder {
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

type (
	thesisSlugAlreadyExistsError struct {
		slug string
	}

	buildThesisError struct {
		slug string
		err  error
	}

	noSuchThesisError struct {
		slug string
	}

	notAllowedStageError struct {
		keyword string
	}
)

func NewThesisSlugAlreadyExistsError(slug string) error {
	return errors.WithStack(thesisSlugAlreadyExistsError{
		slug: slug,
	})
}

func IsThesisSlugAlreadyExistsError(err error) bool {
	var aerr thesisSlugAlreadyExistsError

	return errors.As(err, &aerr)
}

func (e thesisSlugAlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` thesis already exists", e.slug)
}

func NewBuildThesisError(err error, slug string) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildThesisError{
		slug: slug,
		err:  err,
	})
}

func IsBuildThesisError(err error) bool {
	var berr buildThesisError

	return errors.As(err, &berr)
}

func (e buildThesisError) Cause() error {
	return e.err
}

func (e buildThesisError) Unwrap() error {
	return e.err
}

func (e buildThesisError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildThesisError) CommonError() string {
	return fmt.Sprintf("thesis `%s`", e.slug)
}

func (e buildThesisError) Error() string {
	return fmt.Sprintf("thesis `%s`: %s", e.slug, e.err)
}

func NewNoSuchThesisError(slug string) error {
	return errors.WithStack(noSuchThesisError{
		slug: slug,
	})
}

func IsNoSuchThesisError(err error) bool {
	var nerr noSuchThesisError

	return errors.As(err, &nerr)
}

func (e noSuchThesisError) Error() string {
	return fmt.Sprintf("no such thesis `%s`", e.slug)
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

var (
	errThesisEmptySlug         = errors.New("empty thesis slug")
	errNoThesisHTTPOrAssertion = errors.New("no HTTP or assertion")
)

func NewThesisEmptySlugError() error {
	return errThesisEmptySlug
}

func IsThesisEmptySlugError(err error) bool {
	return errors.Is(err, errThesisEmptySlug)
}

func NewNoThesisHTTPOrAssertionError() error {
	return errNoThesisHTTPOrAssertion
}

func IsNoThesisHTTPOrAssertionError(err error) bool {
	return errors.Is(err, errNoThesisHTTPOrAssertion)
}
