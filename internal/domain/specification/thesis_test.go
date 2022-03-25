package specification_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuildThesisSlugging(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		GivenSlug       specification.Slug
		ShouldPanic     bool
		WithExpectedErr error
	}{
		{
			Name:        "foo.bar.baz",
			GivenSlug:   specification.NewThesisSlug("foo", "bar", "baz"),
			ShouldPanic: false,
		},
		{
			Name:            "zero_slug",
			GivenSlug:       specification.Slug{},
			ShouldPanic:     true,
			WithExpectedErr: specification.ErrZeroSlug,
		},
		{
			Name:            "not_thesis_slug",
			GivenSlug:       specification.NewStorySlug("bao"),
			ShouldPanic:     true,
			WithExpectedErr: specification.ErrNotThesisSlug,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewThesisBuilder()

			var thesis specification.Thesis

			buildFn := func() {
				thesis = builder.ErrlessBuild(c.GivenSlug)
			}

			if c.ShouldPanic {
				require.PanicsWithValue(t, c.WithExpectedErr, buildFn)

				return
			}

			require.NotPanics(t, buildFn)

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
		Keyword     specification.Stage
		Behavior    string
		ShouldBeErr bool
	}{
		{
			Name:        "allowed_given",
			Keyword:     specification.Given,
			Behavior:    "hooves delivered to the warehouse",
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_when",
			Keyword:     specification.When,
			Behavior:    "selling hooves",
			ShouldBeErr: false,
		},
		{
			Name:        "allowed_then",
			Keyword:     specification.Then,
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

			var target *specification.NotAllowedStageError

			if c.ShouldBeErr {
				require.ErrorAs(t, err, &target)

				return
			}

			require.False(t, errors.As(err, &target))

			t.Run("stage", func(t *testing.T) {
				require.Equal(t, c.Keyword, thesis.Statement().Stage())
			})

			t.Run("behavior", func(t *testing.T) {
				require.Equal(t, c.Behavior, thesis.Statement().Behavior())
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
				require.ErrorIs(t, err, specification.ErrUselessThesis)

				return
			}

			require.NotErrorIs(t, err, specification.ErrUselessThesis)

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
				require.ErrorIs(t, err, specification.ErrUselessThesis)

				return
			}

			require.NotErrorIs(t, err, specification.ErrUselessThesis)

			require.Equal(t, c.ExpectedHTTP, thesis.HTTP())
		})
	}
}

func TestStageIsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Stage         specification.Stage
		ShouldBeValid bool
	}{
		{
			Stage:         "",
			ShouldBeValid: false,
		},
		{
			Stage:         specification.UnknownStage,
			ShouldBeValid: false,
		},
		{
			Stage:         specification.Given,
			ShouldBeValid: true,
		},
		{
			Stage:         specification.When,
			ShouldBeValid: true,
		},
		{
			Stage:         specification.Then,
			ShouldBeValid: true,
		},
		{
			Stage:         "deploy",
			ShouldBeValid: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeValid, c.Stage.IsValid())
		})
	}
}

func TestBeforeStage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenStage           specification.Stage
		ExpectedBeforeStages []specification.Stage
	}{
		{
			GivenStage:           specification.UnknownStage,
			ExpectedBeforeStages: nil,
		},
		{
			GivenStage:           "foo",
			ExpectedBeforeStages: nil,
		},
		{
			GivenStage:           specification.Given,
			ExpectedBeforeStages: nil,
		},
		{
			GivenStage: specification.When,
			ExpectedBeforeStages: []specification.Stage{
				specification.Given,
			},
		},
		{
			GivenStage: specification.Then,
			ExpectedBeforeStages: []specification.Stage{
				specification.Given,
				specification.When,
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedBeforeStages, c.GivenStage.Before())
		})
	}
}
