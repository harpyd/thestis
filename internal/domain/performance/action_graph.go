package performance

import (
	"strings"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	actionGraph map[string]map[string]Action

	Action struct {
		from          string
		to            string
		thesis        specification.Thesis
		performerType PerformerType
	}
)

func NewAction(from, to string, thesis specification.Thesis, performerType PerformerType) Action {
	return Action{
		from:          from,
		to:            to,
		thesis:        thesis,
		performerType: performerType,
	}
}

func (a Action) From() string {
	return a.from
}

func (a Action) To() string {
	return a.to
}

func (a Action) Thesis() specification.Thesis {
	return a.thesis
}

func (a Action) PerformerType() PerformerType {
	return a.performerType
}

func unmarshalGraph(actions []Action) actionGraph {
	graph := make(actionGraph, len(actions))

	for _, a := range actions {
		initGraphActionsLazy(graph, a.From())

		graph[a.From()][a.To()] = newThesisAction(a.From(), a.To(), a.Thesis())
	}

	return graph
}

func buildGraph(spec *specification.Specification) (actionGraph, error) {
	graph := make(actionGraph)

	stories, _ := spec.Stories()

	for _, story := range stories {
		scenarios, _ := story.Scenarios()

		for _, scenario := range scenarios {
			theses, _ := scenario.Theses()

			addActions(graph, story, scenario, theses)
		}
	}

	if err := checkGraphCycles(graph); err != nil {
		return nil, err
	}

	return graph, nil
}

func addActions(
	graph actionGraph,
	story specification.Story,
	scenario specification.Scenario,
	theses []specification.Thesis,
) {
	var (
		givens = make([]specification.Thesis, 0, len(theses))
		whens  = make([]specification.Thesis, 0, len(theses))
	)

	for _, thesis := range theses {
		if thesis.Statement().Stage() == specification.Given {
			givens = append(givens, thesis)
		} else if thesis.Statement().Stage() == specification.When {
			whens = append(whens, thesis)
		}

		addDependenciesDependentActions(graph, story, scenario, thesis)
		addStageDependentAction(graph, story, scenario, thesis)
	}

	addThesesDependentEmptyActions(graph, story, scenario, givens, specification.When)
	addThesesDependentEmptyActions(graph, story, scenario, whens, specification.Then)
}

func addDependenciesDependentActions(
	graph actionGraph,
	story specification.Story,
	scenario specification.Scenario,
	thesis specification.Thesis,
) {
	for _, dep := range thesis.Dependencies() {
		var (
			from = uniqueThesisName(story.Slug(), scenario.Slug(), dep)
			to   = uniqueThesisName(story.Slug(), scenario.Slug(), thesis.Slug())
		)

		initGraphActionsLazy(graph, from)

		graph[from][to] = newThesisAction(from, to, thesis)
	}
}

func addStageDependentAction(
	graph actionGraph,
	story specification.Story,
	scenario specification.Scenario,
	thesis specification.Thesis,
) {
	var (
		from = uniqueStageName(story.Slug(), scenario.Slug(), thesis.Statement().Stage())
		to   = uniqueThesisName(story.Slug(), scenario.Slug(), thesis.Slug())
	)

	initGraphActionsLazy(graph, from)

	graph[from][to] = newThesisAction(from, to, thesis)
}

func addThesesDependentEmptyActions(
	graph actionGraph,
	story specification.Story,
	scenario specification.Scenario,
	theses []specification.Thesis,
	nextStage specification.Stage,
) {
	for _, thesis := range theses {
		from := uniqueThesisName(story.Slug(), scenario.Slug(), thesis.Slug())
		if len(graph[from]) == 0 {
			initGraphActionsLazy(graph, from)

			to := uniqueStageName(story.Slug(), scenario.Slug(), nextStage)

			graph[from][to] = newEmptyAction(from, to)
		}
	}
}

func uniqueThesisName(storySlug, scenarioSlug, thesisSlug string) string {
	return strings.Join([]string{storySlug, scenarioSlug, thesisSlug}, ".")
}

func uniqueStageName(storySlug, scenarioSlug string, stage specification.Stage) string {
	if stage == specification.Given {
		return givenStageName()
	}

	return strings.Join([]string{"stage", storySlug, scenarioSlug, stage.String()}, ".")
}

func givenStageName() string {
	return strings.Join([]string{"stage", specification.Given.String()}, ".")
}

func thesisPerformerType(thesis specification.Thesis) PerformerType {
	switch {
	case !thesis.HTTP().IsZero():
		return HTTPPerformer
	case !thesis.Assertion().IsZero():
		return AssertionPerformer
	}

	return UnknownPerformer
}

const defaultActionsSize = 1

func initGraphActionsLazy(graph actionGraph, vertex string) {
	if graph[vertex] == nil {
		graph[vertex] = make(map[string]Action, defaultActionsSize)
	}
}

func newThesisAction(from, to string, thesis specification.Thesis) Action {
	return Action{
		from:          from,
		to:            to,
		thesis:        thesis,
		performerType: thesisPerformerType(thesis),
	}
}

func newEmptyAction(from, to string) Action {
	return Action{
		from:          from,
		to:            to,
		performerType: EmptyPerformer,
	}
}

type vertexColor string

const (
	white vertexColor = ""
	gray  vertexColor = "gray"
	black vertexColor = "black"
)

func checkGraphCycles(graph actionGraph) error {
	colors := make(map[string]vertexColor, len(graph))

	return checkGraphCyclesDFS(graph, givenStageName(), colors)
}

func checkGraphCyclesDFS(
	graph actionGraph,
	from string,
	colors map[string]vertexColor,
) error {
	colors[from] = gray

	for to := range graph[from] {
		if c := colors[to]; c == white {
			if err := checkGraphCyclesDFS(graph, to, colors); err != nil {
				return err
			}
		} else if c == gray {
			return NewCyclicGraphError(from, to)
		}
	}

	colors[from] = black

	return nil
}

type lockGraph map[string]map[string]chan struct{}

func (ag actionGraph) toLockGraph() lockGraph {
	lg := make(map[string]map[string]chan struct{}, len(ag))

	for from, as := range ag {
		lg[from] = make(map[string]chan struct{}, len(as))

		for to := range as {
			lg[from][to] = make(chan struct{})
		}
	}

	return lg
}
