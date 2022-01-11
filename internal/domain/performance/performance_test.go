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

func invalidCyclicSpecification(t *testing.T) *specification.Specification {
	t.Helper()

	spec, err := specification.NewBuilder().
		WithStory("story", func(b *specification.StoryBuilder) {
			b.WithScenario("scenario", func(b *specification.ScenarioBuilder) {
				b.
					WithThesis("a", func(b *specification.ThesisBuilder) {
						b.
							WithStatement("given", "a").
							WithDependencies("b").
							WithHTTP(func(b *specification.HTTPBuilder) {
								b.
									WithRequest(func(b *specification.HTTPRequestBuilder) {
										b.
											WithMethod("GET").
											WithURL("https://some-url")
									}).
									WithResponse(func(b *specification.HTTPResponseBuilder) {
										b.WithAllowedCodes([]int{200})
									})
							})
					}).
					WithThesis("b", func(b *specification.ThesisBuilder) {
						b.
							WithAssertion(func(b *specification.AssertionBuilder) {
								b.WithMethod("jsonpath")
								b.WithAssert("some", 1)
							}).
							WithStatement("given", "b").
							WithDependencies("a")
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
				b.WithThesis("thesis", func(b *specification.ThesisBuilder) {
					b.
						WithStatement("then", "check").
						WithAssertion(func(b *specification.AssertionBuilder) {
							b.
								WithMethod("jsonpath").
								WithAssert("some", "some")
						})
				})
			})
		}).
		Build()

	require.NoError(t, err)

	return spec
}
