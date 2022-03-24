package performance

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	// Performance is a ready-to-run pipeline created
	// according to the specification, it can be restarted
	// many times, but if the performance is running, it
	// cannot be started until performing is over.
	Performance struct {
		id      string
		ownerID string
		spec    *specification.Specification

		performers map[PerformerType]Performer

		state lockState
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
	NoPerformer        PerformerType = ""
	UnknownPerformer   PerformerType = "!"
	HTTPPerformer      PerformerType = "HTTP"
	AssertionPerformer PerformerType = "assertion"
)

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
		ID            string
		Specification *specification.Specification
		OwnerID       string
		Started       bool
	}
)

const defaultPerformersSize = 2

// Unmarshal transforms Params to Performance.
// Unmarshal also receives options like FromSpecification.
//
// This function is great for converting
// from a database or using in tests.
//
// You must not use this method in
// business code of domain and app layers.
func Unmarshal(params Params, opts ...Option) *Performance {
	p := &Performance{
		id:         params.ID,
		ownerID:    params.OwnerID,
		spec:       params.Specification,
		performers: make(map[PerformerType]Performer, defaultPerformersSize),
		state:      newLockState(params.Started),
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

// FromSpecification creates new Performance
// from specification.Specification.
//
// FromSpecification receives options that you're
// free to pass or not. You can pass:
// WithHTTP, WithAssertion.
func FromSpecification(
	id string,
	spec *specification.Specification,
	opts ...Option,
) *Performance {
	p := &Performance{
		id:         id,
		ownerID:    "",
		spec:       spec,
		performers: make(map[PerformerType]Performer, defaultPerformersSize),
		state:      unlocked,
	}

	if spec != nil {
		p.ownerID = spec.OwnerID()
	}

	p.applyOpts(opts)

	return p
}

func (p *Performance) applyOpts(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

// ID returns the identifier of the Performance.
func (p *Performance) ID() string {
	return p.id
}

// OwnerID return the identifier of the Performance owner.
func (p *Performance) OwnerID() string {
	return p.ownerID
}

// SpecificationID returns the identifier of the Specification,
// if it isn't nil, else returns empty string.
func (p *Performance) SpecificationID() string {
	if p.spec == nil {
		return ""
	}

	return p.spec.ID()
}

// Started indicates whether the Performance is running.
func (p *Performance) Started() bool {
	return atomic.LoadUint32(&p.state) == locked
}

// WorkingScenarios returns the specification scenarios
// that the Performance will run.
func (p *Performance) WorkingScenarios() []specification.Scenario {
	if p.spec == nil {
		return nil
	}

	return p.spec.Scenarios()
}

// ShouldBeStarted returns ErrNotStarted if
// the Performance is not started.
func (p *Performance) ShouldBeStarted() error {
	if p.Started() {
		return nil
	}

	return ErrNotStarted
}

// Start asynchronously starts performing of the Performance pipeline.
// Start returns non buffered chan of flow Step's. With Step's you can
// build Flow using flow.Reducer.
//
// Only ONE performing can be start at a time. If one goroutine has captured
// performing, then others calls of Start will be return ErrAlreadyStarted.
func (p *Performance) Start(ctx context.Context) (<-chan Step, error) {
	if err := p.lock(); err != nil {
		return nil, err
	}

	steps := make(chan Step)

	go p.run(ctx, steps)

	return steps, nil
}

func (p *Performance) lock() error {
	if !atomic.CompareAndSwapUint32(&p.state, unlocked, locked) {
		return ErrAlreadyStarted
	}

	return nil
}

func (p *Performance) unlock() {
	atomic.StoreUint32(&p.state, unlocked)
}

func (p *Performance) run(ctx context.Context, steps chan<- Step) {
	defer close(steps)
	defer p.unlock()

	p.runScenarios(ctx, steps)
}

func (p *Performance) runScenarios(ctx context.Context, steps chan<- Step) {
	if ctx.Err() != nil {
		steps <- NewScenarioStepWithErr(
			WrapErrorWithTerminated(ctx.Err(), FiredCancel),
			specification.AnyScenarioSlug(),
			FiredCancel,
		)

		return
	}

	var wg sync.WaitGroup

	for _, scenario := range p.WorkingScenarios() {
		wg.Add(1)

		go func(scenario specification.Scenario) {
			defer wg.Done()

			p.runScenario(ctx, steps, scenario)
		}(scenario)
	}

	wg.Wait()
}

const defaultEnvStoreInitialSize = 10

func (p *Performance) runScenario(
	ctx context.Context,
	steps chan<- Step,
	scenario specification.Scenario,
) {
	g, ctx := errgroup.WithContext(ctx)

	var (
		env = NewEnvironment(defaultEnvStoreInitialSize)
		sg  = SyncDependencies(scenario)
	)

	steps <- NewScenarioStep(scenario.Slug(), FiredPerform)

	for _, thesis := range scenario.Theses() {
		g.Go(p.runThesisFn(ctx, steps, env, sg, thesis))
	}

	if err := g.Wait(); err != nil {
		var terr *TerminatedError

		if errors.As(err, &terr) {
			steps <- NewScenarioStepWithErr(err, scenario.Slug(), terr.event)

			return
		}

		steps <- NewScenarioStepWithErr(err, scenario.Slug(), NoEvent)

		return
	}

	steps <- NewScenarioStep(scenario.Slug(), FiredPass)
}

func (p *Performance) runThesisFn(
	ctx context.Context,
	steps chan<- Step,
	env *Environment,
	sg ScenarioSyncGroup,
	thesis specification.Thesis,
) func() error {
	return func() error {
		return p.runThesis(ctx, steps, env, sg, thesis)
	}
}

func (p *Performance) runThesis(
	ctx context.Context,
	steps chan<- Step,
	env *Environment,
	sg ScenarioSyncGroup,
	thesis specification.Thesis,
) error {
	if err := sg.WaitThesisDependencies(ctx, thesis.Slug()); err != nil {
		return err
	}
	defer sg.ThesisDone(thesis.Slug())

	pt := performerType(thesis)

	steps <- NewThesisStep(thesis.Slug(), pt, FiredPerform)

	result := p.performThesis(ctx, env, thesis)

	steps <- NewThesisStepWithErr(result.err, thesis.Slug(), pt, result.event)

	return result.err
}

func (p *Performance) performThesis(
	ctx context.Context,
	env *Environment,
	thesis specification.Thesis,
) Result {
	pt := performerType(thesis)

	performer, ok := p.performers[pt]
	if !ok {
		return Crash(NewRejectedError(pt))
	}

	return performer.Perform(ctx, env, thesis)
}

func performerType(thesis specification.Thesis) PerformerType {
	switch {
	case !thesis.HTTP().IsZero():
		return HTTPPerformer
	case !thesis.Assertion().IsZero():
		return AssertionPerformer
	default:
		return UnknownPerformer
	}
}

var (
	ErrAlreadyStarted = errors.New("performance already started")
	ErrNotStarted     = errors.New("performance not started")
)

type TerminatedError struct {
	err   error
	event Event
}

func WrapErrorWithTerminated(err error, event Event) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(&TerminatedError{
		err:   err,
		event: event,
	})
}

func (e *TerminatedError) Error() string {
	if e == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString("performance has terminated")

	if e.event != NoEvent {
		b.WriteString(" due to `")
		b.WriteString(e.event.String())
		b.WriteString("` event")
	}

	if e.err != nil {
		b.WriteString(": ")
		b.WriteString(e.err.Error())
	}

	return b.String()
}

func (e *TerminatedError) Event() Event {
	return e.event
}

func (e *TerminatedError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.err
}

type RejectedError struct {
	pt PerformerType
}

func NewRejectedError(performerType PerformerType) error {
	return errors.WithStack(&RejectedError{
		pt: performerType,
	})
}

func (e *RejectedError) PerformerType() PerformerType {
	return e.pt
}

func (e *RejectedError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("rejected performer with `%s` type", e.pt)
}
