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
	"github.com/harpyd/thestis/internal/core/entity/specification"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

func TestNewStartPipelineHandlerPanics(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		GivenSpecRepo   service.SpecificationRepository
		GivenPipeRepo   service.PipelineRepository
		GivenMaintainer service.PipelineMaintainer
		ShouldPanic     bool
		PanicMessage    string
	}{
		{
			Name:            "all_dependencies_are_not_nil",
			GivenSpecRepo:   mock.NewSpecificationRepository(),
			GivenPipeRepo:   mock.NewPipelineRepository(),
			GivenMaintainer: mock.NewPipelineMaintainer(false),
			ShouldPanic:     false,
		},
		{
			Name:            "specification_repository_is_nil",
			GivenSpecRepo:   nil,
			GivenPipeRepo:   mock.NewPipelineRepository(),
			GivenMaintainer: mock.NewPipelineMaintainer(false),
			ShouldPanic:     true,
			PanicMessage:    "specification repository is nil",
		},
		{
			Name:            "pipeline_repository_is_nil",
			GivenSpecRepo:   mock.NewSpecificationRepository(),
			GivenPipeRepo:   nil,
			GivenMaintainer: mock.NewPipelineMaintainer(false),
			ShouldPanic:     true,
			PanicMessage:    "pipeline repository is nil",
		},
		{
			Name:            "pipeline_maintainer_is_nil",
			GivenSpecRepo:   mock.NewSpecificationRepository(),
			GivenPipeRepo:   mock.NewPipelineRepository(),
			GivenMaintainer: nil,
			ShouldPanic:     true,
			PanicMessage:    "pipeline maintainer is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = command.NewStartPipelineHandler(
					c.GivenSpecRepo,
					c.GivenPipeRepo,
					c.GivenMaintainer,
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

func TestHandleStartPipeline(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		Command       command.StartPipeline
		Specification *specification.Specification
		ShouldBeErr   bool
		IsErr         func(err error) bool
	}{
		{
			Name: "specification_with_such_test_campaign_id_not_found",
			Command: command.StartPipeline{
				PipelineID:     "74ff7302-4c99-4729-b1b3-94c03db2e1ba",
				TestCampaignID: "68baf422-777f-4a0e-b35a-4fff5858af2d",
				StartedByID:    "d8d1e4ab-8f24-4c79-a1f2-49e24b3f119a",
			},
			Specification: (&specification.Builder{}).
				WithTestCampaignID("d5a7b2ec-c04e-40d8-a2b5-b273d7ad7ffd").
				WithOwnerID("d8d1e4ab-8f24-4c79-a1f2-49e24b3f119a").
				ErrlessBuild(),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, service.ErrSpecificationNotFound)
			},
		},
		{
			Name: "user_cannot_see_specification",
			Command: command.StartPipeline{
				PipelineID:     "144ed969-2703-419f-9664-36e40b9c51fa",
				TestCampaignID: "5ee6228e-5b0b-4d40-b4e5-9a138bef9f84",
				StartedByID:    "fb883739-2c8c-4a4e-bca2-f96b204f4ac8",
			},
			Specification: (&specification.Builder{}).
				WithTestCampaignID("5ee6228e-5b0b-4d40-b4e5-9a138bef9f84").
				WithOwnerID("8ea9dca1-53da-4ed5-8f4b-660c8956ea45").
				ErrlessBuild(),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *user.AccessError

				return errors.As(err, &target)
			},
		},
		{
			Name: "success_pipeline_starting",
			Command: command.StartPipeline{
				PipelineID:     "09193738-a34c-4a43-9471-d5c74af9f11e",
				TestCampaignID: "70c8e87d-395d-4ae6-b53e-3b2f587039a3",
				StartedByID:    "aa584d3d-c790-4ed3-8bfa-19e1b6fed88e",
			},
			Specification: (&specification.Builder{}).
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
				specRepo   = mock.NewSpecificationRepository(c.Specification)
				pipeRepo   = mock.NewPipelineRepository()
				maintainer = mock.NewPipelineMaintainer(false)
				handler    = command.NewStartPipelineHandler(
					specRepo,
					pipeRepo,
					maintainer,
					pipeline.WithHTTP(pipeline.PassingExecutor()),
					pipeline.WithAssertion(pipeline.FailingExecutor()),
				)
			)

			ctx := context.Background()

			err := handler.Handle(ctx, c.Command)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, pipeRepo.PipelinesNumber())

				return
			}

			require.NoError(t, err)

			require.Equal(t, 1, pipeRepo.PipelinesNumber())
		})
	}
}
