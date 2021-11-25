package specification

type AssertionMethod string

const (
	UnknownAssertionMethod AssertionMethod = ""
	JSONPath               AssertionMethod = "JSONPATH"
)
