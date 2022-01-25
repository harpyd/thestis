package app_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
)

func TestIsParsingError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "parsing_error",
			Err:       app.NewParsingError(errors.New("decoding failed")),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("decoding failed"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, app.IsParsingError(c.Err))
		})
	}
}
