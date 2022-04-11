package specification

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type (
	Scenario struct {
		slug        Slug
		description string
		theses      map[string]Thesis
	}

	ScenarioBuilder struct {
		description string
		thesisFns   []thesisFunc
	}

	thesisFunc func(scenarioSlug Slug) Thesis
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
		if staged[thesis.stage] {
			theses = append(theses, thesis)
		}
	}

	return theses
}

var ErrNoScenarioTheses = errors.New("no theses")

func (s Scenario) validate() error {
	var w BuildErrorWrapper

	if len(s.theses) == 0 {
		w.WithError(ErrNoScenarioTheses)
	}

	for _, thesis := range s.theses {
		w.WithError(thesis.validate(s))
	}

	w.WithError(checkCycleDependencies(s))

	return w.SluggedWrap(s.slug)
}

func (b *ScenarioBuilder) Build(slug Slug) Scenario {
	if err := slug.ShouldBeScenarioKind(); err != nil {
		panic(err)
	}

	return Scenario{
		slug:        slug,
		description: b.description,
		theses:      thesesOrNil(slug, b.thesisFns),
	}
}

func thesesOrNil(scenarioSlug Slug, fns []thesisFunc) map[string]Thesis {
	if len(fns) == 0 {
		return nil
	}

	theses := make(map[string]Thesis, len(fns))

	for _, fn := range fns {
		thesis := fn(scenarioSlug)

		theses[thesis.Slug().Thesis()] = thesis
	}

	return theses
}

func (b *ScenarioBuilder) Reset() {
	b.description = ""
	b.thesisFns = nil
}

func (b *ScenarioBuilder) WithDescription(description string) *ScenarioBuilder {
	b.description = description

	return b
}

func (b *ScenarioBuilder) WithThesis(slug string, buildFn func(b *ThesisBuilder)) *ScenarioBuilder {
	var tb ThesisBuilder

	buildFn(&tb)

	b.thesisFns = append(b.thesisFns, func(scenarioSlug Slug) Thesis {
		return tb.Build(NewThesisSlug(scenarioSlug.Story(), scenarioSlug.Scenario(), slug))
	})

	return b
}

func checkCycleDependencies(scenario Scenario) error {
	var w BuildErrorWrapper

	stackTheses := make([]string, 0)
	res := ""

	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for key := range scenario.theses {
		if isCyclic(scenario.theses, key, visited, recStack, stackTheses, &res) {
			w.WithError(NewCycleDependenciesError(scenario.Slug(), res))
			break
		}
	}

	return w.SluggedWrap(scenario.Slug())
}

func isCyclic(theses map[string]Thesis, i string, visited map[string]bool, recStack map[string]bool, path []string, res *string) bool {
	path = append(path, i)

	if recStack[i] {
		*res = strings.Join(path, "->")
		return true
	}

	if visited[i] {
		return false
	}

	visited[i] = true
	recStack[i] = true
	children := theses[i]

	for slug := range children.dependencies {
		if isCyclic(theses, slug.thesis, visited, recStack, path, res) {
			return true
		}
	}

	recStack[i] = false

	return false
}

type CycleDependenciesError struct {
	slug Slug
	path string
}

func NewCycleDependenciesError(slug Slug, path string) error {
	return errors.WithStack(&CycleDependenciesError{
		slug: slug,
		path: path,
	})
}

func (e *CycleDependenciesError) Slug() Slug {
	return e.slug
}

func (e *CycleDependenciesError) Error() string {
	return fmt.Sprintf("cycle dependencies in scenario: %q path:%q", e.slug, e.path)
}
