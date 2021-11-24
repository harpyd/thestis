package specification

type (
	Builder struct {
		Specification
	}

	StoryBuilder struct {
		Story
	}

	ScenarioBuilder struct {
		Scenario
	}

	ThesisBuilder struct {
		Thesis
	}
)

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() *Specification {
	return &b.Specification
}

func (b *Builder) WithAuthor(author string) *Builder {
	b.author = author

	return b
}

func (b *Builder) WithTitle(title string) *Builder {
	b.title = title

	return b
}

func (b *Builder) WithDescription(description string) *Builder {
	b.description = description

	return b
}

func (b *Builder) WithStory(slug string, buildFn func(b *StoryBuilder)) *Builder {
	sb := NewStoryBuilder()
	buildFn(sb)

	story := sb.Build(slug)
	b.Specification.stories[slug] = story

	return b
}

func NewStoryBuilder() *StoryBuilder {
	return &StoryBuilder{}
}

func (b *StoryBuilder) Build(slug string) Story {
	b.Story.slug = slug

	return b.Story
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

	scenario := sb.Build(slug)
	b.Story.scenarios[slug] = scenario

	return b
}

func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{}
}

func (b *ScenarioBuilder) Build(slug string) Scenario {
	b.Scenario.slug = slug

	return b.Scenario
}

func (b *ScenarioBuilder) WithDescription(description string) *ScenarioBuilder {
	b.description = description

	return b
}

func (b *ScenarioBuilder) WithThesis(slug string, buildFn func(b *ThesisBuilder)) *ScenarioBuilder {
	tb := NewThesisBuilder()
	buildFn(tb)

	thesis := tb.Build(slug)
	b.Scenario.theses[slug] = thesis

	return b
}

func NewThesisBuilder() *ThesisBuilder {
	return &ThesisBuilder{}
}

func (b *ThesisBuilder) Build(slug string) Thesis {
	b.Thesis.slug = slug

	return b.Thesis
}
