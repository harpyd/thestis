package query

import (
	"context"

	"github.com/pkg/errors"
)

type SpecificTestCampaign struct {
	TestCampaignID string
	UserID         string
}

type SpecificTestCampaignHandler interface {
	Handle(ctx context.Context, qry SpecificTestCampaign) (SpecificTestCampaignModel, error)
}

type SpecificTestCampaignReadModel interface {
	FindTestCampaign(ctx context.Context, qry SpecificTestCampaign) (SpecificTestCampaignModel, error)
}

type specificTestCampaignHandler struct {
	readModel SpecificTestCampaignReadModel
}

func NewSpecificTestCampaignHandler(readModel SpecificTestCampaignReadModel) SpecificTestCampaignHandler {
	if readModel == nil {
		panic("specific test campaign read model is nil")
	}

	return specificTestCampaignHandler{
		readModel: readModel,
	}
}

func (h specificTestCampaignHandler) Handle(
	ctx context.Context,
	qry SpecificTestCampaign,
) (SpecificTestCampaignModel, error) {
	tc, err := h.readModel.FindTestCampaign(ctx, qry)

	return tc, errors.Wrap(err, "getting test campaign")
}
