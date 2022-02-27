package testcampaign_test

import (
	"testing"
	"time"

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
				ID:        "tc-id",
				ViewName:  "test campaign",
				Summary:   "summary",
				OwnerID:   "user-id",
				CreatedAt: time.Now(),
			},
			ShouldBeErr: false,
		},
		{
			Name: "empty_test_campaign_id",
			Params: testcampaign.Params{
				ID:       "",
				ViewName: "test campaign with empty id",
				Summary:  "",
				OwnerID:  "user-id",
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
			IsErr:       testcampaign.IsEmptyOwnerIDError,
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
			require.Equal(t, c.Params.CreatedAt, tc.CreatedAt())
		})
	}
}

func TestTestCampaignErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "empty_id_error",
			Err:   testcampaign.NewEmptyIDError(),
			IsErr: testcampaign.IsEmptyIDError,
		},
		{
			Name:     "NON_empty_id_error",
			Err:      errors.New("empty id"),
			IsErr:    testcampaign.IsEmptyIDError,
			Reversed: true,
		},
		{
			Name:  "empty_user_id_error",
			Err:   testcampaign.NewEmptyOwnerIDError(),
			IsErr: testcampaign.IsEmptyOwnerIDError,
		},
		{
			Name:     "NON_empty_user_id_error",
			Err:      errors.New("empty owner id"),
			IsErr:    testcampaign.IsEmptyOwnerIDError,
			Reversed: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			if c.Reversed {
				require.False(t, c.IsErr(c.Err))

				return
			}

			require.True(t, c.IsErr(c.Err))
		})
	}
}
