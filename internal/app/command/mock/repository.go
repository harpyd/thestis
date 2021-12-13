package mock

import (
	"context"

	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type TestCampaignsRepository struct {
	campaigns map[string]testcampaign.TestCampaign
}

func NewTestCampaignsRepository(tcs ...*testcampaign.TestCampaign) *TestCampaignsRepository {
	tcm := &TestCampaignsRepository{
		campaigns: make(map[string]testcampaign.TestCampaign, len(tcs)),
	}

	for _, tc := range tcs {
		tcm.campaigns[tc.ID()] = *tc
	}

	return tcm
}

func (m *TestCampaignsRepository) AddTestCampaign(_ context.Context, tc *testcampaign.TestCampaign) error {
	m.campaigns[tc.ID()] = *tc

	return nil
}

func (m *TestCampaignsRepository) TestCampaignsNumber() int {
	return len(m.campaigns)
}
