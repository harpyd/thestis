package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
	"github.com/harpyd/thestis/internal/app/mock"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/user"
)

func TestCancelPerformanceHandler_Handle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Command              app.CancelPerformanceCommand
		Performance          *performance.Performance
		ExpectedPublishCalls int
		ShouldBeErr          bool
		IsErr                func(err error) bool
	}{
		{
			Name: "performance_not_found",
			Command: app.CancelPerformanceCommand{
				PerformanceID: "a64d83e5-4128-4c8b-b5ab-43b77df352ea",
				CanceledByID:  "c89ba386-0976-4671-913d-9252ba29aca4",
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "c89ba386-0976-4671-913d-9252ba29aca4",
				Started: true,
			}, performance.WithID("4abf2481-0546-4f1e-873f-b6859bbe9bf5")),
			ShouldBeErr:          true,
			IsErr:                app.IsPerformanceNotFoundError,
			ExpectedPublishCalls: 0,
		},
		{
			Name: "user_cannot_see_performance",
			Command: app.CancelPerformanceCommand{
				PerformanceID: "1ada8d28-dbdc-425b-b829-dbb45cdae2b3",
				CanceledByID:  "5e1484b4-90ea-4684-bf20-d597446d3eb4",
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "759cf65b-547b-4523-a9f4-9dd4f12188d2",
				Started: true,
			}, performance.WithID("1ada8d28-dbdc-425b-b829-dbb45cdae2b3")),
			ShouldBeErr:          true,
			IsErr:                user.IsCantSeePerformanceError,
			ExpectedPublishCalls: 0,
		},
		{
			Name: "performance_not_started",
			Command: app.CancelPerformanceCommand{
				PerformanceID: "b4e252a1-7b94-46b0-84f0-40f92a6d2ee5",
				CanceledByID:  "93a6224c-3788-49db-a673-ca8683a469ce",
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "93a6224c-3788-49db-a673-ca8683a469ce",
				Started: false,
			}, performance.WithID("b4e252a1-7b94-46b0-84f0-40f92a6d2ee5")),
			ShouldBeErr:          true,
			IsErr:                performance.IsNotStartedError,
			ExpectedPublishCalls: 0,
		},
		{
			Name: "success_performance_cancelation",
			Command: app.CancelPerformanceCommand{
				PerformanceID: "e0c2e511-fc31-4fc4-804b-ceb91de4179f",
				CanceledByID:  "c73e888a-21f2-42c7-84f7-111c4b155be8",
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "c73e888a-21f2-42c7-84f7-111c4b155be8",
				Started: true,
			}, performance.WithID("e0c2e511-fc31-4fc4-804b-ceb91de4179f")),
			ShouldBeErr:          false,
			ExpectedPublishCalls: 1,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				perfsRepo    = mock.NewPerformancesRepository(c.Performance)
				cancelPubsub = mock.NewPerformanceCancelPubsub()
				handler      = command.NewCancelPerformanceHandler(perfsRepo, cancelPubsub)
			)

			err := handler.Handle(context.Background(), c.Command)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, c.ExpectedPublishCalls, cancelPubsub.PublishCalls())
		})
	}
}