package command

import (
	"context"

	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type (
	testCampaignsRepository interface {
		AddTestCampaign(ctx context.Context, tc *testcampaign.TestCampaign) error
	}
)
