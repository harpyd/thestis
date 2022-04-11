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
	return copyAsserts(a.asserts)
}

func (a Assertion) validate() error {
	var w BuildErrorWrapper

	if !a.method.IsValid() {
		w.WithError(NewNotAllowedAssertionMethodError(a.method))
	}

	return w.Wrap("assertion")
}

func copyAsserts(asserts []Assert) []Assert {
	if len(asserts) == 0 {
		return nil
	}

	result := make([]Assert, len(asserts))
	copy(result, asserts)

	return result
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

func (b *AssertionBuilder) Build() Assertion {
	return Assertion{
		method:  b.method,
		asserts: assertsOrNil(b.asserts),
	}
}

func assertsOrNil(asserts []Assert) []Assert {
	if len(asserts) == 0 {
		return nil
	}

	return asserts
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

	return fmt.Sprintf("assertion method %q not allowed", e.method)
}
