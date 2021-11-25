package specification

import (
	"fmt"

	"github.com/pkg/errors"
)

type noElemWithSlugError struct {
	elemName string
	slug     string
}

func NewNoStoryError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: "story",
		slug:     slug,
	})
}

func IsNoStoryError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == "story"
}

func NewNoScenarioError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: "scenario",
		slug:     slug,
	})
}

func IsNoScenarioError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == "scenario"
}

func NewNoThesisError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: "thesis",
		slug:     slug,
	})
}

func IsNoThesisError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == "thesis"
}

func (e noElemWithSlugError) Error() string {
	return fmt.Sprintf("no %s with slug %s", e.elemName, e.slug)
}

type unknownHTTPMethodError struct {
	method string
}

func NewUnknownHTTPMethodError(method string) error {
	return errors.WithStack(unknownHTTPMethodError{
		method: method,
	})
}

func IsUnknownHTTPMethodError(err error) bool {
	var uerr unknownHTTPMethodError

	return errors.As(err, &uerr)
}

func (e unknownHTTPMethodError) Error() string {
	return fmt.Sprintf("unknown HTTP method %s", e.method)
}

type unknownKeywordError struct {
	keyword string
}

func NewUnknownKeywordError(keyword string) error {
	return errors.WithStack(unknownKeywordError{
		keyword: keyword,
	})
}

func IsUnknownKeywordError(err error) bool {
	var uerr unknownKeywordError

	return errors.As(err, &uerr)
}

func (e unknownKeywordError) Error() string {
	return fmt.Sprintf("unknown keyword %s", e.keyword)
}
