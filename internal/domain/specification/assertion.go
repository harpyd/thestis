package specification

import (
	"fmt"

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
		method  AssertionMethod
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
	count := len(a.asserts)

	if count == 0 {
		return nil
	}

	asserts := make([]Assert, count)
	copy(asserts, a.asserts)

	return asserts
}

func (a Assertion) IsZero() bool {
	return a.method == EmptyAssertionMethod && len(a.asserts) == 0
}

func NewAssert(actual string, expected interface{}) Assert {
	return Assert{
		actual:   actual,
		expected: expected,
	}
}

func (a Assert) Actual() string {
	return a.actual
}

func (a Assert) Expected() interface{} {
	return a.expected
}

func (am AssertionMethod) IsValid() bool {
	switch am {
	case EmptyAssertionMethod:
		return true
	case JSONPath:
		return true
	case UnknownAssertionMethod:
		return false
	}

	return false
}

func (am AssertionMethod) String() string {
	return string(am)
}

func NewAssertionBuilder() *AssertionBuilder {
	return &AssertionBuilder{}
}

func (b *AssertionBuilder) Build() (Assertion, error) {
	var err error

	if !b.method.IsValid() {
		err = NewNotAllowedAssertionMethodError(b.method.String())
	}

	assertion := Assertion{
		method:  b.method,
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

func (b *AssertionBuilder) WithMethod(method AssertionMethod) *AssertionBuilder {
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
	var target buildAssertionError

	return errors.As(err, &target)
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
	var target notAllowedAssertionMethodError

	return errors.As(err, &target)
}

func (e notAllowedAssertionMethodError) Error() string {
	return fmt.Sprintf("assertion method `%s` not allowed", e.method)
}
