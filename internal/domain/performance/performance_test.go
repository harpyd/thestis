package performance_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/performance/mock"
	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestFromSpecification(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		Specification *specification.Specification
		ActionsLen    int
		ShouldBeErr   bool
		IsErr         func(err error) bool
	}{
		{
			Name:          "cyclic_performance_graph",
			Specification: invalidCyclicSpecification(t),
			ShouldBeErr:   true,
			IsErr:         performance.IsCyclicGraphError,
		},
		{
			Name:          "valid_performance",
			Specification: validSpecification(t),
			ActionsLen:    5,
			ShouldBeErr:   false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			perf, err := performance.FromSpecification(c.Specification, performance.WithID("perf"))

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)

			actions := perf.Actions()
			require.Len(t, actions, c.ActionsLen)
			require.Equal(t, "perf", perf.ID())
			require.Equal(t, c.Specification.ID(), perf.SpecificationID())
			require.Equal(t, c.Specification.OwnerID(), perf.OwnerID())
		})
	}
}

func TestPerformance_Start(t *testing.T) {
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

	for s := range steps {
		if s.State() == performance.Failed {
			require.Fail(t, "Step with unexpected fail", s.Err())
		}

		if s.State() == performance.Crashed {
			require.Fail(t, "Step with unexpected crash", s.Err())
		}
	}
}

func TestPerformance_Start_one_at_a_time(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(spec)
	require.NoError(t, err)

	_, err = perf.Start(context.Background())
	require.NoError(t, err)

	_, err = perf.Start(context.Background())
	require.True(t, performance.IsAlreadyStartedError(err))
}

func TestPerformance_Start_with_cancel_context_before_steps_reading(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(spec)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	steps, err := perf.Start(ctx)
	require.NoError(t, err)

	cancel()

	requireLastCanceledStep(t, steps)

	_, err = perf.Start(context.Background())
	require.NoError(t, err)
}

func TestPerformance_Start_with_cancel_context_while_steps_reading(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(
		spec,
		performance.WithHTTP(mock.NewPassingPerformer()),
		performance.WithAssertion(mock.NewPassingPerformer()),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	steps, err := perf.Start(ctx)
	require.NoError(t, err)

	<-steps

	cancel()

	requireLastCanceledStep(t, steps)

	_, err = perf.Start(context.Background())
	require.NoError(t, err)
}

func TestPerformance_Start_with_failed_performer(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(
		spec,
		performance.WithHTTP(mock.NewPassingPerformer()),
		performance.WithAssertion(mock.NewFailingPerformer()),
	)
	require.NoError(t, err)

	steps, err := perf.Start(context.Background())
	require.NoError(t, err)

	requireStep(t, steps, performance.Failed)
}

func TestPerformance_Start_with_crashed_performer(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(
		spec,
		performance.WithHTTP(mock.NewCrashingPerformer()),
		performance.WithAssertion(mock.NewPassingPerformer()),
	)
	require.NoError(t, err)

	steps, err := perf.Start(context.Background())
	require.NoError(t, err)

	requireStep(t, steps, performance.Crashed)
}

func TestPerformance_Start_sync_calls_in_a_row(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(spec)
	require.NoError(t, err)

	s1, err := perf.Start(context.Background())
	require.NoError(t, err)

	finish := make(chan bool)

	go func() {
		select {
		case <-finish:
		case <-time.After(3 * time.Second):
			require.Fail(t, "Timeout exceeded, test is not finished")
		}
	}()

	for range s1 {
		// read flow steps of first call
	}

	s2, err := perf.Start(context.Background())
	require.NoError(t, err)

	for range s2 {
		// read flow steps of second call
	}

	finish <- true
}

func TestPerformanceErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "cyclic_performance_error",
			Err:   performance.NewCyclicGraphError("from", "to"),
			IsErr: performance.IsCyclicGraphError,
		},
		{
			Name:     "NON_cyclic_performance_error",
			Err:      errors.New("cyclic performance graph"),
			IsErr:    performance.IsCyclicGraphError,
			Reversed: true,
		},
		{
			Name:  "performance_already_started_error",
			Err:   performance.NewAlreadyStartedError(),
			IsErr: performance.IsAlreadyStartedError,
		},
		{
			Name:     "NON_performance_already_started_error",
			Err:      errors.New("performance already started"),
			IsErr:    performance.IsAlreadyStartedError,
			Reversed: true,
		},
		{
			Name:  "performance_not_started_error",
			Err:   performance.NewNotStartedError(),
			IsErr: performance.IsNotStartedError,
		},
		{
			Name:     "NON_performance_not_started_error",
			Err:      errors.New("performance not started"),
			IsErr:    performance.IsNotStartedError,
			Reversed: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			if c.Reversed {
				require.False(t, c.IsErr(c.Err))

				return
			}

			require.True(t, c.IsErr(c.Err))
		})
	}
}

func requireLastCanceledStep(t *testing.T, steps <-chan performance.Step) {
	t.Helper()

	var step performance.Step

	for s := range steps {
		step = s
	}

	require.Equal(t, performance.Canceled, step.State())
}

func requireStep(t *testing.T, steps <-chan performance.Step, state performance.State) {
	t.Helper()

	for s := range steps {
		if s.State() == state {
			return
		}
	}

	require.Failf(t, "No %s step", state.String())
}

func invalidCyclicSpecification(t *testing.T) *specification.Specification {
	t.Helper()

	spec, err := specification.NewBuilder().
		WithStory("story", func(b *specification.StoryBuilder) {
			b.WithScenario("scenario", func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement("given", "a")
					b.WithDependencies("b")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod("GET")
							b.WithURL("https://some-url")
						})
						b.WithResponse(func(b *specification.HTTPResponseBuilder) {
							b.WithAllowedCodes([]int{200})
						})
					})
				})
				b.WithThesis("b", func(b *specification.ThesisBuilder) {
					b.WithAssertion(func(b *specification.AssertionBuilder) {
						b.WithMethod("jsonpath")
						b.WithAssert("some", 1)
					})
					b.WithStatement("given", "b")
					b.WithDependencies("a")
				})
			})
		}).
		Build()

	require.NoError(t, err)

	return spec
}

func validSpecification(t *testing.T) *specification.Specification {
	t.Helper()

	spec, err := specification.NewBuilder().
		WithStory("story", func(b *specification.StoryBuilder) {
			b.WithScenario("scenario", func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement("given", "state")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod("POST")
							b.WithURL("https://some-api/endpoint")
						})
						b.WithResponse(func(b *specification.HTTPResponseBuilder) {
							b.WithAllowedCodes([]int{201})
							b.WithAllowedContentType("application/json")
						})
					})
				})
				b.WithThesis("b", func(b *specification.ThesisBuilder) {
					b.WithStatement("when", "action")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod("GET")
							b.WithURL("https://some-api/endpoint")
						})
						b.WithResponse(func(b *specification.HTTPResponseBuilder) {
							b.WithAllowedCodes([]int{200})
						})
					})
				})
				b.WithThesis("c", func(b *specification.ThesisBuilder) {
					b.WithStatement("then", "check")
					b.WithAssertion(func(b *specification.AssertionBuilder) {
						b.WithMethod("jsonpath")
						b.WithAssert("some", "some")
					})
				})
			})
		}).
		Build()

	require.NoError(t, err)

	return spec
}
