package specification_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestStoryBuilder_Build_no_scenarios_error(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()

	_, err := builder.Build(specification.NewStorySlug("story"))

	require.True(t, specification.IsNoStoryScenariosError(err))
}

func TestStoryBuilder_Build_slug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Slug        specification.Slug
		ShouldBeErr bool
	}{
		{
			Name:        "build_with_slug",
			Slug:        specification.NewStorySlug("story"),
			ShouldBeErr: false,
		},
		{
			Name:        "dont_build_with_empty_slug",
			Slug:        specification.Slug{},
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewStoryBuilder()

			if c.ShouldBeErr {
				_, err := builder.Build(c.Slug)
				require.True(t, specification.IsEmptySlugError(err))

				return
			}

			story := builder.ErrlessBuild(c.Slug)
			require.Equal(t, c.Slug, story.Slug())
		})
	}
}

func TestStoryBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithDescription("description")

	story := builder.ErrlessBuild(specification.NewStorySlug("someStory"))

	require.Equal(t, "description", story.Description())
}

func TestStoryBuilder_WithAsA(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithAsA("author")

	story := builder.ErrlessBuild(specification.NewStorySlug("story"))

	require.Equal(t, "author", story.AsA())
}

func TestStoryBuilder_WithInOrderTo(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithInOrderTo("to do something")

	story := builder.ErrlessBuild(specification.NewStorySlug("some"))

	require.Equal(t, "to do something", story.InOrderTo())
}

func TestStoryBuilder_WithWantTo(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithWantTo("do work")

	story := builder.ErrlessBuild(specification.NewStorySlug("what"))

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

	story := builder.ErrlessBuild(specification.NewStorySlug("are"))

	expectedFirstScenario := specification.NewScenarioBuilder().
		WithDescription("this is a first scenario").
		ErrlessBuild(specification.NewScenarioSlug("are", "firstScenario"))

	actualFirstScenario, ok := story.Scenario("firstScenario")
	require.True(t, ok)
	require.Equal(t, expectedFirstScenario, actualFirstScenario)

	expectedSecondScenario := specification.NewScenarioBuilder().
		WithDescription("this is a second scenario").
		ErrlessBuild(specification.NewScenarioSlug("are", "secondScenario"))

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

	_, err := builder.Build(specification.NewStorySlug("you"))

	require.True(t, specification.IsScenarioSlugAlreadyExistsError(err))
}

func TestStory_Scenarios(t *testing.T) {
	t.Parallel()

	builder := specification.NewStoryBuilder()
	builder.WithScenario("foo", func(b *specification.ScenarioBuilder) {})
	builder.WithScenario("bar", func(b *specification.ScenarioBuilder) {})

	story := builder.ErrlessBuild(specification.NewStorySlug("baz"))

	expected := []specification.Scenario{
		specification.NewScenarioBuilder().ErrlessBuild(
			specification.NewScenarioSlug("baz", "foo"),
		),
		specification.NewScenarioBuilder().ErrlessBuild(
			specification.NewScenarioSlug("baz", "bar"),
		),
	}

	assert.ElementsMatch(t, expected, story.Scenarios())
}

func TestStoryErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "no_story_scenarios_error",
			Err:   specification.NewNoStoryScenariosError(),
			IsErr: specification.IsNoStoryScenariosError,
		},
		{
			Name:     "NON_no_story_scenarios_error",
			Err:      errors.New("no story scenarios"),
			IsErr:    specification.IsNoStoryScenariosError,
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
