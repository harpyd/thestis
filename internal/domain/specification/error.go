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
	return fmt.Sprintf("no %s with slug `%s`", e.elemName, e.slug)
}

type unknownElemError struct {
	elemName string
	unknown  string
}

func NewUnknownHTTPMethodError(method string) error {
	return errors.WithStack(unknownElemError{
		elemName: "HTTP method",
		unknown:  method,
	})
}

func IsUnknownHTTPMethodError(err error) bool {
	var uerr unknownElemError

	return errors.As(err, &uerr) && uerr.elemName == "HTTP method"
}

func NewUnknownKeywordError(keyword string) error {
	return errors.WithStack(unknownElemError{
		elemName: "keyword",
		unknown:  keyword,
	})
}

func IsUnknownKeywordError(err error) bool {
	var uerr unknownElemError

	return errors.As(err, &uerr) && uerr.elemName == "keyword"
}

func NewUnknownContentTypeError(contentType string) error {
	return errors.WithStack(unknownElemError{
		elemName: "content type",
		unknown:  contentType,
	})
}

func IsUnknownContentTypeError(err error) bool {
	var uerr unknownElemError

	return errors.As(err, &uerr) && uerr.elemName == "content type"
}

func NewUnknownAssertionMethodError(method string) error {
	return errors.WithStack(unknownElemError{
		elemName: "assertion method",
		unknown:  method,
	})
}

func IsUnknownAssertionMethodError(err error) bool {
	var uerr unknownElemError

	return errors.As(err, &uerr) && uerr.elemName == "assertion method"
}

func (e unknownElemError) Error() string {
	return fmt.Sprintf("unknown %s `%s`", e.elemName, e.unknown)
}
