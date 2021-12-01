package specification_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestIsStorySlugAlreadyExistsError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "story_slug_already_exists_error_is_story_slug_already_exists_error",
			Err:       specification.NewStorySlugAlreadyExistsError("story"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_story_slug_already_exists_error",
			Err:       errors.New("story"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsStorySlugAlreadyExistsError(c.Err))
		})
	}
}

func TestIsScenarioSlugAlreadyExistsError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "scenario_slug_already_exists_error_is_scenario_slug_already_exists_error",
			Err:       specification.NewScenarioSlugAlreadyExistsError("scenario"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_scenario_slug_already_exists_error",
			Err:       specification.NewThesisSlugAlreadyExistsError("thesis"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsScenarioSlugAlreadyExistsError(c.Err))
		})
	}
}

func TestIsThesisSlugAlreadyExistsError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "thesis_slug_already_exists_error_is_thesis_slug_already_exists_error",
			Err:       specification.NewThesisSlugAlreadyExistsError("thesis"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_thesis_slug_already_exists_error",
			Err:       errors.New("thesis"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsThesisSlugAlreadyExistsError(c.Err))
		})
	}
}

func TestIsStoryEmptySlugError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "story_empty_slug_error_is_story_empty_slug_error",
			Err:       specification.NewStoryEmptySlugError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_story_empty_slug_error",
			Err:       errors.New("something wrong"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsStoryEmptySlugError(c.Err))
		})
	}
}

func TestIsScenarioEmptySlugError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "scenario_empty_slug_error_is_scenario_empty_slug_error",
			Err:       specification.NewScenarioEmptySlugError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_scenario_empty_slug_error",
			Err:       errors.New("error"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsScenarioEmptySlugError(c.Err))
		})
	}
}

func TestIsThesisEmptySlugError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "thesis_empty_slug_error_is_empty_slug_error",
			Err:       specification.NewThesisEmptySlugError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_thesis_empty_slug_error",
			Err:       errors.New("wrong wrong"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsThesisEmptySlugError(c.Err))
		})
	}
}

func TestIsBuildSpecificationError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "specification_error_is_specification_error",
			Err:       specification.NewBuildSpecificationError(errors.New("badaboom")),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_specification_error",
			Err:       specification.NewNoSuchStoryError("slug"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildSpecificationError(c.Err))
		})
	}
}

func TestIsBuildStoryError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_story_error_is_build_story_error",
			Err:       specification.NewBuildStoryError(errors.New("boom"), "story"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_build_story_error",
			Err:       specification.NewBuildScenarioError(errors.New("boom"), "scenario"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildStoryError(c.Err))
		})
	}
}

func TestIsBuildScenarioError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_scenario_error_is_build_scenario_error",
			Err:       specification.NewBuildScenarioError(errors.New("wrong"), "scenario"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_build_scenario_error",
			Err:       errors.New("wrong"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildScenarioError(c.Err))
		})
	}
}

func TestIsBuildThesisError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_thesis_error_is_build_thesis_error",
			Err:       specification.NewBuildThesisError(errors.New("pew"), "thesis"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_build_thesis_error",
			Err:       errors.New("pew"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsBuildThesisError(c.Err))
		})
	}
}

func TestIsNoSuchStoryError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "no_story_error_is_no_story_error",
			Err:       specification.NewNoSuchStoryError("someStory"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_no_story_error",
			Err:       specification.NewNoSuchThesisError("someThesis"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNoSuchStoryError(c.Err))
		})
	}
}

func TestIsNoSuchScenarioError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "no_scenario_error_is_no_scenario_error",
			Err:       specification.NewNoSuchScenarioError("someScenario"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_no_scenario_error",
			Err:       specification.NewNoSuchThesisError("someThesis"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNoSuchScenarioError(c.Err))
		})
	}
}

func TestIsNoThesisError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "no_thesis_error_is_no_thesis_error",
			Err:       specification.NewNoSuchThesisError("someThesis"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_no_thesis_error",
			Err:       specification.NewNoSuchStoryError("someStory"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNoSuchThesisError(c.Err))
		})
	}
}

func TestIsNotAllowedHTTPMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_http_method_error_is_not_allowed_http_method_error",
			Err:       specification.NewNotAllowedHTTPMethodError("POZT"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_not_allowed_http_method_error",
			Err:       specification.NewNoSuchThesisError("POZT"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedHTTPMethodError(c.Err))
		})
	}
}

func TestIsNotAllowedKeywordError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_keyword_error_is_not_allowed_keyword_error",
			Err:       specification.NewNotAllowedKeywordError("zen"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_not_allowed_keyword_error",
			Err:       specification.NewNotAllowedHTTPMethodError("zen"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedKeywordError(c.Err))
		})
	}
}

func TestIsNotAllowedContentTypeError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_content_type_error_is_not_allowed_content_type_error",
			Err:       specification.NewNotAllowedContentTypeError("some/content"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_not_allowed_content_type_error",
			Err:       specification.NewNoSuchStoryError("some/content"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedContentTypeError(c.Err))
		})
	}
}

func TestIsNotAllowedAssertionMethodError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "not_allowed_assertion_method_error_is_not_allowed_assertion_method_error",
			Err:       specification.NewNotAllowedAssertionMethodError("jzonpad"),
			IsSameErr: true,
		},
		{
			Name:      "another_error_isnt_not_allowed_assertion_method_error",
			Err:       specification.NewNotAllowedKeywordError("jzonpad"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNotAllowedAssertionMethodError(c.Err))
		})
	}
}
