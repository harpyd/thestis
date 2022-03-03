package specification_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func errlessBuildSpecification(
	t *testing.T,
	prepare func(b *specification.Builder),
) *specification.Specification {
	t.Helper()

	builder := specification.NewBuilder()

	prepare(builder)

	return builder.ErrlessBuild()
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
			IsErr:       specification.IsNoSpecificationStoriesError,
		},
		{
			Name: "three_stories",
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {})
				b.WithStory("bar", func(b *specification.StoryBuilder) {})
				b.WithStory("baz", func(b *specification.StoryBuilder) {})
			},
			ExpectedStories: []specification.Story{
				specification.NewStoryBuilder().
					ErrlessBuild(specification.NewStorySlug("foo")),
				specification.NewStoryBuilder().
					ErrlessBuild(specification.NewStorySlug("bar")),
				specification.NewStoryBuilder().
					ErrlessBuild(specification.NewStorySlug("baz")),
			},
			WantThisErr: false,
			IsErr:       specification.IsNoSpecificationStoriesError,
		},
		{
			Name: "story_already_exists",
			Prepare: func(b *specification.Builder) {
				b.WithStory("oops", func(b *specification.StoryBuilder) {})
				b.WithStory("oops", func(b *specification.StoryBuilder) {})
			},
			WantThisErr: true,
			IsErr:       specification.IsStorySlugAlreadyExistsError,
		},
		{
			Name: "no_story_scenarios",
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {})
			},
			WantThisErr: true,
			IsErr:       specification.IsNoStoryScenariosError,
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
				specification.NewStoryBuilder().
					WithScenario("bar", func(b *specification.ScenarioBuilder) {}).
					ErrlessBuild(specification.NewStorySlug("foo")),
				specification.NewStoryBuilder().
					WithScenario("bad", func(b *specification.ScenarioBuilder) {}).
					ErrlessBuild(specification.NewStorySlug("baz")),
			},
			ExpectedScenarios: []specification.Scenario{
				specification.NewScenarioBuilder().
					ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
				specification.NewScenarioBuilder().
					ErrlessBuild(specification.NewScenarioSlug("baz", "bad")),
			},
			WantThisErr: false,
			IsErr:       specification.IsNoStoryScenariosError,
		},
		{
			Name: "not_scenario_theses",
			Prepare: func(b *specification.Builder) {
				b.WithStory("foo", func(b *specification.StoryBuilder) {
					b.WithScenario("bar", func(b *specification.ScenarioBuilder) {})
				})
			},
			WantThisErr: true,
			IsErr:       specification.IsNoScenarioThesesError,
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
				specification.NewStoryBuilder().
					WithScenario("bar", func(b *specification.ScenarioBuilder) {
						b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
					}).
					WithScenario("kap", func(b *specification.ScenarioBuilder) {
						b.WithThesis("dam", func(b *specification.ThesisBuilder) {})
					}).
					ErrlessBuild(specification.NewStorySlug("foo")),
				specification.NewStoryBuilder().
					WithScenario("qyp", func(b *specification.ScenarioBuilder) {
						b.WithThesis("dyq", func(b *specification.ThesisBuilder) {})
					}).
					ErrlessBuild(specification.NewStorySlug("qyz")),
			},
			ExpectedScenarios: []specification.Scenario{
				specification.NewScenarioBuilder().
					WithThesis("baz", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
				specification.NewScenarioBuilder().
					WithThesis("dam", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("foo", "kap")),
				specification.NewScenarioBuilder().
					WithThesis("dyq", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("qyz", "qyp")),
			},
			ExpectedTheses: []specification.Thesis{
				specification.NewThesisBuilder().
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "baz")),
				specification.NewThesisBuilder().
					ErrlessBuild(specification.NewThesisSlug("foo", "kap", "dam")),
				specification.NewThesisBuilder().
					ErrlessBuild(specification.NewThesisSlug("qyz", "qyp", "dyq")),
			},
			WantThisErr: false,
			IsErr:       specification.IsNoScenarioThesesError,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewBuilder()

			c.Prepare(builder)

			spec, err := builder.Build()

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))

			t.Run("stories", func(t *testing.T) {
				assert.ElementsMatch(t, c.ExpectedStories, spec.Stories())
			})

			t.Run("stories_count", func(t *testing.T) {
				assert.Equal(t, len(c.ExpectedStories), spec.StoriesCount())
			})

			t.Run("scenarios", func(t *testing.T) {
				assert.ElementsMatch(t, c.ExpectedScenarios, spec.Scenarios())
			})

			t.Run("scenarios_count", func(t *testing.T) {
				assert.Equal(t, len(c.ExpectedScenarios), spec.ScenariosCount())
			})

			t.Run("theses", func(t *testing.T) {
				assert.ElementsMatch(t, c.ExpectedTheses, spec.Theses())
			})

			t.Run("theses_count", func(t *testing.T) {
				assert.Equal(t, len(c.ExpectedTheses), spec.ThesesCount())
			})
		})
	}
}

func TestSpecificationErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "build_specification_error",
			Err:   specification.NewBuildSpecificationError(errors.New("badaboom")),
			IsErr: specification.IsBuildSpecificationError,
		},
		{
			Name:     "NON_build_specification_error",
			Err:      errors.New("badaboom"),
			IsErr:    specification.IsBuildSpecificationError,
			Reversed: true,
		},
		{
			Name:  "no_specification_stories_error",
			Err:   specification.NewNoSpecificationStoriesError(),
			IsErr: specification.IsNoSpecificationStoriesError,
		},
		{
			Name:     "NON_no_specification_stories_error",
			Err:      errors.New("another"),
			IsErr:    specification.IsNoSpecificationStoriesError,
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
