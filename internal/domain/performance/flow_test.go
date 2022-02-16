package performance_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
)

const (
	from = "from"
	to   = "to"
)

func TestNewTransition(t *testing.T) {
	t.Parallel()

	type params struct {
		State performance.State
		Error string
	}

	testCases := []struct {
		Name    string
		Params  params
		WithErr bool
		IsErr   func(err error) bool
	}{
		{
			Name: "without_error",
			Params: params{
				State: performance.Performing,
			},
			WithErr: false,
		},
		{
			Name: "with_failed_error",
			Params: params{
				State: performance.Failed,
				Error: "performance failed: something wrong",
			},
			WithErr: true,
			IsErr:   performance.IsFailedError,
		},
		{
			Name: "with_crashed_error",
			Params: params{
				State: performance.Crashed,
				Error: "performance crashed: something wrong",
			},
			WithErr: true,
			IsErr:   performance.IsCrashedError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			transition := performance.NewTransition(c.Params.State, from, to, c.Params.Error)

			if c.WithErr {
				require.True(t, c.IsErr(transition.Err()))
				require.EqualError(t, transition.Err(), c.Params.Error)
			}

			require.Equal(t, c.Params.State, transition.State())
			require.Equal(t, from, transition.From())
			require.Equal(t, to, transition.To())
		})
	}
}

func TestFlowFromPerformance(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(spec)
	require.NoError(t, err)

	b := performance.FlowFromPerformance("flow", perf)
	flow := b.Reduce()

	require.Equal(t, performance.NotPerformed, flow.State())
	require.Len(t, flow.Transitions(), len(perf.Actions()))
}

func TestFlowReducer_WithStep_from_performance_start(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(
		spec,
		performance.WithHTTP(passingPerformer(t)),
		performance.WithAssertion(passingPerformer(t)),
	)
	require.NoError(t, err)

	steps, err := perf.Start(context.Background())
	require.NoError(t, err)

	b := performance.FlowFromPerformance("flow", perf)

	for s := range steps {
		requireStepNotError(t, s)
		requireStepNotFailed(t, s)

		b.WithStep(s)

		flow := b.Reduce()
		require.Equal(t, performance.Performing, flow.State())
	}

	flow := b.FinallyReduce()
	require.Equal(t, performance.Passed, flow.State())
}

func TestFlowReducer_WithStep_flow_common_state(t *testing.T) {
	t.Parallel()

	type expectedStates struct {
		CommonState      performance.State
		FinalCommonState performance.State
	}

	testCases := []struct {
		Name             string
		StartCommonState performance.State
		StepState        performance.State
		Expected         expectedStates
	}{
		{
			Name:             "NotPerformed -> NotPerformed => (NotPerformed, NotPerformed)",
			StartCommonState: performance.NotPerformed,
			StepState:        performance.NotPerformed,
			Expected: expectedStates{
				CommonState:      performance.NotPerformed,
				FinalCommonState: performance.NotPerformed,
			},
		},
		{
			Name:             "NotPerformed -> Performing => (Performing, Performing)",
			StartCommonState: performance.NotPerformed,
			StepState:        performance.Performing,
			Expected: expectedStates{
				CommonState:      performance.Performing,
				FinalCommonState: performance.Performing,
			},
		},
		{
			Name:             "NotPerformed -> Passed => (Performing, Passed)",
			StartCommonState: performance.NotPerformed,
			StepState:        performance.Passed,
			Expected: expectedStates{
				CommonState:      performance.Performing,
				FinalCommonState: performance.Passed,
			},
		},
		{
			Name:             "NotPerformed -> Failed => (Failed, Failed)",
			StartCommonState: performance.NotPerformed,
			StepState:        performance.Failed,
			Expected: expectedStates{
				CommonState:      performance.Failed,
				FinalCommonState: performance.Failed,
			},
		},
		{
			Name:             "NotPerformed -> Crashed => (Crashed, Crashed)",
			StartCommonState: performance.NotPerformed,
			StepState:        performance.Crashed,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "NotPerformed -> Canceled => (Canceled, Canceled)",
			StartCommonState: performance.NotPerformed,
			StepState:        performance.Canceled,
			Expected: expectedStates{
				CommonState:      performance.Canceled,
				FinalCommonState: performance.Canceled,
			},
		},
		{
			Name:             "Performing -> NotPerformed => (Performing, Performing)",
			StartCommonState: performance.Performing,
			StepState:        performance.NotPerformed,
			Expected: expectedStates{
				CommonState:      performance.Performing,
				FinalCommonState: performance.Performing,
			},
		},
		{
			Name:             "Performing -> Performing => (Performing, Performing)",
			StartCommonState: performance.Performing,
			StepState:        performance.Performing,
			Expected: expectedStates{
				CommonState:      performance.Performing,
				FinalCommonState: performance.Performing,
			},
		},
		{
			Name:             "Performing -> Passed => (Performing, Passed)",
			StartCommonState: performance.Performing,
			StepState:        performance.Passed,
			Expected: expectedStates{
				CommonState:      performance.Performing,
				FinalCommonState: performance.Passed,
			},
		},
		{
			Name:             "Performing -> Failed => (Failed, Failed)",
			StartCommonState: performance.Performing,
			StepState:        performance.Failed,
			Expected: expectedStates{
				CommonState:      performance.Failed,
				FinalCommonState: performance.Failed,
			},
		},
		{
			Name:             "Performing -> Crashed => (Crashed, Crashed)",
			StartCommonState: performance.Performing,
			StepState:        performance.Crashed,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Performing -> Canceled => (Canceled, Canceled)",
			StartCommonState: performance.Performing,
			StepState:        performance.Canceled,
			Expected: expectedStates{
				CommonState:      performance.Canceled,
				FinalCommonState: performance.Canceled,
			},
		},
		{
			Name:             "Failed -> NotPerformed => (Failed, Failed)",
			StartCommonState: performance.Failed,
			StepState:        performance.NotPerformed,
			Expected: expectedStates{
				CommonState:      performance.Failed,
				FinalCommonState: performance.Failed,
			},
		},
		{
			Name:             "Failed -> Performing => (Failed, Failed)",
			StartCommonState: performance.Failed,
			StepState:        performance.Performing,
			Expected: expectedStates{
				CommonState:      performance.Failed,
				FinalCommonState: performance.Failed,
			},
		},
		{
			Name:             "Failed -> Passed => (Failed, Failed)",
			StartCommonState: performance.Failed,
			StepState:        performance.Passed,
			Expected: expectedStates{
				CommonState:      performance.Failed,
				FinalCommonState: performance.Failed,
			},
		},
		{
			Name:             "Failed -> Failed => (Failed, Failed)",
			StartCommonState: performance.Failed,
			StepState:        performance.Failed,
			Expected: expectedStates{
				CommonState:      performance.Failed,
				FinalCommonState: performance.Failed,
			},
		},
		{
			Name:             "Failed -> Crashed => (Crashed, Crashed)",
			StartCommonState: performance.Failed,
			StepState:        performance.Crashed,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Failed -> Canceled => (Failed, Failed)",
			StartCommonState: performance.Failed,
			StepState:        performance.Canceled,
			Expected: expectedStates{
				CommonState:      performance.Failed,
				FinalCommonState: performance.Failed,
			},
		},
		{
			Name:             "Crashed -> NotPerformed => (Crashed, Crashed)",
			StartCommonState: performance.Crashed,
			StepState:        performance.NotPerformed,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Crashed -> Performing => (Crashed, Crashed)",
			StartCommonState: performance.Crashed,
			StepState:        performance.Performing,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Crashed -> Passed => (Crashed, Crashed)",
			StartCommonState: performance.Crashed,
			StepState:        performance.Passed,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Crashed -> Failed => (Crashed, Crashed)",
			StartCommonState: performance.Crashed,
			StepState:        performance.Failed,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Crashed -> Crashed => (Crashed, Crashed)",
			StartCommonState: performance.Crashed,
			StepState:        performance.Crashed,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Crashed -> Canceled => (Crashed, Crashed)",
			StartCommonState: performance.Crashed,
			StepState:        performance.Canceled,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Canceled -> NotPerformed => (Canceled, Canceled)",
			StartCommonState: performance.Canceled,
			StepState:        performance.NotPerformed,
			Expected: expectedStates{
				CommonState:      performance.Canceled,
				FinalCommonState: performance.Canceled,
			},
		},
		{
			Name:             "Canceled -> Performing => (Canceled, Canceled)",
			StartCommonState: performance.Canceled,
			StepState:        performance.Performing,
			Expected: expectedStates{
				CommonState:      performance.Canceled,
				FinalCommonState: performance.Canceled,
			},
		},
		{
			Name:             "Canceled -> Passed => (Canceled, Canceled)",
			StartCommonState: performance.Canceled,
			StepState:        performance.Passed,
			Expected: expectedStates{
				CommonState:      performance.Canceled,
				FinalCommonState: performance.Canceled,
			},
		},
		{
			Name:             "Canceled -> Failed => (Failed, Failed)",
			StartCommonState: performance.Canceled,
			StepState:        performance.Failed,
			Expected: expectedStates{
				CommonState:      performance.Failed,
				FinalCommonState: performance.Failed,
			},
		},
		{
			Name:             "Canceled -> Crashed => (Crashed, Crashed)",
			StartCommonState: performance.Canceled,
			StepState:        performance.Crashed,
			Expected: expectedStates{
				CommonState:      performance.Crashed,
				FinalCommonState: performance.Crashed,
			},
		},
		{
			Name:             "Canceled -> Canceled => (Canceled, Canceled)",
			StartCommonState: performance.Canceled,
			StepState:        performance.Canceled,
			Expected: expectedStates{
				CommonState:      performance.Canceled,
				FinalCommonState: performance.Canceled,
			},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			b := performance.TestFlowFromState(
				c.StartCommonState,
				performance.NotPerformed,
				from, to,
			)

			b.WithStep(performance.NewTestStep(from, to, c.StepState))

			flow := b.Reduce()
			finalFlow := b.FinallyReduce()

			require.Equal(t, c.Expected.CommonState, flow.State())
			require.Equal(t, c.Expected.FinalCommonState, finalFlow.State())
		})
	}
}

func TestFlowReducer_WithStep_transition_state(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                    string
		StartTransitionState    performance.State
		StepState               performance.State
		ExpectedTransitionState performance.State
	}{
		{
			Name:                    "NotPerformed -> NotPerformed => NotPerformed",
			StartTransitionState:    performance.NotPerformed,
			StepState:               performance.NotPerformed,
			ExpectedTransitionState: performance.NotPerformed,
		},
		{
			Name:                    "NotPerformed -> Performing => Performing",
			StartTransitionState:    performance.NotPerformed,
			StepState:               performance.Performing,
			ExpectedTransitionState: performance.Performing,
		},
		{
			Name:                    "NotPerformed -> Passed => Passed",
			StartTransitionState:    performance.NotPerformed,
			StepState:               performance.Passed,
			ExpectedTransitionState: performance.Passed,
		},
		{
			Name:                    "NotPerformance -> Failed => Failed",
			StartTransitionState:    performance.NotPerformed,
			StepState:               performance.Failed,
			ExpectedTransitionState: performance.Failed,
		},
		{
			Name:                    "NotPerformed -> Crashed => Crashed",
			StartTransitionState:    performance.NotPerformed,
			StepState:               performance.Crashed,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "NotPerformed -> Canceled => NotPerformed",
			StartTransitionState:    performance.NotPerformed,
			StepState:               performance.Canceled,
			ExpectedTransitionState: performance.NotPerformed,
		},
		{
			Name:                    "Performing -> NotPerformed => Performing",
			StartTransitionState:    performance.Performing,
			StepState:               performance.NotPerformed,
			ExpectedTransitionState: performance.Performing,
		},
		{
			Name:                    "Performing -> Performing => Performing",
			StartTransitionState:    performance.Performing,
			StepState:               performance.Performing,
			ExpectedTransitionState: performance.Performing,
		},
		{
			Name:                    "Performing -> Passed => Passed",
			StartTransitionState:    performance.Performing,
			StepState:               performance.Passed,
			ExpectedTransitionState: performance.Passed,
		},
		{
			Name:                    "Performing -> Failed => Failed",
			StartTransitionState:    performance.Performing,
			StepState:               performance.Failed,
			ExpectedTransitionState: performance.Failed,
		},
		{
			Name:                    "Performing -> Crashed => Failed",
			StartTransitionState:    performance.Performing,
			StepState:               performance.Crashed,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "Performing -> Canceled => Performing",
			StartTransitionState:    performance.Performing,
			StepState:               performance.Canceled,
			ExpectedTransitionState: performance.Performing,
		},
		{
			Name:                    "Passed -> NotPerformed => Passed",
			StartTransitionState:    performance.Passed,
			StepState:               performance.NotPerformed,
			ExpectedTransitionState: performance.Passed,
		},
		{
			Name:                    "Passed -> Performing => Passed",
			StartTransitionState:    performance.Passed,
			StepState:               performance.Performing,
			ExpectedTransitionState: performance.Passed,
		},
		{
			Name:                    "Passed -> Passed => Passed",
			StartTransitionState:    performance.Passed,
			StepState:               performance.Passed,
			ExpectedTransitionState: performance.Passed,
		},
		{
			Name:                    "Passed -> Failed => Failed",
			StartTransitionState:    performance.Passed,
			StepState:               performance.Failed,
			ExpectedTransitionState: performance.Failed,
		},
		{
			Name:                    "Passed -> Crashed => Crashed",
			StartTransitionState:    performance.Passed,
			StepState:               performance.Crashed,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "Passed -> Canceled => Passed",
			StartTransitionState:    performance.Passed,
			StepState:               performance.Canceled,
			ExpectedTransitionState: performance.Passed,
		},
		{
			Name:                    "Failed -> NotPerformed => Failed",
			StartTransitionState:    performance.Failed,
			StepState:               performance.NotPerformed,
			ExpectedTransitionState: performance.Failed,
		},
		{
			Name:                    "Failed -> Performing => Failed",
			StartTransitionState:    performance.Failed,
			StepState:               performance.Performing,
			ExpectedTransitionState: performance.Failed,
		},
		{
			Name:                    "Failed -> Passed => Failed",
			StartTransitionState:    performance.Failed,
			StepState:               performance.Passed,
			ExpectedTransitionState: performance.Failed,
		},
		{
			Name:                    "Failed -> Failed => Failed",
			StartTransitionState:    performance.Failed,
			StepState:               performance.Failed,
			ExpectedTransitionState: performance.Failed,
		},
		{
			Name:                    "Failed -> Crashed => Crashed",
			StartTransitionState:    performance.Failed,
			StepState:               performance.Crashed,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "Failed -> Canceled => Failed",
			StartTransitionState:    performance.Failed,
			StepState:               performance.Canceled,
			ExpectedTransitionState: performance.Failed,
		},
		{
			Name:                    "Crashed -> NotPerformed => Crashed",
			StartTransitionState:    performance.Crashed,
			StepState:               performance.NotPerformed,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "Crashed -> Performing => Crashed",
			StartTransitionState:    performance.Crashed,
			StepState:               performance.Performing,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "Crashed -> Passed => Crashed",
			StartTransitionState:    performance.Crashed,
			StepState:               performance.Passed,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "Crashed -> Failed => Crashed",
			StartTransitionState:    performance.Crashed,
			StepState:               performance.Failed,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "Crashed -> Crashed => Crashed",
			StartTransitionState:    performance.Crashed,
			StepState:               performance.Crashed,
			ExpectedTransitionState: performance.Crashed,
		},
		{
			Name:                    "Crashed -> Canceled => Crashed",
			StartTransitionState:    performance.Crashed,
			StepState:               performance.Canceled,
			ExpectedTransitionState: performance.Crashed,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			b := performance.TestFlowFromState(
				performance.NotPerformed,
				c.StartTransitionState,
				from, to,
			)

			b.WithStep(performance.NewTestStep(from, to, c.StepState))

			flow := b.Reduce()
			finalFlow := b.FinallyReduce()

			require.Equal(t, c.ExpectedTransitionState, flow.Transitions()[0].State())
			require.Equal(t, c.ExpectedTransitionState, finalFlow.Transitions()[0].State())
		})
	}
}

func requireStepNotError(t *testing.T, step performance.Step) {
	t.Helper()

	require.NotEqual(t, performance.Crashed, step.State())
	require.NoError(t, step.Err())
}

func requireStepNotFailed(t *testing.T, step performance.Step) {
	t.Helper()

	require.NotEqual(t, performance.Failed, step.State())
	require.NoError(t, step.Err())
}
