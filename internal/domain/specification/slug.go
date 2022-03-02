package specification

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type SlugKind string

const (
	NoSlug       SlugKind = ""
	StorySlug    SlugKind = "story"
	ScenarioSlug SlugKind = "scenario"
	ThesisSlug   SlugKind = "thesis"
)

type Slug struct {
	story    string
	scenario string
	thesis   string

	kind SlugKind
}

func NewStorySlug(slug string) Slug {
	return Slug{
		story: slug,
		kind:  StorySlug,
	}
}

func NewScenarioSlug(storySlug, scenarioSlug string) Slug {
	return Slug{
		story:    storySlug,
		scenario: scenarioSlug,
		kind:     ScenarioSlug,
	}
}

func NewThesisSlug(storySlug, scenarioSlug, thesisSlug string) Slug {
	return Slug{
		story:    storySlug,
		scenario: scenarioSlug,
		thesis:   thesisSlug,
		kind:     ThesisSlug,
	}
}

func (s Slug) Story() string {
	return s.story
}

func (s Slug) Scenario() string {
	return s.scenario
}

func (s Slug) Thesis() string {
	return s.thesis
}

func (s Slug) Kind() SlugKind {
	return s.kind
}

func (s Slug) IsZero() bool {
	return s == Slug{}
}

const (
	emptyReplace   = "*"
	slugsSeparator = "."
)

func (s Slug) String() string {
	switch s.kind {
	case StorySlug:
		return replaceIfEmpty(s.story)
	case ScenarioSlug:
		slugs := mapSlugs([]string{
			s.story,
			s.scenario,
		}, replaceIfEmpty)

		return strings.Join(slugs, slugsSeparator)
	case ThesisSlug:
		slugs := mapSlugs([]string{
			s.story,
			s.scenario,
			s.thesis,
		}, replaceIfEmpty)

		return strings.Join(slugs, slugsSeparator)
	case NoSlug:
	}

	return ""
}

func mapSlugs(slugs []string, fn func(string) string) []string {
	res := make([]string, 0, len(slugs))

	for _, s := range slugs {
		res = append(res, fn(s))
	}

	return res
}

func replaceIfEmpty(s string) string {
	if s == "" {
		return emptyReplace
	}

	return s
}

var errEmptySlug = errors.New("empty slug")

func NewEmptySlugError() error {
	return errEmptySlug
}

func IsEmptySlugError(err error) bool {
	return errors.Is(err, errEmptySlug)
}

type (
	storySlugAlreadyExistsError struct {
		slug string
	}

	scenarioSlugAlreadyExistsError struct {
		slug string
	}

	thesisSlugAlreadyExistsError struct {
		slug string
	}
)

func NewSlugAlreadyExistsError(slug Slug) error {
	switch slug.Kind() {
	case StorySlug:
		return errors.WithStack(storySlugAlreadyExistsError{
			slug: slug.String(),
		})
	case ScenarioSlug:
		return errors.WithStack(scenarioSlugAlreadyExistsError{
			slug: slug.String(),
		})
	case ThesisSlug:
		return errors.WithStack(thesisSlugAlreadyExistsError{
			slug: slug.String(),
		})
	case NoSlug:
	}

	return nil
}

func IsStorySlugAlreadyExistsError(err error) bool {
	var target storySlugAlreadyExistsError

	return errors.As(err, &target)
}

func (e storySlugAlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` story already exists", e.slug)
}

func IsScenarioSlugAlreadyExistsError(err error) bool {
	var target scenarioSlugAlreadyExistsError

	return errors.As(err, &target)
}

func (e scenarioSlugAlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` scenario already exists", e.slug)
}

func IsThesisSlugAlreadyExistsError(err error) bool {
	var target thesisSlugAlreadyExistsError

	return errors.As(err, &target)
}

func (e thesisSlugAlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` thesis already exists", e.slug)
}

type (
	buildStoryError struct {
		slug string
		err  error
	}

	buildScenarioError struct {
		slug string
		err  error
	}

	buildThesisError struct {
		slug string
		err  error
	}
)

func NewBuildSluggedError(err error, slug Slug) error {
	if err == nil {
		return nil
	}

	switch slug.Kind() {
	case StorySlug:
		return errors.WithStack(buildStoryError{
			slug: slug.String(),
			err:  err,
		})
	case ScenarioSlug:
		return errors.WithStack(buildScenarioError{
			slug: slug.String(),
			err:  err,
		})
	case ThesisSlug:
		return errors.WithStack(buildThesisError{
			slug: slug.String(),
			err:  err,
		})
	case NoSlug:
	}

	return nil
}

func IsBuildStoryError(err error) bool {
	var target buildStoryError

	return errors.As(err, &target)
}

func (e buildStoryError) Cause() error {
	return e.err
}

func (e buildStoryError) Unwrap() error {
	return e.err
}

func (e buildStoryError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildStoryError) CommonError() string {
	return fmt.Sprintf("story `%s`", e.slug)
}

func (e buildStoryError) Error() string {
	return fmt.Sprintf("story `%s`: %s", e.slug, e.err)
}

func IsBuildScenarioError(err error) bool {
	var target buildScenarioError

	return errors.As(err, &target)
}

func (e buildScenarioError) Cause() error {
	return e.err
}

func (e buildScenarioError) Unwrap() error {
	return e.err
}

func (e buildScenarioError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildScenarioError) CommonError() string {
	return fmt.Sprintf("scenario `%s`", e.slug)
}

func (e buildScenarioError) Error() string {
	return fmt.Sprintf("scenario `%s`: %s", e.slug, e.err)
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

func (e buildThesisError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildThesisError) CommonError() string {
	return fmt.Sprintf("thesis `%s`", e.slug)
}

func (e buildThesisError) Error() string {
	return fmt.Sprintf("thesis `%s`: %s", e.slug, e.err)
}
