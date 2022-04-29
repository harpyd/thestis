package user_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/domain/performance"
	"github.com/harpyd/thestis/internal/core/domain/specification"
	"github.com/harpyd/thestis/internal/core/domain/testcampaign"
	"github.com/harpyd/thestis/internal/core/domain/user"
)

func TestCanAccessTestCampaign(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name         string
		UserID       string
		TestCampaign *testcampaign.TestCampaign
		Permission   user.Permission
		ShouldBeErr  bool
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
			Permission:  user.Read,
			ShouldBeErr: false,
		},
		{
			Name:   "cannot_see",
			UserID: "3ef42ce3-b4e2-4c64-b41c-92d881d44658",
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:      "c316d7d8-28df-4bce-b28a-80a4364c8c07",
				OwnerID: "f4f560a1-138c-4812-b152-0d7b71236d7f",
			}),
			Permission:  user.Write,
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := user.CanAccessTestCampaign(
				c.UserID,
				c.TestCampaign,
				c.Permission,
			)

			if c.ShouldBeErr {
				var target *user.AccessError

				require.ErrorAs(t, err, &target)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestCanAccessSpecification(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		UserID        string
		Specification *specification.Specification
		Permission    user.Permission
		ShouldBeErr   bool
	}{
		{
			Name:   "can_see",
			UserID: "2e7f4a0b-a23a-4020-9138-756912b705bd",
			Specification: (&specification.Builder{}).
				WithID("17099c59-19c5-4edf-ac0f-b0093fed1ffc").
				WithOwnerID("2e7f4a0b-a23a-4020-9138-756912b705bd").
				ErrlessBuild(),
			Permission:  user.Write,
			ShouldBeErr: false,
		},
		{
			Name:   "cannot_see",
			UserID: "20002cdc-bc58-4e38-a0e0-c46bbfba76e5",
			Specification: (&specification.Builder{}).
				WithID("772f14d3-640a-4e57-82fd-047f426ac623").
				WithOwnerID("4e28e13b-877f-4e53-bc85-0744164b7187").
				ErrlessBuild(),
			ShouldBeErr: true,
			Permission:  user.Read,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := user.CanAccessSpecification(
				c.UserID,
				c.Specification,
				c.Permission,
			)

			if c.ShouldBeErr {
				var target *user.AccessError

				require.ErrorAs(t, err, &target)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestCanAccessPerformance(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		UserID      string
		Performance *performance.Performance
		Permission  user.Permission
		ShouldBeErr bool
	}{
		{
			Name:   "can_see",
			UserID: "e7b9f695-ea8d-4d31-a4ce-cfd5521d52c2",
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "e7b9f695-ea8d-4d31-a4ce-cfd5521d52c2",
			}),
			Permission:  user.Read,
			ShouldBeErr: false,
		},
		{
			Name:   "cannot_see",
			UserID: "732563f6-08a3-4b61-bbfc-818225f58b0b",
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "f36b38e0-829d-4bdf-af2d-de4e8e43b0c0",
			}),
			Permission:  user.Write,
			ShouldBeErr: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := user.CanAccessPerformance(
				c.UserID,
				c.Performance,
				c.Permission,
			)

			if c.ShouldBeErr {
				var target *user.AccessError

				require.ErrorAs(t, err, &target)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestAsAccessError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError         error
		ShouldBeWrapped    bool
		ExpectedUserID     string
		ExpectedResourceID string
		ExpectedResource   user.Resource
		ExpectedPermission user.Permission
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:         &user.AccessError{},
			ShouldBeWrapped:    true,
			ExpectedUserID:     "",
			ExpectedResourceID: "",
			ExpectedResource:   user.NoResource,
			ExpectedPermission: user.NoPermission,
		},
		{
			GivenError: user.NewAccessError(
				"",
				"",
				user.NoResource,
				user.NoPermission,
			),
			ShouldBeWrapped:    true,
			ExpectedUserID:     "",
			ExpectedResourceID: "",
			ExpectedResource:   user.NoResource,
			ExpectedPermission: user.NoPermission,
		},
		{
			GivenError: user.NewAccessError(
				"foo",
				"boo",
				user.Specification,
				user.Read,
			),
			ShouldBeWrapped:    true,
			ExpectedUserID:     "foo",
			ExpectedResourceID: "boo",
			ExpectedResource:   user.Specification,
			ExpectedPermission: user.Read,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *user.AccessError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("user_id", func(t *testing.T) {
					require.Equal(t, c.ExpectedUserID, target.UserID())
				})

				t.Run("resource_id", func(t *testing.T) {
					require.Equal(t, c.ExpectedResourceID, target.ResourceID())
				})

				t.Run("resource", func(t *testing.T) {
					require.Equal(t, c.ExpectedResource, target.Resource())
				})

				t.Run("permission", func(t *testing.T) {
					require.Equal(t, c.ExpectedPermission, target.Permission())
				})
			})
		})
	}
}

func TestFormatAccessError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &user.AccessError{},
			ExpectedErrorString: "can't access",
		},
		{
			GivenError: user.NewAccessError(
				"foo",
				"",
				user.NoResource,
				user.NoPermission,
			),
			ExpectedErrorString: "user #foo can't access",
		},
		{
			GivenError: user.NewAccessError(
				"",
				"bar",
				user.NoResource,
				user.NoPermission,
			),
			ExpectedErrorString: "can't access #bar",
		},
		{
			GivenError: user.NewAccessError(
				"",
				"",
				user.TestCampaign,
				user.NoPermission,
			),
			ExpectedErrorString: "can't access test campaign",
		},
		{
			GivenError: user.NewAccessError(
				"",
				"",
				user.NoResource,
				user.Read,
			),
			ExpectedErrorString: `can't access with "read" permission`,
		},
		{
			GivenError: user.NewAccessError(
				"foo",
				"bar",
				user.TestCampaign,
				user.Read,
			),
			ExpectedErrorString: `user #foo can't access test campaign #bar with "read" permission`,
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
