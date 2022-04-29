package app_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app"
)

func TestAsParseError(t *testing.T) {
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
			GivenError:      app.WrapWithParseError(nil),
			ShouldBeWrapped: false,
		},
		{
			GivenError:        &app.ParseError{},
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: nil,
		},
		{
			GivenError:        app.WrapWithParseError(errors.New("foo")),
			ShouldBeWrapped:   true,
			ExpectedUnwrapped: errors.New("foo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *app.ParseError

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

func TestFormatParseError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &app.ParseError{},
			ExpectedErrorString: "",
		},
		{
			GivenError:          app.WrapWithParseError(errors.New("foo")),
			ExpectedErrorString: "parsing specification: foo",
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
