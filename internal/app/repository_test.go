package app_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
)

func TestAsDatabaseError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError        error
		ShouldBeWrapped   bool
		ExpectedUnwrapped error
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:      app.WrapWithDatabaseError(nil),
			ShouldBeWrapped: false,
		},
		{
			GivenError:        &app.DatabaseError{},
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: nil,
		},
		{
			GivenError:        app.WrapWithDatabaseError(errors.New("foo")),
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: errors.New("foo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *app.DatabaseError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("unwrap", func(t *testing.T) {
					if c.ExpectedUnwrapped != nil {
						require.EqualError(t, target.Unwrap(), c.ExpectedUnwrapped.Error())

						return
					}

					require.NoError(t, target.Unwrap())
				})
			})
		})
	}
}

func TestFormatDatabaseError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &app.DatabaseError{},
			ExpectedErrorString: "",
		},
		{
			GivenError:          app.WrapWithDatabaseError(errors.New("failed")),
			ExpectedErrorString: "database problem: failed",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
		})
	}
}
