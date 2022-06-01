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
		ShouldPanic            bool
		PanicMessage           string
	}{
		{
			Name:                   "all_dependencies_are_not_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          mock.NewEnqueuer(),
			ShouldPanic:            false,
		},
		{
			Name:                   "performance_guard_is_nil",
			GivenGuard:             nil,
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          mock.NewEnqueuer(),
			ShouldPanic:            true,
			PanicMessage:           "performance guard is nil",
		},
		{
			Name:                   "performance_cancel_subscriber_is_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        nil,
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          mock.NewEnqueuer(),
			ShouldPanic:            true,
			PanicMessage:           "performance cancel subscriber is nil",
		},
		{
			Name:                   "steps_policy_is_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: nil,
			GivenEnqueuer:          mock.NewEnqueuer(),
			ShouldPanic:            true,
			PanicMessage:           "performance policy is nil",
		},
		{
			Name:                   "enqueuer_is_nil",
			GivenGuard:             mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:        mock.NewPerformanceCancelPubsub(),
			GivenPerformancePolicy: mock.NewPerformancePolicy(),
			GivenEnqueuer:          nil,
			ShouldPanic:            true,
			PanicMessage:           "enqueuer is nil",
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

	testCases := []struct {
		Name               string
		PerformanceFactory func() *performance.Performance
		Guard              *mock.PerformanceGuard
		ShouldBeErr        bool
		IsErr              func(err error) bool
		ExpectedMessages   []service.Message
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
			Guard:       mock.NewPerformanceGuard(errPerformanceAcquire, nil),
			ShouldBeErr: true,
			IsErr:       func(err error) bool { return errors.Is(err, errPerformanceAcquire) },
		},
		{
			Name: "performance_release_error",
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
			Guard: mock.NewPerformanceGuard(nil, errPerformanceRelease),
			ExpectedMessages: []service.Message{
				service.NewMessageFromStep(performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPerform,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.AssertionPerformer,
					performance.FiredPerform,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.AssertionPerformer,
					performance.FiredPass,
				)),
				service.NewMessageFromStep(performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPass,
				)),
				service.NewMessageFromError(errPerformanceRelease),
			},
			ShouldBeErr: false,
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
			Guard: errlessPerformanceGuard(t),
			ExpectedMessages: []service.Message{
				service.NewMessageFromStep(performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPerform,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "gaz"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "gaz"),
					performance.HTTPPerformer,
					performance.FiredPass,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "gad"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "gad"),
					performance.HTTPPerformer,
					performance.FiredPass,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "waz"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "waz"),
					performance.HTTPPerformer,
					performance.FiredPass,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "taz"),
					performance.AssertionPerformer,
					performance.FiredPerform,
				)),
				service.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "taz"),
					performance.AssertionPerformer,
					performance.FiredPass,
				)),
				service.NewMessageFromStep(performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPass,
				)),
			},
			ShouldBeErr: false,
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
			)

			maintainer := service.NewPerformanceMaintainer(
				c.Guard,
				pubsub,
				policy,
				enqueuer,
				flowTimeout,
			)

			perf := c.PerformanceFactory()

			var actualMessages []service.Message

			done, err := maintainer.MaintainPerformance(context.Background(), perf, func(msg service.Message) {
				actualMessages = append(actualMessages, msg)
			})

			if c.ShouldBeErr {
				t.Run("err", func(t *testing.T) {
					require.True(t, c.IsErr(err))
				})

				return
			}

			require.NoError(t, err)

			<-done

			t.Run("acquire_performance", func(t *testing.T) {
				require.Equal(t, 1, c.Guard.AcquireCalls())
			})

			t.Run("subscribe_performance_canceled", func(t *testing.T) {
				require.Equal(t, 1, pubsub.SubscribeCalls())
			})

			t.Run("performing_enqueued", func(t *testing.T) {
				require.Equal(t, 1, enqueuer.EnqueueCalls())
			})

			t.Run("messages", func(t *testing.T) {
				requireMessagesMatch(t, c.ExpectedMessages, actualMessages)
			})

			t.Run("performance_released", func(t *testing.T) {
				require.Equal(t, 1, c.Guard.ReleaseCalls())
			})
		})
	}
}

func TestCancelWhilePerformanceIsMaintaining(t *testing.T) {
	t.Parallel()

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
		Name                      string
		FlowTimeout               time.Duration
		PublishCancel             bool
		ExpectedContainedMessages []service.Message
	}{
		{
			Name:          "flow_timeout_exceeded",
			FlowTimeout:   5 * time.Millisecond,
			PublishCancel: false,
			ExpectedContainedMessages: []service.Message{
				service.NewMessageFromStep(
					performance.NewScenarioStepWithErr(
						performance.WrapWithTerminatedError(
							context.DeadlineExceeded,
							performance.FiredCancel,
						),
						specification.NewScenarioSlug("foo", "bar"),
						performance.FiredCancel,
					),
				),
			},
		},
		{
			Name:          "cancel_published",
			FlowTimeout:   1 * time.Second,
			PublishCancel: true,
			ExpectedContainedMessages: []service.Message{
				service.NewMessageFromStep(
					performance.NewScenarioStepWithErr(
						performance.WrapWithTerminatedError(
							context.Canceled,
							performance.FiredCancel,
						),
						specification.NewScenarioSlug("foo", "bar"),
						performance.FiredCancel,
					),
				),
			},
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
			)

			maintainer := service.NewPerformanceMaintainer(
				guard,
				pubsub,
				policy,
				enqueuer,
				c.FlowTimeout,
			)

			pass := make(chan struct{})
			defer close(pass)

			perf := performance.Trigger(
				"foo",
				spec,
				performance.WithHTTP(pendingPassPerformer(t, pass)),
			)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var actualMessages []service.Message

			done, err := maintainer.MaintainPerformance(ctx, perf, func(msg service.Message) {
				actualMessages = append(actualMessages, msg)
			})
			require.NoError(t, err)

			if c.PublishCancel {
				err := pubsub.PublishPerformanceCancel("foo")
				require.NoError(t, err)
			}

			<-done

			t.Run("messages", func(t *testing.T) {
				requireMessagesContain(t, actualMessages, c.ExpectedContainedMessages...)
			})

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

func requireMessagesMatch(t *testing.T, expected []service.Message, actual []service.Message) {
	t.Helper()

	require.ElementsMatch(
		t,
		mapMessagesToStrings(expected),
		mapMessagesToStrings(actual),
	)
}

func requireMessagesContain(t *testing.T, messages []service.Message, contain ...service.Message) {
	t.Helper()

	require.Subset(
		t,
		mapMessagesToStrings(messages),
		mapMessagesToStrings(contain),
	)
}

func mapMessagesToStrings(messages []service.Message) []string {
	result := make([]string, 0, len(messages))
	for _, msg := range messages {
		result = append(result, msg.String())
	}

	return result
}
