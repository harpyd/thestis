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

			var b specification.ThesisBuilder

			var thesis specification.Thesis

			buildFn := func() {
				thesis = b.Build(c.GivenSlug)
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

func buildThesis(
	t *testing.T,
	prepare func(b *specification.ThesisBuilder),
) specification.Thesis {
	t.Helper()

	var b specification.ThesisBuilder

	prepare(&b)

	return b.Build(specification.NewThesisSlug("a", "b", "c"))
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
				specification.NewThesisSlug("a", "b", ""),
			},
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithDependency("copy")
				b.WithDependency("copy")
			},
			ExpectedDependencies: []specification.Slug{
				specification.NewThesisSlug("a", "b", "copy"),
			},
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithDependency("pop")
				b.WithDependency("coo")
			},
			ExpectedDependencies: []specification.Slug{
				specification.NewThesisSlug("a", "b", "pop"),
				specification.NewThesisSlug("a", "b", "coo"),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualDeps := buildThesis(t, c.Prepare).Dependencies()

			require.ElementsMatch(t, c.ExpectedDependencies, actualDeps)
		})
	}
}

func TestBuildThesisWithStatement(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare          func(b *specification.ThesisBuilder)
		ExpectedStage    specification.Stage
		ExpectedBehavior string
	}{
		{
			Prepare:          func(b *specification.ThesisBuilder) {},
			ExpectedStage:    specification.NoStage,
			ExpectedBehavior: "",
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithStatement(specification.NoStage, "")
			},
			ExpectedStage:    specification.NoStage,
			ExpectedBehavior: "",
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithStatement(specification.Given, "")
			},
			ExpectedStage:    specification.Given,
			ExpectedBehavior: "",
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithStatement(specification.When, "foo")
			},
			ExpectedStage:    specification.When,
			ExpectedBehavior: "foo",
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithStatement(specification.Then, "bar")
			},
			ExpectedStage:    specification.Then,
			ExpectedBehavior: "bar",
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithStatement("unknown", "hm")
			},
			ExpectedStage:    "unknown",
			ExpectedBehavior: "hm",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			thesis := buildThesis(t, c.Prepare)

			t.Run("stage", func(t *testing.T) {
				actualStage := thesis.Stage()

				require.Equal(t, c.ExpectedStage, actualStage)
			})

			t.Run("behavior", func(t *testing.T) {
				actualBehavior := thesis.Behavior()

				require.Equal(t, c.ExpectedBehavior, actualBehavior)
			})
		})
	}
}

func TestBuildThesisWithAssertion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare           func(b *specification.ThesisBuilder)
		ExpectedAssertion specification.Assertion
	}{
		{
			Prepare:           func(b *specification.ThesisBuilder) {},
			ExpectedAssertion: (&specification.AssertionBuilder{}).Build(),
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithAssertion(func(b *specification.AssertionBuilder) {})
			},
			ExpectedAssertion: (&specification.AssertionBuilder{}).Build(),
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithAssertion(func(b *specification.AssertionBuilder) {
					b.WithMethod("JSONPATH")
				})
			},
			ExpectedAssertion: (&specification.AssertionBuilder{}).
				WithMethod("JSONPATH").
				Build(),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualAssertion := buildThesis(t, c.Prepare).Assertion()

			require.Equal(t, c.ExpectedAssertion, actualAssertion)
		})
	}
}

func TestBuildThesisWithHTTP(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare      func(b *specification.ThesisBuilder)
		ExpectedHTTP specification.HTTP
	}{
		{
			Prepare:      func(b *specification.ThesisBuilder) {},
			ExpectedHTTP: (&specification.HTTPBuilder{}).Build(),
		},
		{
			Prepare: func(b *specification.ThesisBuilder) {
				b.WithHTTP(func(b *specification.HTTPBuilder) {})
			},
			ExpectedHTTP: (&specification.HTTPBuilder{}).Build(),
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
			ExpectedHTTP: (&specification.HTTPBuilder{}).
				WithRequest(func(b *specification.HTTPRequestBuilder) {
					b.WithMethod("GET")
					b.WithURL("https://some-api/v1/endpoint")
				}).
				WithResponse(func(b *specification.HTTPResponseBuilder) {
					b.WithAllowedCodes([]int{200})
					b.WithAllowedContentType("application/json")
				}).
				Build(),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualHTTP := buildThesis(t, c.Prepare).HTTP()

			require.Equal(t, c.ExpectedHTTP, actualHTTP)
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

func TestAsNotAllowedStageError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError      error
		ShouldBeWrapped bool
		ExpectedStage   specification.Stage
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      &specification.NotAllowedStageError{},
			ShouldBeWrapped: true,
			ExpectedStage:   specification.NoStage,
		},
		{
			GivenError:      specification.NewNotAllowedStageError(specification.Given),
			ShouldBeWrapped: true,
			ExpectedStage:   specification.Given,
		},
		{
			GivenError:      specification.NewNotAllowedStageError("else"),
			ShouldBeWrapped: true,
			ExpectedStage:   "else",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *specification.NotAllowedStageError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("stage", func(t *testing.T) {
					require.Equal(t, c.ExpectedStage, target.Stage())
				})
			})
		})
	}
}

func TestFormatNotAllowedStageError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &specification.NotAllowedStageError{},
			ExpectedErrorString: `stage "" not allowed`,
		},
		{
			GivenError:          specification.NewNotAllowedStageError(specification.Given),
			ExpectedErrorString: `stage "given" not allowed`,
		},
		{
			GivenError:          specification.NewNotAllowedStageError("deploy"),
			ExpectedErrorString: `stage "deploy" not allowed`,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
		})
	}
}

func TestFormatUndefinedDependencyError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &specification.UndefinedDependencyError{},
			ExpectedErrorString: `undefined "" dependency`,
		},
		{
			GivenError:          specification.NewUndefinedDependencyError(specification.Slug{}),
			ExpectedErrorString: `undefined "" dependency`,
		},
		{
			GivenError: specification.NewUndefinedDependencyError(
				specification.NewThesisSlug("a", "b", "c")),
			ExpectedErrorString: `undefined "c" dependency`,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
		})
	}
}

func TestAsUndefinedDependencyError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError      error
		ShouldBeWrapped bool
		ExpectedSlug    specification.Slug
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      &specification.UndefinedDependencyError{},
			ShouldBeWrapped: true,
			ExpectedSlug:    specification.Slug{},
		},
		{
			GivenError:      specification.NewUndefinedDependencyError(specification.Slug{}),
			ShouldBeWrapped: true,
			ExpectedSlug:    specification.Slug{},
		},
		{
			GivenError: specification.NewUndefinedDependencyError(
				specification.NewThesisSlug("a", "b", "c"),
			),
			ShouldBeWrapped: true,
			ExpectedSlug:    specification.NewThesisSlug("a", "b", "c"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *specification.UndefinedDependencyError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("slug", func(t *testing.T) {
					require.Equal(t, c.ExpectedSlug, target.Slug())
				})
			})
		})
	}
}
