package pipeline_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

func TestStep(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		StepFactory          func() pipeline.Step
		ShouldPanic          bool
		ExpectedSlug         specification.Slug
		ExpectedExecutorType pipeline.ExecutorType
		ExpectedEvent        pipeline.Event
		ExpectedErr          error
		ExpectedIsZero       bool
		ExpectedString       string
	}{
		{
			StepFactory: func() pipeline.Step {
				return pipeline.Step{}
			},
			ExpectedSlug:         specification.Slug{},
			ExpectedExecutorType: pipeline.NoExecutor,
			ExpectedEvent:        pipeline.NoEvent,
			ExpectedErr:          nil,
			ExpectedIsZero:       true,
			ExpectedString:       "",
		},
		{
			StepFactory: func() pipeline.Step {
				return pipeline.NewScenarioStep(
					specification.Slug{},
					pipeline.FiredFail,
				)
			},
			ShouldPanic: true,
		},
		{
			StepFactory: func() pipeline.Step {
				return pipeline.NewScenarioStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.FiredExecute,
				)
			},
			ShouldPanic: true,
		},
		{
			StepFactory: func() pipeline.Step {
				return pipeline.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredCrash,
				)
			},
			ExpectedSlug:         specification.NewScenarioSlug("foo", "bar"),
			ExpectedExecutorType: pipeline.NoExecutor,
			ExpectedEvent:        pipeline.FiredCrash,
			ExpectedErr:          nil,
			ExpectedIsZero:       false,
			ExpectedString:       "foo.bar: event = crash",
		},
		{
			StepFactory: func() pipeline.Step {
				return pipeline.NewScenarioStepWithErr(
					errors.New("something wrong"),
					specification.NewScenarioSlug("a", "b"),
					pipeline.FiredCancel,
				)
			},
			ExpectedSlug:         specification.NewScenarioSlug("a", "b"),
			ExpectedExecutorType: pipeline.NoExecutor,
			ExpectedEvent:        pipeline.FiredCancel,
			ExpectedErr:          errors.New("something wrong"),
			ExpectedIsZero:       false,
			ExpectedString:       "a.b: event = cancel, err = something wrong",
		},
		{
			StepFactory: func() pipeline.Step {
				return pipeline.NewThesisStep(
					specification.Slug{},
					pipeline.HTTPExecutor,
					pipeline.FiredExecute,
				)
			},
			ShouldPanic: true,
		},
		{
			StepFactory: func() pipeline.Step {
				return pipeline.NewThesisStep(
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.HTTPExecutor,
					pipeline.FiredExecute,
				)
			},
			ShouldPanic: true,
		},
		{
			StepFactory: func() pipeline.Step {
				return pipeline.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.HTTPExecutor,
					pipeline.FiredFail,
				)
			},
			ExpectedSlug:         specification.NewThesisSlug("foo", "bar", "baz"),
			ExpectedExecutorType: pipeline.HTTPExecutor,
			ExpectedEvent:        pipeline.FiredFail,
			ExpectedErr:          nil,
			ExpectedIsZero:       false,
			ExpectedString:       "foo.bar.baz: event = fail, type = HTTP",
		},
		{
			StepFactory: func() pipeline.Step {
				return pipeline.NewThesisStepWithErr(
					errors.New("wrong"),
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.AssertionExecutor,
					pipeline.FiredCrash,
				)
			},
			ExpectedSlug:         specification.NewThesisSlug("foo", "bar", "baz"),
			ExpectedExecutorType: pipeline.AssertionExecutor,
			ExpectedEvent:        pipeline.FiredCrash,
			ExpectedErr:          errors.New("wrong"),
			ExpectedIsZero:       false,
			ExpectedString:       "foo.bar.baz: event = crash, type = assertion, err = wrong",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			if c.ShouldPanic {
				t.Run("panics", func(t *testing.T) {
					require.Panics(t, func() {
						_ = c.StepFactory()
					})
				})

				return
			}

			var step pipeline.Step

			t.Run("not_panics", func(t *testing.T) {
				require.NotPanics(t, func() {
					step = c.StepFactory()
				})
			})

			t.Run("slug", func(t *testing.T) {
				require.Equal(t, c.ExpectedSlug, step.Slug())
			})

			t.Run("executor_type", func(t *testing.T) {
				require.Equal(t, c.ExpectedExecutorType, step.ExecutorType())
			})

			t.Run("event", func(t *testing.T) {
				require.Equal(t, c.ExpectedEvent, step.Event())
			})

			if c.ExpectedErr != nil {
				t.Run("err", func(t *testing.T) {
					require.EqualError(t, step.Err(), c.ExpectedErr.Error())
				})
			} else {
				t.Run("no_err", func(t *testing.T) {
					require.NoError(t, step.Err())
				})
			}

			t.Run("is_zero", func(t *testing.T) {
				require.Equal(t, c.ExpectedIsZero, step.IsZero())
			})

			t.Run("string", func(t *testing.T) {
				require.Equal(t, c.ExpectedString, step.String())
			})
		})
	}
}
