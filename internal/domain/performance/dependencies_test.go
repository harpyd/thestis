package performance_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestSyncDependenciesSnapshotsEqual(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		LeftSnapshot  performance.DependenciesSnapshot
		RightSnapshot performance.DependenciesSnapshot
		ExpectedEqual bool
	}{
		{
			LeftSnapshot:  nil,
			RightSnapshot: nil,
			ExpectedEqual: true,
		},
		{
			LeftSnapshot:  performance.DependenciesSnapshot{},
			RightSnapshot: nil,
			ExpectedEqual: true,
		},
		{
			LeftSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
					specification.NewThesisSlug("c", "c", "c"),
				},
			},
			RightSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
					specification.NewThesisSlug("c", "c", "c"),
				},
			},
			ExpectedEqual: true,
		},
		{
			LeftSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
					specification.NewThesisSlug("c", "c", "c"),
				},
			},
			RightSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("c", "c", "c"),
					specification.NewThesisSlug("b", "b", "b"),
				},
			},
			ExpectedEqual: true,
		},
		{
			LeftSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
					specification.NewThesisSlug("c", "c", "c"),
				},
			},
			RightSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("c", "c", "c"),
					specification.NewThesisSlug("b", "b", "b"),
					specification.NewThesisSlug("b", "b", "b"),
				},
			},
			ExpectedEqual: true,
		},
		{
			LeftSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
					specification.NewThesisSlug("c", "c", "c"),
				},
			},
			RightSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("c", "c", "c"),
					specification.NewThesisSlug("d", "d", "d"),
				},
			},
			ExpectedEqual: false,
		},
		{
			LeftSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
				},
				specification.NewThesisSlug("c", "c", "c"): {
					specification.NewThesisSlug("d", "d", "d"),
				},
			},
			RightSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
					specification.NewThesisSlug("c", "c", "c"),
					specification.NewThesisSlug("d", "d", "d"),
				},
			},
			ExpectedEqual: false,
		},
		{
			LeftSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
				},
				specification.NewThesisSlug("c", "c", "c"): {
					specification.NewThesisSlug("d", "d", "d"),
				},
			},
			RightSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
				},
				specification.NewThesisSlug("g", "g", "g"): {
					specification.NewThesisSlug("d", "d", "d"),
				},
			},
			ExpectedEqual: false,
		},
		{
			LeftSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
					specification.NewThesisSlug("c", "c", "c"),
				},
			},
			RightSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("a", "a", "a"): {
					specification.NewThesisSlug("b", "b", "b"),
				},
			},
			ExpectedEqual: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedEqual, c.LeftSnapshot.Equal(c.RightSnapshot))
		})
	}
}

func TestCollectDependencies(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenScenario    specification.Scenario
		ExpectedSlug     specification.Slug
		ExpectedSnapshot performance.DependenciesSnapshot
	}{
		{
			GivenScenario:    specification.Scenario{},
			ExpectedSlug:     specification.Slug{},
			ExpectedSnapshot: nil,
		},
		{
			GivenScenario: specification.NewScenarioBuilder().
				WithThesis("dak", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "dak")
				}).
				WithThesis("map", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.When, "map")
				}).
				WithThesis("qwe", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Then, "qwe")
				}).
				ErrlessBuild(specification.NewScenarioSlug("foo", "bak")),
			ExpectedSlug: specification.NewScenarioSlug("foo", "bak"),
			ExpectedSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("foo", "bak", "map"): {
					specification.NewThesisSlug("foo", "bak", "dak"),
				},
				specification.NewThesisSlug("foo", "bak", "qwe"): {
					specification.NewThesisSlug("foo", "bak", "dak"),
					specification.NewThesisSlug("foo", "bak", "map"),
				},
			},
		},
		{
			GivenScenario: specification.NewScenarioBuilder().
				WithThesis("qyz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "qyz")
				}).
				WithThesis("qyp", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "qyp")
				}).
				WithThesis("bad", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "bad")
					b.WithDependency("qyz")
					b.WithDependency("qyp")
				}).
				WithThesis("dad", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.When, "dad")
				}).
				WithThesis("tad", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Then, "tad")
				}).
				ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
			ExpectedSlug: specification.NewScenarioSlug("foo", "bar"),
			ExpectedSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("foo", "bar", "bad"): {
					specification.NewThesisSlug("foo", "bar", "qyz"),
					specification.NewThesisSlug("foo", "bar", "qyp"),
				},
				specification.NewThesisSlug("foo", "bar", "dad"): {
					specification.NewThesisSlug("foo", "bar", "qyz"),
					specification.NewThesisSlug("foo", "bar", "qyp"),
					specification.NewThesisSlug("foo", "bar", "bad"),
				},
				specification.NewThesisSlug("foo", "bar", "tad"): {
					specification.NewThesisSlug("foo", "bar", "qyz"),
					specification.NewThesisSlug("foo", "bar", "qyp"),
					specification.NewThesisSlug("foo", "bar", "bad"),
					specification.NewThesisSlug("foo", "bar", "dad"),
				},
			},
		},
		{
			GivenScenario: specification.NewScenarioBuilder().
				WithThesis("baz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "baz")
				}).
				WithThesis("bad", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Then, "bad")
				}).
				ErrlessBuild(specification.NewScenarioSlug("foo", "loo")),
			ExpectedSlug: specification.NewScenarioSlug("foo", "loo"),
			ExpectedSnapshot: performance.DependenciesSnapshot{
				specification.NewThesisSlug("foo", "loo", "bad"): {
					specification.NewThesisSlug("foo", "loo", "baz"),
				},
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			deps := performance.SyncDependencies(c.GivenScenario)

			t.Run("snapshot", func(t *testing.T) {
				actual := deps.Snapshot()

				require.Truef(
					t,
					c.ExpectedSnapshot.Equal(actual),
					"expected : %s\nactual   : %s",
					c.ExpectedSnapshot, actual,
				)
			})

			t.Run("slug", func(t *testing.T) {
				require.Equal(t, c.ExpectedSlug, deps.Slug())
			})
		})
	}
}

func TestWaitThesisDependencies(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenScenario specification.Scenario
		ThesisToWait  specification.Slug
		ThesesToDo    []specification.Slug
		ShouldWait    bool
	}{
		{
			GivenScenario: specification.Scenario{},
			ThesisToWait:  specification.Slug{},
			ShouldWait:    true,
		},
		{
			GivenScenario: specification.Scenario{},
			ThesisToWait:  specification.NewThesisSlug("foo", "bar", "baz"),
			ShouldWait:    true,
		},
		{
			GivenScenario: specification.NewScenarioBuilder().
				WithThesis("poo", func(b *specification.ThesisBuilder) {}).
				ErrlessBuild(specification.NewScenarioSlug("foo", "koo")),
			ThesisToWait: specification.NewThesisSlug("foo", "koo", "poo"),
			ShouldWait:   true,
		},
		{
			GivenScenario: specification.NewScenarioBuilder().
				WithThesis("poo", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "poo")
					b.WithDependency("too")
				}).
				WithThesis("too", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "too")
				}).
				ErrlessBuild(specification.NewScenarioSlug("foo", "koo")),
			ThesisToWait: specification.NewThesisSlug("foo", "koo", "nop"),
			ShouldWait:   true,
		},
		{
			GivenScenario: specification.NewScenarioBuilder().
				WithThesis("baz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "baz")
					b.WithDependency("qyp")
					b.WithDependency("qyz")
				}).
				WithThesis("qyp", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "qyp")
				}).
				WithThesis("qyz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "qyz")
				}).
				ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
			ThesisToWait: specification.NewThesisSlug("foo", "bar", "baz"),
			ThesesToDo: []specification.Slug{
				specification.NewThesisSlug("foo", "bar", "qyp"),
				specification.NewThesisSlug("foo", "bar", "qyz"),
			},
			ShouldWait: true,
		},
		{
			GivenScenario: specification.NewScenarioBuilder().
				WithThesis("baz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "baz")
					b.WithDependency("bad")
				}).
				WithThesis("bad", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "bad")
				}).
				ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
			ThesisToWait: specification.NewThesisSlug("foo", "bar", "baz"),
			ThesesToDo: []specification.Slug{
				specification.NewThesisSlug("foo", "bar", "tad"),
			},
			ShouldWait: false,
		},
		{
			GivenScenario: specification.NewScenarioBuilder().
				WithThesis("baz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "baz")
					b.WithDependency("qyp")
					b.WithDependency("qyz")
					b.WithDependency("pyz")
				}).
				WithThesis("qyp", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "qyp")
				}).
				WithThesis("qyz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "qyz")
				}).
				WithThesis("pyz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "pyz")
				}).
				ErrlessBuild(specification.NewScenarioSlug("foo", "bar")),
			ThesisToWait: specification.NewThesisSlug("foo", "bar", "baz"),
			ThesesToDo: []specification.Slug{
				specification.NewThesisSlug("foo", "bar", "qyp"),
				specification.NewThesisSlug("foo", "bar", "qyz"),
			},
			ShouldWait: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			const timout = 3 * time.Millisecond

			ctx, cancel := context.WithTimeout(context.Background(), timout)
			defer cancel()

			sg := performance.SyncDependencies(c.GivenScenario)

			go func() {
				for _, todo := range c.ThesesToDo {
					sg.ThesisDone(todo)
				}
			}()

			err := sg.WaitThesisDependencies(ctx, c.ThesisToWait)

			if !c.ShouldWait {
				t.Run("cancel_err", func(t *testing.T) {
					var terr *performance.TerminatedError

					require.ErrorAs(t, err, &terr)
					require.Equal(t, performance.FiredCancel, terr.Event)
				})

				return
			}

			t.Run("wait_without_err", func(t *testing.T) {
				require.NoError(t, err)
			})
		})
	}
}
