package performance_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestStepCreation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		StepFactory           func() performance.Step
		ShouldPanic           bool
		ExpectedSlug          specification.Slug
		ExpectedPerformerType performance.PerformerType
		ExpectedEvent         performance.Event
		ExpectedErr           error
		ExpectedIsZero        bool
		ExpectedString        string
	}{
		{
			StepFactory: func() performance.Step {
				return performance.Step{}
			},
			ExpectedSlug:          specification.Slug{},
			ExpectedPerformerType: performance.NoPerformer,
			ExpectedEvent:         performance.NoEvent,
			ExpectedErr:           nil,
			ExpectedIsZero:        true,
			ExpectedString:        "",
		},
		{
			StepFactory: func() performance.Step {
				return performance.NewScenarioStep(
					specification.Slug{},
					performance.FiredFail,
				)
			},
			ShouldPanic: true,
		},
		{
			StepFactory: func() performance.Step {
				return performance.NewScenarioStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.FiredPerform,
				)
			},
			ShouldPanic: true,
		},
		{
			StepFactory: func() performance.Step {
				return performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredCrash,
				)
			},
			ExpectedSlug:          specification.NewScenarioSlug("foo", "bar"),
			ExpectedPerformerType: performance.NoPerformer,
			ExpectedEvent:         performance.FiredCrash,
			ExpectedErr:           nil,
			ExpectedIsZero:        false,
			ExpectedString:        "foo.bar: event = crash",
		},
		{
			StepFactory: func() performance.Step {
				return performance.NewScenarioStepWithErr(
					errors.New("something wrong"),
					specification.AnyScenarioSlug(),
					performance.FiredCancel,
				)
			},
			ExpectedSlug:          specification.AnyScenarioSlug(),
			ExpectedPerformerType: performance.NoPerformer,
			ExpectedEvent:         performance.FiredCancel,
			ExpectedErr:           errors.New("something wrong"),
			ExpectedIsZero:        false,
			ExpectedString:        "*.*: event = cancel, err = something wrong",
		},
		{
			StepFactory: func() performance.Step {
				return performance.NewThesisStep(
					specification.Slug{},
					performance.HTTPPerformer,
					performance.FiredPerform,
				)
			},
			ShouldPanic: true,
		},
		{
			StepFactory: func() performance.Step {
				return performance.NewThesisStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				)
			},
			ShouldPanic: true,
		},
		{
			StepFactory: func() performance.Step {
				return performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.HTTPPerformer,
					performance.FiredFail,
				)
			},
			ExpectedSlug:          specification.NewThesisSlug("foo", "bar", "baz"),
			ExpectedPerformerType: performance.HTTPPerformer,
			ExpectedEvent:         performance.FiredFail,
			ExpectedErr:           nil,
			ExpectedIsZero:        false,
			ExpectedString:        "foo.bar.baz: event = fail, type = HTTP",
		},
		{
			StepFactory: func() performance.Step {
				return performance.NewThesisStepWithErr(
					errors.New("wrong"),
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.AssertionPerformer,
					performance.FiredCrash,
				)
			},
			ExpectedSlug:          specification.NewThesisSlug("foo", "bar", "baz"),
			ExpectedPerformerType: performance.AssertionPerformer,
			ExpectedEvent:         performance.FiredCrash,
			ExpectedErr:           errors.New("wrong"),
			ExpectedIsZero:        false,
			ExpectedString:        "foo.bar.baz: event = crash, type = assertion, err = wrong",
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

			var step performance.Step

			t.Run("not_panics", func(t *testing.T) {
				require.NotPanics(t, func() {
					step = c.StepFactory()
				})
			})

			t.Run("slug", func(t *testing.T) {
				require.Equal(t, c.ExpectedSlug, step.Slug())
			})

			t.Run("performer_type", func(t *testing.T) {
				require.Equal(t, c.ExpectedPerformerType, step.PerformerType())
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
