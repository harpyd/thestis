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
		Step                  performance.Step
		ExpectedSlug          specification.Slug
		ExpectedPerformerType performance.PerformerType
		ExpectedEvent         performance.Event
		ExpectedErr           error
		ExpectedIsZero        bool
		ExpectedString        string
	}{
		{
			Step:                  performance.Step{},
			ExpectedSlug:          specification.Slug{},
			ExpectedPerformerType: performance.NoPerformer,
			ExpectedEvent:         performance.NoEvent,
			ExpectedErr:           nil,
			ExpectedIsZero:        true,
			ExpectedString:        "",
		},
		{
			Step: performance.NewScenarioStep(
				specification.Slug{},
				performance.FiredFail,
			),
			ExpectedSlug:          specification.Slug{},
			ExpectedPerformerType: performance.NoPerformer,
			ExpectedEvent:         performance.FiredFail,
			ExpectedErr:           nil,
			ExpectedIsZero:        false,
			ExpectedString:        ": event = fail",
		},
		{
			Step: performance.NewScenarioStep(
				specification.NewScenarioSlug("foo", "bar"),
				performance.FiredCrash,
			),
			ExpectedSlug:          specification.NewScenarioSlug("foo", "bar"),
			ExpectedPerformerType: performance.NoPerformer,
			ExpectedEvent:         performance.FiredCrash,
			ExpectedErr:           nil,
			ExpectedIsZero:        false,
			ExpectedString:        "foo.bar: event = crash",
		},
		{
			Step: performance.NewScenarioStepWithErr(
				errors.New("something wrong"),
				specification.AnyScenarioSlug(),
				performance.FiredCancel,
			),
			ExpectedSlug:          specification.AnyScenarioSlug(),
			ExpectedPerformerType: performance.NoPerformer,
			ExpectedEvent:         performance.FiredCancel,
			ExpectedErr:           errors.New("something wrong"),
			ExpectedIsZero:        false,
			ExpectedString:        "*.*: event = cancel, err = something wrong",
		},
		{
			Step: performance.NewThesisStep(
				specification.NewThesisSlug("foo", "bar", "baz"),
				performance.HTTPPerformer,
				performance.FiredFail,
			),
			ExpectedSlug:          specification.NewThesisSlug("foo", "bar", "baz"),
			ExpectedPerformerType: performance.HTTPPerformer,
			ExpectedEvent:         performance.FiredFail,
			ExpectedErr:           nil,
			ExpectedIsZero:        false,
			ExpectedString:        "foo.bar.baz: event = fail, type = HTTP",
		},
		{
			Step: performance.NewThesisStepWithErr(
				errors.New("wrong"),
				specification.NewThesisSlug("foo", "bar", "baz"),
				performance.AssertionPerformer,
				performance.FiredCrash,
			),
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

			t.Run("slug", func(t *testing.T) {
				require.Equal(t, c.ExpectedSlug, c.Step.Slug())
			})

			t.Run("performer_type", func(t *testing.T) {
				require.Equal(t, c.ExpectedPerformerType, c.Step.PerformerType())
			})

			t.Run("event", func(t *testing.T) {
				require.Equal(t, c.ExpectedEvent, c.Step.Event())
			})

			if c.ExpectedErr != nil {
				t.Run("err", func(t *testing.T) {
					require.EqualError(t, c.Step.Err(), c.ExpectedErr.Error())
				})
			} else {
				t.Run("no_err", func(t *testing.T) {
					require.NoError(t, c.Step.Err())
				})
			}

			t.Run("is_zero", func(t *testing.T) {
				require.Equal(t, c.ExpectedIsZero, c.Step.IsZero())
			})

			t.Run("string", func(t *testing.T) {
				require.Equal(t, c.ExpectedString, c.Step.String())
			})
		})
	}
}
