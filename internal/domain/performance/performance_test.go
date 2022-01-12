package performance_test

import (
	"testing"

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
			Name: "another_error",
			Err:  errors.New("performance cancelled"),
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
