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

func TestBuilder_Build_no_stories(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()

	_, err := builder.Build()

	require.True(t, specification.IsNoSpecificationStoriesError(err))
}

func TestBuilder_WithStory(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithStory("firstStory", func(b *specification.StoryBuilder) {
		b.WithDescription("this is a first story")
	})
	builder.WithStory("secondStory", func(b *specification.StoryBuilder) {
		b.WithDescription("this is a second story")
	})

	spec := builder.ErrlessBuild()

	expectedFirstStory := specification.NewStoryBuilder().
		WithDescription("this is a first story").
		ErrlessBuild(specification.NewStorySlug("firstStory"))

	actualFirstStory, ok := spec.Story("firstStory")
	require.True(t, ok)
	require.Equal(t, expectedFirstStory, actualFirstStory)

	expectedSecondStory := specification.NewStoryBuilder().
		WithDescription("this is a second story").
		ErrlessBuild(specification.NewStorySlug("secondStory"))

	actualSecondStory, ok := spec.Story("secondStory")
	require.True(t, ok)
	require.Equal(t, expectedSecondStory, actualSecondStory)
}

func TestBuilder_WithStory_when_already_exists(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithStory("story", func(b *specification.StoryBuilder) {
		b.WithDescription("this is a story")
	})
	builder.WithStory("story", func(b *specification.StoryBuilder) {
		b.WithDescription("this is a same story")
	})

	_, err := builder.Build()

	require.True(t, specification.IsStorySlugAlreadyExistsError(err))
}

func TestSpecification_Stories(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithStory("foo", func(b *specification.StoryBuilder) {})
	builder.WithStory("bar", func(b *specification.StoryBuilder) {})

	spec := builder.ErrlessBuild()

	t.Run("match", func(t *testing.T) {
		t.Parallel()

		expected := []specification.Story{
			specification.NewStoryBuilder().ErrlessBuild(
				specification.NewStorySlug("foo"),
			),
			specification.NewStoryBuilder().ErrlessBuild(
				specification.NewStorySlug("bar"),
			),
		}

		assert.ElementsMatch(t, expected, spec.Stories())
	})

	t.Run("count", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, 2, spec.StoriesCount())
	})
}

func TestSpecification_Scenarios(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithStory("foo", func(b *specification.StoryBuilder) {
		b.WithScenario("bar", func(b *specification.ScenarioBuilder) {})
	})
	builder.WithStory("baz", func(b *specification.StoryBuilder) {
		b.WithScenario("qyz", func(b *specification.ScenarioBuilder) {})
	})

	spec := builder.ErrlessBuild()

	t.Run("match", func(t *testing.T) {
		t.Parallel()

		expected := []specification.Scenario{
			specification.NewScenarioBuilder().ErrlessBuild(
				specification.NewScenarioSlug("foo", "bar"),
			),
			specification.NewScenarioBuilder().ErrlessBuild(
				specification.NewScenarioSlug("baz", "qyz"),
			),
		}

		assert.ElementsMatch(t, expected, spec.Scenarios())
	})

	t.Run("count", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, 2, spec.ScenariosCount())
	})
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
