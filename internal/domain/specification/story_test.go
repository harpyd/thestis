package specification_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuildStorySlugging(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		GivenSlug       specification.Slug
		ShouldPanic     bool
		WithExpectedErr error
	}{
		{
			Name:        "foo",
			GivenSlug:   specification.NewStorySlug("foo"),
			ShouldPanic: false,
		},
		{
			Name:            "not_story_slug",
			GivenSlug:       specification.NewScenarioSlug("foo", "bar"),
			ShouldPanic:     true,
			WithExpectedErr: specification.ErrNotStorySlug,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var b specification.StoryBuilder

			var story specification.Story

			buildFn := func() {
				story = b.Build(c.GivenSlug)
			}

			if c.ShouldPanic {
				require.PanicsWithValue(t, c.WithExpectedErr, buildFn)

				return
			}

			require.NotPanics(t, buildFn)

			require.Equal(t, c.GivenSlug, story.Slug())
		})
	}
}

func buildStory(
	t *testing.T,
	slug specification.Slug,
	prepare func(b *specification.StoryBuilder),
) specification.Story {
	t.Helper()

	var b specification.StoryBuilder

	prepare(&b)

	return b.Build(slug)
}

func TestBuildStoryWithDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare             func(b *specification.StoryBuilder)
		ExpectedDescription string
	}{
		{
			Prepare:             func(b *specification.StoryBuilder) {},
			ExpectedDescription: "",
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithDescription("")
			},
			ExpectedDescription: "",
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithDescription("desc")
			},
			ExpectedDescription: "desc",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			slug := specification.NewStorySlug("foo")

			actualDescription := buildStory(t, slug, c.Prepare).Description()

			require.Equal(t, c.ExpectedDescription, actualDescription)
		})
	}
}

func TestBuildStoryWithAsA(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare     func(b *specification.StoryBuilder)
		ExpectedAsA string
	}{
		{
			Prepare:     func(b *specification.StoryBuilder) {},
			ExpectedAsA: "",
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithAsA("")
			},
			ExpectedAsA: "",
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithAsA("boo")
			},
			ExpectedAsA: "boo",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			slug := specification.NewStorySlug("bar")

			actualAsA := buildStory(t, slug, c.Prepare).AsA()

			require.Equal(t, c.ExpectedAsA, actualAsA)
		})
	}
}

func TestBuildStoryWithInOrderTo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare           func(b *specification.StoryBuilder)
		ExpectedInOrderTo string
	}{
		{
			Prepare:           func(b *specification.StoryBuilder) {},
			ExpectedInOrderTo: "",
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithInOrderTo("")
			},
			ExpectedInOrderTo: "",
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithInOrderTo("ord")
			},
			ExpectedInOrderTo: "ord",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			slug := specification.NewStorySlug("aaa")

			actualInOrderTo := buildStory(t, slug, c.Prepare).InOrderTo()

			require.Equal(t, c.ExpectedInOrderTo, actualInOrderTo)
		})
	}
}

func TestBuildStoryWithWantTo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare        func(b *specification.StoryBuilder)
		ExpectedWantTo string
	}{
		{
			Prepare:        func(b *specification.StoryBuilder) {},
			ExpectedWantTo: "",
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithWantTo("")
			},
			ExpectedWantTo: "",
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithWantTo("wanna")
			},
			ExpectedWantTo: "wanna",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			slug := specification.NewStorySlug("buf")

			actualWantTo := buildStory(t, slug, c.Prepare).WantTo()

			require.Equal(t, c.ExpectedWantTo, actualWantTo)
		})
	}
}

func TestBuildStoryWithScenarios(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare           func(b *specification.StoryBuilder)
		ExpectedScenarios []specification.Scenario
	}{
		{
			Prepare:           func(b *specification.StoryBuilder) {},
			ExpectedScenarios: nil,
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithScenario("bar", func(b *specification.ScenarioBuilder) {})
				b.WithScenario("baz", func(b *specification.ScenarioBuilder) {})
			},
			ExpectedScenarios: []specification.Scenario{
				(&specification.ScenarioBuilder{}).
					Build(specification.NewScenarioSlug("a", "bar")),
				(&specification.ScenarioBuilder{}).
					Build(specification.NewScenarioSlug("a", "baz")),
			},
		},
		{
			Prepare: func(b *specification.StoryBuilder) {
				b.WithScenario("wow", func(b *specification.ScenarioBuilder) {})
				b.WithScenario("wow", func(b *specification.ScenarioBuilder) {})
			},
			ExpectedScenarios: []specification.Scenario{
				(&specification.ScenarioBuilder{}).
					Build(specification.NewScenarioSlug("a", "wow")),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			slug := specification.NewStorySlug("a")

			actualScenarios := buildStory(t, slug, c.Prepare).Scenarios()

			require.ElementsMatch(t, c.ExpectedScenarios, actualScenarios)
		})
	}
}

func TestGetStoryScenarioBySlug(t *testing.T) {
	t.Parallel()

	slug := specification.NewStorySlug("foo")

	story := buildStory(t, slug, func(b *specification.StoryBuilder) {
		b.WithScenario("bad", func(b *specification.ScenarioBuilder) {})
		b.WithScenario("baz", func(b *specification.ScenarioBuilder) {})
	})

	var b specification.ScenarioBuilder

	bad, ok := story.Scenario("bad")
	require.True(t, ok)
	require.Equal(
		t,
		b.Build(
			specification.NewScenarioSlug("foo", "bad"),
		),
		bad,
	)

	b.Reset()

	baz, ok := story.Scenario("baz")
	require.True(t, ok)
	require.Equal(
		t,
		b.Build(
			specification.NewScenarioSlug("foo", "baz"),
		),
		baz,
	)

	_, ok = story.Scenario("bak")
	require.False(t, ok)
}
