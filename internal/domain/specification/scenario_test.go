package specification_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuildScenarioSlugging(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		GivenSlug       specification.Slug
		ShouldPanic     bool
		WithExpectedErr error
	}{
		{
			Name:        "foo.bar",
			GivenSlug:   specification.NewScenarioSlug("foo", "bar"),
			ShouldPanic: false,
		},
		{
			Name:            "not_scenario_slug",
			GivenSlug:       specification.NewStorySlug("foo"),
			ShouldPanic:     true,
			WithExpectedErr: specification.ErrNotScenarioSlug,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var b specification.ScenarioBuilder

			var scenario specification.Scenario

			buildFn := func() {
				scenario = b.ErrlessBuild(c.GivenSlug)
			}

			if c.ShouldPanic {
				require.PanicsWithValue(t, c.WithExpectedErr, buildFn)

				return
			}

			require.NotPanics(t, buildFn)

			require.Equal(t, c.GivenSlug, scenario.Slug())
		})
	}
}

func errlessBuildScenario(
	t *testing.T,
	slug specification.Slug,
	prepare func(b *specification.ScenarioBuilder),
) specification.Scenario {
	t.Helper()

	var b specification.ScenarioBuilder

	prepare(&b)

	return b.ErrlessBuild(slug)
}

func TestBuildScenarioWithDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare             func(b *specification.ScenarioBuilder)
		ExpectedDescription string
	}{
		{
			Prepare:             func(b *specification.ScenarioBuilder) {},
			ExpectedDescription: "",
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithDescription("description")
			},
			ExpectedDescription: "description",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var (
				slug        = specification.NewScenarioSlug("foo", "bar")
				description = errlessBuildScenario(t, slug, c.Prepare).Description()
			)

			require.Equal(t, c.ExpectedDescription, description)
		})
	}
}

func TestBuildScenarioWithTheses(t *testing.T) {
	t.Parallel()

	scenarioSlug := specification.NewScenarioSlug("foo", "bar")

	testCases := []struct {
		Name           string
		Prepare        func(b *specification.ScenarioBuilder)
		ExpectedTheses []specification.Thesis
		WantThisErr    bool
		IsErr          func(err error) bool
	}{
		{
			Name:        "no_theses",
			Prepare:     func(b *specification.ScenarioBuilder) {},
			WantThisErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoScenarioTheses)
			},
		},
		{
			Name: "two_theses",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
				b.WithThesis("bad", func(b *specification.ThesisBuilder) {})
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).ErrlessBuild(
					specification.NewThesisSlug(scenarioSlug.Story(), scenarioSlug.Scenario(), "baz"),
				),
				(&specification.ThesisBuilder{}).ErrlessBuild(
					specification.NewThesisSlug(scenarioSlug.Story(), scenarioSlug.Scenario(), "bad"),
				),
			},
			WantThisErr: false,
			IsErr: func(err error) bool {
				return errors.Is(err, specification.ErrNoScenarioTheses)
			},
		},
		{
			Name: "thesis_already_exists",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				var target *specification.DuplicatedError

				return errors.As(err, &target)
			},
		},
		{
			Name: "thesis_has_invalid_dependencies",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {
					b.WithDependency("non-existent")
				})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				var target *specification.InvalidDependenciesError

				return errors.As(err, &target)
			},
		},
		{
			Name: "thesis_has_cyclic_dependencies",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {
					b.WithDependency("baz")
				})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				var target *specification.CyclicDependencyError

				return errors.As(err, &target)
			},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var b specification.ScenarioBuilder

			c.Prepare(&b)

			scenario, err := b.Build(scenarioSlug)

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))

			require.ElementsMatch(t, c.ExpectedTheses, scenario.Theses())
		})
	}
}

func TestGetScenarioThesesByStages(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare        func(b *specification.ScenarioBuilder)
		GivenStages    []specification.Stage
		ExpectedTheses []specification.Thesis
	}{
		{
			Prepare: func(b *specification.ScenarioBuilder) {},
			GivenStages: []specification.Stage{
				specification.Given,
			},
			ExpectedTheses: nil,
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {})
			},
			GivenStages: []specification.Stage{
				specification.Given,
				specification.When,
				specification.Then,
			},
			ExpectedTheses: nil,
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "a")
				})
			},
			GivenStages: []specification.Stage{
				specification.Given,
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).
					WithStatement(specification.Given, "a").
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "a")),
			},
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "a")
				})
			},
			GivenStages: []specification.Stage{
				specification.Given,
				specification.When,
				specification.Then,
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).
					WithStatement(specification.Given, "a").
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "a")),
			},
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "a")
				})
				b.WithThesis("b", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.When, "b")
				})
				b.WithThesis("c", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Then, "c")
				})
			},
			GivenStages: []specification.Stage{
				specification.Given,
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).
					WithStatement(specification.Given, "a").
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "a")),
			},
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "a")
				})
				b.WithThesis("b", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.When, "b")
				})
				b.WithThesis("c", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Then, "c")
				})
			},
			GivenStages: []specification.Stage{
				specification.When,
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).
					WithStatement(specification.When, "b").
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "b")),
			},
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "a")
				})
				b.WithThesis("b", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.When, "b")
				})
				b.WithThesis("c", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Then, "c")
				})
			},
			GivenStages: []specification.Stage{
				specification.Then,
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).
					WithStatement(specification.Then, "c").
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "c")),
			},
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "a")
				})
				b.WithThesis("b", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.When, "b")
				})
				b.WithThesis("c", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Then, "c")
				})
			},
			GivenStages: []specification.Stage{
				specification.Given,
				specification.When,
				specification.Then,
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).
					WithStatement(specification.Given, "a").
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "a")),
				(&specification.ThesisBuilder{}).
					WithStatement(specification.When, "b").
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "b")),
				(&specification.ThesisBuilder{}).
					WithStatement(specification.Then, "c").
					ErrlessBuild(specification.NewThesisSlug("foo", "bar", "c")),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var (
				slug     = specification.NewScenarioSlug("foo", "bar")
				scenario = errlessBuildScenario(t, slug, c.Prepare)
			)

			require.ElementsMatch(t, c.ExpectedTheses, scenario.ThesesByStages(c.GivenStages...))
		})
	}
}

func TestGetScenarioThesisBySlug(t *testing.T) {
	t.Parallel()

	var (
		slug     = specification.NewScenarioSlug("aaa", "bb")
		scenario = errlessBuildScenario(t, slug, func(b *specification.ScenarioBuilder) {
			b.WithThesis("c", func(b *specification.ThesisBuilder) {})
			b.WithThesis("d", func(b *specification.ThesisBuilder) {})
		})
	)

	var b specification.ThesisBuilder

	c, ok := scenario.Thesis("c")
	require.True(t, ok)
	require.Equal(
		t,
		b.ErrlessBuild(
			specification.NewThesisSlug("aaa", "bb", "c"),
		),
		c,
	)

	b.Reset()

	d, ok := scenario.Thesis("d")
	require.True(t, ok)
	require.Equal(
		t,
		b.ErrlessBuild(
			specification.NewThesisSlug("aaa", "bb", "d"),
		),
		d,
	)

	_, ok = scenario.Thesis("f")
	require.False(t, ok)
}
