package performance

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/harpyd/thestis/internal/domain/specification"
)

// Performer carries performing of thesis.
// Performance creators should provide own implementation.
type Performer interface {
	Perform(c *Context, thesis specification.Thesis)
}

type (
	Performance struct {
		started sync.Mutex

		attempts   []Attempt
		performers map[performerType]Performer
		graph      actionGraph
	}

	Option func(p *Performance)
)

type performerType string

const (
	emptyPerformer     performerType = ""
	unknownPerformer   performerType = "!"
	httpPerformer      performerType = "HTTP"
	assertionPerformer performerType = "assertion"
)

type (
	actionGraph map[string]actions

	actions map[string]action

	action struct {
		thesis        specification.Thesis
		performerType performerType

		unlock chan struct{}
	}
)

// WithHTTP registers given Performer as HTTP performer.
func WithHTTP(performer Performer) Option {
	return func(p *Performance) {
		p.performers[httpPerformer] = performer
	}
}

// WithAssertion registers given Performer as assertion performer.
func WithAssertion(performer Performer) Option {
	return func(p *Performance) {
		p.performers[assertionPerformer] = performer
	}
}

func FromSpecification(spec *specification.Specification, opts ...Option) (*Performance, error) {
	graph, err := buildGraph(spec)
	if err != nil {
		return nil, err
	}

	p := &Performance{graph: graph}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
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

	addThesesDependentFakeActions(graph, story, scenario, givens, specification.When)
	addThesesDependentFakeActions(graph, story, scenario, whens, specification.Then)
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

		graph[from][to] = newAction(thesis)
	}
}

func addStageDependentAction(
	graph actionGraph,
	story specification.Story,
	scenario specification.Scenario,
	thesis specification.Thesis,
) {
	var (
		from = thesis.Statement().Stage().String()
		to   = uniqueThesisName(story.Slug(), scenario.Slug(), thesis.Slug())
	)

	initGraphActionsLazy(graph, from)

	graph[from][to] = newAction(thesis)
}

func addThesesDependentFakeActions(
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

			graph[from][nextStage.String()] = newEmptyAction()
		}
	}
}

func uniqueThesisName(storySlug, scenarioSlug, thesisSlug string) string {
	return strings.Join([]string{storySlug, scenarioSlug, thesisSlug}, ".")
}

func thesisPerformerType(thesis specification.Thesis) performerType {
	switch {
	case !thesis.HTTP().IsZero():
		return httpPerformer
	case !thesis.Assertion().IsZero():
		return assertionPerformer
	}

	return unknownPerformer
}

const defaultActionsSize = 1

func initGraphActionsLazy(graph actionGraph, vertex string) {
	if graph[vertex] == nil {
		graph[vertex] = make(actions, defaultActionsSize)
	}
}

func newAction(thesis specification.Thesis) action {
	return action{
		thesis:        thesis,
		performerType: thesisPerformerType(thesis),
		unlock:        make(chan struct{}),
	}
}

func newEmptyAction() action {
	return action{
		performerType: emptyPerformer,
		unlock:        make(chan struct{}),
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

	return checkGraphCyclesDFS(graph, specification.Given.String(), colors)
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

func (p *Performance) Attempts() []Attempt {
	copied := make([]Attempt, len(p.attempts))
	copy(copied, p.attempts)

	return copied
}

func (p *Performance) LastAttempt() Attempt {
	return p.attempts[len(p.attempts)-1]
}

const defaultStreamSize = 1

// Start asynchronously starts performing of Performance action graph.
// Start returns chan of Event with default size equals one.
// Every call of Start creates attempt of performing.
// Only ONE attempt can be start at a time. If one goroutine has captured
// performing, then others will wait for it to complete.
func (p *Performance) Start(ctx context.Context) <-chan Event {
	stream := make(chan Event, defaultStreamSize)

	go p.start(ctx, stream)

	return stream
}

func (p *Performance) start(ctx context.Context, stream chan Event) {
	p.started.Lock()
	defer p.started.Unlock()

	p.attempts = append(p.attempts, newAttempt())

	if err := p.startActions(ctx, stream); err != nil {
		stream <- errEvent{err: err}
	}

	close(stream)
}

func (p *Performance) startActions(ctx context.Context, stream chan Event) error {
	select {
	case <-ctx.Done():
		return NewPerformanceCancelledError()
	default:
	}

	g, ctx := errgroup.WithContext(ctx)

	for from, as := range p.graph {
		for to, a := range as {
			g.Go(p.startActionFn(ctx, stream, from, to, a))
		}
	}

	return g.Wait()
}

func (p *Performance) startActionFn(
	ctx context.Context,
	stream chan Event,
	from, to string,
	a action,
) func() error {
	return func() error {
		return p.startAction(ctx, stream, from, to, a)
	}
}

func (p *Performance) startAction(
	ctx context.Context,
	stream chan Event,
	from, to string,
	a action,
) error {
	if err := p.waitActionLocks(ctx, from); err != nil {
		return err
	}

	p.perform(a)

	stream <- performEvent{
		from:          from,
		to:            to,
		performerType: a.performerType,
	}

	p.unlockAction(from, to)

	return nil
}

func (p *Performance) waitActionLocks(ctx context.Context, to string) error {
	for from := range p.graph {
		select {
		case <-p.graph[from][to].unlock:
		case <-ctx.Done():
			return NewPerformanceCancelledError()
		}
	}

	return nil
}

func (p *Performance) unlockAction(from, to string) {
	close(p.graph[from][to].unlock)
}

func (p *Performance) perform(a action) {
	if a.performerType == emptyPerformer {
		return
	}

	performer, ok := p.performers[a.performerType]
	if !ok {
		return
	}

	performer.Perform(p.LastAttempt().Context(), a.thesis)
}

func (pt performerType) String() string {
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
	return fmt.Sprintf("cyclic performance graph: %s -> %s", e.from, e.to)
}

var errPerformanceCancelled = errors.New("performance cancelled")

func NewPerformanceCancelledError() error {
	return errPerformanceCancelled
}

func IsPerformanceCancelledError(err error) bool {
	return errors.Is(err, errPerformanceCancelled)
}
