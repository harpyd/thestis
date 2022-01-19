package performance_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/performance/mock"
	"github.com/harpyd/thestis/internal/domain/specification"
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

	performer := mock.Performer(func(_ *performance.Environment, _ specification.Thesis) performance.Result {
		return performance.Pass()
	})

	perf, err := performance.FromSpecification(
		spec,
		performance.WithHTTP(performer),
		performance.WithAssertion(performer),
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

func requireStepNotError(t *testing.T, step performance.Step) {
	t.Helper()

	require.NotEqual(t, performance.Crashed, step.State())
	require.NoError(t, step.CrashErr())
}

func requireStepNotFailed(t *testing.T, step performance.Step) {
	t.Helper()

	require.NotEqual(t, performance.Failed, step.State())
	require.NoError(t, step.CrashErr())
}
