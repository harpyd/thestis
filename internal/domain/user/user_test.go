package user_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
	"github.com/harpyd/thestis/internal/domain/user"
)

func TestCanSeeTestCampaign(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name         string
		UserID       string
		TestCampaign *testcampaign.TestCampaign
		ShouldBeErr  bool
		IsErr        func(err error) bool
	}{
		{
			Name:   "can_see",
			UserID: "6e36006d-ec59-40fc-b132-08b2cdd28fc6",
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:       "10fe3333-e6be-498d-a788-5b48aab998cf",
				OwnerID:  "6e36006d-ec59-40fc-b132-08b2cdd28fc6",
				ViewName: "test",
				Summary:  "tests",
			}),
			ShouldBeErr: false,
		},
		{
			Name:   "cannot_see",
			UserID: "3ef42ce3-b4e2-4c64-b41c-92d881d44658",
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:      "c316d7d8-28df-4bce-b28a-80a4364c8c07",
				OwnerID: "f4f560a1-138c-4812-b152-0d7b71236d7f",
			}),
			ShouldBeErr: true,
			IsErr:       user.IsUserCantSeeTestCampaignError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := user.CanSeeTestCampaign(c.UserID, c.TestCampaign)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestCanSeeSpecification(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		UserID        string
		Specification *specification.Specification
		ShouldBeErr   bool
		IsErr         func(err error) bool
	}{
		{
			Name:   "can_see",
			UserID: "2e7f4a0b-a23a-4020-9138-756912b705bd",
			Specification: specification.NewBuilder().
				WithID("17099c59-19c5-4edf-ac0f-b0093fed1ffc").
				WithOwnerID("2e7f4a0b-a23a-4020-9138-756912b705bd").
				ErrlessBuild(),
			ShouldBeErr: false,
		},
		{
			Name:   "cannot_see",
			UserID: "20002cdc-bc58-4e38-a0e0-c46bbfba76e5",
			Specification: specification.NewBuilder().
				WithID("772f14d3-640a-4e57-82fd-047f426ac623").
				WithOwnerID("4e28e13b-877f-4e53-bc85-0744164b7187").
				ErrlessBuild(),
			ShouldBeErr: true,
			IsErr:       user.IsUserCantSeeSpecificationError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := user.CanSeeSpecification(c.UserID, c.Specification)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestIsUserCantSeeTestCampaignError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "user_cant_see_test_campaign_error",
			Err:       user.NewUserCantSeeTestCampaignError("user-id", "owner-id"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       user.NewUserCantSeeSpecificationError("user-id", "owner-id"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, user.IsUserCantSeeTestCampaignError(c.Err))
		})
	}
}

func TestIsUserCantSeeSpecificationError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "user_cant_see_specification_error",
			Err:       user.NewUserCantSeeSpecificationError("user-id", "owner-id"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       user.NewUserCantSeeTestCampaignError("user-id", "owner-id"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, user.IsUserCantSeeSpecificationError(c.Err))
		})
	}
}
