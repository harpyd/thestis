package specification_test

import (
	"errors"
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
			Name:            "zero_slug",
			GivenSlug:       specification.Slug{},
			ShouldPanic:     true,
			WithExpectedErr: specification.ErrZeroSlug,
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

			builder := specification.NewStoryBuilder()

			var story specification.Story

			buildFn := func() {
				story = builder.ErrlessBuild(c.GivenSlug)
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

func errlessBuildStory(
	t *testing.T,
	slug specification.Slug,
	prepare func(b *specification.StoryBuilder),
) specification.Story {
	t.Helper()

	builder := specification.NewStoryBuilder()

	prepare(builder)

	return builder.ErrlessBuild(slug)
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

			var (
				slug        = specification.NewStorySlug("foo")
				description = errlessBuildStory(t, slug, c.Prepare).Description()
			)

			require.Equal(t, c.ExpectedDescription, description)
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

			var (
				slug = specification.NewStorySlug("foo")
				asA  = errlessBuildStory(t, slug, c.Prepare).AsA()
			)

			require.Equal(t, c.ExpectedAsA, asA)
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

			var (
				slug      = specification.NewStorySlug("foo")
				inOrderTo = errlessBuildStory(t, slug, c.Prepare).InOrderTo()
			)

			require.Equal(t, c.ExpectedInOrderTo, inOrderTo)
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

			var (
				slug   = specification.NewStorySlug("foo")
				wantTo = errlessBuildStory(t, slug, c.Prepare).WantTo()
			)

			require.Equal(t, c.ExpectedWantTo, wantTo)
		})
	}
}

func TestBuildStoryWithScenarios(t *testing.T) {
	t.Parallel()

	storySlug := specification.NewStorySlug("foo")

	testCases := []struct {
		Name              string
		Prepare           func(b *specification.StoryBuilder)
		ExpectedScenarios []specification.Scenario
		WantThisErr       bool
		IsErr             func(err error) bool
	}{
		{
			Name:        "no_scenarios",
			Prepare:     func(b *specification.StoryBuilder) {},
			WantThisErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoStoryScenarios)
			},
		},
		{
			Name: "two_scenarios",
			Prepare: func(b *specification.StoryBuilder) {
				b.WithScenario("bar", func(b *specification.ScenarioBuilder) {})
				b.WithScenario("baz", func(b *specification.ScenarioBuilder) {})
			},
			ExpectedScenarios: []specification.Scenario{
				specification.NewScenarioBuilder().
					ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
				specification.NewScenarioBuilder().
					ErrlessBuild(specification.NewScenarioSlug("foo", "baz")),
			},
			WantThisErr: false,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoStoryScenarios)
			},
		},
		{
			Name: "scenario_already_exists",
			Prepare: func(b *specification.StoryBuilder) {
				b.WithScenario("wow", func(b *specification.ScenarioBuilder) {})
				b.WithScenario("wow", func(b *specification.ScenarioBuilder) {})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				var target *specification.DuplicatedError

				return errors.As(err, &target)
			},
		},
		{
			Name: "no_scenario_theses",
			Prepare: func(b *specification.StoryBuilder) {
				b.WithScenario("no", func(b *specification.ScenarioBuilder) {})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoScenarioTheses)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			builder := specification.NewStoryBuilder()

			c.Prepare(builder)

			story, err := builder.Build(storySlug)

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))

			require.ElementsMatch(t, c.ExpectedScenarios, story.Scenarios())
		})
	}
}

func TestGetStoryScenarioBySlug(t *testing.T) {
	t.Parallel()

	var (
		slug  = specification.NewStorySlug("foo")
		story = errlessBuildStory(t, slug, func(b *specification.StoryBuilder) {
			b.WithScenario("bad", func(b *specification.ScenarioBuilder) {})
			b.WithScenario("baz", func(b *specification.ScenarioBuilder) {})
		})
	)

	bad, ok := story.Scenario("bad")
	require.True(t, ok)
	require.Equal(
		t,
		specification.NewScenarioBuilder().ErrlessBuild(
			specification.NewScenarioSlug("foo", "bad"),
		),
		bad,
	)

	baz, ok := story.Scenario("baz")
	require.True(t, ok)
	require.Equal(
		t,
		specification.NewScenarioBuilder().ErrlessBuild(
			specification.NewScenarioSlug("foo", "baz"),
		),
		baz,
	)

	_, ok = story.Scenario("bak")
	require.False(t, ok)
}
