package performance_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/performance/mock"
)

const (
	from = "from"
	to   = "to"
)

func TestNewTransition(t *testing.T) {
	t.Parallel()

	type params struct {
		State        performance.State
		OccurredErrs []string
	}

	testCases := []struct {
		Name     string
		Params   params
		WithErrs bool
		IsErr    func(err error) bool
	}{
		{
			Name: "without_error",
			Params: params{
				State: performance.Performing,
			},
			WithErrs: false,
		},
		{
			Name: "with_failed_error",
			Params: params{
				State: performance.Failed,
				OccurredErrs: []string{
					"performance failed: something wrong",
				},
			},
			WithErrs: true,
			IsErr:    performance.IsFailedError,
		},
		{
			Name: "with_crashed_error",
			Params: params{
				State: performance.Crashed,
				OccurredErrs: []string{
					"performance crashed: something wrong",
				},
			},
			WithErrs: true,
			IsErr:    performance.IsCrashedError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			transition := performance.NewTransition(
				c.Params.State,
				from,
				to,
				c.Params.OccurredErrs...,
			)

			require.ElementsMatch(t, c.Params.OccurredErrs, transition.OccurredErrs())
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
		performance.WithHTTP(mock.NewPassingPerformer()),
		performance.WithAssertion(mock.NewPassingPerformer()),
	)
	require.NoError(t, err)

	steps, err := perf.Start(context.Background())
	require.NoError(t, err)

	fr := performance.FlowFromPerformance("flow", perf)

	for s := range steps {
		requireStepNotError(t, s)
		requireStepNotFailed(t, s)

		fr.WithStep(s)
	}

	flow := fr.Reduce()
	require.Equal(t, performance.Passed, flow.State())
}

func TestFlowReducer_WithStep_flow_common_state(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		StartState          performance.State
		StepState           performance.State
		ExpectedCommonState performance.State
	}{
		{
			Name:                "NotPerformed -> NotPerformed => Crashed",
			StartState:          performance.NotPerformed,
			StepState:           performance.NotPerformed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "NotPerformed -> Performing => Performing",
			StartState:          performance.NotPerformed,
			StepState:           performance.Performing,
			ExpectedCommonState: performance.Performing,
		},
		{
			Name:                "NotPerformed -> Passed => Passed",
			StartState:          performance.NotPerformed,
			StepState:           performance.Passed,
			ExpectedCommonState: performance.Passed,
		},
		{
			Name:                "NotPerformed -> Failed => Failed",
			StartState:          performance.NotPerformed,
			StepState:           performance.Failed,
			ExpectedCommonState: performance.Failed,
		},
		{
			Name:                "NotPerformed -> Crashed => Crashed",
			StartState:          performance.NotPerformed,
			StepState:           performance.Crashed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "NotPerformed -> Canceled => Canceled",
			StartState:          performance.NotPerformed,
			StepState:           performance.Canceled,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Performing -> NotPerformed => Crashed",
			StartState:          performance.Performing,
			StepState:           performance.NotPerformed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Performing -> Performing => Performing",
			StartState:          performance.Performing,
			StepState:           performance.Performing,
			ExpectedCommonState: performance.Performing,
		},
		{
			Name:                "Performing -> Passed => Passed",
			StartState:          performance.Performing,
			StepState:           performance.Passed,
			ExpectedCommonState: performance.Passed,
		},
		{
			Name:                "Performing -> Failed => Failed",
			StartState:          performance.Performing,
			StepState:           performance.Failed,
			ExpectedCommonState: performance.Failed,
		},
		{
			Name:                "Performing -> Crashed => Crashed",
			StartState:          performance.Performing,
			StepState:           performance.Crashed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Performing -> Canceled => Canceled",
			StartState:          performance.Performing,
			StepState:           performance.Canceled,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Failed -> NotPerformed => Crashed",
			StartState:          performance.Failed,
			StepState:           performance.NotPerformed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Failed -> Performing => Failed",
			StartState:          performance.Failed,
			StepState:           performance.Performing,
			ExpectedCommonState: performance.Failed,
		},
		{
			Name:                "Failed -> Passed => Failed",
			StartState:          performance.Failed,
			StepState:           performance.Passed,
			ExpectedCommonState: performance.Failed,
		},
		{
			Name:                "Failed -> Failed => Failed",
			StartState:          performance.Failed,
			StepState:           performance.Failed,
			ExpectedCommonState: performance.Failed,
		},
		{
			Name:                "Failed -> Crashed => Crashed",
			StartState:          performance.Failed,
			StepState:           performance.Crashed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Failed -> Canceled => Canceled",
			StartState:          performance.Failed,
			StepState:           performance.Canceled,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Crashed -> NotPerformed => Crashed",
			StartState:          performance.Crashed,
			StepState:           performance.NotPerformed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Crashed -> Performing => Crashed",
			StartState:          performance.Crashed,
			StepState:           performance.Performing,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Crashed -> Passed => Crashed",
			StartState:          performance.Crashed,
			StepState:           performance.Passed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Crashed -> Failed => Crashed",
			StartState:          performance.Crashed,
			StepState:           performance.Failed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Crashed -> Crashed => Crashed",
			StartState:          performance.Crashed,
			StepState:           performance.Crashed,
			ExpectedCommonState: performance.Crashed,
		},
		{
			Name:                "Crashed -> Canceled => Canceled",
			StartState:          performance.Crashed,
			StepState:           performance.Canceled,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Canceled -> NotPerformed => Canceled",
			StartState:          performance.Canceled,
			StepState:           performance.NotPerformed,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Canceled -> Performing => Canceled",
			StartState:          performance.Canceled,
			StepState:           performance.Performing,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Canceled -> Passed => Canceled",
			StartState:          performance.Canceled,
			StepState:           performance.Passed,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Canceled -> Failed => Canceled",
			StartState:          performance.Canceled,
			StepState:           performance.Failed,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Canceled -> Crashed => Canceled",
			StartState:          performance.Canceled,
			StepState:           performance.Crashed,
			ExpectedCommonState: performance.Canceled,
		},
		{
			Name:                "Canceled -> Canceled => Canceled",
			StartState:          performance.Canceled,
			StepState:           performance.Canceled,
			ExpectedCommonState: performance.Canceled,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			b := performance.TestFlowFromState(
				c.StartState,
				c.StartState,
				from, to,
			)

			b.WithStep(performance.NewTestStep(from, to, c.StepState))

			flow := b.Reduce()

			require.Equal(t, c.ExpectedCommonState, flow.State())
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
			Name:                    "Performing -> NotPerformed => NotPerformed",
			StartTransitionState:    performance.Performing,
			StepState:               performance.NotPerformed,
			ExpectedTransitionState: performance.NotPerformed,
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

			require.Equal(t, c.ExpectedTransitionState, flow.Transitions()[0].State())
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
