package performance

import (
	"context"
	"fmt"

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
		attempts    []Attempt
		performers  map[performerType]Performer
		actionGraph actionGraph

		ready chan bool
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

const defaultPerformersSize = 2

func FromSpecification(spec *specification.Specification, opts ...Option) (*Performance, error) {
	graph, err := buildGraph(spec)
	if err != nil {
		return nil, err
	}

	p := &Performance{
		actionGraph: graph,
		performers:  make(map[performerType]Performer, defaultPerformersSize),
		ready:       make(chan bool, 1),
	}

	for _, opt := range opts {
		opt(p)
	}

	p.ready <- true

	return p, nil
}

func (p *Performance) Attempts() []Attempt {
	copied := make([]Attempt, len(p.attempts))
	copy(copied, p.attempts)

	return copied
}

func (p *Performance) LastAttempt() Attempt {
	return p.attempts[len(p.attempts)-1]
}

// Actions returns flat slice representation of action graph.
func (p *Performance) Actions() []Action {
	copied := make([]Action, 0, len(p.actionGraph))

	for _, as := range p.actionGraph {
		for _, a := range as {
			copied = append(copied, a)
		}
	}

	return copied
}

const defaultStreamSize = 1

// Start asynchronously starts performing of Performance action graph.
// Start returns chan of Event with default size equals one.
// Every call of Start creates attempt of performing.
// Only ONE attempt can be start at a time. If one goroutine has captured
// performing, then others calls of Start will be return error that can
// be detected with method IsPerformanceAlreadyStartedError.
func (p *Performance) Start(ctx context.Context) (<-chan Event, error) {
	select {
	case <-p.ready:
	default:
		return nil, NewPerformanceAlreadyStartedError()
	}

	p.pushNewAttempt()

	stream := make(chan Event, defaultStreamSize)

	go p.start(ctx, stream)

	return stream, nil
}

func (p *Performance) pushNewAttempt() {
	p.attempts = append(p.attempts, newAttempt())
}

func (p *Performance) start(ctx context.Context, stream chan Event) {
	if err := p.startActions(ctx, stream); err != nil {
		stream <- errEvent{err: err}
	}

	close(stream)

	p.ready <- true
}

func (p *Performance) startActions(ctx context.Context, stream chan Event) error {
	select {
	case <-ctx.Done():
		return NewPerformanceCancelledError()
	default:
	}

	g, ctx := errgroup.WithContext(ctx)

	lg := p.actionGraph.toLockGraph()

	for from, as := range p.actionGraph {
		for to, a := range as {
			g.Go(p.startActionFn(ctx, lg, stream, from, to, a))
		}
	}

	return g.Wait()
}

func (p *Performance) startActionFn(
	ctx context.Context,
	lockGraph lockGraph,
	stream chan Event,
	from, to string,
	a Action,
) func() error {
	return func() error {
		return p.startAction(ctx, lockGraph, stream, from, to, a)
	}
}

func (p *Performance) startAction(
	ctx context.Context,
	lockGraph lockGraph,
	stream chan Event,
	from, to string,
	a Action,
) error {
	if err := p.waitActionLocks(ctx, lockGraph, from); err != nil {
		return err
	}

	performed := p.perform(a)

	stream <- actionEvent{
		from:          from,
		to:            to,
		performerType: a.performerType,
		performed:     performed,
	}

	p.unlockAction(lockGraph, from, to)

	return nil
}

func (p *Performance) waitActionLocks(ctx context.Context, lockGraph lockGraph, to string) error {
	for from := range lockGraph {
		lock, ok := lockGraph[from][to]
		if !ok {
			continue
		}

		select {
		case <-lock:
		case <-ctx.Done():
			return NewPerformanceCancelledError()
		}
	}

	return nil
}

func (p *Performance) unlockAction(lockGraph lockGraph, from, to string) {
	close(lockGraph[from][to])
}

func (p *Performance) perform(a Action) (performed bool) {
	if a.performerType == emptyPerformer {
		return
	}

	performer, ok := p.performers[a.performerType]
	if !ok {
		return false
	}

	performer.Perform(p.LastAttempt().Context(), a.thesis)

	return true
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

var (
	errPerformanceCancelled      = errors.New("performance cancelled")
	errPerformanceAlreadyStarted = errors.New("performance already started")
)

func NewPerformanceCancelledError() error {
	return errPerformanceCancelled
}

func IsPerformanceCancelledError(err error) bool {
	return errors.Is(err, errPerformanceCancelled)
}

func NewPerformanceAlreadyStartedError() error {
	return errPerformanceAlreadyStarted
}

func IsPerformanceAlreadyStartedError(err error) bool {
	return errors.Is(err, errPerformanceAlreadyStarted)
}
