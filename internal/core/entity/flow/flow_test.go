package flow_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/entity/flow"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/specification"
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

func TestFulfilledFlow(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		FlowFactory          func() *flow.Flow
		ExpectedFlowID       string
		ExpectedPipelineID   string
		ExpectedStatuses     []*flow.Status
		ExpectedOverallState flow.State
	}{
		{
			FlowFactory: func() *flow.Flow {
				return flow.Fulfill("", &pipeline.Pipeline{})
			},
			ExpectedFlowID:       "",
			ExpectedPipelineID:   "",
			ExpectedStatuses:     nil,
			ExpectedOverallState: flow.NoState,
		},
		{
			FlowFactory: func() *flow.Flow {
				return flow.Fulfill("foo", &pipeline.Pipeline{})
			},
			ExpectedFlowID:       "foo",
			ExpectedPipelineID:   "",
			ExpectedStatuses:     nil,
			ExpectedOverallState: flow.NoState,
		},
		{
			FlowFactory: func() *flow.Flow {
				return flow.Fulfill("bar", pipeline.Unmarshal(pipeline.Params{
					ID: "doo",
				}))
			},
			ExpectedFlowID:       "bar",
			ExpectedPipelineID:   "doo",
			ExpectedStatuses:     nil,
			ExpectedOverallState: flow.NoState,
		},
		{
			FlowFactory: func() *flow.Flow {
				spec := (&specification.Builder{}).
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("koo", func(b *specification.ScenarioBuilder) {
							b.WithThesis("too", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild()

				return flow.Fulfill("rar", pipeline.Trigger("kra", spec))
			},
			ExpectedFlowID:     "rar",
			ExpectedPipelineID: "kra",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "koo"),
					flow.NotExecuted,
					flow.NewThesisStatus("too", flow.NotExecuted),
				),
			},
			ExpectedOverallState: flow.NotExecuted,
		},
		{
			FlowFactory: func() *flow.Flow {
				spec := (&specification.Builder{}).
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
							b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild()

				f := flow.Fulfill("dar", pipeline.Trigger("fla", spec))

				return f.ApplyStep(pipeline.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.HTTPExecutor,
					pipeline.FiredPass,
				))
			},
			ExpectedFlowID:     "dar",
			ExpectedPipelineID: "fla",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.NotExecuted,
					flow.NewThesisStatus("baz", flow.Passed),
				),
			},
			ExpectedOverallState: flow.NotExecuted,
		},
		{
			FlowFactory: func() *flow.Flow {
				spec := (&specification.Builder{}).
					WithStory("doo", func(b *specification.StoryBuilder) {
						b.WithScenario("zoo", func(b *specification.ScenarioBuilder) {
							b.WithThesis("moo", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild()

				return flow.Fulfill("sds", pipeline.Trigger("coo", spec)).
					ApplyStep(pipeline.NewScenarioStepWithErr(
						errors.New("something wrong"),
						specification.NewScenarioSlug("doo", "zoo"),
						pipeline.FiredCrash,
					))
			},
			ExpectedFlowID:     "sds",
			ExpectedPipelineID: "coo",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("doo", "zoo"),
					flow.Crashed,
					flow.NewThesisStatus("moo", flow.NotExecuted),
				),
			},
			ExpectedOverallState: flow.Crashed,
		},
		{
			FlowFactory: func() *flow.Flow {
				return flow.FromStatuses(
					"aba",
					"oba",
					flow.NewStatus(
						specification.NewScenarioSlug("foo", "bar"),
						flow.Executing,
						flow.NewThesisStatus("baz", flow.Passed),
						flow.NewThesisStatus("ban", flow.Executing),
					),
				).ApplyStep(pipeline.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "NOP"),
					pipeline.HTTPExecutor,
					pipeline.FiredFail,
				))
			},
			ExpectedFlowID:     "aba",
			ExpectedPipelineID: "oba",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("foo", "bar"),
					flow.Executing,
					flow.NewThesisStatus("baz", flow.Passed),
					flow.NewThesisStatus("ban", flow.Executing),
				),
			},
			ExpectedOverallState: flow.Executing,
		},
		{
			FlowFactory: func() *flow.Flow {
				return flow.FromStatuses(
					"flow-id",
					"some-pipe-id",
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
			ExpectedFlowID:     "flow-id",
			ExpectedPipelineID: "some-pipe-id",
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
			FlowFactory: func() *flow.Flow {
				return flow.FromStatuses(
					"id",
					"pipe-id",
					flow.NewStatus(
						specification.NewScenarioSlug("a", "b"),
						flow.Executing,
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
			ExpectedFlowID:     "id",
			ExpectedPipelineID: "pipe-id",
			ExpectedStatuses: []*flow.Status{
				flow.NewStatus(
					specification.NewScenarioSlug("a", "b"),
					flow.Executing,
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
			ExpectedOverallState: flow.Executing,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			f := c.FlowFactory()

			t.Run("id", func(t *testing.T) {
				require.Equal(t, c.ExpectedFlowID, f.ID())
			})

			t.Run("pipeline_id", func(t *testing.T) {
				require.Equal(t, c.ExpectedPipelineID, f.PipelineID())
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
