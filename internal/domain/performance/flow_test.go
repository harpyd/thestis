package performance_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/performance/mock"
)

func TestNewFlowBuilder_build_from_new_builder(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(spec)
	require.NoError(t, err)

	b := performance.NewFlowBuilder(perf)
	flow := b.Build()

	require.Equal(t, performance.NotPerformed, flow.State())
	require.Len(t, flow.Transitions(), len(perf.Actions()))
}

func TestFlowBuilder_WithStep_from_valid_performance_start(t *testing.T) {
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

	b := performance.NewFlowBuilder(perf)

	for s := range steps {
		requireStepNotError(t, s)
		requireStepNotFailed(t, s)

		b.WithStep(s)

		flow := b.Build()
		require.Equal(t, performance.Performing, flow.State())
	}

	flow := b.FinallyBuild()
	require.Equal(t, performance.Passed, flow.State())
}

const (
	from = "from"
	to   = "to"
)

func TestFlowBuilder_WithStep_transition_state(t *testing.T) {
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

			b := performance.NewTunedFlowBuilder(
				performance.NotPerformed,
				c.StartTransitionState,
				from, to,
			)

			b.WithStep(step(t, c.StepState))

			flow := b.Build()
			finalFlow := b.FinallyBuild()

			require.Equal(t, c.ExpectedTransitionState, flow.Transitions()[0].State())
			require.Equal(t, c.ExpectedTransitionState, finalFlow.Transitions()[0].State())
		})
	}
}

func step(t *testing.T, state performance.State) performance.Step {
	t.Helper()

	if state == performance.Canceled {
		return mock.NewCanceledStep()
	}

	return mock.NewTransitionStep(state, from, to)
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
