package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/testcampaign"
)

type CreateTestCampaign struct {
	OwnerID  string
	ViewName string
	Summary  string
}

type CreateTestCampaignHandler interface {
	Handle(ctx context.Context, cmd CreateTestCampaign) (string, error)
}

type createTestCampaignHandler struct {
	testCampaignRepo service.TestCampaignRepository
}

func NewCreateTestCampaignHandler(repo service.TestCampaignRepository) CreateTestCampaignHandler {
	if repo == nil {
		panic("test campaign repository is nil")
	}

	return createTestCampaignHandler{testCampaignRepo: repo}
}

func (h createTestCampaignHandler) Handle(
	ctx context.Context,
	cmd CreateTestCampaign,
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
