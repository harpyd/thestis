package specification_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func errlessBuildSpecification(
	t *testing.T,
	prepare func(b *specification.Builder),
) *specification.Specification {
	t.Helper()

	var b specification.Builder

	prepare(&b)

	return b.ErrlessBuild()
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

			id := errlessBuildSpecification(t, c.Prepare).ID()

			require.Equal(t, c.ExpectedID, id)
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

			ownerID := errlessBuildSpecification(t, c.Prepare).OwnerID()

			require.Equal(t, c.ExpectedOwnerID, ownerID)
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

			testCampaignID := errlessBuildSpecification(t, c.Prepare).TestCampaignID()

			require.Equal(t, c.ExpectedTestCampaignID, testCampaignID)
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

			loadedAt := errlessBuildSpecification(t, c.Prepare).LoadedAt()

			require.Equal(t, c.ExpectedLoadedAt, loadedAt)
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

			author := errlessBuildSpecification(t, c.Prepare).Author()

			require.Equal(t, c.ExpectedAuthor, author)
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

			title := errlessBuildSpecification(t, c.Prepare).Title()

			require.Equal(t, c.ExpectedTitle, title)
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

			description := errlessBuildSpecification(t, c.Prepare).Description()

			require.Equal(t, c.ExpectedDescription, description)
		})
	}
}

func TestBuildSpecificationWithStories(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		Prepare           func(b *specification.Builder)
		ExpectedStories   []specification.Story
		ExpectedScenarios []specification.Scenario
		ExpectedTheses    []specification.Thesis
		WantThisErr       bool
		IsErr             func(err error) bool
	}{
		{
			Name:        "no_stories",
			Prepare:     func(b *specification.Builder) {},
			WantThisErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoSpecificationStories)
			},
		},
		{
			Name: "three_stories",
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {})
				b.WithStory("bar", func(b *specification.StoryBuilder) {})
				b.WithStory("baz", func(b *specification.StoryBuilder) {})
			},
			ExpectedStories: []specification.Story{
				(&specification.StoryBuilder{}).
					ErrlessBuild(specification.NewStorySlug("foo")),
				(&specification.StoryBuilder{}).
					ErrlessBuild(specification.NewStorySlug("bar")),
				(&specification.StoryBuilder{}).
					ErrlessBuild(specification.NewStorySlug("baz")),
			},
			WantThisErr: false,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoSpecificationStories)
			},
		},
		{
			Name: "story_already_exists",
			Prepare: func(b *specification.Builder) {
				b.WithStory("oops", func(b *specification.StoryBuilder) {})
				b.WithStory("oops", func(b *specification.StoryBuilder) {})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				var target *specification.DuplicatedError

				return errors.As(err, &target)
			},
		},
		{
			Name: "no_story_scenarios",
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoStoryScenarios)
			},
		},
		{
			Name: "stories_having_scenarios",
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
					ErrlessBuild(specification.NewStorySlug("foo")),
				(&specification.StoryBuilder{}).
					WithScenario("bad", func(b *specification.ScenarioBuilder) {}).
					ErrlessBuild(specification.NewStorySlug("baz")),
			},
			ExpectedScenarios: []specification.Scenario{
				(&specification.ScenarioBuilder{}).
					ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
				(&specification.ScenarioBuilder{}).
					ErrlessBuild(specification.NewScenarioSlug("baz", "bad")),
			},
			WantThisErr: false,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoStoryScenarios)
			},
		},
		{
			Name: "not_scenario_theses",
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {
					b.WithScenario("bar", func(b *specification.ScenarioBuilder) {})
				})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoScenarioTheses)
			},
		},
		{
			Name: "stories_having_scenarios_having_theses",
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
					ErrlessBuild(specification.NewStorySlug("foo")),
				(&specification.StoryBuilder{}).
					WithScenario("qyp", func(b *specification.ScenarioBuilder) {
						b.WithThesis("dyq", func(b *specification.ThesisBuilder) {})
					}).
					ErrlessBuild(specification.NewStorySlug("qyz")),
			},
			ExpectedScenarios: []specification.Scenario{
				(&specification.ScenarioBuilder{}).
					WithThesis("baz", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
				(&specification.ScenarioBuilder{}).
					WithThesis("dam", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("foo", "kap")),
				(&specification.ScenarioBuilder{}).
					WithThesis("dyq", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("qyz", "qyp")),
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "baz")),
				(&specification.ThesisBuilder{}).
					ErrlessBuild(specification.NewThesisSlug("foo", "kap", "dam")),
				(&specification.ThesisBuilder{}).
					ErrlessBuild(specification.NewThesisSlug("qyz", "qyp", "dyq")),
			},
			WantThisErr: false,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoScenarioTheses)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var b specification.Builder

			c.Prepare(&b)

			spec, err := b.Build()

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))

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

	story := errlessBuildSpecification(t, func(b *specification.Builder) {
		b.WithStory("foo", func(b *specification.StoryBuilder) {})
		b.WithStory("bar", func(b *specification.StoryBuilder) {})
	})

	var b specification.StoryBuilder

	foo, ok := story.Story("foo")
	require.True(t, ok)
	require.Equal(
		t,
		b.ErrlessBuild(
			specification.NewStorySlug("foo"),
		),
		foo,
	)

	b.Reset()

	bar, ok := story.Story("bar")
	require.True(t, ok)
	require.Equal(
		t,
		b.ErrlessBuild(
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

	testCases := []struct {
		GivenError      error
		ShouldBeWrapped bool
		ExpectedMessage string
		ExpectedErrors  []error
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
			ExpectedMessage: "foo",
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
			ExpectedMessage: "a.b",
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
				SluggedWrap(specification.AnyThesisSlug()),
			ShouldBeWrapped: true,
			ExpectedMessage: "*.*.*",
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

				t.Run("message", func(t *testing.T) {
					require.Equal(t, c.ExpectedMessage, target.Message())
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
