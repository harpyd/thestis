package specification

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	story    = "story"
	scenario = "scenario"
	thesis   = "thesis"
)

const (
	httpMethod      = "HTTP method"
	keyword         = "keyword"
	contentType     = "content type"
	assertionMethod = "assertion method"
)

type elemSlugAlreadyExistsError struct {
	elemName string
	slug     string
}

func NewStorySlugAlreadyExistsError(slug string) error {
	return errors.WithStack(elemSlugAlreadyExistsError{
		elemName: story,
		slug:     slug,
	})
}

func IsStorySlugAlreadyExistsError(err error) bool {
	var aerr elemSlugAlreadyExistsError

	return errors.As(err, &aerr) && aerr.elemName == story
}

func NewScenarioSlugAlreadyExistsError(slug string) error {
	return errors.WithStack(elemSlugAlreadyExistsError{
		elemName: scenario,
		slug:     slug,
	})
}

func IsScenarioSlugAlreadyExistsError(err error) bool {
	var aerr elemSlugAlreadyExistsError

	return errors.As(err, &aerr) && aerr.elemName == scenario
}

func NewThesisSlugAlreadyExistsError(slug string) error {
	return errors.WithStack(elemSlugAlreadyExistsError{
		elemName: thesis,
		slug:     slug,
	})
}

func IsThesisSlugAlreadyExistsError(err error) bool {
	var aerr elemSlugAlreadyExistsError

	return errors.As(err, &aerr) && aerr.elemName == thesis
}

func (e elemSlugAlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` %s already exists", e.slug, e.elemName)
}

type emptyElemSlugError struct {
	elemName string
}

func NewStoryEmptySlugError() error {
	return errors.WithStack(emptyElemSlugError{
		elemName: story,
	})
}

func IsStoryEmptySlugError(err error) bool {
	var eerr emptyElemSlugError

	return errors.As(err, &eerr) && eerr.elemName == story
}

func NewScenarioEmptySlugError() error {
	return errors.WithStack(emptyElemSlugError{
		elemName: "scenario",
	})
}

func IsScenarioEmptySlugError(err error) bool {
	var eerr emptyElemSlugError

	return errors.As(err, &eerr) && eerr.elemName == scenario
}

func NewThesisEmptySlugError() error {
	return errors.WithStack(emptyElemSlugError{
		elemName: "thesis",
	})
}

func IsThesisEmptySlugError(err error) bool {
	var eerr emptyElemSlugError

	return errors.As(err, &eerr) && eerr.elemName == thesis
}

func (e emptyElemSlugError) Error() string {
	return fmt.Sprintf("empty %s slug", e.elemName)
}

type buildSpecificationError struct {
	err error
}

func NewBuildSpecificationError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildSpecificationError{
		err: err,
	})
}

func IsBuildSpecificationError(err error) bool {
	var berr buildSpecificationError

	return errors.As(err, &berr)
}

func (e buildSpecificationError) Cause() error {
	return e.err
}

func (e buildSpecificationError) Unwrap() error {
	return e.err
}

func (e buildSpecificationError) Error() string {
	return fmt.Sprintf("specification: %s", e.err)
}

type buildSlugElemError struct {
	elemName string
	slug     string
	err      error
}

func NewBuildStoryError(err error, slug string) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildSlugElemError{
		elemName: story,
		slug:     slug,
		err:      err,
	})
}

func IsBuildStoryError(err error) bool {
	var berr buildSlugElemError

	return errors.As(err, &berr) && berr.elemName == story
}

func NewBuildScenarioError(err error, slug string) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildSlugElemError{
		elemName: scenario,
		slug:     slug,
		err:      err,
	})
}

func IsBuildScenarioError(err error) bool {
	var berr buildSlugElemError

	return errors.As(err, &berr) && berr.elemName == scenario
}

func NewBuildThesisError(err error, slug string) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildSlugElemError{
		elemName: thesis,
		slug:     slug,
		err:      err,
	})
}

func IsBuildThesisError(err error) bool {
	var berr buildSlugElemError

	return errors.As(err, &berr) && berr.elemName == thesis
}

func (e buildSlugElemError) Cause() error {
	return e.err
}

func (e buildSlugElemError) Unwrap() error {
	return e.err
}

func (e buildSlugElemError) Error() string {
	return fmt.Sprintf("%s `%s`: %s", e.elemName, e.slug, e.err)
}

type noElemWithSlugError struct {
	elemName string
	slug     string
}

func NewNoStoryError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: story,
		slug:     slug,
	})
}

func IsNoStoryError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == story
}

func NewNoScenarioError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: scenario,
		slug:     slug,
	})
}

func IsNoScenarioError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == scenario
}

func NewNoThesisError(slug string) error {
	return errors.WithStack(noElemWithSlugError{
		elemName: thesis,
		slug:     slug,
	})
}

func IsNoThesisError(err error) bool {
	var nerr noElemWithSlugError

	return errors.As(err, &nerr) && nerr.elemName == thesis
}

func (e noElemWithSlugError) Error() string {
	return fmt.Sprintf("no %s `%s`", e.elemName, e.slug)
}

type notAllowedElemError struct {
	elemName string
	unknown  string
}

func NewNotAllowedHTTPMethodError(method string) error {
	return errors.WithStack(notAllowedElemError{
		elemName: httpMethod,
		unknown:  method,
	})
}

func IsNotAllowedHTTPMethodError(err error) bool {
	var uerr notAllowedElemError

	return errors.As(err, &uerr) && uerr.elemName == httpMethod
}

func NewNotAllowedKeywordError(kw string) error {
	return errors.WithStack(notAllowedElemError{
		elemName: keyword,
		unknown:  kw,
	})
}

func IsNotAllowedKeywordError(err error) bool {
	var uerr notAllowedElemError

	return errors.As(err, &uerr) && uerr.elemName == keyword
}

func NewNotAllowedContentTypeError(ct string) error {
	return errors.WithStack(notAllowedElemError{
		elemName: contentType,
		unknown:  ct,
	})
}

func IsNotAllowedContentTypeError(err error) bool {
	var uerr notAllowedElemError

	return errors.As(err, &uerr) && uerr.elemName == contentType
}

func NewNotAllowedAssertionMethodError(method string) error {
	return errors.WithStack(notAllowedElemError{
		elemName: assertionMethod,
		unknown:  method,
	})
}

func IsNotAllowedAssertionMethodError(err error) bool {
	var uerr notAllowedElemError

	return errors.As(err, &uerr) && uerr.elemName == assertionMethod
}

func (e notAllowedElemError) Error() string {
	if e.unknown == "" {
		return fmt.Sprintf("no %s", e.elemName)
	}

	return fmt.Sprintf("%s `%s` not allowed", e.elemName, e.unknown)
}
