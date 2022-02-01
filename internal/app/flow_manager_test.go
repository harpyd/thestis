package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/mock"
	"github.com/harpyd/thestis/internal/domain/performance"
)

func TestEveryStepSavingFlowManager_ManageFlow(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name               string
		PerformanceFactory func() *performance.Performance
		AddPerformance     bool
		ShouldBeErr        bool
		IsErr              func(err error) bool
		Messages           []string
	}{
		{
			Name: "already_started_performance",
			PerformanceFactory: func() *performance.Performance {
				perf := performance.UnmarshalFromDatabase(performance.Params{
					OwnerID:         "f7b42682-cf52-4699-9bba-f8dac902efb0",
					SpecificationID: "73a7c5f6-f239-4abf-8837-cc4763d59d5f",
					Actions: []performance.Action{
						performance.NewActionWithoutThesis(
							"from",
							"to",
							performance.HTTPPerformer,
						),
					},
				})

				_, err := perf.Start(context.Background())
				require.NoError(t, err)

				return perf
			},
			AddPerformance: true,
			ShouldBeErr:    true,
			IsErr:          performance.IsAlreadyStartedError,
		},
		{
			Name: "exclusive_action_with_performance_failed",
			PerformanceFactory: func() *performance.Performance {
				return performance.UnmarshalFromDatabase(performance.Params{
					OwnerID:         "1baf3001-00ad-4eca-8fea-117ca68d9bc5",
					SpecificationID: "8bc587a9-b7dd-40f8-bf2f-98287518241e",
				})
			},
			AddPerformance: false,
			ShouldBeErr:    true,
			IsErr:          app.IsPerformanceNotFoundError,
		},
		{
			Name: "success_flow_managing",
			PerformanceFactory: func() *performance.Performance {
				return performance.UnmarshalFromDatabase(performance.Params{
					OwnerID:         "d1e0470e-ec44-4d57-b3eb-ef9ed8fe8f01",
					SpecificationID: "e597e3a2-54a2-4076-b1a0-1045e9aeaa7d",
					Actions: []performance.Action{
						performance.NewActionWithoutThesis(
							"first",
							"second",
							performance.HTTPPerformer,
						),
						performance.NewActionWithoutThesis(
							"second",
							"third",
							performance.AssertionPerformer,
						),
					},
				})
			},
			AddPerformance: true,
			ShouldBeErr:    false,
			Messages: []string{
				"Flow step performing `first -(HTTP)-> second`",
				"Flow step not performed `first -(HTTP)-> second`",
				"Flow step performing `second -(assertion)-> third`",
				"Flow step not performed `second -(assertion)-> third`",
			},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			perf := c.PerformanceFactory()

			var (
				perfsRepo = mock.NewPerformancesRepository()
				flowsRepo = mock.NewFlowsRepository()
			)

			if c.AddPerformance {
				err := perfsRepo.AddPerformance(context.Background(), perf)
				require.NoError(t, err)
			}

			fm := app.NewEveryStepSavingFlowManager(perfsRepo, flowsRepo)

			messages, err := fm.ManageFlow(context.Background(), perf)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)

			readMessages := make([]string, 0, len(c.Messages))

			for msg := range messages {
				readMessages = append(readMessages, msg.String())
			}

			require.ElementsMatch(t, c.Messages, readMessages)
		})
	}
}
