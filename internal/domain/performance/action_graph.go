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

// NewAction is factory function for unmarshalling from database.
func NewAction(from, to string, thesis specification.Thesis, performerType PerformerType) Action {
	return Action{
		from:          from,
		to:            to,
		thesis:        thesis,
		performerType: performerType,
	}
}

// NewActionWithoutThesis is factory function for usage in tests.
func NewActionWithoutThesis(from, to string, performerType PerformerType) Action {
	return Action{
		from:          from,
		to:            to,
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

		graph[a.From()][a.To()] = a
	}

	return graph
}

func buildGraph(spec *specification.Specification) (actionGraph, error) {
	graph := make(actionGraph)

	for _, story := range spec.Stories() {
		scenarios, _ := story.Scenarios()

		for _, scenario := range scenarios {
			theses, _ := scenario.Theses()

			addActions(graph, theses)
		}
	}

	if err := checkGraphCycles(graph); err != nil {
		return nil, err
	}

	return graph, nil
}

func addActions(
	graph actionGraph,
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

		addDependenciesDependentActions(graph, thesis)
		addStageDependentAction(graph, thesis)
	}

	addThesesDependentEmptyActions(graph, givens, specification.When)
	addThesesDependentEmptyActions(graph, whens, specification.Then)
}

func addDependenciesDependentActions(
	graph actionGraph,
	thesis specification.Thesis,
) {
	for _, dep := range thesis.Dependencies() {
		var (
			from = dep.String()
			to   = thesis.Slug().String()
		)

		initGraphActionsLazy(graph, from)

		graph[from][to] = newThesisAction(from, to, thesis)
	}
}

func addStageDependentAction(
	graph actionGraph,
	thesis specification.Thesis,
) {
	var (
		from = uniqueStageName(thesis.Slug().Story(), thesis.Slug().Scenario(), thesis.Statement().Stage())
		to   = thesis.Slug().String()
	)

	initGraphActionsLazy(graph, from)

	graph[from][to] = newThesisAction(from, to, thesis)
}

func addThesesDependentEmptyActions(
	graph actionGraph,
	theses []specification.Thesis,
	nextStage specification.Stage,
) {
	for _, thesis := range theses {
		from := thesis.Slug().String()
		if len(graph[from]) == 0 {
			initGraphActionsLazy(graph, from)

			to := uniqueStageName(thesis.Slug().Story(), thesis.Slug().Scenario(), nextStage)

			graph[from][to] = newEmptyAction(from, to)
		}
	}
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
