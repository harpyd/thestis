package performance_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
)

func TestNextState(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		GivenState    performance.State
		GivenEvent    performance.Event
		ExpectedState performance.State
	}{
		{
			Name:          "not_performed-(perform)->performing",
			GivenState:    performance.NotPerformed,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: performance.Performing,
		},
		{
			Name:          "not_performed-(pass)->passed",
			GivenState:    performance.NotPerformed,
			GivenEvent:    performance.FiredPass,
			ExpectedState: performance.Passed,
		},
		{
			Name:          "not_performed-(fail)->failed",
			GivenState:    performance.NotPerformed,
			GivenEvent:    performance.FiredFail,
			ExpectedState: performance.Failed,
		},
		{
			Name:          "not_performed-(crash)->crashed",
			GivenState:    performance.NotPerformed,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "not_performed-(cancel)->canceled",
			GivenState:    performance.NotPerformed,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: performance.Canceled,
		},
		{
			Name:          "performing-(perform)->performing",
			GivenState:    performance.Performing,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: performance.Performing,
		},
		{
			Name:          "performing-(pass)->passed",
			GivenState:    performance.Performing,
			GivenEvent:    performance.FiredPass,
			ExpectedState: performance.Passed,
		},
		{
			Name:          "performing-(fail)->failed",
			GivenState:    performance.Performing,
			GivenEvent:    performance.FiredFail,
			ExpectedState: performance.Failed,
		},
		{
			Name:          "performing-(crash)->crashed",
			GivenState:    performance.Performing,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "performing-(cancel)->canceled",
			GivenState:    performance.Performing,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: performance.Canceled,
		},
		{
			Name:          "passed-(perform)->passed",
			GivenState:    performance.Passed,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: performance.Passed,
		},
		{
			Name:          "passed-(pass)->passed",
			GivenState:    performance.Passed,
			GivenEvent:    performance.FiredPass,
			ExpectedState: performance.Passed,
		},
		{
			Name:          "passed-(fail)->failed",
			GivenState:    performance.Passed,
			GivenEvent:    performance.FiredFail,
			ExpectedState: performance.Failed,
		},
		{
			Name:          "passed-(crash)->crashed",
			GivenState:    performance.Passed,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "passed-(cancel)->canceled",
			GivenState:    performance.Passed,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: performance.Passed,
		},
		{
			Name:          "failed-(perform)->failed",
			GivenState:    performance.Failed,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: performance.Failed,
		},
		{
			Name:          "failed-(pass)->failed",
			GivenState:    performance.Failed,
			GivenEvent:    performance.FiredPass,
			ExpectedState: performance.Failed,
		},
		{
			Name:          "failed-(fail)->failed",
			GivenState:    performance.Failed,
			GivenEvent:    performance.FiredFail,
			ExpectedState: performance.Failed,
		},
		{
			Name:          "failed-(crash)->crashed",
			GivenState:    performance.Failed,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "failed-(cancel)->canceled",
			GivenState:    performance.Failed,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: performance.Failed,
		},
		{
			Name:          "crashed-(perform)->crashed",
			GivenState:    performance.Crashed,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "crashed-(pass)->crashed",
			GivenState:    performance.Crashed,
			GivenEvent:    performance.FiredPass,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "crashed-(fail)->crashed",
			GivenState:    performance.Crashed,
			GivenEvent:    performance.FiredFail,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "crashed-(crash)->crashed",
			GivenState:    performance.Crashed,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "crashed-(cancel)->crashed",
			GivenState:    performance.Crashed,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: performance.Crashed,
		},
		{
			Name:          "canceled-(perform)->canceled",
			GivenState:    performance.Canceled,
			GivenEvent:    performance.FiredPerform,
			ExpectedState: performance.Canceled,
		},
		{
			Name:          "canceled-(pass)->canceled",
			GivenState:    performance.Canceled,
			GivenEvent:    performance.FiredPass,
			ExpectedState: performance.Canceled,
		},
		{
			Name:          "canceled-(fail)->canceled",
			GivenState:    performance.Canceled,
			GivenEvent:    performance.FiredFail,
			ExpectedState: performance.Canceled,
		},
		{
			Name:          "canceled-(crash)->canceled",
			GivenState:    performance.Canceled,
			GivenEvent:    performance.FiredCrash,
			ExpectedState: performance.Canceled,
		},
		{
			Name:          "canceled-(cancel)->canceled",
			GivenState:    performance.Canceled,
			GivenEvent:    performance.FiredCancel,
			ExpectedState: performance.Canceled,
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
