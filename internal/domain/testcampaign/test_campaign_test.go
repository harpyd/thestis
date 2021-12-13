package testcampaign_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type params struct {
		ID       string
		ViewName string
	}

	testCases := []struct {
		Name        string
		Params      params
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "new_without_error",
			Params: params{
				ID:       "tc-id",
				ViewName: "test campaign",
			},
			ShouldBeErr: false,
		},
		{
			Name: "empty_test_campaign_id",
			Params: params{
				ID:       "",
				ViewName: "test campaign with empty id",
			},
			ShouldBeErr: true,
			IsErr:       testcampaign.IsEmptyIDError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			tc, err := testcampaign.New(c.Params.ID, c.Params.ViewName)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, c.Params.ID, tc.ID())
			require.Equal(t, c.Params.ViewName, tc.ViewName())
		})
	}
}

func TestIsEmptyIDError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "empty_id_error",
			Err:       testcampaign.NewEmptyIDError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("some err"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, testcampaign.IsEmptyIDError(c.Err))
		})
	}
}
