package user_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
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
			IsErr:       user.IsCantSeeTestCampaignError,
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
			IsErr:       user.IsCantSeeSpecificationError,
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

func TestCanSeePerformance(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		UserID      string
		Performance *performance.Performance
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name:   "can_see",
			UserID: "e7b9f695-ea8d-4d31-a4ce-cfd5521d52c2",
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "e7b9f695-ea8d-4d31-a4ce-cfd5521d52c2",
			}),
			ShouldBeErr: false,
		},
		{
			Name:   "cannot_see",
			UserID: "732563f6-08a3-4b61-bbfc-818225f58b0b",
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "f36b38e0-829d-4bdf-af2d-de4e8e43b0c0",
			}),
			ShouldBeErr: true,
			IsErr:       user.IsCantSeePerformanceError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := user.CanSeePerformance(c.UserID, c.Performance)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestIsCantSeeTestCampaignError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "cant_see_test_campaign_error",
			Err:       user.NewCantSeeTestCampaignError("user-id", "owner-id"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       user.NewCantSeeSpecificationError("user-id", "owner-id"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, user.IsCantSeeTestCampaignError(c.Err))
		})
	}
}

func TestIsCantSeeSpecificationError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "cant_see_specification_error",
			Err:       user.NewCantSeeSpecificationError("user-id", "owner-id"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       user.NewCantSeeTestCampaignError("user-id", "owner-id"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, user.IsCantSeeSpecificationError(c.Err))
		})
	}
}

func TestIsCantSeePerformanceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name      string
		Err       error
		IsSameErr bool
	}{
		{
			Name:      "cant_see_performance_error",
			Err:       user.NewCantSeePerformanceError("user-id", "owner-id"),
			IsSameErr: true,
		},
		{
			Name:      "another_error",
			Err:       errors.New("user user-id can't see performance user owner-id performance"),
			IsSameErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.IsSameErr, user.IsCantSeePerformanceError(c.Err))
		})
	}
}
