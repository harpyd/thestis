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

var (
	errPerformanceAcquire = errors.New("performance acquire")
	errPerformanceRelease = errors.New("performance release")
)

func TestPerformanceMaintainer_MaintainPerformance(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name               string
		PerformanceFactory func(opts ...performance.Option) *performance.Performance
		Guard              *appMock.PerformanceGuard
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
			Guard:            errlessPerformanceGuard(t),
			StartPerformance: true,
			ShouldBeErr:      true,
			IsErr:            performance.IsAlreadyStartedError,
		},
		{
			Name: "performance_acquire_error",
			PerformanceFactory: func(opts ...performance.Option) *performance.Performance {
				return performance.Unmarshal(performance.Params{
					Actions: []performance.Action{
						performance.NewActionWithoutThesis("b", "c", performance.AssertionPerformer),
					},
				}, opts...)
			},
			Guard:       appMock.NewPerformanceGuard(errPerformanceAcquire, nil),
			ShouldBeErr: true,
			IsErr:       func(err error) bool { return errors.Is(err, errPerformanceAcquire) },
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
			Guard: appMock.NewPerformanceGuard(nil, errPerformanceRelease),
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
			Guard: errlessPerformanceGuard(t),
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
				cancelPubsub = appMock.NewPerformanceCancelPubsub()
				stepsPolicy  = appMock.NewStepsPolicy()
				maintainer   = app.NewPerformanceMaintainer(c.Guard, cancelPubsub, stepsPolicy, flowTimeout)
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

			require.Equal(t, 1, c.Guard.ReleaseCalls())
		})
	}
}

func TestPerformanceMaintainer_MaintainPerformance_cancelation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name             string
		CancelContext    bool
		FlowTimeout      time.Duration
		PublishCancel    bool
		ExpectedMessages []app.Message
		Contains         bool
	}{
		{
			Name:          "context_cancelation",
			CancelContext: true,
			FlowTimeout:   1 * time.Second,
			PublishCancel: false,
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(performance.NewCanceledStep(context.Canceled)),
			},
			Contains: true,
		},
		{
			Name:          "flow_timeout_exceeded",
			CancelContext: false,
			FlowTimeout:   10 * time.Millisecond,
			PublishCancel: false,
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(performance.NewPerformingStep("a", "b", performance.HTTPPerformer)),
				app.NewMessageFromStep(performance.NewCanceledStep(context.DeadlineExceeded)),
			},
		},
		{
			Name:          "cancel_published",
			CancelContext: false,
			FlowTimeout:   1 * time.Second,
			PublishCancel: true,
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(performance.NewCanceledStep(context.Canceled)),
			},
			Contains: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				guard        = errlessPerformanceGuard(t)
				cancelPubsub = appMock.NewPerformanceCancelPubsub()
				stepsPolicy  = appMock.NewStepsPolicy()
				maintainer   = app.NewPerformanceMaintainer(guard, cancelPubsub, stepsPolicy, c.FlowTimeout)
			)

			finish := make(chan struct{})
			defer close(finish)

			performer := perfMock.Performer(func(_ *performance.Environment, _ specification.Thesis) performance.Result {
				<-finish

				return performance.Pass()
			})

			perf := performance.Unmarshal(performance.Params{
				Actions: []performance.Action{
					performance.NewActionWithoutThesis("a", "b", performance.HTTPPerformer),
				},
			}, performance.WithID("id"), performance.WithHTTP(performer))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			messages, err := maintainer.MaintainPerformance(ctx, perf)
			require.NoError(t, err)

			if c.CancelContext {
				cancel()
			}

			if c.PublishCancel {
				err := cancelPubsub.PublishPerformanceCancel("id")
				require.NoError(t, err)
			}

			if c.Contains {
				requireMessagesContains(t, messages, c.ExpectedMessages...)

				return
			}

			requireMessagesEqual(t, c.ExpectedMessages, messages)

			require.Equal(t, 1, guard.ReleaseCalls())
		})
	}
}

func errlessPerformanceGuard(t *testing.T) *appMock.PerformanceGuard {
	t.Helper()

	return appMock.NewPerformanceGuard(nil, nil)
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

	actualMessages := readMessages(t, actual)

	require.ElementsMatch(t, expectedMessages, actualMessages)
}

func requireMessagesContains(t *testing.T, messages <-chan app.Message, contains ...app.Message) {
	t.Helper()

	readMsgs := readMessages(t, messages)

	for _, msg := range contains {
		require.Contains(t, readMsgs, msg.String())
	}
}

func readMessages(t *testing.T, messages <-chan app.Message) []string {
	t.Helper()

	readMsgs := make([]string, 0, len(messages))
	for msg := range messages {
		readMsgs = append(readMsgs, msg.String())
	}

	return readMsgs
}
