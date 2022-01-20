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

func TestFlowBuilder_WithStep(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		StartCommonState     performance.State
		StartTransitionState performance.State
		StepState            performance.State
		CommonState          performance.State
		TransitionState      performance.State
		FinalCommonState     performance.State
	}{
		{
			Name:                 "(NotPerformed, NotPerformed) -> NotPerformed",
			StartCommonState:     performance.NotPerformed,
			StartTransitionState: performance.NotPerformed,
			StepState:            performance.NotPerformed,
			CommonState:          performance.NotPerformed,
			TransitionState:      performance.NotPerformed,
			FinalCommonState:     performance.NotPerformed,
		},
		{
			Name:                 "(NotPerformed, NotPerformed) -> Performing",
			StartCommonState:     performance.NotPerformed,
			StartTransitionState: performance.NotPerformed,
			StepState:            performance.Performing,
			CommonState:          performance.Performing,
			TransitionState:      performance.Performing,
			FinalCommonState:     performance.Performing,
		},
		{
			Name:                 "(NotPerformed, NotPerformed) -> Passed",
			StartCommonState:     performance.NotPerformed,
			StartTransitionState: performance.NotPerformed,
			StepState:            performance.Passed,
			CommonState:          performance.Performing,
			TransitionState:      performance.Passed,
			FinalCommonState:     performance.Passed,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			b := performance.NewTunedFlowBuilder(
				c.StartCommonState,
				c.StartTransitionState,
				"from",
				"to",
			)

			b.WithStep(mock.NewStep(c.StepState, "from", "to"))

			flow := b.Build()

			require.Equal(t, c.CommonState, flow.State())
			require.Equal(t, c.TransitionState, flow.Transitions()[0].State())

			finalFlow := b.FinallyBuild()

			require.Equal(t, c.FinalCommonState, finalFlow.State())
			require.Equal(t, c.TransitionState, finalFlow.Transitions()[0].State())
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
