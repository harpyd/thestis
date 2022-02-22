package performance

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	Performance struct {
		id              string
		ownerID         string
		specificationID string

		performers  map[PerformerType]Performer
		actionGraph actionGraph

		lockState lockState
	}

	Option func(p *Performance)
)

type lockState = uint32

const (
	unlocked lockState = iota
	locked
)

type PerformerType string

const (
	EmptyPerformer     PerformerType = ""
	UnknownPerformer   PerformerType = "!"
	HTTPPerformer      PerformerType = "HTTP"
	AssertionPerformer PerformerType = "assertion"
)

// WithID fills Performance ID with given value.
func WithID(id string) Option {
	return func(p *Performance) {
		p.id = id
	}
}

// WithHTTP registers given Performer as HTTP performer.
func WithHTTP(performer Performer) Option {
	return func(p *Performance) {
		p.performers[HTTPPerformer] = performer
	}
}

// WithAssertion registers given Performer as assertion performer.
func WithAssertion(performer Performer) Option {
	return func(p *Performance) {
		p.performers[AssertionPerformer] = performer
	}
}

type (
	Params struct {
		OwnerID         string
		SpecificationID string
		Actions         []Action
		Started         bool
	}
)

const defaultPerformersSize = 2

func Unmarshal(params Params, opts ...Option) *Performance {
	p := &Performance{
		ownerID:         params.OwnerID,
		specificationID: params.SpecificationID,
		actionGraph:     unmarshalGraph(params.Actions),
		performers:      make(map[PerformerType]Performer, defaultPerformersSize),
		lockState:       newLockState(params.Started),
	}

	p.applyOpts(opts)

	return p
}

func newLockState(started bool) lockState {
	if started {
		return locked
	}

	return unlocked
}

func FromSpecification(spec *specification.Specification, opts ...Option) (*Performance, error) {
	graph, err := buildGraph(spec)
	if err != nil {
		return nil, err
	}

	p := &Performance{
		ownerID:         spec.OwnerID(),
		specificationID: spec.ID(),
		actionGraph:     graph,
		performers:      make(map[PerformerType]Performer, defaultPerformersSize),
		lockState:       unlocked,
	}

	p.applyOpts(opts)

	return p, nil
}

func (p *Performance) applyOpts(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

func (p *Performance) ID() string {
	return p.id
}

func (p *Performance) OwnerID() string {
	return p.ownerID
}

func (p *Performance) SpecificationID() string {
	return p.specificationID
}

// Actions returns flat slice representation of action graph.
func (p *Performance) Actions() []Action {
	actions := make([]Action, 0, len(p.actionGraph))

	for _, as := range p.actionGraph {
		for _, a := range as {
			actions = append(actions, a)
		}
	}

	return actions
}

func (p *Performance) Started() bool {
	return atomic.CompareAndSwapUint32(&p.lockState, locked, locked)
}

func (p *Performance) MustBeStarted() error {
	if p.Started() {
		return nil
	}

	return NewNotStartedError()
}

// Start asynchronously starts performing of Performance action graph.
// Start returns chan of flow Step's. With Step's you can build Flow
// using FlowReducer.
//
// Only ONE performing can be start at a time. If one goroutine has captured
// performing, then others calls of Start will be return error that can
// be detected with method IsAlreadyStartedError.
func (p *Performance) Start(ctx context.Context) (<-chan Step, error) {
	if err := p.lock(); err != nil {
		return nil, err
	}

	steps := make(chan Step)

	go p.start(ctx, steps)

	return steps, nil
}

func (p *Performance) start(ctx context.Context, steps chan Step) {
	defer close(steps)
	defer p.unlock()

	p.startActions(ctx, steps)
}

func (p *Performance) lock() error {
	if !atomic.CompareAndSwapUint32(&p.lockState, unlocked, locked) {
		return NewAlreadyStartedError()
	}

	return nil
}

func (p *Performance) unlock() {
	atomic.StoreUint32(&p.lockState, unlocked)
}

const defaultEnvStoreInitialSize = 10

func (p *Performance) startActions(ctx context.Context, steps chan Step) {
	select {
	case <-ctx.Done():
		steps <- NewCanceledStep(newCanceledError(ctx.Err()))

		return
	default:
	}

	env := NewEnvironment(defaultEnvStoreInitialSize)

	lg := p.actionGraph.toLockGraph()

	g, ctx := errgroup.WithContext(ctx)

	for _, as := range p.actionGraph {
		for _, a := range as {
			g.Go(p.startActionFn(ctx, steps, env, lg, a))
		}
	}

	if err := g.Wait(); IsCanceledError(err) {
		steps <- NewCanceledStep(err)
	}
}

func (p *Performance) startActionFn(
	ctx context.Context,
	steps chan Step,
	env *Environment,
	lockGraph lockGraph,
	a Action,
) func() error {
	return func() error {
		return p.startAction(ctx, steps, env, lockGraph, a)
	}
}

func (p *Performance) startAction(
	ctx context.Context,
	steps chan Step,
	env *Environment,
	lockGraph lockGraph,
	a Action,
) error {
	if err := p.waitActionLocks(ctx, lockGraph, a.from); err != nil {
		return err
	}
	defer p.unlockAction(lockGraph, a.from, a.to)

	steps <- NewPerformingStep(a.from, a.to, a.performerType)

	result := p.perform(ctx, env, a)

	if result.State() == Canceled {
		return result.Err()
	}

	steps <- NewStepFromResult(a.from, a.to, a.performerType, result)

	if result.State() == Failed || result.State() == Crashed {
		return errTerminated
	}

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
			return newCanceledError(ctx.Err())
		}
	}

	return nil
}

func (p *Performance) unlockAction(lockGraph lockGraph, from, to string) {
	close(lockGraph[from][to])
}

func (p *Performance) perform(ctx context.Context, env *Environment, a Action) Result {
	if a.performerType == EmptyPerformer {
		return Pass()
	}

	performer, ok := p.performers[a.performerType]
	if !ok {
		return NotPerform()
	}

	return performer.Perform(ctx, env, a.thesis)
}

type cyclicGraphError struct {
	from string
	to   string
}

func NewCyclicGraphError(from, to string) error {
	return errors.WithStack(cyclicGraphError{
		from: from,
		to:   to,
	})
}

func IsCyclicGraphError(err error) bool {
	var target cyclicGraphError

	return errors.As(err, &target)
}

func (e cyclicGraphError) Error() string {
	return fmt.Sprintf("cyclic performance graph: %s -> %s", e.from, e.to)
}

var (
	errTerminated     = errors.New("performance terminated")
	errAlreadyStarted = errors.New("performance already started")
	errNotStarted     = errors.New("performance not started")
)

func NewAlreadyStartedError() error {
	return errAlreadyStarted
}

func IsAlreadyStartedError(err error) bool {
	return errors.Is(err, errAlreadyStarted)
}

func NewNotStartedError() error {
	return errNotStarted
}

func IsNotStartedError(err error) bool {
	return errors.Is(err, errNotStarted)
}
