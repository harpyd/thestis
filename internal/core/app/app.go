package app

import (
	"github.com/harpyd/thestis/internal/core/app/command"
	"github.com/harpyd/thestis/internal/core/app/query"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		CreateTestCampaign command.CreateTestCampaignHandler
		LoadSpecification  command.LoadSpecificationHandler
		StartPerformance   command.StartPerformanceHandler
		RestartPerformance command.RestartPerformanceHandler
		CancelPerformance  command.CancelPerformanceHandler
	}

	Queries struct {
		SpecificTestCampaign  query.SpecificTestCampaignHandler
		SpecificSpecification query.SpecificSpecificationHandler
		SpecificPerformance   query.SpecificPerformanceHandler
	}
)
