package specification_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuildThesisSlugging(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		GivenSlug   specification.Slug
		WantThisErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "foo.bar.baz",
			GivenSlug:   specification.NewThesisSlug("foo", "bar", "baz"),
			WantThisErr: false,
			IsErr: func(err error) bool {
				return specification.IsEmptySlugError(err) ||
					specification.IsNotThesisSlugError(err)
			},
		},
		{
			Name:        "empty_slug",
			GivenSlug:   specification.Slug{},
			WantThisErr: true,
			IsErr:       specification.IsEmptySlugError,
		},
		{
			Name:        "not_thesis_slug",
			GivenSlug:   specification.NewStorySlug("bao"),
			WantThisErr: true,
			IsErr:       specification.IsNotThesisSlugError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			thesis, err := specification.NewThesisBuilder().Build(c.GivenSlug)

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))

			require.Equal(t, c.GivenSlug, thesis.Slug())
		})
	}
}

func errlessBuildThesis(
	t *testing.T,
	slug specification.Slug,
	prepare func(b *specification.ThesisBuilder),
) specification.Thesis {
	t.Helper()

	builder := specification.NewThesisBuilder()

	prepare(builder)

	return builder.ErrlessBuild(slug)
}

func TestBuildThesisWithDependencies(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare              func(b *specification.ThesisBuilder)
		ExpectedDependencies []specification.Slug
	}{
		{
			Prepare:              func(b *specification.ThesisBuilder) {},
			ExpectedDependencies: nil,
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithDependency("")
			},
			ExpectedDependencies: []specification.Slug{
				specification.NewThesisSlug("foo", "bar", ""),
			},
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithDependency("copy")
				b.WithDependency("copy")
			},
			ExpectedDependencies: []specification.Slug{
				specification.NewThesisSlug("foo", "bar", "copy"),
				specification.NewThesisSlug("foo", "bar", "copy"),
			},
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithDependency("pop")
				b.WithDependency("coo")
			},
			ExpectedDependencies: []specification.Slug{
				specification.NewThesisSlug("foo", "bar", "pop"),
				specification.NewThesisSlug("foo", "bar", "coo"),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var (
				slug = specification.NewThesisSlug("foo", "bar", "qaz")
				deps = errlessBuildThesis(t, slug, c.Prepare).Dependencies()
			)

			require.ElementsMatch(t, c.ExpectedDependencies, deps)
		})
	}
}

func TestBuildThesisWithStatement(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Keyword     string
		Behavior    string
		ShouldBeErr bool
	}{
		{
			Name:        "allowed_given",
			Keyword:     "given",
			Behavior:    "hooves delivered to the warehouse",
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_when",
			Keyword:     "when",
			Behavior:    "selling hooves",
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_then",
			Keyword:     "then",
			Behavior:    "check that hooves are sold",
			ShouldBeErr: false,
		},
		{
			Name:        "not_allowed_zen",
			Keyword:     "zen",
			Behavior:    "zen du dust",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			thesis, err := specification.NewThesisBuilder().
				WithStatement(c.Keyword, c.Behavior).
				Build(specification.NewThesisSlug("foo", "bar", "baz"))

			if c.ShouldBeErr {
				require.True(t, specification.IsNotAllowedStageError(err))

				return
			}

			require.False(t, specification.IsNotAllowedStageError(err))

			t.Run("stage", func(t *testing.T) {
				assert.Equal(t, strings.ToLower(c.Keyword), thesis.Statement().Stage().String())
			})

			t.Run("behavior", func(t *testing.T) {
				assert.Equal(t, c.Behavior, thesis.Statement().Behavior())
			})
		})
	}
}

func TestBuildThesisWithAssertion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare           func(b *specification.ThesisBuilder)
		ExpectedAssertion specification.Assertion
		ShouldBeErr       bool
	}{
		{
			Prepare:     func(b *specification.ThesisBuilder) {},
			ShouldBeErr: true,
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithAssertion(func(b *specification.AssertionBuilder) {})
			},
			ShouldBeErr: true,
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithAssertion(func(b *specification.AssertionBuilder) {
					b.WithMethod("JSONPATH")
				})
			},
			ExpectedAssertion: specification.NewAssertionBuilder().
				WithMethod("JSONPATH").
				ErrlessBuild(),
			ShouldBeErr: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			builder := specification.NewThesisBuilder()

			c.Prepare(builder)

			thesis, err := builder.Build(specification.NewThesisSlug("a", "b", "c"))

			if c.ShouldBeErr {
				require.True(t, specification.IsNoThesisHTTPOrAssertionError(err))

				return
			}

			require.False(t, specification.IsNoThesisHTTPOrAssertionError(err))

			require.Equal(t, c.ExpectedAssertion, thesis.Assertion())
		})
	}
}

func TestBuildThesisWithHTTP(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare      func(b *specification.ThesisBuilder)
		ExpectedHTTP specification.HTTP
		ShouldBeErr  bool
	}{
		{
			Prepare:     func(b *specification.ThesisBuilder) {},
			ShouldBeErr: true,
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithHTTP(func(b *specification.HTTPBuilder) {})
			},
			ShouldBeErr: true,
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithHTTP(func(b *specification.HTTPBuilder) {
					b.WithRequest(func(b *specification.HTTPRequestBuilder) {
						b.WithMethod("GET")
						b.WithURL("https://some-api/v1/endpoint")
					})
					b.WithResponse(func(b *specification.HTTPResponseBuilder) {
						b.WithAllowedCodes([]int{200})
						b.WithAllowedContentType("application/json")
					})
				})
			},
			ExpectedHTTP: specification.NewHTTPBuilder().
				WithRequest(func(b *specification.HTTPRequestBuilder) {
					b.WithMethod("GET")
					b.WithURL("https://some-api/v1/endpoint")
				}).
				WithResponse(func(b *specification.HTTPResponseBuilder) {
					b.WithAllowedCodes([]int{200})
					b.WithAllowedContentType("application/json")
				}).
				ErrlessBuild(),
			ShouldBeErr: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			builder := specification.NewThesisBuilder()

			c.Prepare(builder)

			thesis, err := builder.Build(specification.NewThesisSlug("a", "b", "c"))

			if c.ShouldBeErr {
				require.True(t, specification.IsNoThesisHTTPOrAssertionError(err))

				return
			}

			require.False(t, specification.IsNoThesisHTTPOrAssertionError(err))

			require.Equal(t, c.ExpectedHTTP, thesis.HTTP())
		})
	}
}

func TestThesisErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "not_allowed_stage_error",
			Err:   specification.NewNotAllowedStageError("zen"),
			IsErr: specification.IsNotAllowedStageError,
		},
		{
			Name:     "NON_not_allowed_stage_error",
			Err:      errors.New("zen"),
			IsErr:    specification.IsNotAllowedStageError,
			Reversed: true,
		},
		{
			Name:  "no_thesis_http_or_assertion_error",
			Err:   specification.NewNoThesisHTTPOrAssertionError(),
			IsErr: specification.IsNoThesisHTTPOrAssertionError,
		},
		{
			Name:     "NON_no_thesis_http_or_assertion_error",
			Err:      errors.New("no thesis HTTP or assertion"),
			IsErr:    specification.IsNoThesisHTTPOrAssertionError,
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
