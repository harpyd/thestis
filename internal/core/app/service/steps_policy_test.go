package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/app/service/mock"
	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

func TestPanickingNewEveryStepSavingPolicy(t *testing.T) {
	t.Parallel()

	const saveTimeout = 1 * time.Second

	testCases := []struct {
		Name          string
		GivenFlowRepo service.FlowRepository
		ShouldPanic   bool
		PanicMessage  string
	}{
		{
			Name:          "all_dependencies_are_not_nil",
			GivenFlowRepo: mock.NewFlowRepository(),
			ShouldPanic:   false,
		},
		{
			Name:          "all_dependencies_are_nil",
			GivenFlowRepo: nil,
			ShouldPanic:   true,
			PanicMessage:  "flow repository is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = service.NewEveryStepSavingPolicy(c.GivenFlowRepo, saveTimeout)
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
		GivenInitStatuses []*flow.Status
		GivenSteps        []performance.Step
		ExpectedMessages  []service.Message
		ExpectedStatuses  []*flow.Status
	}{
		{
			Name: "successful_handling_not_performed_to_passed",
			GivenInitStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.NotPerformed,
					flow.NewThesisStatus("dar", flow.NotPerformed),
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
			ExpectedMessages: []service.Message{
				service.NewMessageFromStep(
					performance.NewThesisStep(
						specification.NewThesisSlug("foo", "bar", "dar"),
						performance.HTTPPerformer,
						performance.FiredPass,
					),
				),
				service.NewMessageFromStep(
					performance.NewScenarioStep(
						specification.NewScenarioSlug("foo", "bar"),
						performance.FiredPass,
					),
				),
			},
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.Passed,
					flow.NewThesisStatus("dar", flow.Passed),
				),
			},
		},
		{
			Name: "context_canceled",
			GivenInitStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("a", "b"),
					flow.NotPerformed,
					flow.NewThesisStatus("c", flow.NotPerformed),
				),
			},
			CancelContext: true,
			GivenSteps: []performance.Step{
				performance.NewScenarioStep(
					specification.NewScenarioSlug("a", "b"),
					performance.FiredCancel,
				),
			},
			ExpectedMessages: []service.Message{
				service.NewMessageFromStep(
					performance.NewScenarioStep(
						specification.NewScenarioSlug("a", "b"),
						performance.FiredCancel,
					),
				),
				service.NewMessageFromError(
					service.WrapWithDatabaseError(context.Canceled),
				),
			},
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("a", "b"),
					flow.Canceled,
					flow.NewThesisStatus("c", flow.NotPerformed),
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
				flowRepo = mock.NewFlowRepository()
				policy   = service.NewEveryStepSavingPolicy(flowRepo, saveTimeout)
			)

			var (
				f        = flow.FromStatuses("flow-id", "perf-id", c.GivenInitStatuses...)
				steps    = make(chan performance.Step)
				messages = make(chan service.Message)
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

				policy.HandleSteps(ctx, f, steps, messages)
			}()

			requireMessagesMatch(t, c.ExpectedMessages, messages)

			require.ElementsMatch(t, c.ExpectedStatuses, f.Statuses())
		})
	}
}
