package specification_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestNewSlug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenSlug        specification.Slug
		ExpectedStory    string
		ExpectedScenario string
		ExpectedThesis   string
		ExpectedString   string
		ExpectedKind     specification.SlugKind
		ExpectedPartial  string
	}{
		{
			GivenSlug:        specification.Slug{},
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "",
			ExpectedKind:     specification.NoSlug,
			ExpectedPartial:  "",
		},
		{
			GivenSlug:        specification.AnyStorySlug(),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*",
			ExpectedKind:     specification.StorySlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewStorySlug(""),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*",
			ExpectedKind:     specification.StorySlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewStorySlug("story"),
			ExpectedStory:    "story",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "story",
			ExpectedKind:     specification.StorySlug,
			ExpectedPartial:  "story",
		},
		{
			GivenSlug:        specification.AnyScenarioSlug(),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*.*",
			ExpectedKind:     specification.ScenarioSlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewScenarioSlug("", ""),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*.*",
			ExpectedKind:     specification.ScenarioSlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewScenarioSlug("story", ""),
			ExpectedStory:    "story",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "story.*",
			ExpectedKind:     specification.ScenarioSlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewScenarioSlug("", "scenario"),
			ExpectedStory:    "",
			ExpectedScenario: "scenario",
			ExpectedThesis:   "",
			ExpectedString:   "*.scenario",
			ExpectedKind:     specification.ScenarioSlug,
			ExpectedPartial:  "scenario",
		},
		{
			GivenSlug:        specification.NewScenarioSlug("story", "scenario"),
			ExpectedStory:    "story",
			ExpectedScenario: "scenario",
			ExpectedThesis:   "",
			ExpectedString:   "story.scenario",
			ExpectedKind:     specification.ScenarioSlug,
			ExpectedPartial:  "scenario",
		},
		{
			GivenSlug:        specification.AnyThesisSlug(),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*.*.*",
			ExpectedKind:     specification.ThesisSlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewThesisSlug("", "", ""),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*.*.*",
			ExpectedKind:     specification.ThesisSlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewThesisSlug("story", "", ""),
			ExpectedStory:    "story",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "story.*.*",
			ExpectedKind:     specification.ThesisSlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewThesisSlug("story", "scenario", ""),
			ExpectedStory:    "story",
			ExpectedScenario: "scenario",
			ExpectedThesis:   "",
			ExpectedString:   "story.scenario.*",
			ExpectedKind:     specification.ThesisSlug,
			ExpectedPartial:  "*",
		},
		{
			GivenSlug:        specification.NewThesisSlug("story", "scenario", "thesis"),
			ExpectedStory:    "story",
			ExpectedScenario: "scenario",
			ExpectedThesis:   "thesis",
			ExpectedString:   "story.scenario.thesis",
			ExpectedKind:     specification.ThesisSlug,
			ExpectedPartial:  "thesis",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			t.Run("story", func(t *testing.T) {
				require.Equal(t, c.ExpectedStory, c.GivenSlug.Story())
			})

			t.Run("scenario", func(t *testing.T) {
				require.Equal(t, c.ExpectedScenario, c.GivenSlug.Scenario())
			})

			t.Run("thesis", func(t *testing.T) {
				require.Equal(t, c.ExpectedThesis, c.GivenSlug.Thesis())
			})

			t.Run("string", func(t *testing.T) {
				require.Equal(t, c.ExpectedString, c.GivenSlug.String())
			})

			t.Run("kind", func(t *testing.T) {
				require.Equal(t, c.ExpectedKind, c.GivenSlug.Kind())
			})

			t.Run("partial", func(t *testing.T) {
				require.Equal(t, c.ExpectedPartial, c.GivenSlug.Partial())
			})
		})
	}
}

func TestSlugToStoryKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenSlug    specification.Slug
		ExpectedSlug specification.Slug
	}{
		{
			GivenSlug:    specification.Slug{},
			ExpectedSlug: specification.Slug{},
		},
		{
			GivenSlug:    specification.NewStorySlug("foo"),
			ExpectedSlug: specification.NewStorySlug("foo"),
		},
		{
			GivenSlug:    specification.NewScenarioSlug("foo", "bar"),
			ExpectedSlug: specification.NewStorySlug("foo"),
		},
		{
			GivenSlug:    specification.NewThesisSlug("foo", "bar", "baz"),
			ExpectedSlug: specification.NewStorySlug("foo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedSlug, c.GivenSlug.ToStoryKind())
		})
	}
}

func TestSlugToThesisKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenSlug    specification.Slug
		ExpectedSlug specification.Slug
	}{
		{
			GivenSlug:    specification.Slug{},
			ExpectedSlug: specification.Slug{},
		},
		{
			GivenSlug:    specification.NewStorySlug("foo"),
			ExpectedSlug: specification.NewThesisSlug("foo", "", ""),
		},
		{
			GivenSlug:    specification.NewScenarioSlug("foo", "bar"),
			ExpectedSlug: specification.NewThesisSlug("foo", "bar", ""),
		},
		{
			GivenSlug:    specification.NewThesisSlug("foo", "bar", "baz"),
			ExpectedSlug: specification.NewThesisSlug("foo", "bar", "baz"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedSlug, c.GivenSlug.ToThesisKind())
		})
	}
}

func TestSlugToScenarioKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenSlug    specification.Slug
		ExpectedSlug specification.Slug
	}{
		{
			GivenSlug:    specification.Slug{},
			ExpectedSlug: specification.Slug{},
		},
		{
			GivenSlug:    specification.NewStorySlug("foo"),
			ExpectedSlug: specification.NewScenarioSlug("foo", ""),
		},
		{
			GivenSlug:    specification.NewScenarioSlug("foo", "bar"),
			ExpectedSlug: specification.NewScenarioSlug("foo", "bar"),
		},
		{
			GivenSlug:    specification.NewThesisSlug("foo", "bar", "baz"),
			ExpectedSlug: specification.NewScenarioSlug("foo", "bar"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedSlug, c.GivenSlug.ToScenarioKind())
		})
	}
}

func TestSlugShouldBeOneOfKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		GivenSlug   specification.Slug
		GivenKind   specification.SlugKind
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "story_slug_is_story_slug",
			GivenSlug:   specification.NewStorySlug("foo"),
			GivenKind:   specification.StorySlug,
			ShouldBeErr: false,
		},
		{
			Name:        "story_slug_is_NOT_scenario_slug",
			GivenSlug:   specification.NewStorySlug("foo"),
			GivenKind:   specification.ScenarioSlug,
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNotScenarioSlug)
			},
		},
		{
			Name:        "story_slug_is_NOT_thesis_slug",
			GivenSlug:   specification.NewStorySlug("foo"),
			GivenKind:   specification.ThesisSlug,
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNotThesisSlug)
			},
		},
		{
			Name:        "scenario_slug_is_NOT_story_slug",
			GivenSlug:   specification.NewScenarioSlug("foo", "bar"),
			GivenKind:   specification.StorySlug,
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNotStorySlug)
			},
		},
		{
			Name:        "scenario_slug_is_scenario_slug",
			GivenSlug:   specification.NewScenarioSlug("foo", "bar"),
			GivenKind:   specification.ScenarioSlug,
			ShouldBeErr: false,
		},
		{
			Name:        "scenario_slug_is_NOT_thesis_slug",
			GivenSlug:   specification.NewScenarioSlug("foo", "bar"),
			GivenKind:   specification.ThesisSlug,
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNotThesisSlug)
			},
		},
		{
			Name:        "thesis_slug_is_NOT_story_slug",
			GivenSlug:   specification.NewThesisSlug("foo", "bar", "baz"),
			GivenKind:   specification.StorySlug,
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNotStorySlug)
			},
		},
		{
			Name:        "thesis_slug_is_NOT_scenario_slug",
			GivenSlug:   specification.NewThesisSlug("foo", "bar", "baz"),
			GivenKind:   specification.ScenarioSlug,
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNotScenarioSlug)
			},
		},
		{
			Name:        "thesis_slug_is_thesis_slug",
			GivenSlug:   specification.NewThesisSlug("foo", "bar", "bad"),
			GivenKind:   specification.ThesisSlug,
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var err error

			switch c.GivenKind {
			case specification.StorySlug:
				err = c.GivenSlug.ShouldBeStoryKind()

			case specification.ScenarioSlug:
				err = c.GivenSlug.ShouldBeScenarioKind()

			case specification.ThesisSlug:
				err = c.GivenSlug.ShouldBeThesisKind()

			case specification.NoSlug:
			}

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestAsDuplicatedError(t *testing.T) {
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
			GivenError:      &specification.DuplicatedError{},
			ShouldBeWrapped: true,
			ExpectedSlug:    specification.Slug{},
		},
		{
			GivenError:      specification.NewDuplicatedError(specification.Slug{}),
			ShouldBeWrapped: true,
			ExpectedSlug:    specification.Slug{},
		},
		{
			GivenError: specification.NewDuplicatedError(
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

			var target *specification.DuplicatedError

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

func TestFormatDuplicatedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &specification.DuplicatedError{},
			ExpectedErrorString: "",
		},
		{
			GivenError: specification.NewDuplicatedError(
				specification.NewStorySlug("foo"),
			),
			ExpectedErrorString: `"foo" already exists`,
		},
		{
			GivenError: specification.NewDuplicatedError(
				specification.NewScenarioSlug("a", "b"),
			),
			ExpectedErrorString: `"b" already exists`,
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
