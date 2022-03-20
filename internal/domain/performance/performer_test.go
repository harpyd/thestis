package performance_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

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
			ExpectedErr: &performance.TerminatedError{
				Event: performance.FiredFail,
				Err:   errors.New("foo"),
			},
		},
		{
			GivenResult: performance.Fail(
				&performance.TerminatedError{
					Event: performance.FiredFail,
					Err:   errors.New("bar"),
				},
			),
			ExpectedEvent: performance.FiredFail,
			ExpectedErr: &performance.TerminatedError{
				Event: performance.FiredFail,
				Err:   errors.New("bar"),
			},
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
			ExpectedErr: &performance.TerminatedError{
				Event: performance.FiredCrash,
				Err:   errors.New("boo"),
			},
		},
		{
			GivenResult: performance.Crash(
				&performance.TerminatedError{
					Event: performance.FiredCrash,
					Err:   errors.New("qwe"),
				},
			),
			ExpectedEvent: performance.FiredCrash,
			ExpectedErr: &performance.TerminatedError{
				Event: performance.FiredCrash,
				Err:   errors.New("qwe"),
			},
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
			ExpectedErr: &performance.TerminatedError{
				Event: performance.FiredCancel,
				Err:   errors.New("bar"),
			},
		},
		{
			GivenResult: performance.Cancel(
				&performance.TerminatedError{
					Event: performance.FiredCancel,
					Err:   errors.New("baz"),
				},
			),
			ExpectedEvent: performance.FiredCancel,
			ExpectedErr: &performance.TerminatedError{
				Event: performance.FiredCancel,
				Err:   errors.New("baz"),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			t.Run("event", func(t *testing.T) {
				assert.Equal(t, c.ExpectedEvent, c.GivenResult.Event())
			})

			if c.ExpectedErr != nil {
				t.Run("err", func(t *testing.T) {
					assert.Error(t, c.GivenResult.Err(), c.ExpectedErr.Error())
				})
			} else {
				t.Run("no_err", func(t *testing.T) {
					assert.NoError(t, c.GivenResult.Err())
				})
			}
		})
	}
}
