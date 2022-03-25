package specification

import (
	"fmt"

	"github.com/pkg/errors"
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
	NoAssertionMethod      AssertionMethod = ""
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
	return a.method == NoAssertionMethod && len(a.asserts) == 0
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
	case NoAssertionMethod:
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

func (b *AssertionBuilder) Build() (Assertion, error) {
	var w BuildErrorWrapper

	if !b.method.IsValid() {
		w.WithError(NewNotAllowedAssertionMethodError(b.method))
	}

	assertion := Assertion{
		method:  b.method,
		asserts: make([]Assert, len(b.asserts)),
	}

	copy(assertion.asserts, b.asserts)

	return assertion, w.Wrap("assertion")
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

type NotAllowedAssertionMethodError struct {
	method AssertionMethod
}

func NewNotAllowedAssertionMethodError(method AssertionMethod) error {
	return errors.WithStack(&NotAllowedAssertionMethodError{
		method: method,
	})
}

func (e *NotAllowedAssertionMethodError) Method() AssertionMethod {
	return e.method
}

func (e *NotAllowedAssertionMethodError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("assertion method `%s` not allowed", e.method)
}
