package command

import (
	"context"

	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type (
	testCampaignsRepository interface {
		GetTestCampaign(ctx context.Context, tcID string) (*testcampaign.TestCampaign, error)
		AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error
		UpdateTestCampaign(ctx context.Context, tcID string, updateFn TestCampaignUpdater) error
	}

	TestCampaignUpdater func(
		ctx context.Context,
		tc *testcampaign.TestCampaign,
	) (*testcampaign.TestCampaign, error)
)

type specificationsRepository interface {
	GetSpecification(ctx context.Context, specID string) (*specification.Specification, error)
	AddSpecification(ctx context.Context, spec *specification.Specification) error
}
