package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
	"github.com/harpyd/thestis/internal/app/command/mock"
)

func TestCreateTestCampaignHandler_Handle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Command     app.CreateTestCampaignCommand
		ShouldBeErr bool
		IsErr       func(err error) bool
	}{
		{
			Name: "create_test_campaign",
			Command: app.CreateTestCampaignCommand{
				ViewName: "test campaign",
			},
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				repo    = mock.NewTestCampaignsRepository()
				handler = command.NewCreateTestCampaignHandler(repo)
			)

			tcID, err := handler.Handle(context.Background(), c.Command)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, repo.TestCampaignsNumber())

				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, tcID)
			require.Equal(t, 1, repo.TestCampaignsNumber())
		})
	}
}
