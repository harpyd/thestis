package specification_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestScenarioBuilder_Build_no_theses(t *testing.T) {
	t.Parallel()

	builder := specification.NewScenarioBuilder()

	_, err := builder.Build(specification.NewScenarioSlug("story", "scenario"))

	require.True(t, specification.IsNoScenarioThesesError(err))
}

func TestScenarioBuilder_Build_slug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Slug        specification.Slug
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_slug",
			Slug:        specification.NewScenarioSlug("story", "scenario"),
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_empty_slug",
			Slug:        specification.Slug{},
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewScenarioBuilder()

			if c.ShouldBeErr {
				_, err := builder.Build(c.Slug)
				require.True(t, specification.IsEmptySlugError(err))

				return
			}

			scenario := builder.ErrlessBuild(c.Slug)

			require.Equal(t, c.Slug, scenario.Slug())
		})
	}
}

func TestScenarioBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	builder := specification.NewScenarioBuilder()
	builder.WithDescription("description")

	scenario := builder.ErrlessBuild(specification.NewScenarioSlug("someStory", "someScenario"))

	require.Equal(t, "description", scenario.Description())
}

func TestScenarioBuilder_WithThesis(t *testing.T) {
	t.Parallel()

	builder := specification.NewScenarioBuilder()
	builder.WithThesis("getBeer", func(b *specification.ThesisBuilder) {
		b.WithStatement("when", "get beer")
		b.WithHTTP(func(b *specification.HTTPBuilder) {
			b.WithRequest(func(b *specification.HTTPRequestBuilder) {
				b.WithMethod("GET")
				b.WithURL("https://api/v1/products")
			})
			b.WithResponse(func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedCodes([]int{200})
				b.WithAllowedContentType("application/json")
			})
		})
	})
	builder.WithThesis("checkBeer", func(b *specification.ThesisBuilder) {
		b.WithStatement("then", "check beer")
		b.WithAssertion(func(b *specification.AssertionBuilder) {
			b.WithMethod("JSONPATH")
			b.WithAssert("getSomeBody.response.body.product", "beer")
		})
	})

	scenario := builder.ErrlessBuild(specification.NewScenarioSlug("story", "someScenario"))

	expectedGetBeerThesis := specification.NewThesisBuilder().
		WithStatement("when", "get beer").
		WithHTTP(func(b *specification.HTTPBuilder) {
			b.WithRequest(func(b *specification.HTTPRequestBuilder) {
				b.WithMethod("GET")
				b.WithURL("https://api/v1/products")
			})
			b.WithResponse(func(b *specification.HTTPResponseBuilder) {
				b.WithAllowedCodes([]int{200})
				b.WithAllowedContentType("application/json")
			})
		}).
		ErrlessBuild(specification.NewThesisSlug("story", "someScenario", "getBeer"))

	actualGetBeerThesis, ok := scenario.Thesis("getBeer")
	require.True(t, ok)
	require.Equal(t, expectedGetBeerThesis, actualGetBeerThesis)

	expectedCheckBeerThesis := specification.NewThesisBuilder().
		WithStatement("then", "check beer").
		WithAssertion(func(b *specification.AssertionBuilder) {
			b.WithMethod("jsonpath")
			b.WithAssert("getSomeBody.response.body.product", "beer")
		}).
		ErrlessBuild(specification.NewThesisSlug("story", "someScenario", "checkBeer"))

	actualCheckBeerThesis, ok := scenario.Thesis("checkBeer")
	require.True(t, ok)
	require.Equal(t, expectedCheckBeerThesis, actualCheckBeerThesis)
}

func TestScenarioBuilder_WithThesis_when_already_exists(t *testing.T) {
	t.Parallel()

	builder := specification.NewScenarioBuilder()
	builder.WithThesis("thesis", func(b *specification.ThesisBuilder) {})
	builder.WithThesis("thesis", func(b *specification.ThesisBuilder) {})

	_, err := builder.Build(specification.NewScenarioSlug("someStory", "scenario"))

	require.True(t, specification.IsThesisSlugAlreadyExistsError(err))
}

func TestScenarioErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "no_scenario_theses_error",
			Err:   specification.NewNoScenarioThesesError(),
			IsErr: specification.IsNoScenarioThesesError,
		},
		{
			Name:     "NON_no_scenario_theses_error",
			Err:      errors.New("no scenario theses"),
			IsErr:    specification.IsNoScenarioThesesError,
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
