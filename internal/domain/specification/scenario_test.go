package specification_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestScenarioBuilder_Build_slug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Slug        string
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_slug",
			Slug:        "scenario",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_empty_slug",
			Slug:        "",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewScenarioBuilder()

			scenario, err := builder.Build(c.Slug)

			if c.ShouldBeErr {
				require.True(t, specification.IsScenarioEmptySlugError(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, c.Slug, scenario.Slug())
		})
	}
}

func TestScenarioBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	builder := specification.NewScenarioBuilder()
	builder.WithDescription("description")

	scenario, err := builder.Build("someScenario")

	require.NoError(t, err)
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

	scenario, err := builder.Build("someScenario")

	require.NoError(t, err)

	expectedGetBeerThesis, err := specification.NewThesisBuilder().
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
		Build("getBeer")
	require.NoError(t, err)

	actualGetBeerThesis, ok := scenario.Thesis("getBeer")
	require.True(t, ok)
	require.Equal(t, expectedGetBeerThesis, actualGetBeerThesis)

	expectedCheckBeerThesis, err := specification.NewThesisBuilder().
		WithStatement("then", "check beer").
		WithAssertion(func(b *specification.AssertionBuilder) {
			b.WithMethod("jsonpath")
			b.WithAssert("getSomeBody.response.body.product", "beer")
		}).
		Build("checkBeer")
	require.NoError(t, err)

	actualCheckBeerThesis, ok := scenario.Thesis("checkBeer")
	require.True(t, ok)
	require.Equal(t, expectedCheckBeerThesis, actualCheckBeerThesis)
}

func TestScenarioBuilder_WithThesis_when_already_exists(t *testing.T) {
	t.Parallel()

	builder := specification.NewScenarioBuilder()
	builder.WithThesis("thesis", func(b *specification.ThesisBuilder) {})
	builder.WithThesis("thesis", func(b *specification.ThesisBuilder) {})

	_, err := builder.Build("scenario")

	require.True(t, specification.IsThesisSlugAlreadyExistsError(err))
}

func TestIsScenarioSlugAlreadyExistsError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "scenario_slug_already_exists_error_is_scenario_slug_already_exists_error",
			Err:       specification.NewScenarioSlugAlreadyExistsError("scenario"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_scenario_slug_already_exists_error",
			Err:       specification.NewThesisSlugAlreadyExistsError("thesis"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsScenarioSlugAlreadyExistsError(c.Err))
		})
	}
}

func TestIsScenarioEmptySlugError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "scenario_empty_slug_error_is_scenario_empty_slug_error",
			Err:       specification.NewScenarioEmptySlugError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_scenario_empty_slug_error",
			Err:       errors.New("error"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsScenarioEmptySlugError(c.Err))
		})
	}
}

func TestIsBuildScenarioError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_scenario_error_is_build_scenario_error",
			Err:       specification.NewBuildScenarioError(errors.New("wrong"), "scenario"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_build_scenario_error",
			Err:       errors.New("wrong"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildScenarioError(c.Err))
		})
	}
}

func TestIsNoSuchScenarioError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "no_scenario_error_is_no_scenario_error",
			Err:       specification.NewNoSuchScenarioError("someScenario"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_no_scenario_error",
			Err:       specification.NewNoSuchThesisError("someThesis"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNoSuchScenarioError(c.Err))
		})
	}
}