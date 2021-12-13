package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type CreateTestCampaignHandler struct {
	testCampaignsRepo testCampaignsRepository
}

func NewCreateTestCampaignHandler(repo testCampaignsRepository) CreateTestCampaignHandler {
	if repo == nil {
		panic("test campaigns repository is nil")
	}

	return CreateTestCampaignHandler{testCampaignsRepo: repo}
}

func (h CreateTestCampaignHandler) Handle(
	ctx context.Context,
	cmd app.CreateTestCampaignCommand,
) (testCampaignID string, err error) {
	testCampaignID = uuid.New().String()

	tc := testcampaign.New(testCampaignID, cmd.ViewName)

	if err = h.testCampaignsRepo.AddTestCampaign(ctx, tc); err != nil {
		return "", errors.Wrapf(err, "test campaigns with view name = %s", cmd.ViewName)
	}

	return
}
