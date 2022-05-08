package app

import (
	"context"

	"github.com/harpyd/thestis/internal/core/app/service"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type (
	Commands struct {
		CreateTestCampaign CreateTestCampaignHandler
		LoadSpecification  LoadSpecificationHandler
		StartPerformance   StartPerformanceHandler
		RestartPerformance RestartPerformanceHandler
		CancelPerformance  CancelPerformanceHandler
	}

	CreateTestCampaignHandler interface {
		Handle(ctx context.Context, cmd CreateTestCampaignCommand) (string, error)
	}

	LoadSpecificationHandler interface {
		Handle(ctx context.Context, cmd LoadSpecificationCommand) (string, error)
	}

	StartPerformanceHandler interface {
		Handle(ctx context.Context, cmd StartPerformanceCommand) (string, <-chan service.Message, error)
	}

	RestartPerformanceHandler interface {
		Handle(ctx context.Context, cmd RestartPerformanceCommand) (<-chan service.Message, error)
	}

	CancelPerformanceHandler interface {
		Handle(ctx context.Context, cmd CancelPerformanceCommand) error
	}
)

type (
	Queries struct {
		SpecificTestCampaign  SpecificTestCampaignHandler
		SpecificSpecification SpecificSpecificationHandler
	}

	SpecificTestCampaignHandler interface {
		Handle(ctx context.Context, qry SpecificTestCampaignQuery) (SpecificTestCampaign, error)
	}

	SpecificSpecificationHandler interface {
		Handle(ctx context.Context, qry SpecificSpecificationQuery) (SpecificSpecification, error)
	}
)
