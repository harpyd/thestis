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

	for _, c := range testCases {
		c := c

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

	for _, c := range testCases {
		c := c

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

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsNoThesisError, specification.IsNoThesisError(c.Err))
		})
	}
}

func TestIsNotAllowedHTTPMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name  string
		Err   error
		IsErr bool
	}{
		{
			Name:  "not_allowed_http_method_error_is_not_allowed_http_method_error",
			Err:   specification.NewNotAllowedHTTPMethodError("POZT"),
			IsErr: true,
		},
		{
			Name:  "another_error_isnt_not_allowed_http_method_error",
			Err:   specification.NewNoThesisError("POZT"),
			IsErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsErr, specification.IsNotAllowedHTTPMethodError(c.Err))
		})
	}
}

func TestIsNotAllowedKeywordError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name  string
		Err   error
		IsErr bool
	}{
		{
			Name:  "not_allowed_keyword_error_is_not_allowed_keyword_error",
			Err:   specification.NewNotAllowedKeywordError("zen"),
			IsErr: true,
		},
		{
			Name:  "another_error_isnt_not_allowed_keyword_error",
			Err:   specification.NewNotAllowedHTTPMethodError("zen"),
			IsErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsErr, specification.IsNotAllowedKeywordError(c.Err))
		})
	}
}

func TestIsNotAllowedContentTypeError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name  string
		Err   error
		IsErr bool
	}{
		{
			Name:  "not_allowed_content_type_error_is_not_allowed_content_type_error",
			Err:   specification.NewNotAllowedContentTypeError("some/content"),
			IsErr: true,
		},
		{
			Name:  "another_error_isnt_not_allowed_content_type_error",
			Err:   specification.NewNoStoryError("some/content"),
			IsErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsErr, specification.IsNotAllowedContentTypeError(c.Err))
		})
	}
}

func TestIsNotAllowedAssertionMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name  string
		Err   error
		IsErr bool
	}{
		{
			Name:  "not_allowed_assertion_method_error_is_not_allowed_assertion_method_error",
			Err:   specification.NewNotAllowedAssertionMethodError("jzonpad"),
			IsErr: true,
		},
		{
			Name:  "another_error_isnt_not_allowed_assertion_method_error",
			Err:   specification.NewNotAllowedKeywordError("jzonpad"),
			IsErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsErr, specification.IsNotAllowedAssertionMethodError(c.Err))
		})
	}
}
