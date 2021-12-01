package specification

import (
	"fmt"

	"github.com/pkg/errors"
)

type storySlugAlreadyExistsError struct {
	slug string
}

func NewStorySlugAlreadyExistsError(slug string) error {
	return errors.WithStack(storySlugAlreadyExistsError{
		slug: slug,
	})
}

func IsStorySlugAlreadyExistsError(err error) bool {
	var aerr storySlugAlreadyExistsError

	return errors.As(err, &aerr)
}

func (e storySlugAlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` story already exists", e.slug)
}

type scenarioSlugAlreadyExistsError struct {
	slug string
}

func NewScenarioSlugAlreadyExistsError(slug string) error {
	return errors.WithStack(scenarioSlugAlreadyExistsError{
		slug: slug,
	})
}

func IsScenarioSlugAlreadyExistsError(err error) bool {
	var aerr scenarioSlugAlreadyExistsError

	return errors.As(err, &aerr)
}

func (e scenarioSlugAlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` scenario already exists", e.slug)
}

type thesisSlugAlreadyExistsError struct {
	slug string
}

func NewThesisSlugAlreadyExistsError(slug string) error {
	return errors.WithStack(thesisSlugAlreadyExistsError{
		slug: slug,
	})
}

func IsThesisSlugAlreadyExistsError(err error) bool {
	var aerr thesisSlugAlreadyExistsError

	return errors.As(err, &aerr)
}

func (e thesisSlugAlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` thesis already exists", e.slug)
}

var (
	errStoryEmptySlug    = errors.New("empty story slug")
	errScenarioEmptySlug = errors.New("empty scenario slug")
	errThesisEmptySlug   = errors.New("empty thesis slug")
)

func NewStoryEmptySlugError() error {
	return errors.WithStack(errStoryEmptySlug)
}

func IsStoryEmptySlugError(err error) bool {
	return errors.Is(err, errStoryEmptySlug)
}

func NewScenarioEmptySlugError() error {
	return errors.WithStack(errScenarioEmptySlug)
}

func IsScenarioEmptySlugError(err error) bool {
	return errors.Is(err, errScenarioEmptySlug)
}

func NewThesisEmptySlugError() error {
	return errors.WithStack(errThesisEmptySlug)
}

func IsThesisEmptySlugError(err error) bool {
	return errors.Is(err, errThesisEmptySlug)
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

type buildStoryError struct {
	slug string
	err  error
}

func NewBuildStoryError(err error, slug string) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildStoryError{
		slug: slug,
		err:  err,
	})
}

func IsBuildStoryError(err error) bool {
	var berr buildStoryError

	return errors.As(err, &berr)
}

func (e buildStoryError) Cause() error {
	return e.err
}

func (e buildStoryError) Unwrap() error {
	return e.err
}

func (e buildStoryError) Error() string {
	return fmt.Sprintf("story `%s`: %s", e.slug, e.err)
}

type buildScenarioError struct {
	slug string
	err  error
}

func NewBuildScenarioError(err error, slug string) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildScenarioError{
		slug: slug,
		err:  err,
	})
}

func IsBuildScenarioError(err error) bool {
	var berr buildScenarioError

	return errors.As(err, &berr)
}

func (e buildScenarioError) Cause() error {
	return e.err
}

func (e buildScenarioError) Unwrap() error {
	return e.err
}

func (e buildScenarioError) Error() string {
	return fmt.Sprintf("scenario `%s`: %s", e.slug, e.err)
}

type buildThesisError struct {
	slug string
	err  error
}

func NewBuildThesisError(err error, slug string) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildThesisError{
		slug: slug,
		err:  err,
	})
}

func IsBuildThesisError(err error) bool {
	var berr buildThesisError

	return errors.As(err, &berr)
}

func (e buildThesisError) Cause() error {
	return e.err
}

func (e buildThesisError) Unwrap() error {
	return e.err
}

func (e buildThesisError) Error() string {
	return fmt.Sprintf("thesis `%s`: %s", e.slug, e.err)
}

type noSuchStoryError struct {
	slug string
}

func NewNoSuchStoryError(slug string) error {
	return errors.WithStack(noSuchStoryError{
		slug: slug,
	})
}

func IsNoSuchStoryError(err error) bool {
	var nerr noSuchStoryError

	return errors.As(err, &nerr)
}

func (e noSuchStoryError) Error() string {
	return fmt.Sprintf("no such story `%s`", e.slug)
}

type noSuchScenarioError struct {
	slug string
}

func NewNoSuchScenarioError(slug string) error {
	return errors.WithStack(noSuchScenarioError{
		slug: slug,
	})
}

func IsNoSuchScenarioError(err error) bool {
	var nerr noSuchScenarioError

	return errors.As(err, &nerr)
}

func (e noSuchScenarioError) Error() string {
	return fmt.Sprintf("no such scenario `%s`", e.slug)
}

type noSuchThesisError struct {
	slug string
}

func NewNoSuchThesisError(slug string) error {
	return errors.WithStack(noSuchThesisError{
		slug: slug,
	})
}

func IsNoSuchThesisError(err error) bool {
	var nerr noSuchThesisError

	return errors.As(err, &nerr)
}

func (e noSuchThesisError) Error() string {
	return fmt.Sprintf("no such thesis `%s`", e.slug)
}

type notAllowedHTTPMethodError struct {
	method string
}

func NewNotAllowedHTTPMethodError(method string) error {
	return errors.WithStack(notAllowedHTTPMethodError{
		method: method,
	})
}

func IsNotAllowedHTTPMethodError(err error) bool {
	var nerr notAllowedHTTPMethodError

	return errors.As(err, &nerr)
}

func (e notAllowedHTTPMethodError) Error() string {
	return fmt.Sprintf("HTTP method `%s` not allowed", e.method)
}

type notAllowedKeywordError struct {
	keyword string
}

func NewNotAllowedKeywordError(keyword string) error {
	return errors.WithStack(notAllowedKeywordError{
		keyword: keyword,
	})
}

func IsNotAllowedKeywordError(err error) bool {
	var nerr notAllowedKeywordError

	return errors.As(err, &nerr)
}

func (e notAllowedKeywordError) Error() string {
	if e.keyword == "" {
		return "no keyword"
	}

	return fmt.Sprintf("keyword `%s` not allowed", e.keyword)
}

type notAllowedContentTypeError struct {
	contentType string
}

func NewNotAllowedContentTypeError(contentType string) error {
	return errors.WithStack(notAllowedContentTypeError{
		contentType: contentType,
	})
}

func IsNotAllowedContentTypeError(err error) bool {
	var nerr notAllowedContentTypeError

	return errors.As(err, &nerr)
}

func (e notAllowedContentTypeError) Error() string {
	return fmt.Sprintf("content type `%s` not allowed", e.contentType)
}

type notAllowedAssertionMethodError struct {
	method string
}

func NewNotAllowedAssertionMethodError(method string) error {
	return errors.WithStack(notAllowedAssertionMethodError{
		method: method,
	})
}

func IsNotAllowedAssertionMethodError(err error) bool {
	var nerr notAllowedAssertionMethodError

	return errors.As(err, &nerr)
}

func (e notAllowedAssertionMethodError) Error() string {
	return fmt.Sprintf("assertion method `%s` not allowed", e.method)
}
