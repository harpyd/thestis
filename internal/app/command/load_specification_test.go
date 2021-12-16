package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/app/command"
	"github.com/harpyd/thestis/internal/app/command/mock"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

const spec = `
---
author: Djerys
title: horns-and-hooves API test
description: declarative auto tests for horns-and-hooves API

...
`

func TestLoadSpecificationHandler_Handle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		Command             app.LoadSpecificationCommand
		TestCampaignFactory func() *testcampaign.TestCampaign
		ParseWithErr        bool
		ShouldBeErr         bool
		IsErr               func(err error) bool
	}{
		{
			Name: "load_valid_specification",
			Command: app.LoadSpecificationCommand{
				TestCampaignID: "35474763-28f4-43a6-a184-e8894f50cba8",
				Content:        []byte(spec),
			},
			TestCampaignFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New("35474763-28f4-43a6-a184-e8894f50cba8", "view name", "summary")
				require.NoError(t, err)

				return tc
			},
			ParseWithErr: false,
			ShouldBeErr:  false,
		},
		{
			Name: "load_invalid_specification",
			Command: app.LoadSpecificationCommand{
				TestCampaignID: "f18fdd19-d69c-4afe-a639-8bcefd6c4af9",
				Content:        []byte(spec),
			},
			TestCampaignFactory: func() *testcampaign.TestCampaign {
				tc, err := testcampaign.New("f18fdd19-d69c-4afe-a639-8bcefd6c4af9", "view", "summary")
				require.NoError(t, err)

				return tc
			},
			ParseWithErr: true,
			ShouldBeErr:  true,
			IsErr:        specification.IsBuildSpecificationError,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			var (
				specRepo = mock.NewSpecificationsRepository()
				tcRepo   = mock.NewTestCampaignsRepository(c.TestCampaignFactory())
				parser   = mock.NewSpecificationParserService(c.ParseWithErr)
				handler  = command.NewLoadSpecificationHandler(specRepo, tcRepo, parser)
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

			tc, err := tcRepo.GetTestCampaign(ctx, c.Command.TestCampaignID)
			require.NoError(t, err)
			require.Equal(t, specID, tc.ActiveSpecificationID())
		})
	}
}
