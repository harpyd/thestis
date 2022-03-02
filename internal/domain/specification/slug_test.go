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

func TestSlugMustBeOneOfKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		GivenSlug       specification.Slug
		GivenMustBeKind specification.SlugKind
		ShouldBeErr     bool
		IsErr           func(err error) bool
	}{
		{
			Name:            "story_slug_is_story_slug",
			GivenSlug:       specification.NewStorySlug("foo"),
			GivenMustBeKind: specification.StorySlug,
			ShouldBeErr:     false,
		},
		{
			Name:            "story_slug_is_NOT_scenario_slug",
			GivenSlug:       specification.NewStorySlug("foo"),
			GivenMustBeKind: specification.ScenarioSlug,
			ShouldBeErr:     true,
			IsErr:           specification.IsNotScenarioSlugError,
		},
		{
			Name:            "story_slug_is_NOT_thesis_slug",
			GivenSlug:       specification.NewStorySlug("foo"),
			GivenMustBeKind: specification.ThesisSlug,
			ShouldBeErr:     true,
			IsErr:           specification.IsNotThesisSlugError,
		},
		{
			Name:            "scenario_slug_is_NOT_story_slug",
			GivenSlug:       specification.NewScenarioSlug("foo", "bar"),
			GivenMustBeKind: specification.StorySlug,
			ShouldBeErr:     true,
			IsErr:           specification.IsNotStorySlugError,
		},
		{
			Name:            "scenario_slug_is_scenario_slug",
			GivenSlug:       specification.NewScenarioSlug("foo", "bar"),
			GivenMustBeKind: specification.ScenarioSlug,
			ShouldBeErr:     false,
		},
		{
			Name:            "scenario_slug_is_NOT_thesis_slug",
			GivenSlug:       specification.NewScenarioSlug("foo", "bar"),
			GivenMustBeKind: specification.ThesisSlug,
			ShouldBeErr:     true,
			IsErr:           specification.IsNotThesisSlugError,
		},
		{
			Name:            "thesis_slug_is_NOT_story_slug",
			GivenSlug:       specification.NewThesisSlug("foo", "bar", "baz"),
			GivenMustBeKind: specification.StorySlug,
			ShouldBeErr:     true,
			IsErr:           specification.IsNotStorySlugError,
		},
		{
			Name:            "thesis_slug_is_NOT_scenario_slug",
			GivenSlug:       specification.NewThesisSlug("foo", "bar", "baz"),
			GivenMustBeKind: specification.ScenarioSlug,
			ShouldBeErr:     true,
			IsErr:           specification.IsNotScenarioSlugError,
		},
		{
			Name:            "thesis_slug_is_thesis_slug",
			GivenSlug:       specification.NewThesisSlug("foo", "bar", "bad"),
			GivenMustBeKind: specification.ThesisSlug,
			ShouldBeErr:     false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var err error

			switch c.GivenMustBeKind {
			case specification.StorySlug:
				err = c.GivenSlug.MustBeStoryKind()
			case specification.ScenarioSlug:
				err = c.GivenSlug.MustBeScenarioKind()
			case specification.ThesisSlug:
				err = c.GivenSlug.MustBeThesisKind()
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

func TestSlugErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "not_story_slug_error",
			Err:   specification.NewNotStorySlugError(),
			IsErr: specification.IsNotStorySlugError,
		},
		{
			Name:     "NON_not_story_slug_error",
			Err:      errors.New("not story slug"),
			IsErr:    specification.IsNotStorySlugError,
			Reversed: true,
		},
		{
			Name:  "not_scenario_slug_error",
			Err:   specification.NewNotScenarioSlugError(),
			IsErr: specification.IsNotScenarioSlugError,
		},
		{
			Name:     "NON_not_scenario_slug_error",
			Err:      errors.New("not scenario slug"),
			IsErr:    specification.IsNotScenarioSlugError,
			Reversed: true,
		},
		{
			Name:  "not_thesis_slug_error",
			Err:   specification.NewNotThesisSlugError(),
			IsErr: specification.IsNotThesisSlugError,
		},
		{
			Name:     "NON_not_thesis_slug_error",
			Err:      errors.New("not thesis slug"),
			IsErr:    specification.IsNotThesisSlugError,
			Reversed: true,
		},
		{
			Name:  "empty_slug_error",
			Err:   specification.NewEmptySlugError(),
			IsErr: specification.IsEmptySlugError,
		},
		{
			Name:     "NON_empty_slug_error",
			Err:      errors.New("empty slug"),
			IsErr:    specification.IsEmptySlugError,
			Reversed: true,
		},
		{
			Name: "story_slug_already_exists_error",
			Err: specification.NewSlugAlreadyExistsError(
				specification.NewStorySlug("story"),
			),
			IsErr: specification.IsStorySlugAlreadyExistsError,
		},
		{
			Name:     "NON_story_slug_already_exists_error",
			Err:      errors.New("story"),
			IsErr:    specification.IsStorySlugAlreadyExistsError,
			Reversed: true,
		},
		{
			Name: "scenario_slug_already_exists_error",
			Err: specification.NewSlugAlreadyExistsError(
				specification.NewScenarioSlug("story", "scenario"),
			),
			IsErr: specification.IsScenarioSlugAlreadyExistsError,
		},
		{
			Name:     "NON_scenario_slug_already_exists_error",
			Err:      errors.New("scenario"),
			IsErr:    specification.IsScenarioSlugAlreadyExistsError,
			Reversed: true,
		},
		{
			Name: "thesis_slug_already_exists_error",
			Err: specification.NewSlugAlreadyExistsError(
				specification.NewThesisSlug("story", "scenario", "thesis"),
			),
			IsErr: specification.IsThesisSlugAlreadyExistsError,
		},
		{
			Name:     "NON_thesis_slug_already_exists_error",
			Err:      errors.New("thesis"),
			IsErr:    specification.IsThesisSlugAlreadyExistsError,
			Reversed: true,
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
			Name:     "NON_build_story_error",
			Err:      errors.New("foo"),
			IsErr:    specification.IsBuildStoryError,
			Reversed: true,
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
			Name:     "NON_build_scenario_error",
			Err:      errors.New("bar"),
			IsErr:    specification.IsBuildScenarioError,
			Reversed: true,
		},
		{
			Name: "build_thesis_error",
			Err: specification.NewBuildSluggedError(
				errors.New("baz"),
				specification.NewThesisSlug("story", "scenario", "thesis"),
			),
			IsErr: specification.IsBuildThesisError,
		},
		{
			Name:     "NON_build_thesis_error",
			Err:      errors.New("baz"),
			IsErr:    specification.IsBuildThesisError,
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
