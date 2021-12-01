package specification

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	Story struct {
		slug        string
		description string
		asA         string
		inOrderTo   string
		wantTo      string
		scenarios   map[string]Scenario
	}

	StoryBuilder struct {
		description       string
		asA               string
		inOrderTo         string
		wantTo            string
		scenarioFactories []scenarioFactory
	}

	scenarioFactory func() (Scenario, error)
)

func (s Story) Slug() string {
	return s.slug
}

func (s Story) Description() string {
	return s.description
}

func (s Story) AsA() string {
	return s.asA
}

func (s Story) InOrderTo() string {
	return s.inOrderTo
}

func (s Story) WantTo() string {
	return s.wantTo
}

func (s Story) Scenarios(slugs ...string) ([]Scenario, error) {
	if shouldGetAll(slugs) {
		return s.allScenarios(), nil
	}

	return s.filteredScenarios(slugs)
}

func (s Story) allScenarios() []Scenario {
	scenarios := make([]Scenario, 0, len(s.scenarios))

	for _, scenario := range s.scenarios {
		scenarios = append(scenarios, scenario)
	}

	return scenarios
}

func (s Story) filteredScenarios(slugs []string) ([]Scenario, error) {
	scenarios := make([]Scenario, 0, len(slugs))

	var err error

	for _, slug := range slugs {
		if scenario, ok := s.Scenario(slug); ok {
			scenarios = append(scenarios, scenario)
		} else {
			err = multierr.Append(err, NewNoSuchScenarioError(slug))
		}
	}

	return scenarios, err
}

func (s Story) Scenario(slug string) (scenario Scenario, ok bool) {
	scenario, ok = s.scenarios[slug]

	return
}

func NewStoryBuilder() *StoryBuilder {
	return &StoryBuilder{}
}

func (b *StoryBuilder) Build(slug string) (Story, error) {
	if slug == "" {
		return Story{}, NewStoryEmptySlugError()
	}

	stry := Story{
		slug:        slug,
		description: b.description,
		asA:         b.asA,
		inOrderTo:   b.inOrderTo,
		wantTo:      b.wantTo,
		scenarios:   make(map[string]Scenario),
	}

	var err error

	for _, scnFactory := range b.scenarioFactories {
		scn, scnErr := scnFactory()
		if _, ok := stry.scenarios[scn.Slug()]; ok {
			err = multierr.Append(err, NewScenarioSlugAlreadyExistsError(scn.Slug()))

			continue
		}

		err = multierr.Append(err, scnErr)

		stry.scenarios[scn.Slug()] = scn
	}

	return stry, NewBuildStoryError(err, slug)
}

func (b *StoryBuilder) Reset() {
	b.description = ""
	b.asA = ""
	b.inOrderTo = ""
	b.wantTo = ""
	b.scenarioFactories = nil
}

func (b *StoryBuilder) WithDescription(description string) *StoryBuilder {
	b.description = description

	return b
}

func (b *StoryBuilder) WithAsA(asA string) *StoryBuilder {
	b.asA = asA

	return b
}

func (b *StoryBuilder) WithInOrderTo(inOrderTo string) *StoryBuilder {
	b.inOrderTo = inOrderTo

	return b
}

func (b *StoryBuilder) WithWantTo(wantTo string) *StoryBuilder {
	b.wantTo = wantTo

	return b
}

func (b *StoryBuilder) WithScenario(slug string, buildFn func(b *ScenarioBuilder)) *StoryBuilder {
	sb := NewScenarioBuilder()
	buildFn(sb)

	b.scenarioFactories = append(b.scenarioFactories, func() (Scenario, error) {
		return sb.Build(slug)
	})

	return b
}

type (
	storySlugAlreadyExistsError struct {
		slug string
	}

	buildStoryError struct {
		slug string
		err  error
	}

	noSuchStoryError struct {
		slug string
	}
)

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

var errStoryEmptySlug = errors.New("empty story slug")

func NewStoryEmptySlugError() error {
	return errStoryEmptySlug
}

func IsStoryEmptySlugError(err error) bool {
	return errors.Is(err, errStoryEmptySlug)
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
