package specification

import "github.com/pkg/errors"

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

func (s Scenario) ThesesByStages(stages ...Stage) []Thesis {
	theses := make([]Thesis, 0, len(s.theses))

	staged := make(map[Stage]bool, len(stages))
	for _, stage := range stages {
		staged[stage] = true
	}

	for _, thesis := range s.theses {
		if staged[thesis.statement.stage] {
			theses = append(theses, thesis)
		}
	}

	return theses
}

var ErrNoScenarioTheses = errors.New("no theses")

func (b *ScenarioBuilder) Build(slug Slug) (Scenario, error) {
	if err := slug.ShouldBeScenarioKind(); err != nil {
		panic(err)
	}

	scenario := Scenario{
		slug:        slug,
		description: b.description,
		theses:      make(map[string]Thesis),
	}

	var w BuildErrorWrapper

	if len(b.thesisFactories) == 0 {
		w.WithError(ErrNoScenarioTheses)
	}

	for _, thesisFry := range b.thesisFactories {
		thesis, err := thesisFry(slug)
		w.WithError(err)

		if _, ok := scenario.theses[thesis.Slug().Thesis()]; ok {
			w.WithError(NewDuplicatedError(thesis.Slug()))

			continue
		}

		scenario.theses[thesis.Slug().Thesis()] = thesis
	}

	checkThesesDependencies(&w, scenario.theses)
	checkCycleDependencies(&w, scenario)

	return scenario, w.SluggedWrap(slug)
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
	var tb ThesisBuilder

	buildFn(&tb)

	b.thesisFactories = append(b.thesisFactories, func(scenarioSlug Slug) (Thesis, error) {
		return tb.Build(NewThesisSlug(scenarioSlug.Story(), scenarioSlug.Scenario(), slug))
	})

	return b
}

func checkThesesDependencies(w *BuildErrorWrapper, theses map[string]Thesis) {
	for _, thesis := range theses {
		for _, dependency := range thesis.dependencies {
			if _, ok := theses[dependency.Thesis()]; !ok {
				w.WithError(NewUndefinedDependencyError(thesis.slug))
			}
		}
	}
}

func checkCycleDependencies(w *BuildErrorWrapper, scenario Scenario) {
	g := NewGraph(scenario.theses)
	if g.IsCyclic() {
		w.WithError(NewCycleDependenciesError(scenario.Slug()))
	}
}

type Graph struct {
	adj map[string][]string
}

func NewGraph(theses map[string]Thesis) *Graph {
	adj := make(map[string][]string)

	for key, value := range theses {
		for _, dependency := range value.dependencies {
			adj[key] = append(adj[key], dependency.thesis)
		}
	}

	return &Graph{
		adj: adj,
	}
}

func (g *Graph) isCyclicUtil(i string, visited map[string]bool, recStack map[string]bool) bool {
	if recStack[i] {
		return true
	}

	if visited[i] {
		return false
	}

	visited[i] = true
	recStack[i] = true
	children := g.adj[i]

	for _, c := range children {
		if g.isCyclicUtil(c, visited, recStack) {
			return true
		}
	}

	recStack[i] = false

	return false
}

func (g *Graph) IsCyclic() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for key := range g.adj {
		if g.isCyclicUtil(key, visited, recStack) {
			return true
		}
	}

	return false
}
