package app

import "context"

type Application struct {
	Commands Commands
}

type (
	Commands struct {
		CreateTestCampaign createTestCampaignHandler
	}

	createTestCampaignHandler interface {
		Handle(ctx context.Context, cmd CreateTestCampaignCommand) (string, error)
	}
)
