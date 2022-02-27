package specification_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestSlugIsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name string

		GivenSlug specification.Slug

		ExpectedStory    string
		ExpectedScenario string
		ExpectedThesis   string
		ExpectedString   string
		ExpectedKind     specification.SlugKind
	}{
		{
			GivenSlug:        specification.Slug{},
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "",
			ExpectedKind:     specification.NoSlug,
		},
		{
			GivenSlug:        specification.NewStorySlug(""),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*",
			ExpectedKind:     specification.StorySlug,
		},
		{
			GivenSlug:        specification.NewStorySlug("story"),
			ExpectedStory:    "story",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "story",
			ExpectedKind:     specification.StorySlug,
		},
		{
			GivenSlug:        specification.NewScenarioSlug("", ""),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*.*",
			ExpectedKind:     specification.ScenarioSlug,
		},
		{
			GivenSlug:        specification.NewScenarioSlug("story", ""),
			ExpectedStory:    "story",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "story.*",
			ExpectedKind:     specification.ScenarioSlug,
		},
		{
			GivenSlug:        specification.NewScenarioSlug("", "scenario"),
			ExpectedStory:    "",
			ExpectedScenario: "scenario",
			ExpectedThesis:   "",
			ExpectedString:   "*.scenario",
			ExpectedKind:     specification.ScenarioSlug,
		},
		{
			GivenSlug:        specification.NewScenarioSlug("story", "scenario"),
			ExpectedStory:    "story",
			ExpectedScenario: "scenario",
			ExpectedThesis:   "",
			ExpectedString:   "story.scenario",
			ExpectedKind:     specification.ScenarioSlug,
		},
		{
			GivenSlug:        specification.NewThesisSlug("", "", ""),
			ExpectedStory:    "",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "*.*.*",
			ExpectedKind:     specification.ThesisSlug,
		},
		{
			GivenSlug:        specification.NewThesisSlug("story", "", ""),
			ExpectedStory:    "story",
			ExpectedScenario: "",
			ExpectedThesis:   "",
			ExpectedString:   "story.*.*",
			ExpectedKind:     specification.ThesisSlug,
		},
		{
			GivenSlug:        specification.NewThesisSlug("story", "scenario", ""),
			ExpectedStory:    "story",
			ExpectedScenario: "scenario",
			ExpectedThesis:   "",
			ExpectedString:   "story.scenario.*",
			ExpectedKind:     specification.ThesisSlug,
		},
		{
			GivenSlug:        specification.NewThesisSlug("story", "scenario", "thesis"),
			ExpectedStory:    "story",
			ExpectedScenario: "scenario",
			ExpectedThesis:   "thesis",
			ExpectedString:   "story.scenario.thesis",
			ExpectedKind:     specification.ThesisSlug,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.ExpectedString, func(t *testing.T) {
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
		})
	}
}

func TestSlugErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name  string
		Err   error
		IsErr func(err error) bool
	}{
		{
			Name:  "empty_slug_error",
			Err:   specification.NewEmptySlugError(),
			IsErr: specification.IsEmptySlugError,
		},
		{
			Name: "story_slug_already_exists_error",
			Err: specification.NewSlugAlreadyExistsError(
				specification.NewStorySlug("story"),
			),
			IsErr: specification.IsStorySlugAlreadyExistsError,
		},
		{
			Name: "scenario_slug_already_exists_error",
			Err: specification.NewSlugAlreadyExistsError(
				specification.NewScenarioSlug("story", "scenario"),
			),
			IsErr: specification.IsScenarioSlugAlreadyExistsError,
		},
		{
			Name: "thesis_slug_already_exists_error",
			Err: specification.NewSlugAlreadyExistsError(
				specification.NewThesisSlug("story", "scenario", "thesis"),
			),
			IsErr: specification.IsThesisSlugAlreadyExistsError,
		},
		{
			Name: "no_such_story_slug_error",
			Err: specification.NewNoSuchSlugError(
				specification.NewStorySlug("story"),
			),
			IsErr: specification.IsNoSuchStorySlugError,
		},
		{
			Name: "no_such_scenario_slug_error",
			Err: specification.NewNoSuchSlugError(
				specification.NewScenarioSlug("story", "scenario"),
			),
			IsErr: specification.IsNoSuchScenarioSlugError,
		},
		{
			Name: "no_such_thesis_slug_error",
			Err: specification.NewNoSuchSlugError(
				specification.NewThesisSlug("story", "scenario", "thesis"),
			),
			IsErr: specification.IsNoSuchThesisSlugError,
		},
		{
			Name: "build_story_error",
			Err: specification.NewBuildSluggedError(
				errors.New("foo"),
				specification.NewStorySlug("story"),
			),
			IsErr: specification.IsBuildStoryError,
		},
		{
			Name: "build_scenario_error",
			Err: specification.NewBuildSluggedError(
				errors.New("bar"),
				specification.NewScenarioSlug("story", "scenario"),
			),
			IsErr: specification.IsBuildScenarioError,
		},
		{
			Name: "build_thesis_error",
			Err: specification.NewBuildSluggedError(
				errors.New("baz"),
				specification.NewThesisSlug("story", "scenario", "thesis"),
			),
			IsErr: specification.IsBuildThesisError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.True(t, c.IsErr(c.Err))
		})
	}
}
