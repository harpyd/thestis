package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/mock"
	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestPanickingNewEveryStepSavingPolicy(t *testing.T) {
	t.Parallel()

	const saveTimeout = 1 * time.Second

	testCases := []struct {
		Name           string
		GivenFlowsRepo app.FlowsRepository
		ShouldPanic    bool
		PanicMessage   string
	}{
		{
			Name:           "all_dependencies_are_not_nil",
			GivenFlowsRepo: mock.NewFlowsRepository(),
			ShouldPanic:    false,
		},
		{
			Name:           "all_dependencies_are_nil",
			GivenFlowsRepo: nil,
			ShouldPanic:    true,
			PanicMessage:   "flows repository is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = app.NewEveryStepSavingPolicy(c.GivenFlowsRepo, saveTimeout)
			}

			if !c.ShouldPanic {
				require.NotPanics(t, init)

				return
			}

			require.PanicsWithValue(t, c.PanicMessage, init)
		})
	}
}

func TestHandleEveryStepSavingPolicy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		CancelContext     bool
		GivenInitStatuses []flow.Status
		GivenSteps        []performance.Step
		ExpectedMessages  []app.Message
		ExpectedStatuses  []flow.Status
	}{
		{
			Name: "successful_handling_not_performed_to_passed",
			GivenInitStatuses: []flow.Status{
				flow.NewStatus(
					specification.NewThesisSlug("foo", "bar", "dar"),
					flow.NotPerformed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.NotPerformed,
				),
			},
			GivenSteps: []performance.Step{
				performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "dar"),
					performance.HTTPPerformer,
					performance.FiredPass,
				),
				performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPass,
				),
			},
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(
					performance.NewThesisStep(
						specification.NewThesisSlug("foo", "bar", "dar"),
						performance.HTTPPerformer,
						performance.FiredPass,
					),
				),
				app.NewMessageFromStep(
					performance.NewScenarioStep(
						specification.NewScenarioSlug("foo", "bar"),
						performance.FiredPass,
					),
				),
			},
			ExpectedStatuses: []flow.Status{
				flow.NewStatus(
					specification.NewThesisSlug("foo", "bar", "dar"),
					flow.Passed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.Passed,
				),
			},
		},
		{
			Name: "context_canceled",
			GivenInitStatuses: []flow.Status{
				flow.NewStatus(
					specification.NewThesisSlug("a", "b", "c"),
					flow.NotPerformed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("a", "b"),
					flow.NotPerformed,
				),
			},
			CancelContext: true,
			GivenSteps: []performance.Step{
				performance.NewScenarioStep(
					specification.AnyScenarioSlug(),
					performance.FiredCancel,
				),
			},
			ExpectedMessages: []app.Message{
				app.NewMessageFromStep(
					performance.NewScenarioStep(
						specification.AnyScenarioSlug(),
						performance.FiredCancel,
					),
				),
				app.NewMessageFromError(
					app.WrapWithDatabaseError(context.Canceled),
				),
			},
			ExpectedStatuses: []flow.Status{
				flow.NewStatus(
					specification.NewThesisSlug("a", "b", "c"),
					flow.NotPerformed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("a", "b"),
					flow.NotPerformed,
				),
			},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			const saveTimeout = 3 * time.Second

			var (
				flowsRepo = mock.NewFlowsRepository()
				policy    = app.NewEveryStepSavingPolicy(flowsRepo, saveTimeout)
			)

			var (
				fr       = flow.FromStatuses("flow-id", "perf-id", c.GivenInitStatuses...)
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

			f := fr.Reduce()

			require.ElementsMatch(t, c.ExpectedStatuses, f.Statuses())
		})
	}
}
