package specification

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	Thesis struct {
		slug      string
		after     []string
		statement Statement
		http      HTTP
		assertion Assertion
	}

	Statement struct {
		keyword  Keyword
		behavior string
	}

	ThesisBuilder struct {
		after            []string
		keyword          string
		behavior         string
		httpBuilder      *HTTPBuilder
		assertionBuilder *AssertionBuilder
	}

	Keyword string
)

const (
	UnknownKeyword Keyword = "!"
	Given          Keyword = "given"
	When           Keyword = "when"
	Then           Keyword = "then"
)

func (t Thesis) Slug() string {
	return t.slug
}

func (t Thesis) After() []string {
	return t.after
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

func (s Statement) Keyword() Keyword {
	return s.keyword
}

func (s Statement) Behavior() string {
	return s.behavior
}

func newKeywordFromString(keyword string) (Keyword, error) {
	switch strings.ToLower(keyword) {
	case "given":
		return Given, nil
	case "when":
		return When, nil
	case "then":
		return Then, nil
	}

	return UnknownKeyword, NewNotAllowedKeywordError(keyword)
}

func (k Keyword) String() string {
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

	kw, keywordErr := newKeywordFromString(b.keyword)
	http, httpErr := b.httpBuilder.Build()
	assertion, assertionErr := b.assertionBuilder.Build()

	err := multierr.Combine(httpErr, assertionErr)
	if err == nil && http.IsZero() && assertion.IsZero() {
		err = NewNoThesisHTTPOrAssertionError()
	}

	thsis := Thesis{
		slug:  slug,
		after: make([]string, len(b.after)),
		statement: Statement{
			keyword:  kw,
			behavior: b.behavior,
		},
		http:      http,
		assertion: assertion,
	}

	copy(thsis.after, b.after)

	return thsis, NewBuildThesisError(multierr.Combine(keywordErr, err), slug)
}

func (b *ThesisBuilder) ErrlessBuild(slug string) Thesis {
	t, _ := b.Build(slug)

	return t
}

func (b *ThesisBuilder) Reset() {
	b.after = nil
	b.keyword = ""
	b.behavior = ""
	b.assertionBuilder.Reset()
	b.httpBuilder.Reset()
}

func (b *ThesisBuilder) WithAfter(after string) *ThesisBuilder {
	b.after = append(b.after, after)

	return b
}

func (b *ThesisBuilder) WithStatement(keyword string, behavior string) *ThesisBuilder {
	b.keyword = keyword
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

	notAllowedKeywordError struct {
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

func NewNotAllowedKeywordError(keyword string) error {
	return errors.WithStack(notAllowedKeywordError{
		keyword: keyword,
	})
}

func IsNotAllowedKeywordError(err error) bool {
	var nerr notAllowedKeywordError

	return errors.As(err, &nerr)
}

func (e notAllowedKeywordError) Error() string {
	if e.keyword == "" {
		return "no keyword"
	}

	return fmt.Sprintf("keyword `%s` not allowed", e.keyword)
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
