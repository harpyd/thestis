package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app/command"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/app/service/mock"
)

func TestNewCreateTestCampaignHandlerPanics(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                  string
		GivenTestCampaignRepo service.TestCampaignRepository
		ShouldPanic           bool
		PanicMessage          string
	}{
		{
			Name:                  "all_dependencies_are_not_nil",
			GivenTestCampaignRepo: mock.NewTestCampaignRepository(),
			ShouldPanic:           false,
		},
		{
			Name:                  "all_dependencies_are_nil",
			GivenTestCampaignRepo: nil,
			ShouldPanic:           true,
			PanicMessage:          "test campaign repository is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = command.NewCreateTestCampaignHandler(c.GivenTestCampaignRepo)
			}

			if !c.ShouldPanic {
				require.NotPanics(t, init)

				return
			}

			require.PanicsWithValue(t, c.PanicMessage, init)
		})
	}
}

func TestHandleCreateTestCampaign(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Command     command.CreateTestCampaign
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "create_test_campaign",
			Command: command.CreateTestCampaign{
				TestCampaignID: "774fcdd9-3b8d-41c6-b7ee-a544a4247e30",
				OwnerID:        "61fcde9c-b729-4ae1-9c86-a80d706eda6c",
				ViewName:       "test campaign",
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				repo    = mock.NewTestCampaignRepository()
				handler = command.NewCreateTestCampaignHandler(repo)
			)

			err := handler.Handle(context.Background(), c.Command)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, repo.TestCampaignsNumber())

				return
			}

			require.NoError(t, err)
			require.Equal(t, 1, repo.TestCampaignsNumber())
		})
	}
}
