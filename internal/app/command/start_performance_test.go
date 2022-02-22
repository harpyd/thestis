package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
	appMock "github.com/harpyd/thestis/internal/app/mock"
	perfMock "github.com/harpyd/thestis/internal/domain/performance/mock"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/user"
)

func TestStartPerformanceHandler_Handle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		Command       app.StartPerformanceCommand
		Specification *specification.Specification
		ShouldBeErr   bool
		IsErr         func(err error) bool
	}{
		{
			Name: "specification_with_such_test_campaign_id_not_found",
			Command: app.StartPerformanceCommand{
				TestCampaignID: "68baf422-777f-4a0e-b35a-4fff5858af2d",
				StartedByID:    "d8d1e4ab-8f24-4c79-a1f2-49e24b3f119a",
			},
			Specification: specification.NewBuilder().
				WithTestCampaignID("d5a7b2ec-c04e-40d8-a2b5-b273d7ad7ffd").
				WithOwnerID("d8d1e4ab-8f24-4c79-a1f2-49e24b3f119a").
				ErrlessBuild(),
			ShouldBeErr: true,
			IsErr:       app.IsSpecificationNotFoundError,
		},
		{
			Name: "user_cannot_see_specification",
			Command: app.StartPerformanceCommand{
				TestCampaignID: "5ee6228e-5b0b-4d40-b4e5-9a138bef9f84",
				StartedByID:    "fb883739-2c8c-4a4e-bca2-f96b204f4ac8",
			},
			Specification: specification.NewBuilder().
				WithTestCampaignID("5ee6228e-5b0b-4d40-b4e5-9a138bef9f84").
				WithOwnerID("8ea9dca1-53da-4ed5-8f4b-660c8956ea45").
				ErrlessBuild(),
			ShouldBeErr: true,
			IsErr:       user.IsCantSeeSpecificationError,
		},
		{
			Name: "success_performance_starting",
			Command: app.StartPerformanceCommand{
				TestCampaignID: "70c8e87d-395d-4ae6-b53e-3b2f587039a3",
				StartedByID:    "aa584d3d-c790-4ed3-8bfa-19e1b6fed88e",
			},
			Specification: specification.NewBuilder().
				WithTestCampaignID("70c8e87d-395d-4ae6-b53e-3b2f587039a3").
				WithOwnerID("aa584d3d-c790-4ed3-8bfa-19e1b6fed88e").
				ErrlessBuild(),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				specsRepo  = appMock.NewSpecificationsRepository(c.Specification)
				perfsRepo  = appMock.NewPerformancesRepository()
				maintainer = appMock.NewPerformanceMaintainer(false)
				handler    = command.NewStartPerformanceHandler(
					specsRepo,
					perfsRepo,
					maintainer,
					app.WithHTTPPerformer(perfMock.NewPassingPerformer()),
					app.WithAssertionPerformer(perfMock.NewFailingPerformer()),
				)
			)

			ctx := context.Background()

			perfID, messages, err := handler.Handle(ctx, c.Command)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, perfsRepo.PerformancesNumber())

				return
			}

			require.NoError(t, err)

			require.NotEmpty(t, perfID)
			require.NotNil(t, messages)
			require.Equal(t, 1, perfsRepo.PerformancesNumber())
		})
	}
}
