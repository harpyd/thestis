package specification_test

import (
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
				scenario = b.Build(c.GivenSlug)
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

func buildScenario(
	t *testing.T,
	slug specification.Slug,
	prepare func(b *specification.ScenarioBuilder),
) specification.Scenario {
	t.Helper()

	var b specification.ScenarioBuilder

	prepare(&b)

	return b.Build(slug)
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

			slug := specification.NewScenarioSlug("foo", "bar")

			actualDescription := buildScenario(t, slug, c.Prepare).Description()

			require.Equal(t, c.ExpectedDescription, actualDescription)
		})
	}
}

func TestBuildScenarioWithTheses(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Prepare        func(b *specification.ScenarioBuilder)
		ExpectedTheses []specification.Thesis
	}{
		{
			Prepare:        func(b *specification.ScenarioBuilder) {},
			ExpectedTheses: nil,
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
				b.WithThesis("bad", func(b *specification.ThesisBuilder) {})
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).Build(
					specification.NewThesisSlug("a", "b", "baz"),
				),
				(&specification.ThesisBuilder{}).Build(
					specification.NewThesisSlug("a", "b", "bad"),
				),
			},
		},
		{
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
			},
			ExpectedTheses: []specification.Thesis{
				(&specification.ThesisBuilder{}).Build(
					specification.NewThesisSlug("a", "b", "baz"),
				),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			slug := specification.NewScenarioSlug("a", "b")

			actualTheses := buildScenario(t, slug, c.Prepare).Theses()

			require.ElementsMatch(t, c.ExpectedTheses, actualTheses)
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
					Build(specification.NewThesisSlug("a", "bb", "a")),
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
					Build(specification.NewThesisSlug("a", "bb", "a")),
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
					Build(specification.NewThesisSlug("a", "bb", "a")),
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
					Build(specification.NewThesisSlug("a", "bb", "b")),
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
					Build(specification.NewThesisSlug("a", "bb", "c")),
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
					Build(specification.NewThesisSlug("a", "bb", "a")),
				(&specification.ThesisBuilder{}).
					WithStatement(specification.When, "b").
					Build(specification.NewThesisSlug("a", "bb", "b")),
				(&specification.ThesisBuilder{}).
					WithStatement(specification.Then, "c").
					Build(specification.NewThesisSlug("a", "bb", "c")),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			slug := specification.NewScenarioSlug("a", "bb")

			actualTheses := buildScenario(t, slug, c.Prepare).ThesesByStages(c.GivenStages...)

			require.ElementsMatch(t, c.ExpectedTheses, actualTheses)
		})
	}
}

func TestGetScenarioThesisBySlug(t *testing.T) {
	t.Parallel()

	slug := specification.NewScenarioSlug("a", "b")

	scenario := buildScenario(t, slug, func(b *specification.ScenarioBuilder) {
		b.WithThesis("c", func(b *specification.ThesisBuilder) {})
		b.WithThesis("d", func(b *specification.ThesisBuilder) {})
	})

	var b specification.ThesisBuilder

	c, ok := scenario.Thesis("c")
	require.True(t, ok)
	require.Equal(
		t,
		b.Build(
			specification.NewThesisSlug("a", "b", "c"),
		),
		c,
	)

	b.Reset()

	d, ok := scenario.Thesis("d")
	require.True(t, ok)
	require.Equal(
		t,
		b.Build(
			specification.NewThesisSlug("a", "b", "d"),
		),
		d,
	)

	_, ok = scenario.Thesis("f")
	require.False(t, ok)
}

func TestCheckCycleInThesesDependencies(t *testing.T) {
	t.Parallel()

	scenarioSlug := specification.NewScenarioSlug("foo", "bar")

	testCases := []struct {
		Name        string
		Prepare     func(b *specification.ScenarioBuilder)
		WantThisErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "thesis_has_cycle_dependencies",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("t1", func(b *specification.ThesisBuilder) {
					b.WithDependency("t2")
				})
				b.WithThesis("t2", func(b *specification.ThesisBuilder) {
					b.WithDependency("t3")
				})
				b.WithThesis("t3", func(b *specification.ThesisBuilder) {
					b.WithDependency("t1")
				})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				var target *specification.CycleDependenciesError

				return errors.As(err, &target)
			},
		},
		{
			Name: "thesis_without_cycle_dependencies",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("t1", func(b *specification.ThesisBuilder) {
					b.WithDependency("t2")
				})
				b.WithThesis("t2", func(b *specification.ThesisBuilder) {
					b.WithDependency("t3")
				})
			},
			WantThisErr: false,
			IsErr: func(err error) bool {
				var target *specification.CycleDependenciesError

				return errors.As(err, &target)
			},
		},
		{
			Name: "thesis_void_dependencies",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("t1", func(b *specification.ThesisBuilder) {})
				b.WithThesis("t2", func(b *specification.ThesisBuilder) {})
			},
			WantThisErr: false,
			IsErr: func(err error) bool {
				var target *specification.CycleDependenciesError

				return errors.As(err, &target)
			},
		},
		{
			Name:        "void_thesis",
			Prepare:     func(b *specification.ScenarioBuilder) {},
			WantThisErr: false,
			IsErr: func(err error) bool {
				var target *specification.CycleDependenciesError

				return errors.As(err, &target)
			},
		},

		{
			Name: "thesis_has_hard_cycle_dependencies",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("t1", func(b *specification.ThesisBuilder) {
					b.WithDependency("t2")
					b.WithDependency("t3")
					b.WithDependency("t4")
				})
				b.WithThesis("t2", func(b *specification.ThesisBuilder) {
					b.WithDependency("t5")
				})
				b.WithThesis("t3", func(b *specification.ThesisBuilder) {
					b.WithDependency("t2")
				})
				b.WithThesis("t4", func(b *specification.ThesisBuilder) {})
				b.WithThesis("t5", func(b *specification.ThesisBuilder) {
					b.WithDependency("t3")
				})
			},
			WantThisErr: true,
			IsErr: func(err error) bool {
				var target *specification.CycleDependenciesError

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

			_, err := b.Build(scenarioSlug)

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))
		})
	}
}
