package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
	appMock "github.com/harpyd/thestis/internal/app/mock"
	"github.com/harpyd/thestis/internal/domain/performance"
	perfMock "github.com/harpyd/thestis/internal/domain/performance/mock"
	"github.com/harpyd/thestis/internal/domain/specification"
)

var errPerformanceRelease = errors.New("performance release")

func TestPerformanceMaintainer_MaintainPerformance(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name               string
		PerformanceFactory func(opts ...performance.Option) *performance.Performance
		Releaser           app.PerformanceReleaser
		StartPerformance   bool
		ShouldBeErr        bool
		IsErr              func(err error) bool
		ExpectedMessages   []app.Message
	}{
		{
			Name: "already_started_performance",
			PerformanceFactory: func(opts ...performance.Option) *performance.Performance {
				return performance.Unmarshal(performance.Params{
					Actions: []performance.Action{
						performance.NewActionWithoutThesis("from", "to", performance.HTTPPerformer),
					},
				}, opts...)
			},
			Releaser:         emptyPerformanceReleaser(t),
			StartPerformance: true,
			ShouldBeErr:      true,
			IsErr:            performance.IsAlreadyStartedError,
		},
		{
			Name: "performance_release_error",
			PerformanceFactory: func(opts ...performance.Option) *performance.Performance {
				return performance.Unmarshal(performance.Params{
					Actions: []performance.Action{
						performance.NewActionWithoutThesis("a", "b", performance.HTTPPerformer),
					},
				}, opts...)
			},
			Releaser: appMock.PerformanceReleaser(func(ctx context.Context, perfID string) error {
				return errPerformanceRelease
			}),
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(
					performance.NewPerformingStep("a", "b", performance.HTTPPerformer),
				),
				app.NewMessageFromStep(
					performance.NewStepFromResult("a", "b", performance.HTTPPerformer, performance.Pass()),
				),
				app.NewMessageFromError(errPerformanceRelease),
			},
		},
		{
			Name: "successfully_maintain_performance",
			PerformanceFactory: func(opts ...performance.Option) *performance.Performance {
				return performance.Unmarshal(performance.Params{
					Actions: []performance.Action{
						performance.NewActionWithoutThesis("a", "c", performance.HTTPPerformer),
						performance.NewActionWithoutThesis("b", "c", performance.HTTPPerformer),
						performance.NewActionWithoutThesis("c", "d", performance.AssertionPerformer),
					},
				}, opts...)
			},
			Releaser: emptyPerformanceReleaser(t),
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(
					performance.NewPerformingStep("a", "c", performance.HTTPPerformer),
				),
				app.NewMessageFromStep(
					performance.NewStepFromResult("a", "c", performance.HTTPPerformer, performance.Pass()),
				),
				app.NewMessageFromStep(
					performance.NewPerformingStep("b", "c", performance.HTTPPerformer),
				),
				app.NewMessageFromStep(
					performance.NewStepFromResult("b", "c", performance.HTTPPerformer, performance.Pass()),
				),
				app.NewMessageFromStep(
					performance.NewPerformingStep("c", "d", performance.AssertionPerformer),
				),
				app.NewMessageFromStep(
					performance.NewStepFromResult("c", "d", performance.AssertionPerformer, performance.Pass()),
				),
			},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			const flowTimeout = 5 * time.Second

			var (
				stepsPolicy = appMock.NewStepsPolicy()
				maintainer  = app.NewPerformanceMaintainer(c.Releaser, stepsPolicy, flowTimeout)
			)

			perf := c.PerformanceFactory(
				performance.WithHTTP(passedPerformer(t)),
				performance.WithAssertion(passedPerformer(t)),
			)

			if c.StartPerformance {
				_, err := perf.Start(context.Background())

				require.NoError(t, err)
			}

			messages, err := maintainer.MaintainPerformance(context.Background(), perf)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)

			requireMessagesEqual(t, c.ExpectedMessages, messages)
		})
	}
}

func TestPerformanceMaintainer_MaintainPerformance_cancelation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name             string
		ContextCanceled  bool
		FlowTimeout      time.Duration
		ExpectedMessages []app.Message
	}{
		{
			Name:            "context_cancelation",
			ContextCanceled: true,
			FlowTimeout:     1 * time.Second,
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(performance.NewCanceledStep(context.Canceled)),
			},
		},
		{
			Name:            "flow_timeout_exceeded",
			ContextCanceled: false,
			FlowTimeout:     1 * time.Millisecond,
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(performance.NewPerformingStep("a", "b", performance.HTTPPerformer)),
				app.NewMessageFromStep(performance.NewCanceledStep(context.DeadlineExceeded)),
			},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				releaser    = emptyPerformanceReleaser(t)
				stepsPolicy = appMock.NewStepsPolicy()
				maintainer  = app.NewPerformanceMaintainer(releaser, stepsPolicy, c.FlowTimeout)
			)

			finish := make(chan struct{})
			defer close(finish)

			performer := perfMock.Performer(func(_ *performance.Environment, _ specification.Thesis) performance.Result {
				<-finish

				return performance.Pass()
			})

			perf := performanceWithHTTPPerformer(t, "a", "b", performer)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if c.ContextCanceled {
				cancel()
			}

			messages, err := maintainer.MaintainPerformance(ctx, perf)
			require.NoError(t, err)

			requireMessagesEqual(t, c.ExpectedMessages, messages)
		})
	}
}

func emptyPerformanceReleaser(t *testing.T) app.PerformanceReleaser {
	t.Helper()

	return appMock.PerformanceReleaser(func(ctx context.Context, perfID string) error {
		return nil
	})
}

func performanceWithHTTPPerformer(t *testing.T, from, to string, performer performance.Performer) *performance.Performance {
	t.Helper()

	return performance.Unmarshal(performance.Params{
		Actions: []performance.Action{
			performance.NewActionWithoutThesis(from, to, performance.HTTPPerformer),
		},
	}, performance.WithHTTP(performer))
}

func passedPerformer(t *testing.T) performance.Performer {
	t.Helper()

	return perfMock.Performer(func(_ *performance.Environment, _ specification.Thesis) performance.Result {
		return performance.Pass()
	})
}

func requireMessagesEqual(t *testing.T, expected []app.Message, actual <-chan app.Message) {
	t.Helper()

	expectedMessages := make([]string, 0, len(expected))
	for _, msg := range expected {
		expectedMessages = append(expectedMessages, msg.String())
	}

	readMessages := make([]string, 0, len(expected))
	for msg := range actual {
		readMessages = append(readMessages, msg.String())
	}

	require.ElementsMatch(t, expectedMessages, readMessages)
}
