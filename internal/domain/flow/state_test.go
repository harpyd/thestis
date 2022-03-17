package flow_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/performance"
)

func TestNextState(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		GivenState    flow.State
		GivenEvent    performance.Event
		ExpectedState flow.State
	}{
		{
			Name:          "no_state-(perform)->no_state",
			GivenState:    flow.NoState,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: flow.NoState,
		},
		{
			Name:          "performing-(no_event)->no_state",
			GivenState:    flow.Performing,
			GivenEvent:    performance.NoEvent,
			ExpectedState: flow.NoState,
		},
		{
			Name:          "not_performed-(perform)->performing",
			GivenState:    flow.NotPerformed,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: flow.Performing,
		},
		{
			Name:          "not_performed-(pass)->passed",
			GivenState:    flow.NotPerformed,
			GivenEvent:    performance.FiredPass,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "not_performed-(fail)->failed",
			GivenState:    flow.NotPerformed,
			GivenEvent:    performance.FiredFail,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "not_performed-(crash)->crashed",
			GivenState:    flow.NotPerformed,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "not_performed-(cancel)->canceled",
			GivenState:    flow.NotPerformed,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "performing-(perform)->performing",
			GivenState:    flow.Performing,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: flow.Performing,
		},
		{
			Name:          "performing-(pass)->passed",
			GivenState:    flow.Performing,
			GivenEvent:    performance.FiredPass,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "performing-(fail)->failed",
			GivenState:    flow.Performing,
			GivenEvent:    performance.FiredFail,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "performing-(crash)->crashed",
			GivenState:    flow.Performing,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "performing-(cancel)->canceled",
			GivenState:    flow.Performing,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "passed-(perform)->passed",
			GivenState:    flow.Passed,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "passed-(pass)->passed",
			GivenState:    flow.Passed,
			GivenEvent:    performance.FiredPass,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "passed-(fail)->failed",
			GivenState:    flow.Passed,
			GivenEvent:    performance.FiredFail,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "passed-(crash)->crashed",
			GivenState:    flow.Passed,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "passed-(cancel)->canceled",
			GivenState:    flow.Passed,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: flow.Passed,
		},
		{
			Name:          "failed-(perform)->failed",
			GivenState:    flow.Failed,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "failed-(pass)->failed",
			GivenState:    flow.Failed,
			GivenEvent:    performance.FiredPass,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "failed-(fail)->failed",
			GivenState:    flow.Failed,
			GivenEvent:    performance.FiredFail,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "failed-(crash)->crashed",
			GivenState:    flow.Failed,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "failed-(cancel)->canceled",
			GivenState:    flow.Failed,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: flow.Failed,
		},
		{
			Name:          "crashed-(perform)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "crashed-(pass)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    performance.FiredPass,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "crashed-(fail)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    performance.FiredFail,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "crashed-(crash)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "crashed-(cancel)->crashed",
			GivenState:    flow.Crashed,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: flow.Crashed,
		},
		{
			Name:          "canceled-(perform)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "canceled-(pass)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    performance.FiredPass,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "canceled-(fail)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    performance.FiredFail,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "canceled-(crash)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: flow.Canceled,
		},
		{
			Name:          "canceled-(cancel)->canceled",
			GivenState:    flow.Canceled,
			GivenEvent:    performance.FiredCancel,
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
