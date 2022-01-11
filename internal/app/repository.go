package app

import (
	"context"

	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type (
	TestCampaignsRepository interface {
		GetTestCampaign(ctx context.Context, tcID string) (*testcampaign.TestCampaign, error)
		AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error
		UpdateTestCampaign(ctx context.Context, tcID string, updateFn TestCampaignUpdater) error
	}

	TestCampaignUpdater func(
		ctx context.Context,
		tc *testcampaign.TestCampaign,
	) (*testcampaign.TestCampaign, error)

	SpecificTestCampaignReadModel interface {
		FindTestCampaign(ctx context.Context, qry SpecificTestCampaignQuery) (SpecificTestCampaign, error)
	}
)

type (
	SpecificationsRepository interface {
		GetSpecification(ctx context.Context, specID string) (*specification.Specification, error)
		AddSpecification(ctx context.Context, spec *specification.Specification) error
	}

	SpecificSpecificationReadModel interface {
		FindSpecification(ctx context.Context, qry SpecificSpecificationQuery) (SpecificSpecification, error)
	}
)