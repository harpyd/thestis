package specification

import "github.com/pkg/errors"

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

var ErrNoStoryScenarios = errors.New("no scenarios")

func (b *StoryBuilder) Build(slug Slug) (Story, error) {
	if err := slug.ShouldBeStoryKind(); err != nil {
		panic(err)
	}

	story := Story{
		slug:        slug,
		description: b.description,
		asA:         b.asA,
		inOrderTo:   b.inOrderTo,
		wantTo:      b.wantTo,
		scenarios:   make(map[string]Scenario),
	}

	var w BuildErrorWrapper

	if len(b.scenarioFactories) == 0 {
		w.WithError(ErrNoStoryScenarios)
	}

	for _, scenarioFry := range b.scenarioFactories {
		scenario, err := scenarioFry(slug)
		w.WithError(err)

		if _, ok := story.scenarios[scenario.Slug().Scenario()]; ok {
			w.WithError(NewDuplicatedError(scenario.Slug()))

			continue
		}

		story.scenarios[scenario.Slug().Scenario()] = scenario
	}

	return story, w.SluggedWrap(slug)
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
	var sb ScenarioBuilder

	buildFn(&sb)

	b.scenarioFactories = append(b.scenarioFactories, func(storySlug Slug) (Scenario, error) {
		return sb.Build(NewScenarioSlug(storySlug.Story(), slug))
	})

	return b
}
