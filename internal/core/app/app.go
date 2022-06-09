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
		StartPipeline      command.StartPipelineHandler
		RestartPipeline    command.RestartPipelineHandler
		CancelPipeline     command.CancelPipelineHandler
	}

	Queries struct {
		SpecificTestCampaign  query.SpecificTestCampaignHandler
		SpecificSpecification query.SpecificSpecificationHandler
		SpecificPipeline      query.SpecificPipelineHandler
	}
)
