package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/app/service/mock"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/specification"
	"github.com/harpyd/thestis/pkg/correlationid"
)

var (
	errPerformanceAcquire = errors.New("performance acquire")
	errPerformanceRelease = errors.New("performance release")
)

func TestNewPerformanceMaintainerPanics(t *testing.T) {
	t.Parallel()

	const flowTimeout = 1 * time.Second

	testCases := []struct {
		Name                   string
		GivenGuard             service.PerformanceGuard
		GivenSubscriber        service.PerformanceCancelSubscriber
		GivenPerformancePolicy service.PerformancePolicy
		GivenEnqueuer          service.Enqueuer
		GivenLogger            service.Logger
		ShouldPanic            bool
		PanicMessage           string
	}{
		{
			Name:                   "all_dependencies_are_not_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          mock.NewEnqueuer(),
			GivenLogger:            mock.NewMemoryLogger(),
			ShouldPanic:            false,
		},
		{
			Name:                   "performance_guard_is_nil",
			GivenGuard:             nil,
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          mock.NewEnqueuer(),
			GivenLogger:            mock.NewMemoryLogger(),
			ShouldPanic:            true,
			PanicMessage:           "performance guard is nil",
		},
		{
			Name:                   "performance_cancel_subscriber_is_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        nil,
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          mock.NewEnqueuer(),
			GivenLogger:            mock.NewMemoryLogger(),
			ShouldPanic:            true,
			PanicMessage:           "performance cancel subscriber is nil",
		},
		{
			Name:                   "steps_policy_is_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: nil,
			GivenEnqueuer:          mock.NewEnqueuer(),
			GivenLogger:            mock.NewMemoryLogger(),
			ShouldPanic:            true,
			PanicMessage:           "performance policy is nil",
		},
		{
			Name:                   "enqueuer_is_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          nil,
			GivenLogger:            mock.NewMemoryLogger(),
			ShouldPanic:            true,
			PanicMessage:           "enqueuer is nil",
		},
		{
			Name:                   "logger_is_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          mock.NewEnqueuer(),
			GivenLogger:            nil,
			ShouldPanic:            true,
			PanicMessage:           "logger is nil",
		},
		{
			Name:                   "all_dependencies_are_nil",
			GivenGuard:             nil,
			GivenSubscriber:        nil,
			GivenPerformancePolicy: nil,
			GivenEnqueuer:          nil,
			ShouldPanic:            true,
			PanicMessage:           "performance guard is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = service.NewPerformanceMaintainer(
					c.GivenGuard,
					c.GivenSubscriber,
					c.GivenPerformancePolicy,
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

func TestMaintainPerformance(t *testing.T) {
	t.Parallel()

	const correlationID = "corr"

	testCases := []struct {
		Name                   string
		PerformanceFactory     func() *performance.Performance
		Guard                  *mock.PerformanceGuard
		ShouldBeErr            bool
		IsErr                  func(err error) bool
		ExpectedAcquireCalls   int
		ExpectedSubscribeCalls int
		ExpectedReleaseCalls   int
	}{
		{
			Name: "performance_acquire_error",
			PerformanceFactory: func() *performance.Performance {
				return performance.Trigger(
					"que",
					(&specification.Builder{}).
						WithID("due").
						ErrlessBuild(),
				)
			},
			Guard:                mock.NewPerformanceGuard(errPerformanceAcquire, nil),
			ShouldBeErr:          true,
			IsErr:                func(err error) bool { return errors.Is(err, errPerformanceAcquire) },
			ExpectedAcquireCalls: 1,
		},
		{
			Name: "performance_is_not_released",
			PerformanceFactory: func() *performance.Performance {
				return performance.Trigger(
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
					performance.WithAssertion(performance.PassingPerformer()),
				)
			},
			Guard:                  mock.NewPerformanceGuard(nil, errPerformanceRelease),
			ShouldBeErr:            false,
			ExpectedAcquireCalls:   1,
			ExpectedReleaseCalls:   1,
			ExpectedSubscribeCalls: 1,
		},
		{
			Name: "successfully_maintain_performance",
			PerformanceFactory: func() *performance.Performance {
				return performance.Trigger(
					"perf",
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
					performance.WithHTTP(performance.PassingPerformer()),
					performance.WithAssertion(performance.PassingPerformer()),
				)
			},
			Guard:                  errlessPerformanceGuard(t),
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
				pubsub   = mock.NewPerformanceCancelPubsub()
				policy   = mock.NewPerformancePolicy()
				enqueuer = mock.NewEnqueuer()
				logger   = mock.NewMemoryLogger()
			)

			maintainer := service.NewPerformanceMaintainer(
				c.Guard, pubsub,
				policy, enqueuer,
				logger, flowTimeout,
			)

			perf := c.PerformanceFactory()

			ctx := correlationid.AssignToCtx(context.Background(), correlationID)

			done, err := maintainer.MaintainPerformance(ctx, perf)

			t.Run("performance_acquired", func(t *testing.T) {
				require.Equal(t, c.ExpectedAcquireCalls, c.Guard.AcquireCalls())
			})

			t.Run("performance_cancellation_subscribed", func(t *testing.T) {
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

			t.Run("performance_enqueued", func(t *testing.T) {
				require.Equal(t, 1, enqueuer.EnqueueCalls())
			})

			t.Run("performance_released", func(t *testing.T) {
				require.Equal(t, c.ExpectedReleaseCalls, c.Guard.ReleaseCalls())
			})
		})
	}
}

func TestCancelWhilePerformanceIsMaintaining(t *testing.T) {
	t.Parallel()

	const (
		correlationID = "corr"
		performanceID = "perf"
	)

	spec := (&specification.Builder{}).
		WithID("perf").
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
				guard    = errlessPerformanceGuard(t)
				pubsub   = mock.NewPerformanceCancelPubsub()
				policy   = mock.NewPerformancePolicy()
				enqueuer = mock.NewEnqueuer()
				logger   = mock.NewMemoryLogger()
			)

			maintainer := service.NewPerformanceMaintainer(
				guard, pubsub,
				policy, enqueuer,
				logger, c.FlowTimeout,
			)

			pass := make(chan struct{})
			defer close(pass)

			perf := performance.Trigger(
				performanceID,
				spec,
				performance.WithHTTP(pendingPassPerformer(t, pass)),
			)

			ctx := correlationid.AssignToCtx(context.Background(), correlationID)

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			done, err := maintainer.MaintainPerformance(ctx, perf)
			require.NoError(t, err)

			if c.PublishCancel {
				err := pubsub.PublishPerformanceCancel(performanceID)
				require.NoError(t, err)
			}

			<-done

			t.Run("release_performance", func(t *testing.T) {
				require.Equal(t, 1, guard.ReleaseCalls())
			})
		})
	}
}

func pendingPassPerformer(t *testing.T, pass <-chan struct{}) performance.Performer {
	t.Helper()

	return performance.PerformerFunc(func(
		ctx context.Context,
		_ *performance.Environment,
		_ specification.Thesis,
	) performance.Result {
		select {
		case <-pass:
		case <-ctx.Done():
			return performance.Cancel(ctx.Err())
		}

		return performance.Pass()
	})
}

func errlessPerformanceGuard(t *testing.T) *mock.PerformanceGuard {
	t.Helper()

	return mock.NewPerformanceGuard(nil, nil)
}
