package testcampaign_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Params      testcampaign.Params
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "new_without_error",
			Params: testcampaign.Params{
				ID:       "tc-id",
				ViewName: "test campaign",
				Summary:  "summary",
				UserID:   "user-id",
			},
			ShouldBeErr: false,
		},
		{
			Name: "empty_test_campaign_id",
			Params: testcampaign.Params{
				ID:       "",
				ViewName: "test campaign with empty id",
				Summary:  "",
				UserID:   "user-id",
			},
			ShouldBeErr: true,
			IsErr:       testcampaign.IsEmptyIDError,
		},
		{
			Name: "empty_user_id",
			Params: testcampaign.Params{
				ID:       "tc-id",
				ViewName: "some name",
				Summary:  "",
			},
			ShouldBeErr: true,
			IsErr:       testcampaign.IsEmptyUserIDError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			tc, err := testcampaign.New(c.Params)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, c.Params.ID, tc.ID())
			require.Equal(t, c.Params.ViewName, tc.ViewName())
			require.Equal(t, c.Params.Summary, tc.Summary())
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

func TestIsEmptyUserIDError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "empty_user_id_error",
			Err:       testcampaign.NewEmptyUserIDError(),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("another error"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, testcampaign.IsEmptyUserIDError(c.Err))
		})
	}
}
