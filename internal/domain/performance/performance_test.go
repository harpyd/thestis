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

			perf, err := performance.FromSpecification(c.Specification)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)

			actions := perf.Actions()
			require.Len(t, actions, c.ActionsLen)
		})
	}
}

func TestPerformance_Start(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	http := mock.Performer(func(env *performance.Environment, t specification.Thesis) performance.Result {
		env.Store(t.Slug(), "HTTP")

		return performance.Pass()
	})

	assertion := mock.Performer(func(env *performance.Environment, t specification.Thesis) performance.Result {
		env.Store(t.Slug(), "assertion")

		return performance.Pass()
	})

	perf, err := performance.FromSpecification(
		spec,
		performance.WithHTTP(http),
		performance.WithAssertion(assertion),
	)
	require.NoError(t, err)

	steps, err := perf.Start(context.Background())
	require.NoError(t, err)

	for s := range steps {
		if s.State() == performance.Failed {
			require.Fail(t, "Step with unexpected fail", s.Err())
		}

		if s.State() == performance.Crashed {
			require.Fail(t, "Step with unexpected error", s.Err())
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

func TestPerformance_Start_with_cancel_context(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(spec)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	steps, err := perf.Start(ctx)
	require.NoError(t, err)

	cancel()

	requireCanceledStep(t, steps)

	_, err = perf.Start(context.Background())
	require.NoError(t, err)
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

func TestIsCyclicGraphError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "cyclic_performance_error",
			Err:       performance.NewCyclicGraphError("from", "to"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("from to"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, performance.IsCyclicGraphError(c.Err))
		})
	}
}

func TestIsAlreadyStartedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "performance_already_started_error",
			Err:       performance.NewAlreadyStartedError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("performance already started"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, performance.IsAlreadyStartedError(c.Err))
		})
	}
}

func requireCanceledStep(t *testing.T, steps <-chan performance.Step) {
	t.Helper()

	for s := range steps {
		if s.State() == performance.Canceled {
			return
		}
	}

	require.Fail(t, "No canceled event")
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
