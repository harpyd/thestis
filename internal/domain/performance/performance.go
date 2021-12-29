package performance

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type PerformerType string

const (
	EmptyPerformer     PerformerType = ""
	UnknownPerformer   PerformerType = "!"
	HTTPPerformer      PerformerType = "HTTP"
	AssertionPerformer PerformerType = "assertion"
)

type (
	Performance struct {
		context    *Context
		performers map[PerformerType]Performer
		graph      actionGraph
	}

	actionGraph map[string]actions

	actions map[string]action

	Performer interface {
		Perform(c *Context, thesis specification.Thesis)
	}

	EventStream chan Event

	Event struct {
		from          string
		to            string
		performerType PerformerType
	}

	Option struct {
		performer     Performer
		performerType PerformerType
	}

	action struct {
		thesis        specification.Thesis
		performerType PerformerType

		fake bool

		unlock chan bool
	}
)

func WithHTTPPerformer(performer Performer) Option {
	return Option{
		performer:     performer,
		performerType: HTTPPerformer,
	}
}

func WithAssertionPerformer(performer Performer) Option {
	return Option{
		performer:     performer,
		performerType: AssertionPerformer,
	}
}

func FromSpecification(spec *specification.Specification, opts ...Option) (*Performance, error) {
	graph, err := buildGraph(spec)
	if err != nil {
		return nil, err
	}

	return &Performance{
		graph:      graph,
		performers: buildPerformers(opts),
	}, nil
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
		if thesis.Statement().Keyword() == specification.Given {
			givens = append(givens, thesis)
		} else if thesis.Statement().Keyword() == specification.When {
			whens = append(whens, thesis)
		}

		added := addActionsDependingOnAfter(graph, story, scenario, thesis)
		if !added {
			addActionDependingOnBDDStageStart(graph, story, scenario, thesis)
		}
	}

	addActionsDependingOnTheses(graph, story, scenario, givens, specification.When)
	addActionsDependingOnTheses(graph, story, scenario, whens, specification.Then)
}

func addActionsDependingOnAfter(
	graph actionGraph,
	story specification.Story,
	scenario specification.Scenario,
	thesis specification.Thesis,
) bool {
	if len(thesis.After()) == 0 {
		return false
	}

	for _, after := range thesis.After() {
		var (
			from = uniqueThesisName(story.Slug(), scenario.Slug(), after)
			to   = uniqueThesisName(story.Slug(), scenario.Slug(), thesis.Slug())
		)

		if graph[from] == nil {
			graph[from] = make(actions, 1)
		}

		graph[from][to] = newAction(thesis)
	}

	return true
}

func addActionDependingOnBDDStageStart(
	graph actionGraph,
	story specification.Story,
	scenario specification.Scenario,
	thesis specification.Thesis,
) {
	var (
		from = thesis.Statement().Keyword().String()
		to   = uniqueThesisName(story.Slug(), scenario.Slug(), thesis.Slug())
	)

	if graph[from] == nil {
		graph[from] = make(actions, 1)
	}

	graph[from][to] = newAction(thesis)
}

func addActionsDependingOnTheses(
	graph actionGraph,
	story specification.Story,
	scenario specification.Scenario,
	theses []specification.Thesis,
	nextStage specification.Keyword,
) {
	for _, thesis := range theses {
		from := uniqueThesisName(story.Slug(), scenario.Slug(), thesis.Slug())
		if len(graph[from]) == 0 {
			if graph[from] == nil {
				graph[from] = make(actions, 1)
			}

			graph[from][nextStage.String()] = newFakeAction()
		}
	}
}

func newAction(thesis specification.Thesis) action {
	return action{
		thesis:        thesis,
		performerType: thesisPerformerType(thesis),
		unlock:        make(chan bool, 1),
	}
}

func newFakeAction() action {
	return action{
		fake:   true,
		unlock: make(chan bool, 1),
	}
}

type vertexColor string

const (
	white vertexColor = ""
	gray  vertexColor = "gray"
	black vertexColor = "black"
)

func checkGraphCycles(unsortedGraph actionGraph) error {
	colors := make(map[string]vertexColor, len(unsortedGraph))

	return checkGraphCyclesDFS(unsortedGraph, specification.Given.String(), colors)
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
			return NewCyclicPerformanceGraphError(from, to)
		}
	}

	colors[from] = black

	return nil
}

func uniqueThesisName(storySlug, scenarioSlug, thesisSlug string) string {
	return strings.Join([]string{storySlug, scenarioSlug, thesisSlug}, ".")
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

func buildPerformers(opts []Option) map[PerformerType]Performer {
	performers := make(map[PerformerType]Performer, len(opts))

	for _, opt := range opts {
		performers[opt.performerType] = opt.performer
	}

	return performers
}

func (p *Performance) Start(stream EventStream) {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		p.unlockGraph(specification.Given.String())
	}()

	for from, as := range p.graph {
		for to, a := range as {
			wg.Add(1)

			go func(from, to string, a action) {
				defer wg.Done()

				p.startAction(stream, from, to, a)
			}(from, to, a)
		}
	}

	wg.Wait()

	close(stream)
}

func (p *Performance) startAction(stream EventStream, from, to string, a action) {
	p.waitGraphLocks(from)

	p.perform(a)

	stream <- Event{
		from:          from,
		to:            to,
		performerType: a.performerType,
	}

	p.unlockGraph(to)
}

func (p *Performance) unlockGraph(from string) {
	for to := range p.graph[from] {
		p.graph[from][to].unlock <- true
	}
}

func (p *Performance) waitGraphLocks(to string) {
	for from := range p.graph {
		if dep, ok := p.graph[from][to]; ok {
			<-dep.unlock
		}
	}
}

func (p *Performance) perform(a action) {
	if a.fake {
		return
	}

	performer, ok := p.performers[a.performerType]
	if !ok {
		return
	}

	performer.Perform(p.context, a.thesis)
}

func (e Event) From() string {
	return e.from
}

func (e Event) To() string {
	return e.to
}

func (e Event) PerformerType() PerformerType {
	return e.performerType
}

func (e Event) String() string {
	return fmt.Sprintf("Performance event `%s -(%s)-> %s`", e.from, e.performerType, e.to)
}

func (pt PerformerType) String() string {
	return string(pt)
}

type cyclicPerformanceGraphError struct {
	from string
	to   string
}

func NewCyclicPerformanceGraphError(from, to string) error {
	return errors.WithStack(cyclicPerformanceGraphError{
		from: from,
		to:   to,
	})
}

func IsCyclicPerformanceGraphError(err error) bool {
	var target cyclicPerformanceGraphError

	return errors.As(err, &target)
}

func (e cyclicPerformanceGraphError) Error() string {
	return fmt.Sprintf("cyclic performance graph: (%s, %s)", e.from, e.to)
}
