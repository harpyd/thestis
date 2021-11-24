package specification_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestIsNoStoryError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		Err            error
		IsNoStoryError bool
	}{
		{
			Name:           "no_story_error_is_no_story_error",
			Err:            specification.NewNoStoryError("someStory"),
			IsNoStoryError: true,
		},
		{
			Name:           "another_error_isnt_no_story_error",
			Err:            specification.NewNoThesisError("someThesis"),
			IsNoStoryError: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsNoStoryError, specification.IsNoStoryError(c.Err))
		})
	}
}

func TestIsNoScenarioError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		Err               error
		IsNoScenarioError bool
	}{
		{
			Name:              "no_scenario_error_is_no_scenario_error",
			Err:               specification.NewNoScenarioError("someScenario"),
			IsNoScenarioError: true,
		},
		{
			Name:              "another_error_isnt_no_scenario_error",
			Err:               specification.NewNoThesisError("someThesis"),
			IsNoScenarioError: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsNoScenarioError, specification.IsNoScenarioError(c.Err))
		})
	}
}

func TestIsNoThesisError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		Err             error
		IsNoThesisError bool
	}{
		{
			Name:            "no_thesis_error_is_no_thesis_error",
			Err:             specification.NewNoThesisError("someThesis"),
			IsNoThesisError: true,
		},
		{
			Name:            "another_error_isnt_no_thesis_error",
			Err:             specification.NewNoStoryError("someStory"),
			IsNoThesisError: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsNoThesisError, specification.IsNoThesisError(c.Err))
		})
	}
}
