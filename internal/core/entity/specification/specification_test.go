package specification_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/entity/specification"
)

func errlessBuildSpec(
	t *testing.T,
	prepare func(b *specification.Builder),
) *specification.Specification {
	t.Helper()

	var b specification.Builder

	prepare(&b)

	return b.ErrlessBuild()
}

func TestValidateSpecificationBuilding(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare     func(b *specification.Builder)
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Prepare:     func(b *specification.Builder) {},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoSpecificationStories)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoStoryScenarios)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {
					b.WithScenario("bar", func(b *specification.ScenarioBuilder) {})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoScenarioTheses)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("aaa", func(b *specification.StoryBuilder) {
					b.WithScenario("bb", func(b *specification.ScenarioBuilder) {
						b.WithThesis("c", func(b *specification.ThesisBuilder) {
							b.WithStatement(specification.When, "useless")
						})
					})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrUselessThesis)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("a", func(b *specification.StoryBuilder) {
					b.WithScenario("b", func(b *specification.ScenarioBuilder) {
						b.WithThesis("c", func(b *specification.ThesisBuilder) {
							b.WithStatement(specification.Given, "oops")
							b.WithDependency("undefined")
						})
					})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *specification.UndefinedDependencyError

				return errors.As(err, &target)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {
					b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
						b.WithThesis("baz", func(b *specification.ThesisBuilder) {
							b.WithStatement(specification.Given, "none")
							b.WithAssertion(func(b *specification.AssertionBuilder) {
								b.WithMethod("unknown")
							})
						})
					})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *specification.NotAllowedAssertionMethodError

				return errors.As(err, &target)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("a", func(b *specification.StoryBuilder) {
					b.WithScenario("b", func(b *specification.ScenarioBuilder) {
						b.WithThesis("c", func(b *specification.ThesisBuilder) {
							b.WithStatement(specification.Given, "some behavior")
							b.WithHTTP(func(b *specification.HTTPBuilder) {
								b.WithResponse(func(b *specification.HTTPResponseBuilder) {
									b.WithAllowedContentType(specification.ApplicationJSON)
								})
							})
						})
					})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoHTTPRequest)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("a", func(b *specification.StoryBuilder) {
					b.WithScenario("b", func(b *specification.ScenarioBuilder) {
						b.WithThesis("c", func(b *specification.ThesisBuilder) {
							b.WithStatement(specification.Then, "then")
							b.WithHTTP(func(b *specification.HTTPBuilder) {
								b.WithRequest(func(b *specification.HTTPRequestBuilder) {
									b.WithMethod("unknown")
								})
							})
						})
					})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *specification.NotAllowedHTTPMethodError

				return errors.As(err, &target)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("a", func(b *specification.StoryBuilder) {
					b.WithScenario("b", func(b *specification.ScenarioBuilder) {
						b.WithThesis("c", func(b *specification.ThesisBuilder) {
							b.WithStatement(specification.When, "when")
							b.WithHTTP(func(b *specification.HTTPBuilder) {
								b.WithRequest(func(b *specification.HTTPRequestBuilder) {
									b.WithContentType("unknown/content")
								})
							})
						})
					})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *specification.NotAllowedContentTypeError

				return errors.As(err, &target)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {
					b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
						b.WithThesis("bad", func(b *specification.ThesisBuilder) {
							b.WithStatement(specification.Given, "given")
							b.WithHTTP(func(b *specification.HTTPBuilder) {
								b.WithResponse(func(b *specification.HTTPResponseBuilder) {
									b.WithAllowedContentType("bad/content")
								})
							})
						})
					})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *specification.NotAllowedContentTypeError

				return errors.As(err, &target)
			},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("story", func(b *specification.StoryBuilder) {
					b.WithScenario("scenario", func(b *specification.ScenarioBuilder) {
						b.WithThesis("thesis", func(b *specification.ThesisBuilder) {})
					})
				})
			},
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *specification.NotAllowedStageError

				return errors.As(err, &target)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var b specification.Builder

			c.Prepare(&b)

			_, err := b.Build()

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestBuildSpecificationWithID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare    func(b *specification.Builder)
		ExpectedID string
	}{
		{
			Prepare:    func(b *specification.Builder) {},
			ExpectedID: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithID("")
			},
			ExpectedID: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithID("some-id")
			},
			ExpectedID: "some-id",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualID := errlessBuildSpec(t, c.Prepare).ID()

			require.Equal(t, c.ExpectedID, actualID)
		})
	}
}

func TestBuildSpecificationWithOwnerID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare         func(b *specification.Builder)
		ExpectedOwnerID string
	}{
		{
			Prepare:         func(b *specification.Builder) {},
			ExpectedOwnerID: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithOwnerID("")
			},
			ExpectedOwnerID: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithOwnerID("owner-id")
			},
			ExpectedOwnerID: "owner-id",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualOwnerID := errlessBuildSpec(t, c.Prepare).OwnerID()

			require.Equal(t, c.ExpectedOwnerID, actualOwnerID)
		})
	}
}

func TestBuildSpecificationWithTestCampaignID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare                func(b *specification.Builder)
		ExpectedTestCampaignID string
	}{
		{
			Prepare:                func(b *specification.Builder) {},
			ExpectedTestCampaignID: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithTestCampaignID("")
			},
			ExpectedTestCampaignID: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithTestCampaignID("test-campaign-id")
			},
			ExpectedTestCampaignID: "test-campaign-id",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualTestCampaignID := errlessBuildSpec(t, c.Prepare).TestCampaignID()

			require.Equal(t, c.ExpectedTestCampaignID, actualTestCampaignID)
		})
	}
}

func TestBuildSpecificationWithLoadedAt(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	testCases := []struct {
		Prepare          func(b *specification.Builder)
		ExpectedLoadedAt time.Time
	}{
		{
			Prepare:          func(b *specification.Builder) {},
			ExpectedLoadedAt: time.Time{},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithLoadedAt(time.Time{})
			},
			ExpectedLoadedAt: time.Time{},
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithLoadedAt(now)
			},
			ExpectedLoadedAt: now,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualLoadedAt := errlessBuildSpec(t, c.Prepare).LoadedAt()

			require.Equal(t, c.ExpectedLoadedAt, actualLoadedAt)
		})
	}
}

func TestBuildSpecificationWithAuthor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare        func(b *specification.Builder)
		ExpectedAuthor string
	}{
		{
			Prepare:        func(b *specification.Builder) {},
			ExpectedAuthor: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithAuthor("")
			},
			ExpectedAuthor: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithAuthor("djerys")
			},
			ExpectedAuthor: "djerys",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualAuthor := errlessBuildSpec(t, c.Prepare).Author()

			require.Equal(t, c.ExpectedAuthor, actualAuthor)
		})
	}
}

func TestBuildSpecificationWithTitle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare       func(b *specification.Builder)
		ExpectedTitle string
	}{
		{
			Prepare:       func(b *specification.Builder) {},
			ExpectedTitle: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithTitle("")
			},
			ExpectedTitle: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithTitle("foo")
			},
			ExpectedTitle: "foo",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualTitle := errlessBuildSpec(t, c.Prepare).Title()

			require.Equal(t, c.ExpectedTitle, actualTitle)
		})
	}
}

func TestBuildSpecificationWithDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare             func(b *specification.Builder)
		ExpectedDescription string
	}{
		{
			Prepare:             func(b *specification.Builder) {},
			ExpectedDescription: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithDescription("")
			},
			ExpectedDescription: "",
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithDescription("desc")
			},
			ExpectedDescription: "desc",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			actualDescription := errlessBuildSpec(t, c.Prepare).Description()

			require.Equal(t, c.ExpectedDescription, actualDescription)
		})
	}
}

func TestBuildSpecificationWithStories(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare           func(b *specification.Builder)
		ExpectedStories   []specification.Story
		ExpectedScenarios []specification.Scenario
		ExpectedTheses    []specification.Thesis
	}{
		{
			Prepare:           func(b *specification.Builder) {},
			ExpectedStories:   nil,
			ExpectedScenarios: nil,
			ExpectedTheses:    nil,
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {})
				b.WithStory("bar", func(b *specification.StoryBuilder) {})
				b.WithStory("baz", func(b *specification.StoryBuilder) {})
			},
			ExpectedStories: []specification.Story{
				(&specification.StoryBuilder{}).
					Build(specification.NewStorySlug("foo")),
				(&specification.StoryBuilder{}).
					Build(specification.NewStorySlug("bar")),
				(&specification.StoryBuilder{}).
					Build(specification.NewStorySlug("baz")),
			},
			ExpectedScenarios: nil,
			ExpectedTheses:    nil,
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("oops", func(b *specification.StoryBuilder) {})
				b.WithStory("oops", func(b *specification.StoryBuilder) {})
			},
			ExpectedStories: []specification.Story{
				(&specification.StoryBuilder{}).
					Build(specification.NewStorySlug("oops")),
			},
			ExpectedScenarios: nil,
			ExpectedTheses:    nil,
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {
					b.WithScenario("bar", func(b *specification.ScenarioBuilder) {})
				})
				b.WithStory("baz", func(b *specification.StoryBuilder) {
					b.WithScenario("bad", func(b *specification.ScenarioBuilder) {})
				})
			},
			ExpectedStories: []specification.Story{
				(&specification.StoryBuilder{}).
					WithScenario("bar", func(b *specification.ScenarioBuilder) {}).
					Build(specification.NewStorySlug("foo")),
				(&specification.StoryBuilder{}).
					WithScenario("bad", func(b *specification.ScenarioBuilder) {}).
					Build(specification.NewStorySlug("baz")),
			},
			ExpectedScenarios: []specification.Scenario{
				(&specification.ScenarioBuilder{}).
					Build(specification.NewScenarioSlug("foo", "bar")),
				(&specification.ScenarioBuilder{}).
					Build(specification.NewScenarioSlug("baz", "bad")),
			},
			ExpectedTheses: nil,
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {
					b.WithScenario("bar", func(b *specification.ScenarioBuilder) {})
				})
			},
			ExpectedStories: []specification.Story{
				(&specification.StoryBuilder{}).
					WithScenario("bar", func(b *specification.ScenarioBuilder) {}).
					Build(specification.NewStorySlug("foo")),
			},
			ExpectedScenarios: []specification.Scenario{
				(&specification.ScenarioBuilder{}).
					Build(specification.NewScenarioSlug("foo", "bar")),
			},
			ExpectedTheses: nil,
		},
		{
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {
					b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
						b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
					})

					b.WithScenario("kap", func(b *specification.ScenarioBuilder) {
						b.WithThesis("dam", func(b *specification.ThesisBuilder) {})
					})
				})
				b.WithStory("qyz", func(b *specification.StoryBuilder) {
					b.WithScenario("qyp", func(b *specification.ScenarioBuilder) {
						b.WithThesis("dyq", func(b *specification.ThesisBuilder) {})
					})
				})
			},
			ExpectedStories: []specification.Story{
				(&specification.StoryBuilder{}).
					WithScenario("bar", func(b *specification.ScenarioBuilder) {
						b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
					}).
					WithScenario("kap", func(b *specification.ScenarioBuilder) {
						b.WithThesis("dam", func(b *specification.ThesisBuilder) {})
					}).
					Build(specification.NewStorySlug("foo")),
				(&specification.StoryBuilder{}).
					WithScenario("qyp", func(b *specification.ScenarioBuilder) {
						b.WithThesis("dyq", func(b *specification.ThesisBuilder) {})
					}).
					Build(specification.NewStorySlug("qyz")),
			},
			ExpectedScenarios: []specification.Scenario{
				(&specification.ScenarioBuilder{}).
					WithThesis("baz", func(b *specification.ThesisBuilder) {}).
					Build(specification.NewScenarioSlug("foo", "bar")),
				(&specification.ScenarioBuilder{}).
					WithThesis("dam", func(b *specification.ThesisBuilder) {}).
					Build(specification.NewScenarioSlug("foo", "kap")),
				(&specification.ScenarioBuilder{}).
					WithThesis("dyq", func(b *specification.ThesisBuilder) {}).
					Build(specification.NewScenarioSlug("qyz", "qyp")),
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).
					Build(specification.NewThesisSlug("foo", "bar", "baz")),
				(&specification.ThesisBuilder{}).
					Build(specification.NewThesisSlug("foo", "kap", "dam")),
				(&specification.ThesisBuilder{}).
					Build(specification.NewThesisSlug("qyz", "qyp", "dyq")),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			spec := errlessBuildSpec(t, c.Prepare)

			t.Run("stories", func(t *testing.T) {
				require.ElementsMatch(t, c.ExpectedStories, spec.Stories())
			})

			t.Run("stories_count", func(t *testing.T) {
				require.Equal(t, len(c.ExpectedStories), spec.StoriesCount())
			})

			t.Run("scenarios", func(t *testing.T) {
				require.ElementsMatch(t, c.ExpectedScenarios, spec.Scenarios())
			})

			t.Run("scenarios_count", func(t *testing.T) {
				require.Equal(t, len(c.ExpectedScenarios), spec.ScenariosCount())
			})

			t.Run("theses", func(t *testing.T) {
				require.ElementsMatch(t, c.ExpectedTheses, spec.Theses())
			})

			t.Run("theses_count", func(t *testing.T) {
				require.Equal(t, len(c.ExpectedTheses), spec.ThesesCount())
			})
		})
	}
}

func TestGetSpecificationStoryBySlug(t *testing.T) {
	t.Parallel()

	story := errlessBuildSpec(t, func(b *specification.Builder) {
		b.WithStory("foo", func(b *specification.StoryBuilder) {})
		b.WithStory("bar", func(b *specification.StoryBuilder) {})
	})

	var b specification.StoryBuilder

	foo, ok := story.Story("foo")
	require.True(t, ok)
	require.Equal(
		t,
		b.Build(
			specification.NewStorySlug("foo"),
		),
		foo,
	)

	b.Reset()

	bar, ok := story.Story("bar")
	require.True(t, ok)
	require.Equal(
		t,
		b.Build(
			specification.NewStorySlug("bar"),
		),
		bar,
	)

	_, ok = story.Story("baz")
	require.False(t, ok)
}

var errTest = errors.New("foo")

func TestIsWrappedInBuildError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError error
		ExpectedIs bool
	}{
		{
			GivenError: nil,
			ExpectedIs: false,
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(errors.New("bar")).
				Wrap("foo"),
			ExpectedIs: false,
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(errTest).
				Wrap("foo"),
			ExpectedIs: true,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedIs, errors.Is(c.GivenError, errTest))
		})
	}
}

func TestAsWrappedInBuildError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError error
		ExpectedAs bool
	}{
		{
			GivenError: nil,
			ExpectedAs: false,
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(errors.New("baz")).
				Wrap("bar"),
			ExpectedAs: false,
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(testError{}).
				Wrap("bar"),
			ExpectedAs: true,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target testError

			require.Equal(t, c.ExpectedAs, errors.As(c.GivenError, &target))
		})
	}
}

func TestAsBuildError(t *testing.T) {
	t.Parallel()

	type stringContext struct {
		Value string
		Ok    bool
	}

	type slugContext struct {
		Value specification.Slug
		Ok    bool
	}

	testCases := []struct {
		GivenError            error
		ShouldBeWrapped       bool
		ExpectedStringContext stringContext
		ExpectedSlugContext   slugContext
		ExpectedErrors        []error
	}{
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(nil).
				WithError(nil).
				Wrap("bob"),
			ShouldBeWrapped: false,
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(errors.New("bar")).
				Wrap("foo"),
			ShouldBeWrapped: true,
			ExpectedStringContext: stringContext{
				Value: "foo",
				Ok:    true,
			},
			ExpectedErrors: []error{
				errors.New("bar"),
			},
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(errors.New("foo")).
				WithError(errors.New("bar")).
				SluggedWrap(specification.NewScenarioSlug("a", "b")),
			ShouldBeWrapped: true,
			ExpectedSlugContext: slugContext{
				Value: specification.NewScenarioSlug("a", "b"),
				Ok:    true,
			},
			ExpectedErrors: []error{
				errors.New("foo"),
				errors.New("bar"),
			},
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(errors.New("foo")).
				WithError(nil).
				WithError(errors.New("bar")).
				SluggedWrap(specification.NewStorySlug("story")),
			ShouldBeWrapped: true,
			ExpectedSlugContext: slugContext{
				Value: specification.NewStorySlug("story"),
				Ok:    true,
			},
			ExpectedErrors: []error{
				errors.New("foo"),
				errors.New("bar"),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *specification.BuildError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("string_context", func(t *testing.T) {
					msg, ok := target.StringContext()

					require.Equal(t, c.ExpectedStringContext.Ok, ok)
					require.Equal(t, c.ExpectedStringContext.Value, msg)
				})

				t.Run("slug_context", func(t *testing.T) {
					slug, ok := target.SlugContext()

					require.Equal(t, c.ExpectedSlugContext.Ok, ok)
					require.Equal(t, c.ExpectedSlugContext.Value, slug)
				})

				t.Run("errors", func(t *testing.T) {
					require.ElementsMatch(t, c.ExpectedErrors, target.Errors())
				})
			})
		})
	}
}

func TestFormatBuildError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(errors.New("bar")).
				WithError(errors.New("baz")).
				Wrap("foo"),
			ExpectedErrorString: "foo: [bar; baz]",
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(errors.New("ba\nbaba")).
				WithError(errors.New("foo")).
				SluggedWrap(specification.NewThesisSlug("a", "b", "c")),
			ExpectedErrorString: "a.b.c: [ba\nbaba; foo]",
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(
					(&specification.BuildErrorWrapper{}).
						WithError(errors.New("doo")).
						WithError(errors.New("qoo")).
						Wrap("bar"),
				).
				Wrap("foo"),
			ExpectedErrorString: "foo: [bar: [doo; qoo]]",
		},
		{
			GivenError: (&specification.BuildErrorWrapper{}).
				WithError(
					(&specification.BuildErrorWrapper{}).
						WithError(errors.New("c")).
						WithError(errors.New("d")).
						Wrap("b"),
				).
				WithError(
					(&specification.BuildErrorWrapper{}).
						WithError(errors.New("f")).
						WithError(
							(&specification.BuildErrorWrapper{}).
								WithError(errors.New("h")).
								Wrap("g"),
						).
						Wrap("e"),
				).
				Wrap("a"),
			ExpectedErrorString: "a: [b: [c; d]; e: [f; g: [h]]]",
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

type testError struct{}

func (e testError) Error() string {
	return "test"
}
