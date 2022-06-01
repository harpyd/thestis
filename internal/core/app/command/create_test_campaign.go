package command

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/testcampaign"
)

type CreateTestCampaign struct {
	TestCampaignID string
	OwnerID        string
	ViewName       string
	Summary        string
}

type CreateTestCampaignHandler interface {
	Handle(ctx context.Context, cmd CreateTestCampaign) error
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
) (err error) {
	defer func() {
		err = errors.Wrap(err, "test campaign creation")
	}()

	tc, err := testcampaign.New(testcampaign.Params{
		ID:        cmd.TestCampaignID,
		OwnerID:   cmd.OwnerID,
		ViewName:  cmd.ViewName,
		Summary:   cmd.Summary,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return h.testCampaignRepo.AddTestCampaign(ctx, tc)
}
