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
		Messages           []app.Message
	}{
		{
			Name: "already_started_performance",
			PerformanceFactory: func(opts ...performance.Option) *performance.Performance {
				return performance.Unmarshal(performance.Params{
					OwnerID:         "f7b42682-cf52-4699-9bba-f8dac902efb0",
					SpecificationID: "73a7c5f6-f239-4abf-8837-cc4763d59d5f",
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
					OwnerID:         "3fe69ab2-ebd0-4890-babf-4970a5d4d1d1",
					SpecificationID: "d62b068f-39cc-4ac1-9245-f74aba314d78",
					Actions: []performance.Action{
						performance.NewActionWithoutThesis("a", "b", performance.HTTPPerformer),
					},
				}, opts...)
			},
			Releaser: appMock.PerformanceReleaser(func(ctx context.Context, perfID string) error {
				return errPerformanceRelease
			}),
			Messages: []app.Message{
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
					OwnerID:         "52f42fbb-572c-47f3-9774-52107dffec03",
					SpecificationID: "6b28440c-fe30-4f6f-8290-bf783c8e36c2",
					Actions: []performance.Action{
						performance.NewActionWithoutThesis("a", "c", performance.HTTPPerformer),
						performance.NewActionWithoutThesis("b", "c", performance.HTTPPerformer),
						performance.NewActionWithoutThesis("c", "d", performance.AssertionPerformer),
					},
				}, opts...)
			},
			Releaser: emptyPerformanceReleaser(t),
			Messages: []app.Message{
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
				performance.WithID("e148737a-5825-4a39-bca2-9c671f2e0386"),
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

			requireMessagesEqual(t, c.Messages, messages)
		})
	}
}

func TestPerformanceMaintainer_MaintainPerformance_timeout_exceeded(t *testing.T) {
	t.Parallel()

	const flowTimeout = 1 * time.Millisecond

	var (
		releaser    = emptyPerformanceReleaser(t)
		stepsPolicy = appMock.NewStepsPolicy()
		maintainer  = app.NewPerformanceMaintainer(releaser, stepsPolicy, flowTimeout)
	)

	finish := make(chan struct{})
	defer close(finish)

	performer := perfMock.Performer(func(_ *performance.Environment, _ specification.Thesis) performance.Result {
		<-finish

		return performance.Pass()
	})

	perf := performanceWithHTTPPerformer(t, "a", "b", performer)

	messages, err := maintainer.MaintainPerformance(context.Background(), perf)
	require.NoError(t, err)

	requireMessagesEqual(
		t,
		[]app.Message{
			app.NewMessageFromStep(performance.NewPerformingStep("a", "b", performance.HTTPPerformer)),
			app.NewMessageFromStep(performance.NewCanceledStep(context.DeadlineExceeded)),
		},
		messages,
	)
}

func TestPerformanceMaintainer_MaintainPerformance_context_canceled(t *testing.T) {
	t.Parallel()

	const flowTimeout = 3 * time.Second

	var (
		releaser    = emptyPerformanceReleaser(t)
		stepsPolicy = appMock.NewStepsPolicy()
		maintainer  = app.NewPerformanceMaintainer(releaser, stepsPolicy, flowTimeout)
	)

	finish := make(chan struct{})
	defer close(finish)

	performer := perfMock.Performer(func(_ *performance.Environment, _ specification.Thesis) performance.Result {
		<-finish

		return performance.Pass()
	})

	perf := performanceWithHTTPPerformer(t, "a", "b", performer)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	messages, err := maintainer.MaintainPerformance(ctx, perf)
	require.NoError(t, err)

	requireMessagesEqual(
		t,
		[]app.Message{app.NewMessageFromStep(performance.NewCanceledStep(context.Canceled))},
		messages,
	)
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
