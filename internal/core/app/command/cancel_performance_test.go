package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/app/command"
	"github.com/harpyd/thestis/internal/core/app/mock"
	"github.com/harpyd/thestis/internal/core/domain/performance"
	"github.com/harpyd/thestis/internal/core/domain/user"
)

func TestPanickingNewCancelPerformanceHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		GivenPerfRepo  app.PerformanceRepository
		GivenPublisher app.PerformanceCancelPublisher
		ShouldPanic    bool
		PanicMessage   string
	}{
		{
			Name:           "all_dependencies_are_not_nil",
			GivenPerfRepo:  mock.NewPerformanceRepository(),
			GivenPublisher: mock.NewPerformanceCancelPubsub(),
			ShouldPanic:    false,
		},
		{
			Name:           "performance_repository_is_nil",
			GivenPerfRepo:  nil,
			GivenPublisher: mock.NewPerformanceCancelPubsub(),
			ShouldPanic:    true,
			PanicMessage:   "performance repository is nil",
		},
		{
			Name:           "performance_cancel_publisher_is_nil",
			GivenPerfRepo:  mock.NewPerformanceRepository(),
			GivenPublisher: nil,
			ShouldPanic:    true,
			PanicMessage:   "performance cancel publisher is nil",
		},
		{
			Name:           "all_dependencies_are_nil",
			GivenPerfRepo:  nil,
			GivenPublisher: nil,
			ShouldPanic:    true,
			PanicMessage:   "performance repository is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = command.NewCancelPerformanceHandler(
					c.GivenPerfRepo,
					c.GivenPublisher,
				)
			}

			if !c.ShouldPanic {
				require.NotPanics(t, init)

				return
			}

			require.PanicsWithValue(t, c.PanicMessage, init)
		})
	}
}

func TestHandleCancelPerformance(t *testing.T) {
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
				ID:      "4abf2481-0546-4f1e-873f-b6859bbe9bf5",
				OwnerID: "c89ba386-0976-4671-913d-9252ba29aca4",
				Started: true,
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrPerformanceNotFound)
			},
			ExpectedPublishCalls: 0,
		},
		{
			Name: "user_cannot_see_performance",
			Command: app.CancelPerformanceCommand{
				PerformanceID: "1ada8d28-dbdc-425b-b829-dbb45cdae2b3",
				CanceledByID:  "5e1484b4-90ea-4684-bf20-d597446d3eb4",
			},
			Performance: performance.Unmarshal(performance.Params{
				ID:      "1ada8d28-dbdc-425b-b829-dbb45cdae2b3",
				OwnerID: "759cf65b-547b-4523-a9f4-9dd4f12188d2",
				Started: true,
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *user.AccessError

				return errors.As(err, &target)
			},
			ExpectedPublishCalls: 0,
		},
		{
			Name: "performance_not_started",
			Command: app.CancelPerformanceCommand{
				PerformanceID: "b4e252a1-7b94-46b0-84f0-40f92a6d2ee5",
				CanceledByID:  "93a6224c-3788-49db-a673-ca8683a469ce",
			},
			Performance: performance.Unmarshal(performance.Params{
				ID:      "b4e252a1-7b94-46b0-84f0-40f92a6d2ee5",
				OwnerID: "93a6224c-3788-49db-a673-ca8683a469ce",
				Started: false,
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, performance.ErrNotStarted)
			},
			ExpectedPublishCalls: 0,
		},
		{
			Name: "success_performance_cancelation",
			Command: app.CancelPerformanceCommand{
				PerformanceID: "e0c2e511-fc31-4fc4-804b-ceb91de4179f",
				CanceledByID:  "c73e888a-21f2-42c7-84f7-111c4b155be8",
			},
			Performance: performance.Unmarshal(performance.Params{
				ID:      "e0c2e511-fc31-4fc4-804b-ceb91de4179f",
				OwnerID: "c73e888a-21f2-42c7-84f7-111c4b155be8",
				Started: true,
			}),
			ShouldBeErr:          false,
			ExpectedPublishCalls: 1,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				perfRepo     = mock.NewPerformanceRepository(c.Performance)
				cancelPubsub = mock.NewPerformanceCancelPubsub()
				handler      = command.NewCancelPerformanceHandler(perfRepo, cancelPubsub)
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
