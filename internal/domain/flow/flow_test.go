package flow_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/harpyd/thestis/internal/domain/flow"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestNewStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Slug         specification.Slug
		State        flow.State
		OccurredErrs []string
	}{
		{
			Slug:  specification.Slug{},
			State: flow.NoState,
		},
		{
			Slug:  specification.NewThesisSlug("foo", "bar", "zar"),
			State: flow.Canceled,
		},
		{
			Slug:  specification.NewThesisSlug("foo", "pam", "par"),
			State: flow.Failed,
			OccurredErrs: []string{
				"some error",
				"other error",
				"another error",
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			s := flow.NewStatus(c.Slug, c.State, c.OccurredErrs...)

			t.Run("slug", func(t *testing.T) {
				assert.Equal(t, c.Slug, s.Slug())
			})

			t.Run("state", func(t *testing.T) {
				assert.Equal(t, c.State, s.State())
			})

			t.Run("occurred_errs", func(t *testing.T) {
				assert.ElementsMatch(t, c.OccurredErrs, s.OccurredErrs())
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
		ExpectedStatuses      []flow.Status
	}{
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromPerformance("", &performance.Performance{})
			},
			ExpectedFlowID:        "",
			ExpectedPerformanceID: "",
			ExpectedStatuses:      nil,
		},
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromPerformance("foo", &performance.Performance{})
			},
			ExpectedFlowID:        "foo",
			ExpectedPerformanceID: "",
			ExpectedStatuses:      nil,
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
		},
		{
			FlowReducer: func() *flow.Reducer {
				spec := specification.NewBuilder().
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
			ExpectedStatuses: []flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "koo"),
					flow.NotPerformed,
				),
				flow.NewStatus(
					specification.NewThesisSlug("foo", "koo", "too"),
					flow.NotPerformed,
				),
			},
		},
		{
			FlowReducer: func() *flow.Reducer {
				spec := specification.NewBuilder().
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
							b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild()

				return flow.FromPerformance("dar", performance.FromSpecification("fla", spec)).
					WithStep(performance.NewThesisStep(
						specification.NewThesisSlug("foo", "bar", "baz"),
						performance.HTTPPerformer,
						performance.FiredPass,
					))
			},
			ExpectedFlowID:        "dar",
			ExpectedPerformanceID: "fla",
			ExpectedStatuses: []flow.Status{
				flow.NewStatus(
					specification.NewThesisSlug("foo", "bar", "baz"),
					flow.Passed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.NotPerformed,
				),
			},
		},
		{
			FlowReducer: func() *flow.Reducer {
				spec := specification.NewBuilder().
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
			ExpectedStatuses: []flow.Status{
				flow.NewStatus(
					specification.NewThesisSlug("doo", "zoo", "moo"),
					flow.NotPerformed,
				),
				flow.NewStatus(
					specification.NewScenarioSlug("doo", "zoo"),
					flow.Crashed,
					"something wrong",
				),
			},
		},
		{
			FlowReducer: func() *flow.Reducer {
				return flow.FromStatuses(
					"aba",
					"oba",
					flow.NewStatus(
						specification.NewThesisSlug("foo", "bar", "baz"),
						flow.NotPerformed,
					),
				).WithStep(performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "NOP"),
					performance.HTTPPerformer,
					performance.FiredFail,
				))
			},
			ExpectedFlowID:        "aba",
			ExpectedPerformanceID: "oba",
			ExpectedStatuses: []flow.Status{
				flow.NewStatus(
					specification.NewThesisSlug("foo", "bar", "baz"),
					flow.NotPerformed,
				),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			f := c.FlowReducer().Reduce()

			t.Run("id", func(t *testing.T) {
				assert.Equal(t, c.ExpectedFlowID, f.ID())
			})

			t.Run("performance_id", func(t *testing.T) {
				assert.Equal(t, c.ExpectedPerformanceID, f.PerformanceID())
			})

			t.Run("statuses", func(t *testing.T) {
				assert.ElementsMatch(t, c.ExpectedStatuses, f.Statuses())
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
				Statuses: []flow.Status{
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "bar"),
						flow.NotPerformed,
					),
				},
			},
		},
		{
			FlowParams: flow.Params{
				ID:            "flow-id",
				PerformanceID: "perf-id",
				Statuses: []flow.Status{
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "doo"),
						flow.Performing,
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
				assert.Equal(t, c.FlowParams.ID, f.ID())
			})

			t.Run("performance_id", func(t *testing.T) {
				assert.Equal(t, c.FlowParams.PerformanceID, f.PerformanceID())
			})

			t.Run("statuses", func(t *testing.T) {
				assert.ElementsMatch(t, c.FlowParams.Statuses, f.Statuses())
			})
		})
	}
}
