package performance_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestFromSpecification(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		Specification *specification.Specification
		ShouldBeErr   bool
		IsErr         func(err error) bool
	}{
		{
			Name:          "cyclic_performance_graph",
			Specification: invalidCyclicSpecification(t),
			ShouldBeErr:   true,
			IsErr:         performance.IsCyclicPerformanceGraphError,
		},
		{
			Name:          "valid_performance",
			Specification: validSpecification(t),
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
			require.NotNil(t, perf)
		})
	}
}

func TestPerformance_Start(t *testing.T) {
	t.Parallel()
}

func TestPerformance_Start_one_at_a_time(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(spec)
	require.NoError(t, err)

	_, err = perf.Start(context.Background())
	require.NoError(t, err)

	_, err = perf.Start(context.Background())
	require.True(t, performance.IsPerformanceAlreadyStartedError(err))
}

func TestPerformance_Start_with_cancel_context(t *testing.T) {
	t.Parallel()

	spec := validSpecification(t)

	perf, err := performance.FromSpecification(spec)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	s, err := perf.Start(ctx)
	require.NoError(t, err)

	cancel()

	requireCancelledEvent(t, s)

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
		// read event stream of first call
	}

	s2, err := perf.Start(context.Background())
	require.NoError(t, err)

	for range s2 {
		// read event stream of second call
	}

	finish <- true
}

func TestIsCyclicPerformanceGraphError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "cyclic_performance_error",
			Err:       performance.NewCyclicPerformanceGraphError("from", "to"),
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

			require.Equal(t, c.IsSameErr, performance.IsCyclicPerformanceGraphError(c.Err))
		})
	}
}

func TestIsPerformanceCancelledError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "performance_cancelled_error",
			Err:       performance.NewPerformanceCancelledError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("performance cancelled"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, performance.IsPerformanceCancelledError(c.Err))
		})
	}
}

func TestIsPerformanceAlreadyStartedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "performance_already_started_error",
			Err:       performance.NewPerformanceAlreadyStartedError(),
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

			require.Equal(t, c.IsSameErr, performance.IsPerformanceAlreadyStartedError(c.Err))
		})
	}
}

func requireCancelledEvent(t *testing.T, stream <-chan performance.Event) {
	t.Helper()

	for e := range stream {
		if performance.IsPerformanceCancelledError(e.Err()) {
			return
		}
	}

	require.Fail(t, "No cancelled event")
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
