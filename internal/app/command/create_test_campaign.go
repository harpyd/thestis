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
	defer func() {
		err = errors.Wrap(err, "test campaign creation")
	}()

	testCampaignID = uuid.New().String()

	tc, err := testcampaign.New(testcampaign.Params{
		ID:       testCampaignID,
		OwnerID:  cmd.OwnerID,
		ViewName: cmd.ViewName,
		Summary:  cmd.Summary,
	})
	if err != nil {
		return "", err
	}

	if err = h.testCampaignsRepo.AddTestCampaign(ctx, tc); err != nil {
		return "", err
	}

	return
}
