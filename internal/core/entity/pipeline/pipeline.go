package pipeline

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/harpyd/thestis/internal/core/entity/specification"
)

type (
	// Pipeline is a ready-to-run pipeline created
	// according to the specification, it can be restarted
	// many times, but if the pipeline is running, it
	// cannot be started until executing is over.
	Pipeline struct {
		id      string
		ownerID string
		spec    *specification.Specification

		executors map[ExecutorType]Executor

		state lockState
	}

	ExecutorRegistrar func(p *Pipeline)
)

type lockState = uint32

const (
	unlocked lockState = iota
	locked
)

type ExecutorType string

const (
	NoExecutor        ExecutorType = ""
	UnknownExecutor   ExecutorType = "!"
	HTTPExecutor      ExecutorType = "HTTP"
	AssertionExecutor ExecutorType = "assertion"
)

// WithHTTP registers given Executor as HTTP.
func WithHTTP(executor Executor) ExecutorRegistrar {
	return func(p *Pipeline) {
		p.executors[HTTPExecutor] = executor
	}
}

// WithAssertion registers given Executor as assertion.
func WithAssertion(executor Executor) ExecutorRegistrar {
	return func(p *Pipeline) {
		p.executors[AssertionExecutor] = executor
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

const defaultExecutorsSize = 2

// Unmarshal transforms Params to Pipeline.
// Unmarshal also receives options like Trigger.
//
// This function is great for converting
// from a database or using in tests.
//
// You must not use this method in
// business code of domain and app layers.
func Unmarshal(params Params, registrars ...ExecutorRegistrar) *Pipeline {
	p := &Pipeline{
		id:        params.ID,
		ownerID:   params.OwnerID,
		spec:      params.Specification,
		executors: make(map[ExecutorType]Executor, defaultExecutorsSize),
		state:     newLockState(params.Started),
	}

	p.applyOpts(registrars)

	return p
}

func newLockState(started bool) lockState {
	if started {
		return locked
	}

	return unlocked
}

// Trigger creates new Pipeline
// from specification.Specification.
//
// Trigger receives options that you're
// free to pass or not. You can pass:
// WithHTTP, WithAssertion.
func Trigger(
	id string,
	spec *specification.Specification,
	registrars ...ExecutorRegistrar,
) *Pipeline {
	p := &Pipeline{
		id:        id,
		ownerID:   "",
		spec:      spec,
		executors: make(map[ExecutorType]Executor, defaultExecutorsSize),
		state:     unlocked,
	}

	if spec != nil {
		p.ownerID = spec.OwnerID()
	}

	p.applyOpts(registrars)

	return p
}

func (p *Pipeline) applyOpts(opts []ExecutorRegistrar) {
	for _, opt := range opts {
		opt(p)
	}
}

// ID returns the identifier of the Pipeline.
func (p *Pipeline) ID() string {
	return p.id
}

// OwnerID return the identifier of the Pipeline owner.
func (p *Pipeline) OwnerID() string {
	return p.ownerID
}

// SpecificationID returns the identifier of the Specification,
// if it isn't nil, else returns empty string.
func (p *Pipeline) SpecificationID() string {
	if p.spec == nil {
		return ""
	}

	return p.spec.ID()
}

// Started indicates whether the Pipeline is running.
func (p *Pipeline) Started() bool {
	return atomic.LoadUint32(&p.state) == locked
}

// WorkingScenarios returns the specification scenarios
// that the Pipeline will run.
func (p *Pipeline) WorkingScenarios() []specification.Scenario {
	if p.spec == nil {
		return nil
	}

	return p.spec.Scenarios()
}

// ShouldBeStarted returns ErrNotStarted if
// the Pipeline is not started.
func (p *Pipeline) ShouldBeStarted() error {
	if p.Started() {
		return nil
	}

	return ErrNotStarted
}

// Start asynchronously starts executing of the Pipeline.
// Start returns non buffered chan of flow Step's. With Step's you can
// build Flow using flow.Reducer.
//
// Only ONE executing can be start at a time. If one goroutine has captured
// executing, then others calls of Start will be return ErrAlreadyStarted.
func (p *Pipeline) Start(ctx context.Context) (<-chan Step, error) {
	if err := p.lock(); err != nil {
		return nil, err
	}

	steps := make(chan Step)

	go p.run(ctx, steps)

	return steps, nil
}

// MustStart is similar to Start, but instead of the error it panics.
func (p *Pipeline) MustStart(ctx context.Context) <-chan Step {
	steps, err := p.Start(ctx)
	if err != nil {
		panic(err)
	}

	return steps
}

func (p *Pipeline) lock() error {
	if !atomic.CompareAndSwapUint32(&p.state, unlocked, locked) {
		return ErrAlreadyStarted
	}

	return nil
}

func (p *Pipeline) unlock() {
	atomic.StoreUint32(&p.state, unlocked)
}

func (p *Pipeline) run(ctx context.Context, steps chan<- Step) {
	defer close(steps)
	defer p.unlock()

	p.runScenarios(ctx, steps)
}

func (p *Pipeline) runScenarios(ctx context.Context, steps chan<- Step) {
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

func (p *Pipeline) runScenario(
	ctx context.Context,
	steps chan<- Step,
	scenario specification.Scenario,
) {
	g, ctx := errgroup.WithContext(ctx)

	var (
		env = NewEnvironment(defaultEnvStoreInitialSize)
		sg  = SyncDependencies(scenario)
	)

	steps <- NewScenarioStep(scenario.Slug(), FiredExecute)

	for _, thesis := range scenario.Theses() {
		g.Go(p.runThesisFn(ctx, steps, env, sg, thesis))
	}

	if err := g.Wait(); err != nil {
		var terr *TerminatedError

		if errors.As(err, &terr) {
			steps <- NewScenarioStepWithErr(err, scenario.Slug(), terr.Event())

			return
		}

		steps <- NewScenarioStepWithErr(err, scenario.Slug(), NoEvent)

		return
	}

	steps <- NewScenarioStep(scenario.Slug(), FiredPass)
}

func (p *Pipeline) runThesisFn(
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

func (p *Pipeline) runThesis(
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

	pt := executorType(thesis)

	steps <- NewThesisStep(thesis.Slug(), pt, FiredExecute)

	result := p.executeThesis(ctx, env, thesis)

	steps <- NewThesisStepWithErr(result.err, thesis.Slug(), pt, result.event)

	return result.err
}

func (p *Pipeline) executeThesis(
	ctx context.Context,
	env *Environment,
	thesis specification.Thesis,
) Result {
	pt := executorType(thesis)

	exec, ok := p.executors[pt]
	if !ok {
		return Crash(NewUndefinedExecutorError(pt))
	}

	return exec.Execute(ctx, env, thesis)
}

func executorType(thesis specification.Thesis) ExecutorType {
	switch {
	case !thesis.HTTP().IsZero():
		return HTTPExecutor
	case !thesis.Assertion().IsZero():
		return AssertionExecutor
	default:
		return UnknownExecutor
	}
}

var (
	ErrAlreadyStarted = errors.New("pipeline already started")
	ErrNotStarted     = errors.New("pipeline not started")
)

type TerminatedError struct {
	err   error
	event Event
}

func WrapWithTerminatedError(err error, event Event) error {
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

	b.WriteString("pipeline has terminated")

	if e.event != NoEvent {
		_, _ = fmt.Fprintf(&b, " due to %q event", e.event)
	}

	if e.err != nil {
		_, _ = fmt.Fprintf(&b, ": %s", e.err)
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

type UndefinedExecutorError struct {
	pt ExecutorType
}

func NewUndefinedExecutorError(executorType ExecutorType) error {
	return errors.WithStack(&UndefinedExecutorError{
		pt: executorType,
	})
}

func (e *UndefinedExecutorError) ExecutorType() ExecutorType {
	return e.pt
}

func (e *UndefinedExecutorError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("undefined executor with `%s` type", e.pt)
}
