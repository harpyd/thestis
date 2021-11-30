package specification

import "strings"

type AssertionMethod string

const (
	UnknownAssertionMethod AssertionMethod = "!"
	EmptyAssertionMethod   AssertionMethod = ""
	JSONPath               AssertionMethod = "jsonpath"
)

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
