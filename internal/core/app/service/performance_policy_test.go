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
		GivenLogger         service.Logger
		ShouldPanic         bool
		PanicMessage        string
	}{
		{
			Name:                "all_dependencies_are_not_nil",
			GivenFlowRepository: mock.NewFlowRepository(),
			GivenLogger:         mock.NewMemoryLogger(),
			ShouldPanic:         false,
		},
		{
			Name:                "flow_repository_is_nil",
			GivenFlowRepository: nil,
			GivenLogger:         mock.NewMemoryLogger(),
			ShouldPanic:         true,
			PanicMessage:        "flow repository is nil",
		},
		{
			Name:                "logger_is_nil",
			GivenFlowRepository: mock.NewFlowRepository(),
			GivenLogger:         nil,
			ShouldPanic:         true,
			PanicMessage:        "logger is nil",
		},
		{
			Name:                "all_dependencies_are_nil",
			GivenFlowRepository: nil,
			GivenLogger:         nil,
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
					c.GivenLogger,
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

	testCases := []struct {
		Name             string
		CancelContext    bool
		InitSaveTimeout  time.Duration
		GivenPerformance *performance.Performance
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
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				flowRepo = mock.NewFlowRepository()
				logger   = mock.NewMemoryLogger()
				policy   = service.NewSavePerStepPolicy(
					flowRepo,
					logger,
					c.InitSaveTimeout,
				)
			)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if c.CancelContext {
				cancel()
			}

			policy.ConsumePerformance(ctx, c.GivenPerformance)
		})
	}
}
