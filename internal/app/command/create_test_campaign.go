package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type CreateTestCampaignHandler struct {
	testCampaignRepo app.TestCampaignRepository
}

func NewCreateTestCampaignHandler(repo app.TestCampaignRepository) CreateTestCampaignHandler {
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
