package testcampaign_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

func TestTestCampaign(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Params      testcampaign.Params
		ShouldBeErr bool
		ExpectedErr error
	}{
		{
			Name: "without_error",
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
			ExpectedErr: testcampaign.ErrEmptyID,
		},
		{
			Name: "empty_user_id",
			Params: testcampaign.Params{
				ID:       "tc-id",
				ViewName: "some name",
				Summary:  "",
			},
			ShouldBeErr: true,
			ExpectedErr: testcampaign.ErrEmptyOwnerID,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			tc, err := testcampaign.New(c.Params)

			must := func() {
				_ = testcampaign.MustNew(c.Params)
			}

			if c.ShouldBeErr {
				t.Run("err", func(t *testing.T) {
					t.Run("is", func(t *testing.T) {
						require.ErrorIs(t, err, c.ExpectedErr)
					})

					t.Run("panic", func(t *testing.T) {
						require.PanicsWithValue(t, c.ExpectedErr, must)
					})
				})

				return
			}

			t.Run("no_err", func(t *testing.T) {
				require.NoError(t, err)
				require.NotPanics(t, must)

				t.Run("id", func(t *testing.T) {
					require.Equal(t, c.Params.ID, tc.ID())
				})

				t.Run("owner_id", func(t *testing.T) {
					require.Equal(t, c.Params.OwnerID, tc.OwnerID())
				})

				t.Run("view_name", func(t *testing.T) {
					require.Equal(t, c.Params.ViewName, tc.ViewName())
				})

				t.Run("summary", func(t *testing.T) {
					require.Equal(t, c.Params.Summary, tc.Summary())
				})

				t.Run("created_at", func(t *testing.T) {
					require.Equal(t, c.Params.CreatedAt, tc.CreatedAt())
				})
			})
		})
	}
}

func TestSetTestCampaignViewName(t *testing.T) {
	t.Parallel()

	tc := testcampaign.MustNew(testcampaign.Params{
		ID:       "id",
		OwnerID:  "owner-id",
		ViewName: "foo",
	})

	require.Equal(t, "foo", tc.ViewName())

	tc.SetViewName("bar")

	require.Equal(t, "bar", tc.ViewName())
}

func TestSetTestCampaignSummary(t *testing.T) {
	t.Parallel()

	tc := testcampaign.MustNew(testcampaign.Params{
		ID:      "id",
		OwnerID: "owner-id",
		Summary: "doo",
	})

	require.Equal(t, "doo", tc.Summary())

	tc.SetSummary("qoo")

	require.Equal(t, "qoo", tc.Summary())
}
