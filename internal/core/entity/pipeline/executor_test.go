package pipeline_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/entity/pipeline"
)

func TestResult(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenResult   pipeline.Result
		ExpectedEvent pipeline.Event
		ExpectedErr   error
	}{
		{
			GivenResult:   pipeline.Pass(),
			ExpectedEvent: pipeline.FiredPass,
		},
		{
			GivenResult:   pipeline.Fail(nil),
			ExpectedEvent: pipeline.FiredFail,
			ExpectedErr:   nil,
		},
		{
			GivenResult: pipeline.Fail(
				errors.New("foo"),
			),
			ExpectedEvent: pipeline.FiredFail,
			ExpectedErr: pipeline.WrapWithTerminatedError(
				errors.New("foo"),
				pipeline.FiredFail,
			),
		},
		{
			GivenResult: pipeline.Fail(
				pipeline.WrapWithTerminatedError(
					errors.New("foo"),
					pipeline.FiredFail,
				),
			),
			ExpectedEvent: pipeline.FiredFail,
			ExpectedErr: pipeline.WrapWithTerminatedError(
				errors.New("bar"),
				pipeline.FiredFail,
			),
		},
		{
			GivenResult:   pipeline.Crash(nil),
			ExpectedEvent: pipeline.FiredCrash,
			ExpectedErr:   nil,
		},
		{
			GivenResult: pipeline.Crash(
				errors.New("boo"),
			),
			ExpectedEvent: pipeline.FiredCrash,
			ExpectedErr: pipeline.WrapWithTerminatedError(
				errors.New("boo"),
				pipeline.FiredCrash,
			),
		},
		{
			GivenResult: pipeline.Crash(
				pipeline.WrapWithTerminatedError(
					errors.New("qwe"),
					pipeline.FiredCrash,
				),
			),
			ExpectedEvent: pipeline.FiredCrash,
			ExpectedErr: pipeline.WrapWithTerminatedError(
				errors.New("qwe"),
				pipeline.FiredCrash,
			),
		},
		{
			GivenResult:   pipeline.Cancel(nil),
			ExpectedEvent: pipeline.FiredCancel,
			ExpectedErr:   nil,
		},
		{
			GivenResult: pipeline.Cancel(
				errors.New("bar"),
			),
			ExpectedEvent: pipeline.FiredCancel,
			ExpectedErr: pipeline.WrapWithTerminatedError(
				errors.New("bar"),
				pipeline.FiredCancel,
			),
		},
		{
			GivenResult: pipeline.Cancel(
				pipeline.WrapWithTerminatedError(
					errors.New("baz"),
					pipeline.FiredCancel,
				),
			),
			ExpectedEvent: pipeline.FiredCancel,
			ExpectedErr: pipeline.WrapWithTerminatedError(
				errors.New("baz"),
				pipeline.FiredCancel,
			),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			t.Run("event", func(t *testing.T) {
				require.Equal(t, c.ExpectedEvent, c.GivenResult.Event())
			})

			if c.ExpectedErr != nil {
				t.Run("err", func(t *testing.T) {
					require.Error(t, c.GivenResult.Err(), c.ExpectedErr.Error())
				})
			} else {
				t.Run("no_err", func(t *testing.T) {
					require.NoError(t, c.GivenResult.Err())
				})
			}
		})
	}
}
