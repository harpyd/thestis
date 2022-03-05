package specification

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	Scenario struct {
		slug        Slug
		description string
		theses      map[string]Thesis
	}

	ScenarioBuilder struct {
		description     string
		thesisFactories []thesisFactory
	}

	thesisFactory func(scenarioSlug Slug) (Thesis, error)
)

func (s Scenario) Slug() Slug {
	return s.slug
}

func (s Scenario) Description() string {
	return s.description
}

func (s Scenario) Theses() []Thesis {
	theses := make([]Thesis, 0, len(s.theses))

	for _, thesis := range s.theses {
		theses = append(theses, thesis)
	}

	return theses
}

func (s Scenario) Thesis(slug string) (thesis Thesis, ok bool) {
	thesis, ok = s.theses[slug]

	return
}

func (s Scenario) ThesesByStage(stage Stage) []Thesis {
	theses := make([]Thesis, 0, len(s.theses))

	for _, thesis := range s.theses {
		if thesis.statement.stage == stage {
			theses = append(theses, thesis)
		}
	}

	return theses
}

func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{}
}

func (b *ScenarioBuilder) Build(slug Slug) (Scenario, error) {
	if slug.IsZero() {
		return Scenario{}, NewEmptySlugError()
	}

	if err := slug.MustBeScenarioKind(); err != nil {
		return Scenario{}, err
	}

	scenario := Scenario{
		slug:        slug,
		description: b.description,
		theses:      make(map[string]Thesis),
	}

	if len(b.thesisFactories) == 0 {
		return scenario, NewBuildSluggedError(NewNoScenarioThesesError(), slug)
	}

	var err error

	for _, thesisFry := range b.thesisFactories {
		thesis, thesisErr := thesisFry(slug)
		if _, ok := scenario.theses[thesis.Slug().Thesis()]; ok {
			err = multierr.Append(err, NewSlugAlreadyExistsError(thesis.Slug()))

			continue
		}

		err = multierr.Append(err, thesisErr)

		scenario.theses[thesis.Slug().Thesis()] = thesis
	}

	return scenario, NewBuildSluggedError(err, slug)
}

func (b *ScenarioBuilder) ErrlessBuild(slug Slug) Scenario {
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

	b.thesisFactories = append(b.thesisFactories, func(scenarioSlug Slug) (Thesis, error) {
		return tb.Build(NewThesisSlug(scenarioSlug.Story(), scenarioSlug.Scenario(), slug))
	})

	return b
}

var errNoScenarioTheses = errors.New("no theses")

func NewNoScenarioThesesError() error {
	return errNoScenarioTheses
}

func IsNoScenarioThesesError(err error) bool {
	return errors.Is(err, errNoScenarioTheses)
}
