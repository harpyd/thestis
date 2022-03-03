package specification_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuildStorySlugging(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		GivenSlug   specification.Slug
		WantThisErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "foo",
			GivenSlug:   specification.NewStorySlug("foo"),
			WantThisErr: false,
			IsErr: func(err error) bool {
				return specification.IsEmptySlugError(err) ||
					specification.IsNotStorySlugError(err)
			},
		},
		{
			Name:        "empty_slug",
			GivenSlug:   specification.Slug{},
			WantThisErr: true,
			IsErr:       specification.IsEmptySlugError,
		},
		{
			Name:        "not_story_slug",
			GivenSlug:   specification.NewScenarioSlug("foo", "bar"),
			WantThisErr: true,
			IsErr:       specification.IsNotStorySlugError,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			builder := specification.NewStoryBuilder()

			story, err := builder.Build(c.GivenSlug)

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))

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
			IsErr:       specification.IsNoStoryScenariosError,
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
			IsErr:       specification.IsNoStoryScenariosError,
		},
		{
			Name: "scenario_already_exists",
			Prepare: func(b *specification.StoryBuilder) {
				b.WithScenario("wow", func(b *specification.ScenarioBuilder) {})
				b.WithScenario("wow", func(b *specification.ScenarioBuilder) {})
			},
			WantThisErr: true,
			IsErr:       specification.IsScenarioSlugAlreadyExistsError,
		},
		{
			Name: "no_scenario_theses",
			Prepare: func(b *specification.StoryBuilder) {
				b.WithScenario("no", func(b *specification.ScenarioBuilder) {})
			},
			WantThisErr: true,
			IsErr:       specification.IsNoScenarioThesesError,
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
