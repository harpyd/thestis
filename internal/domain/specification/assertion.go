package specification

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	Assertion struct {
		method  AssertionMethod
		asserts []Assert
	}

	Assert struct {
		actual   string
		expected interface{}
	}

	AssertionBuilder struct {
		method  string
		asserts []Assert
	}

	AssertionMethod string
)

const (
	UnknownAssertionMethod AssertionMethod = "!"
	EmptyAssertionMethod   AssertionMethod = ""
	JSONPath               AssertionMethod = "jsonpath"
)

func (a Assertion) Method() AssertionMethod {
	return a.method
}

func (a Assertion) Asserts() []Assert {
	asserts := make([]Assert, len(a.asserts))
	copy(asserts, a.asserts)

	return asserts
}

func (a Assertion) IsZero() bool {
	return a.method == EmptyAssertionMethod && len(a.asserts) == 0
}

func (a Assert) Actual() string {
	return a.actual
}

func (a Assert) Expected() interface{} {
	return a.expected
}

func newAssertionMethodFromString(method string) (AssertionMethod, error) {
	switch strings.ToLower(method) {
	case "":
		return EmptyAssertionMethod, nil
	case "jsonpath":
		return JSONPath, nil
	}

	return UnknownAssertionMethod, NewNotAllowedAssertionMethodError(method)
}

func (am AssertionMethod) String() string {
	return string(am)
}

func NewAssertionBuilder() *AssertionBuilder {
	return &AssertionBuilder{}
}

func (b *AssertionBuilder) Build() (Assertion, error) {
	method, err := newAssertionMethodFromString(b.method)

	assertion := Assertion{
		method:  method,
		asserts: make([]Assert, len(b.asserts)),
	}

	copy(assertion.asserts, b.asserts)

	return assertion, NewBuildAssertionError(err)
}

func (b *AssertionBuilder) ErrlessBuild() Assertion {
	a, _ := b.Build()

	return a
}

func (b *AssertionBuilder) Reset() {
	b.method = ""
	b.asserts = nil
}

func (b *AssertionBuilder) WithMethod(method string) *AssertionBuilder {
	b.method = method

	return b
}

func (b *AssertionBuilder) WithAssert(actual string, expected interface{}) *AssertionBuilder {
	b.asserts = append(b.asserts, Assert{
		actual:   actual,
		expected: expected,
	})

	return b
}

type (
	buildAssertionError struct {
		err error
	}

	notAllowedAssertionMethodError struct {
		method string
	}
)

func NewBuildAssertionError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildAssertionError{
		err: err,
	})
}

func IsBuildAssertionError(err error) bool {
	var berr buildAssertionError

	return errors.As(err, &berr)
}

func (e buildAssertionError) Cause() error {
	return e.err
}

func (e buildAssertionError) Unwrap() error {
	return e.err
}

func (e buildAssertionError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildAssertionError) CommonError() string {
	return "assertion"
}

func (e buildAssertionError) Error() string {
	return fmt.Sprintf("assertion: %s", e.err)
}

func NewNotAllowedAssertionMethodError(method string) error {
	return errors.WithStack(notAllowedAssertionMethodError{
		method: method,
	})
}

func IsNotAllowedAssertionMethodError(err error) bool {
	var nerr notAllowedAssertionMethodError

	return errors.As(err, &nerr)
}

func (e notAllowedAssertionMethodError) Error() string {
	return fmt.Sprintf("assertion method `%s` not allowed", e.method)
}
