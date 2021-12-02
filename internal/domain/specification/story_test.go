package specification_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestStoryBuilder_Build_slug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Slug        string
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_slug",
			Slug:        "story",
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_empty_slug",
			Slug:        "",
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewStoryBuilder()

			story, err := builder.Build(c.Slug)

			if c.ShouldBeErr {
				require.True(t, specification.IsStoryEmptySlugError(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, c.Slug, story.Slug())
		})
	}
}

func TestStoryBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithDescription("description")

	story, err := builder.Build("someStory")

	require.NoError(t, err)
	require.Equal(t, "description", story.Description())
}

func TestStoryBuilder_WithAsA(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithAsA("author")

	story, err := builder.Build("someStory")

	require.NoError(t, err)
	require.Equal(t, "author", story.AsA())
}

func TestStoryBuilder_WithInOrderTo(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithInOrderTo("to do something")

	story, err := builder.Build("someStory")

	require.NoError(t, err)
	require.Equal(t, "to do something", story.InOrderTo())
}

func TestStoryBuilder_WithWantTo(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithWantTo("do work")

	story, err := builder.Build("someStory")

	require.NoError(t, err)
	require.Equal(t, "do work", story.WantTo())
}

func TestStoryBuilder_WithScenario(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithScenario("firstScenario", func(b *specification.ScenarioBuilder) {
		b.WithDescription("this is a first scenario")
	})
	builder.WithScenario("secondScenario", func(b *specification.ScenarioBuilder) {
		b.WithDescription("this is a second scenario")
	})

	story, err := builder.Build("someStory")

	require.NoError(t, err)

	expectedFirstScenario, err := specification.NewScenarioBuilder().
		WithDescription("this is a first scenario").
		Build("firstScenario")
	require.NoError(t, err)

	actualFirstScenario, ok := story.Scenario("firstScenario")
	require.True(t, ok)
	require.Equal(t, expectedFirstScenario, actualFirstScenario)

	expectedSecondScenario, err := specification.NewScenarioBuilder().
		WithDescription("this is a second scenario").
		Build("secondScenario")
	require.NoError(t, err)

	actualSecondScenario, ok := story.Scenario("secondScenario")
	require.True(t, ok)
	require.Equal(t, expectedSecondScenario, actualSecondScenario)
}

func TestStoryBuilder_WithScenario_when_already_exists(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithScenario("scenario", func(b *specification.ScenarioBuilder) {
		b.WithDescription("this is a scenario")
	})
	builder.WithScenario("scenario", func(b *specification.ScenarioBuilder) {
		b.WithDescription("this is a same scenario")
	})

	_, err := builder.Build("story")

	require.True(t, specification.IsScenarioSlugAlreadyExistsError(err))
}

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