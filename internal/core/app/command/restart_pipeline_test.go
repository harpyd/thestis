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

func TestNewRestartPipelineHandlerPanics(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		GivenPipeRepo   service.PipelineRepository
		GivenSpecGetter service.SpecificationGetter
		GivenMaintainer service.PipelineMaintainer
		ShouldPanic     bool
		PanicMessage    string
	}{
		{
			Name:            "all_dependencies_are_not_nil",
			GivenPipeRepo:   mock.NewPipelineRepository(),
			GivenSpecGetter: service.WithoutSpecification(),
			GivenMaintainer: mock.NewPipelineMaintainer(false),
			ShouldPanic:     false,
		},
		{
			Name:            "pipeline_repository_is_nil",
			GivenPipeRepo:   nil,
			GivenSpecGetter: service.WithoutSpecification(),
			GivenMaintainer: mock.NewPipelineMaintainer(false),
			ShouldPanic:     true,
			PanicMessage:    "pipeline repository is nil",
		},
		{
			Name:            "specification_getter_is_nil",
			GivenPipeRepo:   mock.NewPipelineRepository(),
			GivenSpecGetter: nil,
			GivenMaintainer: mock.NewPipelineMaintainer(false),
			ShouldPanic:     true,
			PanicMessage:    "specification getter is nil",
		},
		{
			Name:            "pipeline_maintainer_is_nil",
			GivenPipeRepo:   mock.NewPipelineRepository(),
			GivenSpecGetter: service.WithoutSpecification(),
			GivenMaintainer: nil,
			ShouldPanic:     true,
			PanicMessage:    "pipeline maintainer is nil",
		},
		{
			Name:            "all_dependencies_are_nil",
			GivenPipeRepo:   nil,
			GivenSpecGetter: service.WithoutSpecification(),
			GivenMaintainer: mock.NewPipelineMaintainer(false),
			ShouldPanic:     true,
			PanicMessage:    "pipeline repository is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = command.NewRestartPipelineHandler(
					c.GivenPipeRepo,
					c.GivenSpecGetter,
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

func TestHandleRestartPipeline(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		Command                command.RestartPipeline
		Pipeline               *pipeline.Pipeline
		PipelineAlreadyStarted bool
		ShouldBeErr            bool
		IsErr                  func(err error) bool
	}{
		{
			Name: "pipeline_not_found",
			Command: command.RestartPipeline{
				PipelineID:  "9524edfe-9a47-40d5-9c40-5c575de1368b",
				StartedByID: "7177997a-b63d-4e1b-9288-0a581f7ff03a",
			},
			Pipeline: pipeline.Unmarshal(pipeline.Params{
				ID:      "ac622aa2-59e9-4389-97cc-2309271f1f36",
				OwnerID: "7177997a-b63d-4e1b-9288-0a581f7ff03a",
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				return errors.Is(err, service.ErrPipelineNotFound)
			},
		},
		{
			Name: "user_cannot_see_pipeline",
			Command: command.RestartPipeline{
				PipelineID:  "50b59781-6932-422e-a7e6-f7424b5f5d36",
				StartedByID: "be746c97-fe6e-4795-820a-97337e8d98b2",
			},
			Pipeline: pipeline.Unmarshal(pipeline.Params{
				ID:      "50b59781-6932-422e-a7e6-f7424b5f5d36",
				OwnerID: "0a5782ec-7580-4526-b4e4-c3a7489ca512",
			}),
			ShouldBeErr: true,
			IsErr: func(err error) bool {
				var target *user.AccessError

				return errors.As(err, &target)
			},
		},
		{
			Name: "pipeline_already_started",
			Command: command.RestartPipeline{
				PipelineID:  "6f112cf1-3dd5-4f14-a5ef-7ef18dfb8921",
				StartedByID: "960f7ba1-b16c-43eb-9f87-d367ec9e0ba9",
			},
			Pipeline: pipeline.Unmarshal(pipeline.Params{
				ID:      "6f112cf1-3dd5-4f14-a5ef-7ef18dfb8921",
				OwnerID: "960f7ba1-b16c-43eb-9f87-d367ec9e0ba9",
			}),
			PipelineAlreadyStarted: true,
			ShouldBeErr:            true,
			IsErr: func(err error) bool {
				return errors.Is(err, pipeline.ErrAlreadyStarted)
			},
		},
		{
			Name: "success_pipeline_restarting",
			Command: command.RestartPipeline{
				PipelineID:  "fc2f14b3-6125-47fa-a343-5fabcac9abd1",
				StartedByID: "5da02570-a192-4a9a-9180-1a2704732b06",
			},
			Pipeline: pipeline.Unmarshal(pipeline.Params{
				ID:      "fc2f14b3-6125-47fa-a343-5fabcac9abd1",
				OwnerID: "5da02570-a192-4a9a-9180-1a2704732b06",
			}),
			ShouldBeErr: false,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				pipeRepo   = mock.NewPipelineRepository(c.Pipeline)
				maintainer = mock.NewPipelineMaintainer(c.PipelineAlreadyStarted)
				handler    = command.NewRestartPipelineHandler(
					pipeRepo,
					service.WithoutSpecification(),
					maintainer,
					pipeline.WithHTTP(pipeline.PassingExecutor()),
					pipeline.WithAssertion(pipeline.FailingExecutor()),
				)
			)

			ctx := context.Background()

			err := handler.Handle(ctx, c.Command)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))

				return
			}

			require.NoError(t, err)
		})
	}
}
