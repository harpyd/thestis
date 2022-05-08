package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/testcampaign"
)

type CreateTestCampaignHandler struct {
	testCampaignRepo service.TestCampaignRepository
}

func NewCreateTestCampaignHandler(repo service.TestCampaignRepository) CreateTestCampaignHandler {
	if repo == nil {
		panic("test campaign repository is nil")
	}

	return CreateTestCampaignHandler{testCampaignRepo: repo}
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
		ID:        testCampaignID,
		OwnerID:   cmd.OwnerID,
		ViewName:  cmd.ViewName,
		Summary:   cmd.Summary,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return "", err
	}

	if err = h.testCampaignRepo.AddTestCampaign(ctx, tc); err != nil {
		return "", err
	}

	return
}
