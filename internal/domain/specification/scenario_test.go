package specification_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestBuildScenarioSlugging(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		GivenSlug   specification.Slug
		WantThisErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:        "foo.bar",
			GivenSlug:   specification.NewScenarioSlug("foo", "bar"),
			WantThisErr: false,
			IsErr: func(err error) bool {
				return specification.IsEmptySlugError(err) ||
					specification.IsNotScenarioSlugError(err)
			},
		},
		{
			Name:        "empty_slug",
			GivenSlug:   specification.Slug{},
			WantThisErr: true,
			IsErr:       specification.IsEmptySlugError,
		},
		{
			Name:        "not_scenario_slug",
			GivenSlug:   specification.NewStorySlug("foo"),
			WantThisErr: true,
			IsErr:       specification.IsNotScenarioSlugError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewScenarioBuilder()

			scenario, err := builder.Build(c.GivenSlug)

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))

			require.Equal(t, c.GivenSlug, scenario.Slug())
		})
	}
}

func errlessBuildScenario(
	t *testing.T,
	prepare func(b *specification.ScenarioBuilder),
) specification.Scenario {
	t.Helper()

	builder := specification.NewScenarioBuilder()

	prepare(builder)

	return builder.ErrlessBuild(
		specification.NewScenarioSlug("story", "scenario"),
	)
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

			description := errlessBuildScenario(t, c.Prepare).Description()

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
			IsErr:       specification.IsNoScenarioThesesError,
		},
		{
			Name: "two_theses",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
				b.WithThesis("bad", func(b *specification.ThesisBuilder) {})
			},
			ExpectedTheses: []specification.Thesis{
				specification.NewThesisBuilder().ErrlessBuild(
					specification.NewThesisSlug(scenarioSlug.Story(), scenarioSlug.Scenario(), "baz"),
				),
				specification.NewThesisBuilder().ErrlessBuild(
					specification.NewThesisSlug(scenarioSlug.Story(), scenarioSlug.Scenario(), "bad"),
				),
			},
			WantThisErr: false,
			IsErr:       specification.IsNoScenarioThesesError,
		},
		{
			Name: "thesis_already_exists",
			Prepare: func(b *specification.ScenarioBuilder) {
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
				b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
			},
			WantThisErr: true,
			IsErr:       specification.IsThesisSlugAlreadyExistsError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			builder := specification.NewScenarioBuilder()

			c.Prepare(builder)

			scenario, err := builder.Build(scenarioSlug)

			if c.WantThisErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.False(t, c.IsErr(err))

			require.ElementsMatch(t, c.ExpectedTheses, scenario.Theses())
		})
	}
}

func TestGetScenarioThesisBySlug(t *testing.T) {
	t.Parallel()

	scenario := errlessBuildScenario(t, func(b *specification.ScenarioBuilder) {
		b.WithThesis("c", func(b *specification.ThesisBuilder) {})
		b.WithThesis("d", func(b *specification.ThesisBuilder) {})
	})

	c, ok := scenario.Thesis("c")
	require.True(t, ok)
	require.Equal(
		t,
		specification.NewThesisBuilder().ErrlessBuild(
			specification.NewThesisSlug("aaa", "bb", "c"),
		),
		c,
	)

	d, ok := scenario.Thesis("d")
	require.True(t, ok)
	require.Equal(
		t,
		specification.NewThesisBuilder().ErrlessBuild(
			specification.NewThesisSlug("aaa", "bb", "d"),
		),
		d,
	)

	_, ok = scenario.Thesis("f")
	require.False(t, ok)
}

func TestScenarioErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "no_scenario_theses_error",
			Err:   specification.NewNoScenarioThesesError(),
			IsErr: specification.IsNoScenarioThesesError,
		},
		{
			Name:     "NON_no_scenario_theses_error",
			Err:      errors.New("no scenario theses"),
			IsErr:    specification.IsNoScenarioThesesError,
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
