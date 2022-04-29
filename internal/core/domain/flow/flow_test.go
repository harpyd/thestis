package flow_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/domain/flow"
	"github.com/harpyd/thestis/internal/core/domain/performance"
	"github.com/harpyd/thestis/internal/core/domain/specification"
)

func TestNewStatus(t *testing.T) {
	t.Parallel()

	type params struct {
		Slug           specification.Slug
		State          flow.State
		ThesisStatuses []*flow.ThesisStatus
	}

	testCases := []struct {
		Given       params
		Expected    params
		ShouldPanic bool
	}{
		{
			Given: params{
				Slug:  specification.Slug{},
				State: flow.NoState,
			},
			ShouldPanic: true,
		},
		{
			Given: params{
				Slug:  specification.NewThesisSlug("a", "b", "c"),
				State: flow.Crashed,
			},
			ShouldPanic: true,
		},
		{
			Given: params{
				Slug:  specification.NewScenarioSlug("foo", "bar"),
				State: flow.Canceled,
				ThesisStatuses: []*flow.ThesisStatus{
					flow.NewThesisStatus("baz", flow.Passed),
					flow.NewThesisStatus("bam", flow.Canceled),
				},
			},
			Expected: params{
				Slug:  specification.NewScenarioSlug("foo", "bar"),
				State: flow.Canceled,
				ThesisStatuses: []*flow.ThesisStatus{
					flow.NewThesisStatus("baz", flow.Passed),
					flow.NewThesisStatus("bam", flow.Canceled),
				},
			},
			ShouldPanic: false,
		},
		{
			Given: params{
				Slug:  specification.NewScenarioSlug("foo", "pam"),
				State: flow.Failed,
				ThesisStatuses: []*flow.ThesisStatus{
					nil,
					flow.NewThesisStatus("tam", flow.Crashed),
					nil,
					flow.NewThesisStatus("ram", flow.Failed),
				},
			},
			Expected: params{
				Slug:  specification.NewScenarioSlug("foo", "pam"),
				State: flow.Failed,
				ThesisStatuses: []*flow.ThesisStatus{
					flow.NewThesisStatus("tam", flow.Crashed),
					flow.NewThesisStatus("ram", flow.Failed),
				},
			},
			ShouldPanic: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var s *flow.Status

			newFn := func() {
				s = flow.NewStatus(c.Given.Slug, c.Given.State, c.Given.ThesisStatuses...)
			}

			if c.ShouldPanic {
				t.Run("panics", func(t *testing.T) {
					require.PanicsWithValue(t, specification.ErrNotScenarioSlug, newFn)
				})

				return
			}

			t.Run("not_panics", func(t *testing.T) {
				require.NotPanics(t, newFn)
			})

			t.Run("slug", func(t *testing.T) {
				require.Equal(t, c.Expected.Slug, s.Slug())
			})

			t.Run("state", func(t *testing.T) {
				require.Equal(t, c.Expected.State, s.State())
			})

			t.Run("thesis_statuses", func(t *testing.T) {
				require.ElementsMatch(t, c.Expected.ThesisStatuses, s.ThesisStatuses())
			})
		})
	}
}

func TestNewThesisStatus(t *testing.T) {
	t.Parallel()

	type params struct {
		ThesisSlug   string
		State        flow.State
		OccurredErrs []string
	}

	testCases := []struct {
		Given    params
		Expected params
	}{
		{
			Given: params{
				ThesisSlug: "",
				State:      flow.NoState,
			},
			Expected: params{
				ThesisSlug: "",
				State:      flow.NoState,
			},
		},
		{
			Given: params{
				ThesisSlug:   "c",
				State:        flow.Failed,
				OccurredErrs: []string{"some err"},
			},
			Expected: params{
				ThesisSlug:   "c",
				State:        flow.Failed,
				OccurredErrs: []string{"some err"},
			},
		},
		{
			Given: params{
				ThesisSlug: "qyp",
				State:      flow.Crashed,
				OccurredErrs: []string{
					"",
					"bar",
					"foo",
				},
			},
			Expected: params{
				ThesisSlug: "qyp",
				State:      flow.Crashed,
				OccurredErrs: []string{
					"",
					"bar",
					"foo",
				},
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			s := flow.NewThesisStatus(c.Given.ThesisSlug, c.Given.State, c.Given.OccurredErrs...)

			t.Run("thesis_slug", func(t *testing.T) {
				require.Equal(t, c.Expected.ThesisSlug, s.ThesisSlug())
			})

			t.Run("state", func(t *testing.T) {
				require.Equal(t, c.Expected.State, s.State())
			})

			t.Run("occurred_errs", func(t *testing.T) {
				require.Equal(t, c.Expected.OccurredErrs, s.OccurredErrs())
			})
		})
	}
}

func TestReduceFlow(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		FlowReducer           func() *flow.Reducer
		ExpectedFlowID        string
		ExpectedPerformanceID string
		ExpectedStatuses      []*flow.Status
		ExpectedOverallState  flow.State
	}{
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromPerformance("", &performance.Performance{})
			},
			ExpectedFlowID:        "",
			ExpectedPerformanceID: "",
			ExpectedStatuses:      nil,
			ExpectedOverallState:  flow.NoState,
		},
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromPerformance("foo", &performance.Performance{})
			},
			ExpectedFlowID:        "foo",
			ExpectedPerformanceID: "",
			ExpectedStatuses:      nil,
			ExpectedOverallState:  flow.NoState,
		},
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromPerformance("bar", performance.Unmarshal(performance.Params{
					ID: "doo",
				}))
			},
			ExpectedFlowID:        "bar",
			ExpectedPerformanceID: "doo",
			ExpectedStatuses:      nil,
			ExpectedOverallState:  flow.NoState,
		},
		{
			FlowReducer: func() *flow.Reducer {
				spec := (&specification.Builder{}).
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("koo", func(b *specification.ScenarioBuilder) {
							b.WithThesis("too", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild()

				return flow.FromPerformance("rar", performance.FromSpecification("kra", spec))
			},
			ExpectedFlowID:        "rar",
			ExpectedPerformanceID: "kra",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "koo"),
					flow.NotPerformed,
					flow.NewThesisStatus("too", flow.NotPerformed),
				),
			},
			ExpectedOverallState: flow.NotPerformed,
		},
		{
			FlowReducer: func() *flow.Reducer {
				spec := (&specification.Builder{}).
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
							b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild()

				f := flow.FromPerformance("dar", performance.FromSpecification("fla", spec))

				return f.WithStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.HTTPPerformer,
					performance.FiredPass,
				))
			},
			ExpectedFlowID:        "dar",
			ExpectedPerformanceID: "fla",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.NotPerformed,
					flow.NewThesisStatus("baz", flow.Passed),
				),
			},
			ExpectedOverallState: flow.NotPerformed,
		},
		{
			FlowReducer: func() *flow.Reducer {
				spec := (&specification.Builder{}).
					WithStory("doo", func(b *specification.StoryBuilder) {
						b.WithScenario("zoo", func(b *specification.ScenarioBuilder) {
							b.WithThesis("moo", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild()

				return flow.FromPerformance("sds", performance.FromSpecification("coo", spec)).
					WithStep(performance.NewScenarioStepWithErr(
						errors.New("something wrong"),
						specification.NewScenarioSlug("doo", "zoo"),
						performance.FiredCrash,
					))
			},
			ExpectedFlowID:        "sds",
			ExpectedPerformanceID: "coo",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("doo", "zoo"),
					flow.Crashed,
					flow.NewThesisStatus("moo", flow.NotPerformed),
				),
			},
			ExpectedOverallState: flow.Crashed,
		},
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromStatuses(
					"aba",
					"oba",
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "bar"),
						flow.Performing,
						flow.NewThesisStatus("baz", flow.Passed),
						flow.NewThesisStatus("ban", flow.Performing),
					),
				).WithStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "NOP"),
					performance.HTTPPerformer,
					performance.FiredFail,
				))
			},
			ExpectedFlowID:        "aba",
			ExpectedPerformanceID: "oba",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.Performing,
					flow.NewThesisStatus("baz", flow.Passed),
					flow.NewThesisStatus("ban", flow.Performing),
				),
			},
			ExpectedOverallState: flow.Performing,
		},
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromStatuses(
					"flow-id",
					"some-perf-id",
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "bar"),
						flow.Passed,
					),
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "baz"),
						flow.Passed,
					),
				)
			},
			ExpectedFlowID:        "flow-id",
			ExpectedPerformanceID: "some-perf-id",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.Passed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "baz"),
					flow.Passed,
				),
			},
			ExpectedOverallState: flow.Passed,
		},
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromStatuses(
					"id",
					"perf-id",
					flow.NewStatus(
						specification.NewScenarioSlug("a", "b"),
						flow.Performing,
					),
					flow.NewStatus(
						specification.NewScenarioSlug("a", "d"),
						flow.Passed,
					),
					flow.NewStatus(
						specification.NewScenarioSlug("b", "c"),
						flow.Failed,
					),
					flow.NewStatus(
						specification.NewScenarioSlug("b", "d"),
						flow.Crashed,
					),
				)
			},
			ExpectedFlowID:        "id",
			ExpectedPerformanceID: "perf-id",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("a", "b"),
					flow.Performing,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("a", "d"),
					flow.Passed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("b", "c"),
					flow.Failed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("b", "d"),
					flow.Crashed,
				),
			},
			ExpectedOverallState: flow.Performing,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			f := c.FlowReducer().Reduce()

			t.Run("id", func(t *testing.T) {
				require.Equal(t, c.ExpectedFlowID, f.ID())
			})

			t.Run("performance_id", func(t *testing.T) {
				require.Equal(t, c.ExpectedPerformanceID, f.PerformanceID())
			})

			t.Run("statuses", func(t *testing.T) {
				require.ElementsMatch(t, c.ExpectedStatuses, f.Statuses())
			})

			t.Run("overall_state", func(t *testing.T) {
				require.Equal(t, c.ExpectedOverallState, f.OverallState())
			})
		})
	}
}

func TestUnmarshalFlow(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		FlowParams flow.Params
	}{
		{
			FlowParams: flow.Params{},
		},
		{
			FlowParams: flow.Params{
				ID: "flow-id",
			},
		},
		{
			FlowParams: flow.Params{
				PerformanceID: "perf-id",
			},
		},
		{
			FlowParams: flow.Params{
				Statuses: []*flow.Status{
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "bar"),
						flow.NotPerformed,
					),
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "dar"),
						flow.Performing,
					),
				},
			},
		},
		{
			FlowParams: flow.Params{
				ID:            "flow-id",
				PerformanceID: "perf-id",
				Statuses: []*flow.Status{
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "doo"),
						flow.Performing,
						flow.NewThesisStatus("boo", flow.Performing),
					),
					nil,
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "doo"),
						flow.Failed,
						flow.NewThesisStatus("zoo", flow.Failed),
					),
				},
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			f := flow.Unmarshal(c.FlowParams)

			t.Run("id", func(t *testing.T) {
				require.Equal(t, c.FlowParams.ID, f.ID())
			})

			t.Run("performance_id", func(t *testing.T) {
				require.Equal(t, c.FlowParams.PerformanceID, f.PerformanceID())
			})

			t.Run("statuses", func(t *testing.T) {
				require.ElementsMatch(t, c.FlowParams.Statuses, f.Statuses())
			})
		})
	}
}
