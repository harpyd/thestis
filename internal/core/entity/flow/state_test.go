package flow_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

func TestNextState(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		GivenState    flow.State
		GivenEvent    pipeline.Event
		ExpectedState flow.State
	}{
		{
			Name:          "no_state-(execute)->no_state",
			GivenState:    flow.NoState,
			GivenEvent:    pipeline.FiredExecute,
			ExpectedState: flow.NoState,
		},
		{
			Name:          "executing-(no_event)->executing",
			GivenState:    flow.Executing,
			GivenEvent:    pipeline.NoEvent,
			ExpectedState: flow.Executing,
		},
		{
			Name:          "not_executed-(execute)->executing",
			GivenState:    flow.NotExecuted,
			GivenEvent:    pipeline.FiredExecute,
			ExpectedState: flow.Executing,
		},
		{
			Name:          "not_executed-(pass)->passed",
			GivenState:    flow.NotExecuted,
			GivenEvent:    pipeline.FiredPass,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "not_executed-(fail)->failed",
			GivenState:    flow.NotExecuted,
			GivenEvent:    pipeline.FiredFail,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "not_executed-(crash)->crashed",
			GivenState:    flow.NotExecuted,
			GivenEvent:    pipeline.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "not_executed-(cancel)->canceled",
			GivenState:    flow.NotExecuted,
			GivenEvent:    pipeline.FiredCancel,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "executing-(execute)->executing",
			GivenState:    flow.Executing,
			GivenEvent:    pipeline.FiredExecute,
			ExpectedState: flow.Executing,
		},
		{
			Name:          "executing-(pass)->passed",
			GivenState:    flow.Executing,
			GivenEvent:    pipeline.FiredPass,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "executing-(fail)->failed",
			GivenState:    flow.Executing,
			GivenEvent:    pipeline.FiredFail,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "executing-(crash)->crashed",
			GivenState:    flow.Executing,
			GivenEvent:    pipeline.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "executing-(cancel)->canceled",
			GivenState:    flow.Executing,
			GivenEvent:    pipeline.FiredCancel,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "passed-(execute)->passed",
			GivenState:    flow.Passed,
			GivenEvent:    pipeline.FiredExecute,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "passed-(pass)->passed",
			GivenState:    flow.Passed,
			GivenEvent:    pipeline.FiredPass,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "passed-(fail)->failed",
			GivenState:    flow.Passed,
			GivenEvent:    pipeline.FiredFail,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "passed-(crash)->crashed",
			GivenState:    flow.Passed,
			GivenEvent:    pipeline.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "passed-(cancel)->canceled",
			GivenState:    flow.Passed,
			GivenEvent:    pipeline.FiredCancel,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "failed-(execute)->failed",
			GivenState:    flow.Failed,
			GivenEvent:    pipeline.FiredExecute,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "failed-(pass)->failed",
			GivenState:    flow.Failed,
			GivenEvent:    pipeline.FiredPass,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "failed-(fail)->failed",
			GivenState:    flow.Failed,
			GivenEvent:    pipeline.FiredFail,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "failed-(crash)->crashed",
			GivenState:    flow.Failed,
			GivenEvent:    pipeline.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "failed-(cancel)->canceled",
			GivenState:    flow.Failed,
			GivenEvent:    pipeline.FiredCancel,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "crashed-(execute)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    pipeline.FiredExecute,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "crashed-(pass)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    pipeline.FiredPass,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "crashed-(fail)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    pipeline.FiredFail,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "crashed-(crash)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    pipeline.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "crashed-(cancel)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    pipeline.FiredCancel,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "canceled-(execute)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    pipeline.FiredExecute,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "canceled-(pass)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    pipeline.FiredPass,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "canceled-(fail)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    pipeline.FiredFail,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "canceled-(crash)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    pipeline.FiredCrash,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "canceled-(cancel)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    pipeline.FiredCancel,
			ExpectedState: flow.Canceled,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			actual := c.GivenState.Next(c.GivenEvent)
			require.Equal(t, c.ExpectedState, actual)
		})
	}
}

type states []flow.State

func (s states) Len() int {
	return len(s)
}

func (s states) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s states) Less(i, j int) bool {
	return s[i].Precedence() < s[j].Precedence()
}

func TestStatePrecedenceOrder(t *testing.T) {
	t.Parallel()

	sortedStates := states{
		flow.NoState,
		flow.Executing,
		flow.Passed,
		flow.Failed,
		flow.Crashed,
		flow.Canceled,
		flow.NotExecuted,
	}

	sort.Sort(sortedStates)

	expectedStates := states{
		flow.NoState,
		flow.Passed,
		flow.NotExecuted,
		flow.Canceled,
		flow.Failed,
		flow.Crashed,
		flow.Executing,
	}

	require.Equal(t, expectedStates, sortedStates)
}
