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

type notAllowedElemError struct {
	elemName string
	unknown  string
}

func NewNotAllowedHTTPMethodError(method string) error {
	return errors.WithStack(notAllowedElemError{
		elemName: "HTTP method",
		unknown:  method,
	})
}

func IsNotAllowedHTTPMethodError(err error) bool {
	var uerr notAllowedElemError

	return errors.As(err, &uerr) && uerr.elemName == "HTTP method"
}

func NewNotAllowedKeywordError(keyword string) error {
	return errors.WithStack(notAllowedElemError{
		elemName: "keyword",
		unknown:  keyword,
	})
}

func IsNotAllowedKeywordError(err error) bool {
	var uerr notAllowedElemError

	return errors.As(err, &uerr) && uerr.elemName == "keyword"
}

func NewNotAllowedContentTypeError(contentType string) error {
	return errors.WithStack(notAllowedElemError{
		elemName: "content type",
		unknown:  contentType,
	})
}

func IsNotAllowedContentTypeError(err error) bool {
	var uerr notAllowedElemError

	return errors.As(err, &uerr) && uerr.elemName == "content type"
}

func NewNotAllowedAssertionMethodError(method string) error {
	return errors.WithStack(notAllowedElemError{
		elemName: "assertion method",
		unknown:  method,
	})
}

func IsNotAllowedAssertionMethodError(err error) bool {
	var uerr notAllowedElemError

	return errors.As(err, &uerr) && uerr.elemName == "assertion method"
}

func (e notAllowedElemError) Error() string {
	return fmt.Sprintf("%s `%s` not allowed", e.elemName, e.unknown)
}
