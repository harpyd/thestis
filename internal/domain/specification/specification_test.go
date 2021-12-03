package specification_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuilder_WithAuthor(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithAuthor("author")

	spec := builder.ErrlessBuild()

	require.Equal(t, "author", spec.Author())
}

func TestBuilder_WithTitle(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithTitle("specification")

	spec := builder.ErrlessBuild()

	require.Equal(t, "specification", spec.Title())
}

func TestBuilder_WithDescription(t *testing.T) {
	t.Parallel()

	builder := specification.NewBuilder()
	builder.WithDescription("description")

	spec := builder.ErrlessBuild()

	require.Equal(t, "description", spec.Description())
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
		ErrlessBuild("firstStory")

	actualFirstStory, ok := spec.Story("firstStory")
	require.True(t, ok)
	require.Equal(t, expectedFirstStory, actualFirstStory)

	expectedSecondStory := specification.NewStoryBuilder().
		WithDescription("this is a second story").
		ErrlessBuild("secondStory")

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

func TestIsBuildSpecificationError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "build_specification_error",
			Err:       specification.NewBuildSpecificationError(errors.New("badaboom")),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
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

func TestIsNoSpecificationStoriesError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "no_specification_stories_error",
			Err:       specification.NewNoSpecificationStoriesError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("another"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, specification.IsNoSpecificationStoriesError(c.Err))
		})
	}
}
