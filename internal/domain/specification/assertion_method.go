package specification

import "strings"

type AssertionMethod string

const (
	UnknownAssertionMethod AssertionMethod = ""
	JSONPath               AssertionMethod = "jsonpath"
)

func newAssertionMethodFromString(method string) (AssertionMethod, error) {
	if strings.ToLower(method) == "jsonpath" {
		return JSONPath, nil
	}

	return UnknownAssertionMethod, NewUnknownAssertionMethodError(method)
}
