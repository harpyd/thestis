package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/app/service/mock"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/specification"
	"github.com/harpyd/thestis/pkg/correlationid"
)

var (
	errPipelineAcquire = errors.New("pipeline acquire")
	errPipelineRelease = errors.New("pipeline release")
)

func TestNewPipelineMaintainerPanics(t *testing.T) {
	t.Parallel()

	const flowTimeout = 1 * time.Second

	testCases := []struct {
		Name            string
		GivenGuard      service.PipelineGuard
		GivenSubscriber service.PipelineCancelSubscriber
		GivenPipePolicy service.PipelinePolicy
		GivenEnqueuer   service.Enqueuer
		GivenLogger     service.Logger
		ShouldPanic     bool
		PanicMessage    string
	}{
		{
			Name:            "all_dependencies_are_not_nil",
			GivenGuard:      mock.NewPipelineGuard(nil, nil),
			GivenSubscriber: mock.NewPipelineCancelPubsub(),
			GivenPipePolicy: mock.NewPipelinePolicy(),
			GivenEnqueuer:   mock.NewEnqueuer(),
			GivenLogger:     mock.NewMemoryLogger(),
			ShouldPanic:     false,
		},
		{
			Name:            "pipeline_guard_is_nil",
			GivenGuard:      nil,
			GivenSubscriber: mock.NewPipelineCancelPubsub(),
			GivenPipePolicy: mock.NewPipelinePolicy(),
			GivenEnqueuer:   mock.NewEnqueuer(),
			GivenLogger:     mock.NewMemoryLogger(),
			ShouldPanic:     true,
			PanicMessage:    "pipeline guard is nil",
		},
		{
			Name:            "pipeline_cancel_subscriber_is_nil",
			GivenGuard:      mock.NewPipelineGuard(nil, nil),
			GivenSubscriber: nil,
			GivenPipePolicy: mock.NewPipelinePolicy(),
			GivenEnqueuer:   mock.NewEnqueuer(),
			GivenLogger:     mock.NewMemoryLogger(),
			ShouldPanic:     true,
			PanicMessage:    "pipeline cancel subscriber is nil",
		},
		{
			Name:            "steps_policy_is_nil",
			GivenGuard:      mock.NewPipelineGuard(nil, nil),
			GivenSubscriber: mock.NewPipelineCancelPubsub(),
			GivenPipePolicy: nil,
			GivenEnqueuer:   mock.NewEnqueuer(),
			GivenLogger:     mock.NewMemoryLogger(),
			ShouldPanic:     true,
			PanicMessage:    "pipeline policy is nil",
		},
		{
			Name:            "enqueuer_is_nil",
			GivenGuard:      mock.NewPipelineGuard(nil, nil),
			GivenSubscriber: mock.NewPipelineCancelPubsub(),
			GivenPipePolicy: mock.NewPipelinePolicy(),
			GivenEnqueuer:   nil,
			GivenLogger:     mock.NewMemoryLogger(),
			ShouldPanic:     true,
			PanicMessage:    "enqueuer is nil",
		},
		{
			Name:            "logger_is_nil",
			GivenGuard:      mock.NewPipelineGuard(nil, nil),
			GivenSubscriber: mock.NewPipelineCancelPubsub(),
			GivenPipePolicy: mock.NewPipelinePolicy(),
			GivenEnqueuer:   mock.NewEnqueuer(),
			GivenLogger:     nil,
			ShouldPanic:     true,
			PanicMessage:    "logger is nil",
		},
		{
			Name:            "all_dependencies_are_nil",
			GivenGuard:      nil,
			GivenSubscriber: nil,
			GivenPipePolicy: nil,
			GivenEnqueuer:   nil,
			ShouldPanic:     true,
			PanicMessage:    "pipeline guard is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = service.NewPipelineMaintainer(
					c.GivenGuard,
					c.GivenSubscriber,
					c.GivenPipePolicy,
					c.GivenEnqueuer,
					c.GivenLogger,
					flowTimeout,
				)
			}

			if !c.ShouldPanic {
				require.NotPanics(t, init)

				return
			}

			require.PanicsWithValue(t, c.PanicMessage, init)
		})
	}
}

func TestMaintainPipeline(t *testing.T) {
	t.Parallel()

	const correlationID = "corr"

	testCases := []struct {
		Name                   string
		PipelineFactory        func() *pipeline.Pipeline
		Guard                  *mock.PipelineGuard
		ShouldBeErr            bool
		IsErr                  func(err error) bool
		ExpectedAcquireCalls   int
		ExpectedSubscribeCalls int
		ExpectedReleaseCalls   int
	}{
		{
			Name: "pipeline_acquire_error",
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
					"que",
					(&specification.Builder{}).
						WithID("due").
						ErrlessBuild(),
				)
			},
			Guard:                mock.NewPipelineGuard(errPipelineAcquire, nil),
			ShouldBeErr:          true,
			IsErr:                func(err error) bool { return errors.Is(err, errPipelineAcquire) },
			ExpectedAcquireCalls: 1,
		},
		{
			Name: "pipeline_is_not_released",
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
					"suu",
					(&specification.Builder{}).
						WithID("quu").
						WithStory("foo", func(b *specification.StoryBuilder) {
							b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
								b.WithThesis("baz", func(b *specification.ThesisBuilder) {
									b.WithStatement(specification.Given, "baz")
									b.WithAssertion(func(b *specification.AssertionBuilder) {
										b.WithMethod(specification.JSONPath)
									})
								})
							})
						}).
						ErrlessBuild(),
					pipeline.WithAssertion(pipeline.PassingExecutor()),
				)
			},
			Guard:                  mock.NewPipelineGuard(nil, errPipelineRelease),
			ShouldBeErr:            false,
			ExpectedAcquireCalls:   1,
			ExpectedReleaseCalls:   1,
			ExpectedSubscribeCalls: 1,
		},
		{
			Name: "successfully_maintain_pipeline",
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
					"pipe",
					(&specification.Builder{}).
						WithID("spec").
						WithStory("foo", func(b *specification.StoryBuilder) {
							b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
								b.WithThesis("gaz", func(b *specification.ThesisBuilder) {
									b.WithStatement(specification.Given, "gaz")
									b.WithHTTP(func(b *specification.HTTPBuilder) {
										b.WithRequest(func(b *specification.HTTPRequestBuilder) {
											b.WithURL("https://prepare.com")
											b.WithMethod(specification.POST)
										})
									})
								})
								b.WithThesis("gad", func(b *specification.ThesisBuilder) {
									b.WithStatement(specification.Given, "gad")
									b.WithHTTP(func(b *specification.HTTPBuilder) {
										b.WithRequest(func(b *specification.HTTPRequestBuilder) {
											b.WithURL("https://prepareee.com")
											b.WithMethod(specification.POST)
										})
									})
								})
								b.WithThesis("waz", func(b *specification.ThesisBuilder) {
									b.WithStatement(specification.When, "was")
									b.WithHTTP(func(b *specification.HTTPBuilder) {
										b.WithRequest(func(b *specification.HTTPRequestBuilder) {
											b.WithURL("https://localhost:8000/endpooint")
											b.WithMethod(specification.GET)
										})
									})
								})
								b.WithThesis("taz", func(b *specification.ThesisBuilder) {
									b.WithStatement(specification.Then, "taz")
									b.WithAssertion(func(b *specification.AssertionBuilder) {
										b.WithMethod(specification.JSONPath)
									})
								})
							})
						}).
						ErrlessBuild(),
					pipeline.WithHTTP(pipeline.PassingExecutor()),
					pipeline.WithAssertion(pipeline.PassingExecutor()),
				)
			},
			Guard:                  errlessPipelineGuard(t),
			ShouldBeErr:            false,
			ExpectedAcquireCalls:   1,
			ExpectedSubscribeCalls: 1,
			ExpectedReleaseCalls:   1,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			const flowTimeout = 5 * time.Second

			var (
				pubsub   = mock.NewPipelineCancelPubsub()
				policy   = mock.NewPipelinePolicy()
				enqueuer = mock.NewEnqueuer()
				logger   = mock.NewMemoryLogger()
			)

			maintainer := service.NewPipelineMaintainer(
				c.Guard, pubsub,
				policy, enqueuer,
				logger, flowTimeout,
			)

			pipe := c.PipelineFactory()

			ctx := correlationid.AssignToCtx(context.Background(), correlationID)

			done, err := maintainer.MaintainPipeline(ctx, pipe)

			t.Run("pipeline_acquired", func(t *testing.T) {
				require.Equal(t, c.ExpectedAcquireCalls, c.Guard.AcquireCalls())
			})

			t.Run("pipeline_cancellation_subscribed", func(t *testing.T) {
				require.Equal(t, c.ExpectedSubscribeCalls, pubsub.SubscribeCalls())
			})

			if c.ShouldBeErr {
				t.Run("err", func(t *testing.T) {
					require.True(t, c.IsErr(err))
				})

				return
			}

			require.NoError(t, err)

			<-done

			t.Run("pipeline_enqueued", func(t *testing.T) {
				require.Equal(t, 1, enqueuer.EnqueueCalls())
			})

			t.Run("pipeline_released", func(t *testing.T) {
				require.Equal(t, c.ExpectedReleaseCalls, c.Guard.ReleaseCalls())
			})
		})
	}
}

func TestCancelWhilePipelineIsMaintaining(t *testing.T) {
	t.Parallel()

	const (
		correlationID = "corr"
		pipelineID    = "pipe"
	)

	spec := (&specification.Builder{}).
		WithID("pipe").
		WithStory("foo", func(b *specification.StoryBuilder) {
			b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "baz")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod(specification.POST)
						})
					})
				})
				b.WithThesis("bad", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.When, "bad")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod(specification.GET)
						})
					})
				})
			})
		}).
		ErrlessBuild()

	testCases := []struct {
		Name          string
		FlowTimeout   time.Duration
		PublishCancel bool
	}{
		{
			Name:          "flow_timeout_exceeded",
			FlowTimeout:   5 * time.Millisecond,
			PublishCancel: false,
		},
		{
			Name:          "cancel_published",
			FlowTimeout:   1 * time.Second,
			PublishCancel: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				guard    = errlessPipelineGuard(t)
				pubsub   = mock.NewPipelineCancelPubsub()
				policy   = mock.NewPipelinePolicy()
				enqueuer = mock.NewEnqueuer()
				logger   = mock.NewMemoryLogger()
			)

			maintainer := service.NewPipelineMaintainer(
				guard, pubsub,
				policy, enqueuer,
				logger, c.FlowTimeout,
			)

			pass := make(chan struct{})
			defer close(pass)

			pipe := pipeline.Trigger(
				pipelineID,
				spec,
				pipeline.WithHTTP(pendingPassExecutor(t, pass)),
			)

			ctx := correlationid.AssignToCtx(context.Background(), correlationID)

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			done, err := maintainer.MaintainPipeline(ctx, pipe)
			require.NoError(t, err)

			if c.PublishCancel {
				err := pubsub.PublishPipelineCancel(pipelineID)
				require.NoError(t, err)
			}

			<-done

			t.Run("release_pipeline", func(t *testing.T) {
				require.Equal(t, 1, guard.ReleaseCalls())
			})
		})
	}
}

func pendingPassExecutor(t *testing.T, pass <-chan struct{}) pipeline.Executor {
	t.Helper()

	return pipeline.ExecutorFunc(func(
		ctx context.Context,
		_ *pipeline.Environment,
		_ specification.Thesis,
	) pipeline.Result {
		select {
		case <-pass:
		case <-ctx.Done():
			return pipeline.Cancel(ctx.Err())
		}

		return pipeline.Pass()
	})
}

func errlessPipelineGuard(t *testing.T) *mock.PipelineGuard {
	t.Helper()

	return mock.NewPipelineGuard(nil, nil)
}
