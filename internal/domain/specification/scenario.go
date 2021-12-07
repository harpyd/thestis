package specification

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	Scenario struct {
		slug        string
		description string
		theses      map[string]Thesis
	}

	ScenarioBuilder struct {
		description     string
		thesisFactories []thesisFactory
	}

	thesisFactory func() (Thesis, error)
)

func (s Scenario) Slug() string {
	return s.slug
}

func (s Scenario) Description() string {
	return s.description
}

func (s Scenario) Theses(slugs ...string) ([]Thesis, error) {
	if shouldGetAll(slugs) {
		return s.allTheses(), nil
	}

	return s.filteredTheses(slugs)
}

func (s Scenario) allTheses() []Thesis {
	theses := make([]Thesis, 0, len(s.theses))

	for _, thesis := range s.theses {
		theses = append(theses, thesis)
	}

	return theses
}

func (s Scenario) filteredTheses(slugs []string) ([]Thesis, error) {
	theses := make([]Thesis, 0, len(slugs))

	var err error

	for _, slug := range slugs {
		if thesis, ok := s.Thesis(slug); ok {
			theses = append(theses, thesis)
		} else {
			err = multierr.Append(err, NewNoSuchThesisError(slug))
		}
	}

	return theses, err
}

func (s Scenario) Thesis(slug string) (thesis Thesis, ok bool) {
	thesis, ok = s.theses[slug]

	return
}

func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{}
}

func (b *ScenarioBuilder) Build(slug string) (Scenario, error) {
	if slug == "" {
		return Scenario{}, NewScenarioEmptySlugError()
	}

	scn := Scenario{
		slug:        slug,
		description: b.description,
		theses:      make(map[string]Thesis),
	}

	if len(b.thesisFactories) == 0 {
		return scn, NewBuildScenarioError(NewNoScenarioThesesError(), slug)
	}

	var err error

	for _, thsisFactory := range b.thesisFactories {
		thsis, thsisErr := thsisFactory()
		if _, ok := scn.theses[thsis.Slug()]; ok {
			err = multierr.Append(err, NewThesisSlugAlreadyExistsError(thsis.Slug()))

			continue
		}

		err = multierr.Append(err, thsisErr)

		scn.theses[thsis.Slug()] = thsis
	}

	return scn, NewBuildScenarioError(err, slug)
}

func (b *ScenarioBuilder) ErrlessBuild(slug string) Scenario {
	s, _ := b.Build(slug)

	return s
}

func (b *ScenarioBuilder) Reset() {
	b.description = ""
	b.thesisFactories = nil
}

func (b *ScenarioBuilder) WithDescription(description string) *ScenarioBuilder {
	b.description = description

	return b
}

func (b *ScenarioBuilder) WithThesis(slug string, buildFn func(b *ThesisBuilder)) *ScenarioBuilder {
	tb := NewThesisBuilder()
	buildFn(tb)

	b.thesisFactories = append(b.thesisFactories, func() (Thesis, error) {
		return tb.Build(slug)
	})

	return b
}

type (
	scenarioSlugAlreadyExistsError struct {
		slug string
	}

	buildScenarioError struct {
		slug string
		err  error
	}

	noSuchScenarioError struct {
		slug string
	}
)

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

func (e buildScenarioError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildScenarioError) CommonError() string {
	return fmt.Sprintf("scenario `%s`", e.slug)
}

func (e buildScenarioError) Error() string {
	return fmt.Sprintf("scenario `%s`: %s", e.slug, e.err)
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

var (
	errScenarioEmptySlug = errors.New("empty scenario slug")
	errNoScenarioTheses  = errors.New("no theses")
)

func NewScenarioEmptySlugError() error {
	return errScenarioEmptySlug
}

func IsScenarioEmptySlugError(err error) bool {
	return errors.Is(err, errScenarioEmptySlug)
}

func NewNoScenarioThesesError() error {
	return errNoScenarioTheses
}

func IsNoScenarioThesesError(err error) bool {
	return errors.Is(err, errNoScenarioTheses)
}
