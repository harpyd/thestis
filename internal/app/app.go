package app

import "context"

type Application struct {
	Commands Commands
}

type (
	Commands struct {
		CreateTestCampaign createTestCampaignHandler
		LoadSpecification  loadSpecificationHandler
	}

	createTestCampaignHandler interface {
		Handle(ctx context.Context, cmd CreateTestCampaignCommand) (string, error)
	}

	loadSpecificationHandler interface {
		Handle(ctx context.Context, cmd LoadSpecificationCommand) (string, error)
	}
)
