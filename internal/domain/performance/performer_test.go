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
			ExpectedErr: performance.NewFailedError(
				errors.New("foo"),
			),
		},
		{
			GivenResult: performance.Fail(
				performance.NewFailedError(
					errors.New("foo"),
				),
			),
			ExpectedEvent: performance.FiredFail,
			ExpectedErr: performance.NewFailedError(
				errors.New("foo"),
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
			ExpectedErr: performance.NewCrashedError(
				errors.New("foo"),
			),
		},
		{
			GivenResult: performance.Crash(
				performance.NewCrashedError(
					errors.New("boo"),
				),
			),
			ExpectedEvent: performance.FiredCrash,
			ExpectedErr: performance.NewCrashedError(
				errors.New("boo"),
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
			ExpectedErr: performance.NewCanceledError(
				errors.New("bar"),
			),
		},
		{
			GivenResult: performance.Cancel(
				performance.NewCanceledError(
					errors.New("bar"),
				),
			),
			ExpectedEvent: performance.FiredCancel,
			ExpectedErr: performance.NewCanceledError(
				errors.New("bar"),
			),
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
