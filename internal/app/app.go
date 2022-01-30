package app

import "context"

type Application struct {
	Commands Commands
	Queries  Queries
}

type (
	Commands struct {
		CreateTestCampaign  createTestCampaignHandler
		LoadSpecification   loadSpecificationHandler
		StartNewPerformance startNewPerformanceHandler
	}

	createTestCampaignHandler interface {
		Handle(ctx context.Context, cmd CreateTestCampaignCommand) (string, error)
	}

	loadSpecificationHandler interface {
		Handle(ctx context.Context, cmd LoadSpecificationCommand) (string, error)
	}

	startNewPerformanceHandler interface {
		Handle(ctx context.Context, cmd StartNewPerformanceCommand) (string, <-chan Message, error)
	}
)

type (
	Queries struct {
		SpecificTestCampaign  specificTestCampaignHandler
		SpecificSpecification specificSpecificationHandler
	}

	specificTestCampaignHandler interface {
		Handle(ctx context.Context, qry SpecificTestCampaignQuery) (SpecificTestCampaign, error)
	}

	specificSpecificationHandler interface {
		Handle(ctx context.Context, qry SpecificSpecificationQuery) (SpecificSpecification, error)
	}
)
