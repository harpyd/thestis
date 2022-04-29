package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/app/command"
	"github.com/harpyd/thestis/internal/core/app/mock"
	"github.com/harpyd/thestis/internal/core/domain/specification"
	"github.com/harpyd/thestis/internal/core/domain/testcampaign"
	"github.com/harpyd/thestis/internal/core/domain/user"
)

func TestPanickingNewLoadSpecificationHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                  string
		GivenSpecRepo         app.SpecificationRepository
		GivenTestCampaignRepo app.TestCampaignRepository
		GivenSpecParser       app.SpecificationParser
		ShouldPanic           bool
		PanicMessage          string
	}{
		{
			Name:                  "all_dependencies_are_not_nil",
			GivenSpecRepo:         mock.NewSpecificationRepository(),
			GivenTestCampaignRepo: mock.NewTestCampaignRepository(),
			GivenSpecParser:       mock.NewSpecificationParserService(false),
			ShouldPanic:           false,
		},
		{
			Name:                  "specification_repository_is_nil",
			GivenSpecRepo:         nil,
			GivenTestCampaignRepo: mock.NewTestCampaignRepository(),
			GivenSpecParser:       mock.NewSpecificationParserService(false),
			ShouldPanic:           true,
			PanicMessage:          "specification repository is nil",
		},
		{
			Name:                  "test_campaign_repository_is_nil",
			GivenSpecRepo:         mock.NewSpecificationRepository(),
			GivenTestCampaignRepo: nil,
			GivenSpecParser:       mock.NewSpecificationParserService(false),
			ShouldPanic:           true,
			PanicMessage:          "test campaign repository is nil",
		},
		{
			Name:                  "specification_parser_is_nil",
			GivenSpecRepo:         mock.NewSpecificationRepository(),
			GivenTestCampaignRepo: mock.NewTestCampaignRepository(),
			GivenSpecParser:       nil,
			ShouldPanic:           true,
			PanicMessage:          "specification parser is nil",
		},
		{
			Name:                  "all_dependencies_are_nil",
			GivenSpecRepo:         nil,
			GivenTestCampaignRepo: nil,
			GivenSpecParser:       nil,
			ShouldPanic:           true,
			PanicMessage:          "specification repository is nil",
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			init := func() {
				_ = command.NewLoadSpecificationHandler(
					c.GivenSpecRepo,
					c.GivenTestCampaignRepo,
					c.GivenSpecParser,
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

const spec = `
---
author: Djerys
title: horns-and-hooves API test
description: declarative auto tests for horns-and-hooves API

...
`

func TestHandleLoadSpecification(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name         string
		Command      app.LoadSpecificationCommand
		TestCampaign *testcampaign.TestCampaign
		ParseWithErr bool
		ShouldBeErr  bool
		IsErr        func(err error) bool
	}{
		{
			Name: "load_valid_specification",
			Command: app.LoadSpecificationCommand{
				TestCampaignID: "35474763-28f4-43a6-a184-e8894f50cba8",
				LoadedByID:     "cb39a8e2-8f79-484b-bc48-51f83a8e8c33",
				Content:        []byte(spec),
			},
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:       "35474763-28f4-43a6-a184-e8894f50cba8",
				ViewName: "view name",
				Summary:  "summary",
				OwnerID:  "cb39a8e2-8f79-484b-bc48-51f83a8e8c33",
			}),
			ParseWithErr: false,
			ShouldBeErr:  false,
		},
		{
			Name: "load_invalid_specification",
			Command: app.LoadSpecificationCommand{
				TestCampaignID: "f18fdd19-d69c-4afe-a639-8bcefd6c4af9",
				LoadedByID:     "dc0479de-33ed-4631-b9a4-2834c3efb7b1",
				Content:        []byte(spec),
			},
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:       "f18fdd19-d69c-4afe-a639-8bcefd6c4af9",
				ViewName: "view",
				Summary:  "summary",
				OwnerID:  "dc0479de-33ed-4631-b9a4-2834c3efb7b1",
			}),
			ParseWithErr: true,
			ShouldBeErr:  true,
			IsErr: func(err error) bool {
				var target *specification.BuildError

				return errors.As(err, &target)
			},
		},
		{
			Name: "user_cant_see_test_campaign",
			Command: app.LoadSpecificationCommand{
				TestCampaignID: "e7c57ccf-3bff-402b-ada5-71990e3ab0cd",
				LoadedByID:     "1dccc358-2f91-427b-b2a8-f46169fc3a04",
				Content:        []byte(spec),
			},
			TestCampaign: testcampaign.MustNew(testcampaign.Params{
				ID:      "e7c57ccf-3bff-402b-ada5-71990e3ab0cd",
				OwnerID: "98fd29f8-442b-420c-9cf8-e4d3c1f105e8",
			}),
			ParseWithErr: false,
			ShouldBeErr:  true,
			IsErr: func(err error) bool {
				var target *user.AccessError

				return errors.As(err, &target)
			},
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				specRepo = mock.NewSpecificationRepository()
				tcsRepo  = mock.NewTestCampaignRepository(c.TestCampaign)
				parser   = mock.NewSpecificationParserService(c.ParseWithErr)
				handler  = command.NewLoadSpecificationHandler(specRepo, tcsRepo, parser)
			)

			ctx := context.Background()

			specID, err := handler.Handle(ctx, c.Command)

			if c.ShouldBeErr {
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, specRepo.SpecificationsNumber())

				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, specID)
			require.Equal(t, 1, specRepo.SpecificationsNumber())
		})
	}
}
