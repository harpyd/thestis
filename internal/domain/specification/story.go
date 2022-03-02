package specification

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	Story struct {
		slug        Slug
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

	scenarioFactory func(storySlug Slug) (Scenario, error)
)

func (s Story) Slug() Slug {
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

func (s Story) Scenarios() []Scenario {
	scenarios := make([]Scenario, 0, len(s.scenarios))

	for _, scenario := range s.scenarios {
		scenarios = append(scenarios, scenario)
	}

	return scenarios
}

func (s Story) Scenario(slug string) (scenario Scenario, ok bool) {
	scenario, ok = s.scenarios[slug]

	return
}

func NewStoryBuilder() *StoryBuilder {
	return &StoryBuilder{}
}

func (b *StoryBuilder) Build(slug Slug) (Story, error) {
	if slug.IsZero() {
		return Story{}, NewEmptySlugError()
	}

	stry := Story{
		slug:        slug,
		description: b.description,
		asA:         b.asA,
		inOrderTo:   b.inOrderTo,
		wantTo:      b.wantTo,
		scenarios:   make(map[string]Scenario),
	}

	if len(b.scenarioFactories) == 0 {
		return stry, NewBuildSluggedError(NewNoStoryScenariosError(), slug)
	}

	var err error

	for _, scnFactory := range b.scenarioFactories {
		scn, scnErr := scnFactory(slug)
		if _, ok := stry.scenarios[scn.Slug().Scenario()]; ok {
			err = multierr.Append(err, NewSlugAlreadyExistsError(scn.Slug()))

			continue
		}

		err = multierr.Append(err, scnErr)

		stry.scenarios[scn.Slug().Scenario()] = scn
	}

	return stry, NewBuildSluggedError(err, slug)
}

func (b *StoryBuilder) ErrlessBuild(slug Slug) Story {
	s, _ := b.Build(slug)

	return s
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

	b.scenarioFactories = append(b.scenarioFactories, func(storySlug Slug) (Scenario, error) {
		return sb.Build(NewScenarioSlug(storySlug.Story(), slug))
	})

	return b
}

var errNoStoryScenarios = errors.New("no scenarios")

func NewNoStoryScenariosError() error {
	return errNoStoryScenarios
}

func IsNoStoryScenariosError(err error) bool {
	return errors.Is(err, errNoStoryScenarios)
}
