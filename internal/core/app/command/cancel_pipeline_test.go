package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app/command"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/app/service/mock"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

func TestNewCancelPipelineHandlerPanics(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		GivenPipeRepo  service.PipelineRepository
		GivenPublisher service.PipelineCancelPublisher
		ShouldPanic    bool
		PanicMessage   string
	}{
		{
			Name:           "all_dependencies_are_not_nil",
			GivenPipeRepo:  mock.NewPipelineRepository(),
			GivenPublisher: mock.NewPipelineCancelPubsub(),
			ShouldPanic:    false,
		},
		{
			Name:           "pipeline_repository_is_nil",
			GivenPipeRepo:  nil,
			GivenPublisher: mock.NewPipelineCancelPubsub(),
			ShouldPanic:    true,
			PanicMessage:   "pipeline repository is nil",
		},
		{
			Name:           "pipeline_cancel_publisher_is_nil",
			GivenPipeRepo:  mock.NewPipelineRepository(),
			GivenPublisher: nil,
			ShouldPanic:    true,
			PanicMessage:   "pipeline cancel publisher is nil",
		},
		{
			Name:           "all_dependencies_are_nil",
			GivenPipeRepo:  nil,
			GivenPublisher: nil,
			ShouldPanic:    true,
			PanicMessage:   "pipeline repository is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = command.NewCancelPipelineHandler(
					c.GivenPipeRepo,
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

func TestHandleCancelPipeline(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Command              command.CancelPipeline
		Pipeline             *pipeline.Pipeline
		ExpectedPublishCalls int
		ShouldBeErr          bool
		IsErr                func(err error) bool
	}{
		{
			Name: "pipeline_not_found",
			Command: command.CancelPipeline{
				PipelineID:   "a64d83e5-4128-4c8b-b5ab-43b77df352ea",
				CanceledByID: "c89ba386-0976-4671-913d-9252ba29aca4",
			},
			Pipeline: pipeline.Unmarshal(pipeline.Params{
				ID:      "4abf2481-0546-4f1e-873f-b6859bbe9bf5",
				OwnerID: "c89ba386-0976-4671-913d-9252ba29aca4",
				Started: true,
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, service.ErrPipelineNotFound)
			},
			ExpectedPublishCalls: 0,
		},
		{
			Name: "user_cannot_see_pipeline",
			Command: command.CancelPipeline{
				PipelineID:   "1ada8d28-dbdc-425b-b829-dbb45cdae2b3",
				CanceledByID: "5e1484b4-90ea-4684-bf20-d597446d3eb4",
			},
			Pipeline: pipeline.Unmarshal(pipeline.Params{
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
			Name: "pipeline_not_started",
			Command: command.CancelPipeline{
				PipelineID:   "b4e252a1-7b94-46b0-84f0-40f92a6d2ee5",
				CanceledByID: "93a6224c-3788-49db-a673-ca8683a469ce",
			},
			Pipeline: pipeline.Unmarshal(pipeline.Params{
				ID:      "b4e252a1-7b94-46b0-84f0-40f92a6d2ee5",
				OwnerID: "93a6224c-3788-49db-a673-ca8683a469ce",
				Started: false,
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, pipeline.ErrNotStarted)
			},
			ExpectedPublishCalls: 0,
		},
		{
			Name: "success_pipeline_cancellation",
			Command: command.CancelPipeline{
				PipelineID:   "e0c2e511-fc31-4fc4-804b-ceb91de4179f",
				CanceledByID: "c73e888a-21f2-42c7-84f7-111c4b155be8",
			},
			Pipeline: pipeline.Unmarshal(pipeline.Params{
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
				pipeRepo     = mock.NewPipelineRepository(c.Pipeline)
				cancelPubsub = mock.NewPipelineCancelPubsub()
				handler      = command.NewCancelPipelineHandler(pipeRepo, cancelPubsub)
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
