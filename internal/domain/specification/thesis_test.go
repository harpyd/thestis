package specification_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestThesisBuilder_WithAfter(t *testing.T) {
	t.Parallel()

	builder := specification.NewThesisBuilder()
	builder.WithStatement("when", "something")
	builder.WithAfter("anotherOneThesis")
	builder.WithAfter("anotherTwoThesis")

	thesis, err := builder.Build("thesis")

	require.NoError(t, err)
	require.ElementsMatch(t, []string{"anotherOneThesis", "anotherTwoThesis"}, thesis.After())
}

func TestThesisBuilder_Build_slug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Slug        string
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_slug",
			Slug:        "thesis",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_empty_slug",
			Slug:        "",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewThesisBuilder()
			builder.WithStatement("when", "do something")

			thesis, err := builder.Build(c.Slug)

			if c.ShouldBeErr {
				require.True(t, specification.IsThesisEmptySlugError(err))

				return
			}

			require.NoError(t, err)
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
			Name:        "build_with_allowed_given_keyword",
			Keyword:     "given",
			Behavior:    "hooves delivered to the warehouse",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_when_keyword",
			Keyword:     "when",
			Behavior:    "selling hooves",
			ShouldBeErr: false,
		},
		{
			Name:        "build_with_allowed_then_keyword",
			Keyword:     "then",
			Behavior:    "check that hooves are sold",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_not_allowed_keyword",
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

			thesis, err := builder.Build("sellHooves")

			if c.ShouldBeErr {
				require.True(t, specification.IsNotAllowedKeywordError(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, strings.ToLower(c.Keyword), thesis.Statement().Keyword().String())
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

	thesis, err := builder.Build("someThesis")

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

	thesis, err := builder.Build("someThesis")

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
			Name:      "thesis_slug_already_exists_error_is_thesis_slug_already_exists_error",
			Err:       specification.NewThesisSlugAlreadyExistsError("thesis"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_thesis_slug_already_exists_error",
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
			Name:      "thesis_empty_slug_error_is_empty_slug_error",
			Err:       specification.NewThesisEmptySlugError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_thesis_empty_slug_error",
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
			Name:      "build_thesis_error_is_build_thesis_error",
			Err:       specification.NewBuildThesisError(errors.New("pew"), "thesis"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_build_thesis_error",
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
			Name:      "no_thesis_error_is_no_thesis_error",
			Err:       specification.NewNoSuchThesisError("someThesis"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_no_thesis_error",
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

func TestIsNotAllowedKeywordError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_keyword_error_is_not_allowed_keyword_error",
			Err:       specification.NewNotAllowedKeywordError("zen"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_not_allowed_keyword_error",
			Err:       specification.NewNotAllowedHTTPMethodError("zen"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedKeywordError(c.Err))
		})
	}
}
