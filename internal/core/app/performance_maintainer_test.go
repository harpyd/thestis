package app_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/app/mock"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Message        app.Message
		ExpectedEvent  performance.Event
		ExpectedErr    error
		ExpectedString string
	}{
		{
			Message:        app.Message{},
			ExpectedEvent:  performance.NoEvent,
			ExpectedErr:    nil,
			ExpectedString: "",
		},
		{
			Message:        app.NewMessageFromError(errors.New("foo")),
			ExpectedEvent:  performance.NoEvent,
			ExpectedErr:    errors.New("foo"),
			ExpectedString: "foo",
		},
		{
			Message: app.NewMessageFromStep(performance.NewThesisStep(
				specification.NewThesisSlug("foo", "bar", "baz"),
				performance.HTTPPerformer,
				performance.FiredPerform,
			)),
			ExpectedEvent:  performance.FiredPerform,
			ExpectedErr:    nil,
			ExpectedString: "foo.bar.baz: event = perform, type = HTTP",
		},
		{
			Message: app.NewMessageFromStep(performance.NewScenarioStepWithErr(
				errors.New("something wrong"),
				specification.NewScenarioSlug("foo", "bar"),
				performance.FiredCrash,
			)),
			ExpectedEvent:  performance.FiredCrash,
			ExpectedErr:    errors.New("something wrong"),
			ExpectedString: "foo.bar: event = crash, err = something wrong",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			t.Run("event", func(t *testing.T) {
				require.Equal(t, c.ExpectedEvent, c.Message.Event())
			})

			if c.ExpectedErr != nil {
				t.Run("err", func(t *testing.T) {
					require.EqualError(t, c.Message.Err(), c.ExpectedErr.Error())
				})
			} else {
				t.Run("no_err", func(t *testing.T) {
					require.NoError(t, c.Message.Err())
				})
			}

			t.Run("string", func(t *testing.T) {
				require.Equal(t, c.ExpectedString, c.Message.String())
			})
		})
	}
}

var (
	errPerformanceAcquire = errors.New("performance acquire")
	errPerformanceRelease = errors.New("performance release")
)

func TestPanickingNewPerformanceMaintainer(t *testing.T) {
	t.Parallel()

	const flowTimeout = 1 * time.Second

	testCases := []struct {
		Name             string
		GivenGuard       app.PerformanceGuard
		GivenSubscriber  app.PerformanceCancelSubscriber
		GivenStepsPolicy app.StepsPolicy
		GivenEnqueuer    app.Enqueuer
		ShouldPanic      bool
		PanicMessage     string
	}{
		{
			Name:             "all_dependencies_are_not_nil",
			GivenGuard:       mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:  mock.NewPerformanceCancelPubsub(),
			GivenStepsPolicy: mock.NewStepsPolicy(),
			GivenEnqueuer:    mock.NewEnqueuer(),
			ShouldPanic:      false,
		},
		{
			Name:             "performance_guard_is_nil",
			GivenGuard:       nil,
			GivenSubscriber:  mock.NewPerformanceCancelPubsub(),
			GivenStepsPolicy: mock.NewStepsPolicy(),
			GivenEnqueuer:    mock.NewEnqueuer(),
			ShouldPanic:      true,
			PanicMessage:     "performance guard is nil",
		},
		{
			Name:             "performance_cancel_subscriber_is_nil",
			GivenGuard:       mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:  nil,
			GivenStepsPolicy: mock.NewStepsPolicy(),
			GivenEnqueuer:    mock.NewEnqueuer(),
			ShouldPanic:      true,
			PanicMessage:     "performance cancel subscriber is nil",
		},
		{
			Name:             "steps_policy_is_nil",
			GivenGuard:       mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:  mock.NewPerformanceCancelPubsub(),
			GivenStepsPolicy: nil,
			GivenEnqueuer:    mock.NewEnqueuer(),
			ShouldPanic:      true,
			PanicMessage:     "steps policy is nil",
		},
		{
			Name:             "enqueuer_is_nil",
			GivenGuard:       mock.NewPerformanceGuard(nil, nil),
			GivenSubscriber:  mock.NewPerformanceCancelPubsub(),
			GivenStepsPolicy: mock.NewStepsPolicy(),
			GivenEnqueuer:    nil,
			ShouldPanic:      true,
			PanicMessage:     "enqueuer is nil",
		},
		{
			Name:             "all_dependencies_are_nil",
			GivenGuard:       nil,
			GivenSubscriber:  nil,
			GivenStepsPolicy: nil,
			GivenEnqueuer:    nil,
			ShouldPanic:      true,
			PanicMessage:     "performance guard is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = app.NewPerformanceMaintainer(
					c.GivenGuard,
					c.GivenSubscriber,
					c.GivenStepsPolicy,
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
		ExpectedMessages   []app.Message
	}{
		{
			Name: "performance_acquire_error",
			PerformanceFactory: func() *performance.Performance {
				return performance.FromSpecification(
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
				return performance.FromSpecification(
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
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPerform,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.AssertionPerformer,
					performance.FiredPerform,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.AssertionPerformer,
					performance.FiredPass,
				)),
				app.NewMessageFromStep(performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPass,
				)),
				app.NewMessageFromError(errPerformanceRelease),
			},
			ShouldBeErr: false,
		},
		{
			Name: "successfully_maintain_performance",
			PerformanceFactory: func() *performance.Performance {
				return performance.FromSpecification(
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
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPerform,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "gaz"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "gaz"),
					performance.HTTPPerformer,
					performance.FiredPass,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "gad"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "gad"),
					performance.HTTPPerformer,
					performance.FiredPass,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "waz"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "waz"),
					performance.HTTPPerformer,
					performance.FiredPass,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "taz"),
					performance.AssertionPerformer,
					performance.FiredPerform,
				)),
				app.NewMessageFromStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "taz"),
					performance.AssertionPerformer,
					performance.FiredPass,
				)),
				app.NewMessageFromStep(performance.NewScenarioStep(
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
				cancelPubsub = mock.NewPerformanceCancelPubsub()
				stepsPolicy  = mock.NewStepsPolicy()
				enqueuer     = mock.NewEnqueuer()
			)

			maintainer := app.NewPerformanceMaintainer(
				c.Guard,
				cancelPubsub,
				stepsPolicy,
				enqueuer,
				flowTimeout,
			)

			perf := c.PerformanceFactory()

			messages, err := maintainer.MaintainPerformance(context.Background(), perf)

			t.Run("acquire_performance", func(t *testing.T) {
				require.Equal(t, 1, c.Guard.AcquireCalls())
			})

			if c.ShouldBeErr {
				t.Run("err", func(t *testing.T) {
					require.True(t, c.IsErr(err))
				})

				return
			}

			t.Run("subscribe_performance_canceled", func(t *testing.T) {
				require.Equal(t, 1, cancelPubsub.SubscribeCalls())
			})

			t.Run("performing_enqueued", func(t *testing.T) {
				require.Equal(t, 1, enqueuer.EnqueueCalls())
			})

			t.Run("messages", func(t *testing.T) {
				require.NoError(t, err)

				requireMessagesMatch(t, c.ExpectedMessages, messages)
			})

			t.Run("release_performance", func(t *testing.T) {
				require.Equal(t, 1, c.Guard.ReleaseCalls())
			})
		})
	}
}

func TestCancelMaintainPerformance(t *testing.T) {
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
		Name                     string
		CancelContext            bool
		FlowTimeout              time.Duration
		PublishCancel            bool
		ExpectedIncludedMessages []app.Message
	}{
		{
			Name:          "context_canceled",
			CancelContext: true,
			FlowTimeout:   1 * time.Second,
			PublishCancel: false,
			ExpectedIncludedMessages: []app.Message{
				app.NewMessageFromStep(
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
		{
			Name:          "flow_timeout_exceeded",
			CancelContext: false,
			FlowTimeout:   5 * time.Millisecond,
			PublishCancel: false,
			ExpectedIncludedMessages: []app.Message{
				app.NewMessageFromStep(
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
			CancelContext: false,
			FlowTimeout:   1 * time.Second,
			PublishCancel: true,
			ExpectedIncludedMessages: []app.Message{
				app.NewMessageFromStep(
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
				guard        = errlessPerformanceGuard(t)
				cancelPubsub = mock.NewPerformanceCancelPubsub()
				stepsPolicy  = mock.NewStepsPolicy()
				enqueuer     = mock.NewEnqueuer()
			)

			maintainer := app.NewPerformanceMaintainer(
				guard,
				cancelPubsub,
				stepsPolicy,
				enqueuer,
				c.FlowTimeout,
			)

			pass := make(chan struct{})
			defer close(pass)

			perf := performance.FromSpecification(
				"foo",
				spec,
				performance.WithHTTP(pendingPassPerformer(t, pass)),
			)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			messages, err := maintainer.MaintainPerformance(ctx, perf)
			require.NoError(t, err)

			if c.CancelContext {
				cancel()
			}

			if c.PublishCancel {
				err := cancelPubsub.PublishPerformanceCancel("foo")
				require.NoError(t, err)
			}

			t.Run("cancel_messages", func(t *testing.T) {
				requireMessagesContain(t, messages, c.ExpectedIncludedMessages...)
			})

			t.Run("release_performance", func(t *testing.T) {
				require.Equal(t, 1, guard.ReleaseCalls())
			})
		})
	}
}

func TestAsPublishCancelError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError        error
		ShouldBeWrapped   bool
		ExpectedUnwrapped error
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      app.WrapWithPublishCancelError(nil),
			ShouldBeWrapped: false,
		},
		{
			GivenError:        &app.PublishCancelError{},
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: nil,
		},
		{
			GivenError:        app.WrapWithPublishCancelError(errors.New("foo")),
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: errors.New("foo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *app.PublishCancelError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("unwrap", func(t *testing.T) {
					if c.ExpectedUnwrapped != nil {
						require.EqualError(t, target.Unwrap(), c.ExpectedUnwrapped.Error())

						return
					}

					require.NoError(t, target.Unwrap())
				})
			})
		})
	}
}

func TestFormatPublishCancelError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &app.PublishCancelError{},
			ExpectedErrorString: "",
		},
		{
			GivenError:          app.WrapWithPublishCancelError(errors.New("failed")),
			ExpectedErrorString: "publish cancel: failed",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
		})
	}
}

func TestAsSubscribeCancelError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError        error
		ShouldBeWrapped   bool
		ExpectedUnwrapped error
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      app.WrapWithSubscribeCancelError(nil),
			ShouldBeWrapped: false,
		},
		{
			GivenError:        &app.SubscribeCancelError{},
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: nil,
		},
		{
			GivenError:        app.WrapWithSubscribeCancelError(errors.New("qoo")),
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: errors.New("qoo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *app.SubscribeCancelError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("unwrap", func(t *testing.T) {
					if c.ExpectedUnwrapped != nil {
						require.EqualError(t, target.Unwrap(), c.ExpectedUnwrapped.Error())

						return
					}

					require.NoError(t, target.Unwrap())
				})
			})
		})
	}
}

func TestFormatSubscribeCancelError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &app.SubscribeCancelError{},
			ExpectedErrorString: "",
		},
		{
			GivenError:          app.WrapWithSubscribeCancelError(errors.New("wrong")),
			ExpectedErrorString: "subscribe cancel: wrong",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
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

func requireMessagesMatch(t *testing.T, expected []app.Message, actual <-chan app.Message) {
	t.Helper()

	require.ElementsMatch(
		t,
		mapMessagesSliceToStrings(expected),
		mapMessagesChanToStrings(actual, len(expected)),
	)
}

func requireMessagesContain(t *testing.T, messages <-chan app.Message, contain ...app.Message) {
	t.Helper()

	require.Subset(
		t,
		mapMessagesChanToStrings(messages, len(contain)),
		mapMessagesSliceToStrings(contain),
	)
}

func mapMessagesSliceToStrings(messages []app.Message) []string {
	result := make([]string, 0, len(messages))
	for _, msg := range messages {
		result = append(result, msg.String())
	}

	return result
}

func mapMessagesChanToStrings(messages <-chan app.Message, capacity int) []string {
	result := make([]string, 0, capacity)
	for msg := range messages {
		result = append(result, msg.String())
	}

	return result
}
