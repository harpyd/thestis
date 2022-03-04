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

func TestHandleRestartPerformance(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                      string
		Command                   app.RestartPerformanceCommand
		Performance               *performance.Performance
		PerformanceAlreadyStarted bool
		ShouldBeErr               bool
		IsErr                     func(err error) bool
	}{
		{
			Name: "performance_not_found",
			Command: app.RestartPerformanceCommand{
				PerformanceID: "9524edfe-9a47-40d5-9c40-5c575de1368b",
				StartedByID:   "7177997a-b63d-4e1b-9288-0a581f7ff03a",
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID:         "7177997a-b63d-4e1b-9288-0a581f7ff03a",
				SpecificationID: "64be3dee-a11b-451c-91f3-fddcc4ab7ced",
			}, performance.WithID("ac622aa2-59e9-4389-97cc-2309271f1f36")),
			ShouldBeErr: true,
			IsErr:       app.IsPerformanceNotFoundError,
		},
		{
			Name: "user_cannot_see_performance",
			Command: app.RestartPerformanceCommand{
				PerformanceID: "50b59781-6932-422e-a7e6-f7424b5f5d36",
				StartedByID:   "be746c97-fe6e-4795-820a-97337e8d98b2",
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "0a5782ec-7580-4526-b4e4-c3a7489ca512",
			}, performance.WithID("50b59781-6932-422e-a7e6-f7424b5f5d36")),
			ShouldBeErr: true,
			IsErr:       user.IsCantSeePerformanceError,
		},
		{
			Name: "performance_already_started",
			Command: app.RestartPerformanceCommand{
				PerformanceID: "6f112cf1-3dd5-4f14-a5ef-7ef18dfb8921",
				StartedByID:   "960f7ba1-b16c-43eb-9f87-d367ec9e0ba9",
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "960f7ba1-b16c-43eb-9f87-d367ec9e0ba9",
			}, performance.WithID("6f112cf1-3dd5-4f14-a5ef-7ef18dfb8921")),
			PerformanceAlreadyStarted: true,
			ShouldBeErr:               true,
			IsErr:                     performance.IsAlreadyStartedError,
		},
		{
			Name: "success_performance_restarting",
			Command: app.RestartPerformanceCommand{
				PerformanceID: "fc2f14b3-6125-47fa-a343-5fabcac9abd1",
				StartedByID:   "5da02570-a192-4a9a-9180-1a2704732b06",
			},
			Performance: performance.Unmarshal(performance.Params{
				OwnerID: "5da02570-a192-4a9a-9180-1a2704732b06",
			}, performance.WithID("fc2f14b3-6125-47fa-a343-5fabcac9abd1")),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				perfsRepo  = mock.NewPerformancesRepository(c.Performance)
				maintainer = mock.NewPerformanceMaintainer(c.PerformanceAlreadyStarted)
				handler    = command.NewRestartPerformanceHandler(
					perfsRepo,
					maintainer,
					app.WithHTTPPerformer(performance.PassingPerformer()),
					app.WithAssertionPerformer(performance.FailingPerformer()),
				)
			)

			ctx := context.Background()

			messages, err := handler.Handle(ctx, c.Command)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)

			require.NotNil(t, messages)
		})
	}
}
