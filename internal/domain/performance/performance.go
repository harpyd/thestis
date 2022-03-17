package performance

import (
	"context"
	"fmt"
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

func (p *Performance) ID() string {
	return p.id
}

func (p *Performance) OwnerID() string {
	return p.ownerID
}

func (p *Performance) SpecificationID() string {
	if p.spec == nil {
		return ""
	}

	return p.spec.ID()
}

func (p *Performance) Started() bool {
	return atomic.LoadUint32(&p.state) == locked
}

func (p *Performance) WorkingScenarios() []specification.Scenario {
	if p.spec == nil {
		return nil
	}

	return p.spec.Scenarios()
}

func (p *Performance) MustBeStarted() error {
	if p.Started() {
		return nil
	}

	return NewNotStartedError()
}

// Start asynchronously starts performing of the Performance pipeline.
// Start returns non buffered chan of flow Step's. With Step's you can
// build Flow using flow.Reducer.
//
// Only ONE performing can be start at a time. If one goroutine has captured
// performing, then others calls of Start will be return error that can
// be detected with method IsAlreadyStartedError.
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
		return NewAlreadyStartedError()
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
			NewCanceledError(ctx.Err()),
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

	err := g.Wait()

	switch {
	case IsCanceledError(err):
		steps <- NewScenarioStepWithErr(err, scenario.Slug(), FiredCancel)
	case IsFailedError(err):
		steps <- NewScenarioStepWithErr(err, scenario.Slug(), FiredFail)
	case IsCrashedError(err):
		steps <- NewScenarioStepWithErr(err, scenario.Slug(), FiredCrash)
	case err != nil:
		steps <- NewScenarioStepWithErr(err, scenario.Slug(), NoEvent)
	default:
		steps <- NewScenarioStep(scenario.Slug(), FiredPass)
	}
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
		return Crash(NewNoSatisfyingPerformerError(pt))
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

type (
	canceledError struct {
		err error
	}

	failedError struct {
		err error
	}

	crashedError struct {
		err error
	}
)

func NewCanceledError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(canceledError{err: err})
}

func IsCanceledError(err error) bool {
	var target canceledError

	return errors.As(err, &target)
}

func (e canceledError) Error() string {
	return fmt.Sprintf("performance canceled: %s", e.err)
}

func (e canceledError) Cause() error {
	return e.err
}

func (e canceledError) Unwrap() error {
	return e.err
}

func NewFailedError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(failedError{err: err})
}

func IsFailedError(err error) bool {
	var target failedError

	return errors.As(err, &target)
}

func (e failedError) Error() string {
	return fmt.Sprintf("performance failed: %s", e.err)
}

func (e failedError) Cause() error {
	return e.err
}

func (e failedError) Unwrap() error {
	return e.err
}

func NewCrashedError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(crashedError{err: err})
}

func IsCrashedError(err error) bool {
	var target crashedError

	return errors.As(err, &target)
}

func (e crashedError) Error() string {
	return fmt.Sprintf("performance crashed: %s", e.err)
}

func (e crashedError) Cause() error {
	return e.err
}

func (e crashedError) Unwrap() error {
	return e.err
}

type noSatisfyingPerformerError struct {
	pt PerformerType
}

func NewNoSatisfyingPerformerError(pt PerformerType) error {
	return errors.WithStack(noSatisfyingPerformerError{
		pt: pt,
	})
}

func IsNoSatisfyingPerformerError(err error) bool {
	var target noSatisfyingPerformerError

	return errors.As(err, &target)
}

func (e noSatisfyingPerformerError) Error() string {
	return fmt.Sprintf("no satisfying performer for %s", e.pt)
}
