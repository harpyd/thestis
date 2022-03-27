package performance_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
)

func TestResult(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenResult   performance.Result
		ExpectedEvent performance.Event
		ExpectedErr   error
	}{
		{
			GivenResult:   performance.Pass(),
			ExpectedEvent: performance.FiredPass,
		},
		{
			GivenResult:   performance.Fail(nil),
			ExpectedEvent: performance.FiredFail,
			ExpectedErr:   nil,
		},
		{
			GivenResult: performance.Fail(
				errors.New("foo"),
			),
			ExpectedEvent: performance.FiredFail,
			ExpectedErr: performance.WrapWithTerminatedError(
				errors.New("foo"),
				performance.FiredFail,
			),
		},
		{
			GivenResult: performance.Fail(
				performance.WrapWithTerminatedError(
					errors.New("foo"),
					performance.FiredFail,
				),
			),
			ExpectedEvent: performance.FiredFail,
			ExpectedErr: performance.WrapWithTerminatedError(
				errors.New("bar"),
				performance.FiredFail,
			),
		},
		{
			GivenResult:   performance.Crash(nil),
			ExpectedEvent: performance.FiredCrash,
			ExpectedErr:   nil,
		},
		{
			GivenResult: performance.Crash(
				errors.New("boo"),
			),
			ExpectedEvent: performance.FiredCrash,
			ExpectedErr: performance.WrapWithTerminatedError(
				errors.New("boo"),
				performance.FiredCrash,
			),
		},
		{
			GivenResult: performance.Crash(
				performance.WrapWithTerminatedError(
					errors.New("qwe"),
					performance.FiredCrash,
				),
			),
			ExpectedEvent: performance.FiredCrash,
			ExpectedErr: performance.WrapWithTerminatedError(
				errors.New("qwe"),
				performance.FiredCrash,
			),
		},
		{
			GivenResult:   performance.Cancel(nil),
			ExpectedEvent: performance.FiredCancel,
			ExpectedErr:   nil,
		},
		{
			GivenResult: performance.Cancel(
				errors.New("bar"),
			),
			ExpectedEvent: performance.FiredCancel,
			ExpectedErr: performance.WrapWithTerminatedError(
				errors.New("bar"),
				performance.FiredCancel,
			),
		},
		{
			GivenResult: performance.Cancel(
				performance.WrapWithTerminatedError(
					errors.New("baz"),
					performance.FiredCancel,
				),
			),
			ExpectedEvent: performance.FiredCancel,
			ExpectedErr: performance.WrapWithTerminatedError(
				errors.New("baz"),
				performance.FiredCancel,
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
