package query

import (
	"context"

	"github.com/pkg/errors"
)

type TestCampaign struct {
	TestCampaignID string
	UserID         string
}

type TestCampaignHandler interface {
	Handle(ctx context.Context, qry TestCampaign) (TestCampaignModel, error)
}

type TestCampaignReadModel interface {
	FindTestCampaign(ctx context.Context, qry TestCampaign) (TestCampaignModel, error)
}

type testCampaignHandler struct {
	readModel TestCampaignReadModel
}

func NewTestCampaignHandler(readModel TestCampaignReadModel) TestCampaignHandler {
	if readModel == nil {
		panic("test campaign read model is nil")
	}

	return testCampaignHandler{
		readModel: readModel,
	}
}

func (h testCampaignHandler) Handle(
	ctx context.Context,
	qry TestCampaign,
) (TestCampaignModel, error) {
	tc, err := h.readModel.FindTestCampaign(ctx, qry)

	return tc, errors.Wrap(err, "getting test campaign")
}
