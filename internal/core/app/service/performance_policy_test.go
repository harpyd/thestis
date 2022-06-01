package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/app/service/mock"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

func TestNewSavePerStepPolicyPanics(t *testing.T) {
	t.Parallel()

	const saveTimeout = 1 * time.Second

	testCases := []struct {
		Name                string
		GivenFlowRepository service.FlowRepository
		ShouldPanic         bool
		PanicMessage        string
	}{
		{
			Name:                "all_dependencies_are_not_nil",
			GivenFlowRepository: mock.NewFlowRepository(),
			ShouldPanic:         false,
		},
		{
			Name:                "all_dependencies_are_nil",
			GivenFlowRepository: nil,
			ShouldPanic:         true,
			PanicMessage:        "flow repository is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = service.NewSavePerStepPolicy(
					c.GivenFlowRepository,
					saveTimeout,
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

func TestConsumePerformanceSavePerStep(t *testing.T) {
	t.Parallel()

	type operationType int

	const (
		match operationType = iota + 1
		contain
	)

	testCases := []struct {
		Name                string
		CancelContext       bool
		InitSaveTimeout     time.Duration
		GivenPerformance    *performance.Performance
		ExpectedMessages    []service.Message
		AssertOperationType operationType
	}{
		{
			Name:            "save_timeout_exceeded",
			CancelContext:   false,
			InitSaveTimeout: 0,
			GivenPerformance: performance.Trigger(
				"id",
				(&specification.Builder{}).
					WithStory("a", func(b *specification.StoryBuilder) {
						b.WithScenario("b", func(b *specification.ScenarioBuilder) {
							b.WithThesis("c", func(b *specification.ThesisBuilder) {
								b.WithHTTP(func(b *specification.HTTPBuilder) {
									b.WithRequest(func(b *specification.HTTPRequestBuilder) {
										b.WithURL("https://some-api/v1/resources")
									})
								})
							})
						})
					}).
					ErrlessBuild(),
				performance.WithHTTP(performance.PassingPerformer()),
			),
			ExpectedMessages: []service.Message{
				service.NewMessageFromError(
					service.WrapWithDatabaseError(context.DeadlineExceeded),
				),
			},
			AssertOperationType: contain,
		},
		{
			Name:            "context_canceled",
			CancelContext:   true,
			InitSaveTimeout: 1 * time.Second,
			GivenPerformance: performance.Trigger(
				"id",
				(&specification.Builder{}).
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
							b.WithThesis("baz", func(b *specification.ThesisBuilder) {
								b.WithAssertion(func(b *specification.AssertionBuilder) {
									b.WithMethod(specification.JSONPath)
								})
							})
						})
					}).
					ErrlessBuild(),
				performance.WithAssertion(performance.PassingPerformer()),
			),
			ExpectedMessages: []service.Message{
				service.NewMessageFromStep(
					performance.NewScenarioStepWithErr(
						performance.WrapWithTerminatedError(context.Canceled, performance.FiredCancel),
						specification.NewScenarioSlug("foo", "bar"),
						performance.FiredCancel,
					),
				),
				service.NewMessageFromStep(
					performance.NewThesisStepWithErr(
						performance.WrapWithTerminatedError(context.Canceled, performance.FiredCancel),
						specification.NewThesisSlug("foo", "bar", "baz"),
						performance.AssertionPerformer,
						performance.FiredCancel,
					)),
				service.NewMessageFromError(
					service.WrapWithDatabaseError(context.Canceled),
				),
			},
			AssertOperationType: contain,
		},
		{
			Name:            "performance_consumed_wo_context_cancellation",
			CancelContext:   false,
			InitSaveTimeout: 1 * time.Second,
			GivenPerformance: performance.Trigger(
				"id",
				(&specification.Builder{}).
					WithStory("a", func(b *specification.StoryBuilder) {
						b.WithScenario("b", func(b *specification.ScenarioBuilder) {
							b.WithThesis("c", func(b *specification.ThesisBuilder) {
								b.WithAssertion(func(b *specification.AssertionBuilder) {
									b.WithMethod(specification.JSONPath)
								})
							})
						})
					}).
					ErrlessBuild(),
				performance.WithAssertion(performance.PassingPerformer()),
			),
			ExpectedMessages: []service.Message{
				service.NewMessageFromStep(
					performance.NewScenarioStep(
						specification.NewScenarioSlug("a", "b"),
						performance.FiredPerform,
					),
				),
				service.NewMessageFromStep(
					performance.NewThesisStep(
						specification.NewThesisSlug("a", "b", "c"),
						performance.AssertionPerformer,
						performance.FiredPerform,
					),
				),
				service.NewMessageFromStep(
					performance.NewThesisStep(
						specification.NewThesisSlug("a", "b", "c"),
						performance.AssertionPerformer,
						performance.FiredPass,
					),
				),
				service.NewMessageFromStep(
					performance.NewScenarioStep(
						specification.NewScenarioSlug("a", "b"),
						performance.FiredPass,
					),
				),
			},
			AssertOperationType: match,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				flowRepo = mock.NewFlowRepository()
				policy   = service.NewSavePerStepPolicy(flowRepo, c.InitSaveTimeout)
			)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if c.CancelContext {
				cancel()
			}

			var actualMessages []service.Message

			policy.ConsumePerformance(ctx, c.GivenPerformance, func(msg service.Message) {
				actualMessages = append(actualMessages, msg)
			})

			switch c.AssertOperationType {
			case match:
				requireMessagesMatch(t, c.ExpectedMessages, actualMessages)
			case contain:
				requireMessagesContain(t, actualMessages, c.ExpectedMessages...)
			default:
				require.Fail(t, "unexpected assert operation type value")
			}
		})
	}
}
