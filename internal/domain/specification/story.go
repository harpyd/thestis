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
		description string
		asA         string
		inOrderTo   string
		wantTo      string
		scenarioFns []scenarioFunc
	}

	scenarioFunc func(storySlug Slug) Scenario
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

func (s Story) validate() error {
	var w BuildErrorWrapper

	if len(s.scenarios) == 0 {
		w.WithError(ErrNoStoryScenarios)
	}

	for _, scenario := range s.scenarios {
		w.WithError(scenario.validate())
	}

	return w.SluggedWrap(s.slug)
}

func (b *StoryBuilder) Build(slug Slug) Story {
	if err := slug.ShouldBeStoryKind(); err != nil {
		panic(err)
	}

	return Story{
		slug:        slug,
		description: b.description,
		asA:         b.asA,
		inOrderTo:   b.inOrderTo,
		wantTo:      b.wantTo,
		scenarios:   scenariosOrNil(slug, b.scenarioFns),
	}
}

func scenariosOrNil(storySlug Slug, fns []scenarioFunc) map[string]Scenario {
	if len(fns) == 0 {
		return nil
	}

	scenarios := make(map[string]Scenario, len(fns))

	for _, fn := range fns {
		scenario := fn(storySlug)

		scenarios[scenario.Slug().Scenario()] = scenario
	}

	return scenarios
}

func (b *StoryBuilder) Reset() {
	b.description = ""
	b.asA = ""
	b.inOrderTo = ""
	b.wantTo = ""
	b.scenarioFns = nil
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

	b.scenarioFns = append(b.scenarioFns, func(storySlug Slug) Scenario {
		return sb.Build(NewScenarioSlug(storySlug.Story(), slug))
	})

	return b
}
