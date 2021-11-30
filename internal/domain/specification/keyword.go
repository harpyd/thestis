package specification

import "strings"

type Keyword string

const (
	UnknownKeyword Keyword = "!"
	Given          Keyword = "given"
	When           Keyword = "when"
	Then           Keyword = "then"
)

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
