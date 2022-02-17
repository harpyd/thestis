package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/mock"
	"github.com/harpyd/thestis/internal/domain/performance"
)

type states struct {
	CommonState     performance.State
	TransitionState performance.State
}

func TestEveryStepSavingPolicy_HandleSteps(t *testing.T) {
	t.Parallel()

	const (
		from = "a"
		to   = "b"
	)

	testCases := []struct {
		Name             string
		CancelContext    bool
		GivenSteps       []performance.Step
		GivenStates      states
		ExpectedMessages []app.Message
		ExpectedStates   states
	}{
		{
			Name: "successful_handling_not_performed_to_passed",
			GivenSteps: []performance.Step{
				performance.NewTestStep(from, to, performance.Performing),
				performance.NewTestStep(from, to, performance.Passed),
			},
			GivenStates: states{
				CommonState:     performance.NotPerformed,
				TransitionState: performance.NotPerformed,
			},
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(
					performance.NewTestStep(from, to, performance.Performing),
				),
				app.NewMessageFromStep(
					performance.NewTestStep(from, to, performance.Passed),
				),
			},
			ExpectedStates: states{
				CommonState:     performance.Passed,
				TransitionState: performance.Passed,
			},
		},
		{
			Name:          "context_canceled",
			CancelContext: true,
			GivenSteps: []performance.Step{
				performance.NewTestStep(from, to, performance.Canceled),
			},
			GivenStates: states{
				CommonState:     performance.Performing,
				TransitionState: performance.Performing,
			},
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(
					performance.NewTestStep(from, to, performance.Canceled),
				),
				app.NewMessageFromError(
					app.NewDatabaseError(context.Canceled),
				),
			},
			ExpectedStates: states{
				CommonState:     performance.Canceled,
				TransitionState: performance.Performing,
			},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			const saveTImeout = 3 * time.Second

			var (
				flowsRepo = mock.NewFlowsRepository()
				policy    = app.NewEveryStepSavingPolicy(flowsRepo, saveTImeout)
			)

			var (
				fr = performance.TestFlowFromState(
					c.GivenStates.CommonState,
					c.GivenStates.TransitionState,
					from,
					to,
				)
				steps    = make(chan performance.Step)
				messages = make(chan app.Message)
			)

			go func() {
				defer close(steps)

				for _, s := range c.GivenSteps {
					steps <- s
				}
			}()

			go func() {
				defer close(messages)

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				if c.CancelContext {
					cancel()
				}

				policy.HandleSteps(ctx, fr, steps, messages)
			}()

			requireMessagesEqual(t, c.ExpectedMessages, messages)

			flow := fr.Reduce()

			require.Equal(t, c.ExpectedStates.CommonState, flow.State())
			require.Equal(t, c.ExpectedStates.TransitionState, flow.Transitions()[0].State())
		})
	}
}
