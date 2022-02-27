package specification_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestThesisBuilder_Build_no_http_or_assertion(t *testing.T) {
	t.Parallel()

	builder := specification.NewThesisBuilder()

	_, err := builder.Build(specification.NewThesisSlug("someStory", "foo", "bar"))

	require.True(t, specification.IsNoThesisHTTPOrAssertionError(err))
}

func TestThesisBuilder_WithDependencies(t *testing.T) {
	t.Parallel()

	builder := specification.NewThesisBuilder()
	builder.WithStatement("when", "something")
	builder.WithDependencies("anotherOneThesis")
	builder.WithDependencies("anotherTwoThesis")

	expectedDependencies := []specification.Slug{
		specification.NewThesisSlug("story", "scenario", "anotherOneThesis"),
		specification.NewThesisSlug("story", "scenario", "anotherTwoThesis"),
	}

	thesis := builder.ErrlessBuild(specification.NewThesisSlug("story", "scenario", "thesis"))

	require.ElementsMatch(t, expectedDependencies, thesis.Dependencies())
}

func TestThesisBuilder_Build_slug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Slug        specification.Slug
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_slug",
			Slug:        specification.NewThesisSlug("story", "scenario", "thesis"),
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_empty_slug",
			Slug:        specification.Slug{},
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewThesisBuilder()
			builder.WithStatement("when", "do something")

			if c.ShouldBeErr {
				_, err := builder.Build(c.Slug)
				require.True(t, specification.IsThesisEmptySlugError(err))

				return
			}

			thesis := builder.ErrlessBuild(c.Slug)
			require.Equal(t, c.Slug, thesis.Slug())
		})
	}
}

func TestThesisBuilder_WithStatement(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Keyword     string
		Behavior    string
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_allowed_given_stage",
			Keyword:     "given",
			Behavior:    "hooves delivered to the warehouse",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_when_stage",
			Keyword:     "when",
			Behavior:    "selling hooves",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_then_stage",
			Keyword:     "then",
			Behavior:    "check that hooves are sold",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_stage",
			Keyword:     "zen",
			Behavior:    "zen du dust",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewThesisBuilder()
			builder.WithStatement(c.Keyword, c.Behavior)

			if c.ShouldBeErr {
				_, err := builder.Build(specification.NewThesisSlug("story", "scenario", "sellHooves"))
				require.True(t, specification.IsNotAllowedStageError(err))

				return
			}

			thesis := builder.ErrlessBuild(specification.NewThesisSlug("story", "scenario", "sellHooves"))
			require.Equal(t, strings.ToLower(c.Keyword), thesis.Statement().Stage().String())
			require.Equal(t, c.Behavior, thesis.Statement().Behavior())
		})
	}
}

func TestThesisBuilder_WithAssertion(t *testing.T) {
	t.Parallel()

	builder := specification.NewThesisBuilder()
	builder.WithStatement("when", "something wrong")
	builder.WithAssertion(func(b *specification.AssertionBuilder) {
		b.WithMethod("JSONPATH")
		b.WithAssert("getSomeBody.response.body.type", "product")
	})

	thesis, err := builder.Build(specification.NewThesisSlug("story", "scenario", "someThesis"))

	require.NoError(t, err)
	expectedAssertion, err := specification.NewAssertionBuilder().
		WithMethod("JSONPATH").
		WithAssert("getSomeBody.response.body.type", "product").
		Build()
	require.NoError(t, err)
	require.Equal(t, expectedAssertion, thesis.Assertion())
}

func TestThesisBuilder_WithHTTP(t *testing.T) {
	t.Parallel()

	builder := specification.NewThesisBuilder()
	builder.WithStatement("given", "some state")
	builder.WithHTTP(func(b *specification.HTTPBuilder) {
		b.WithRequest(func(b *specification.HTTPRequestBuilder) {
			b.WithMethod("GET")
			b.WithURL("https://some-api/v1/endpoint")
		})
		b.WithResponse(func(b *specification.HTTPResponseBuilder) {
			b.WithAllowedCodes([]int{200})
			b.WithAllowedContentType("application/json")
		})
	})

	thesis, err := builder.Build(specification.NewThesisSlug("story", "scenario", "thesis"))

	require.NoError(t, err)
	expectedHTTP, err := specification.NewHTTPBuilder().
		WithRequest(func(b *specification.HTTPRequestBuilder) {
			b.WithMethod("GET")
			b.WithURL("https://some-api/v1/endpoint")
		}).
		WithResponse(func(b *specification.HTTPResponseBuilder) {
			b.WithAllowedCodes([]int{200})
			b.WithAllowedContentType("application/json")
		}).
		Build()
	require.NoError(t, err)
	require.Equal(t, expectedHTTP, thesis.HTTP())
}

func TestIsThesisSlugAlreadyExistsError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name: "thesis_slug_already_exists_error",
			Err: specification.NewSlugAlreadyExistsError(
				specification.NewThesisSlug("story", "scenario", "thesis"),
			),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("thesis"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsThesisSlugAlreadyExistsError(c.Err))
		})
	}
}

func TestIsThesisEmptySlugError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "thesis_empty_slug_error",
			Err:       specification.NewThesisEmptySlugError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("wrong wrong"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsThesisEmptySlugError(c.Err))
		})
	}
}

func TestIsBuildThesisError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name: "build_thesis_error",
			Err: specification.NewBuildThesisError(
				errors.New("pew"),
				specification.NewThesisSlug("foo", "bar", "baz"),
			),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("pew"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildThesisError(c.Err))
		})
	}
}

func TestIsNoSuchThesisError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "no_thesis_error",
			Err:       specification.NewNoSuchThesisError("someThesis"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       specification.NewNoSuchStoryError("someStory"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNoSuchThesisError(c.Err))
		})
	}
}

func TestIsNotAllowedStageError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_stage_error",
			Err:       specification.NewNotAllowedStageError("zen"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       specification.NewNotAllowedHTTPMethodError("zen"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedStageError(c.Err))
		})
	}
}

func TestIsNoThesisHTTPOrAssertionError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "no_thesis_http_or_assertion_error",
			Err:       specification.NewNoThesisHTTPOrAssertionError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.Wrap(specification.NewNoStoryScenariosError(), "wrap"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNoThesisHTTPOrAssertionError(c.Err))
		})
	}
}
