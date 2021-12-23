package query

import (
	"context"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
)

type SpecificTestCampaignHandler struct {
	readModel app.SpecificTestCampaignReadModel
}

func NewSpecificTestCampaignHandler(readModel app.SpecificTestCampaignReadModel) SpecificTestCampaignHandler {
	if readModel == nil {
		panic("specific test campaign read model is nil")
	}

	return SpecificTestCampaignHandler{
		readModel: readModel,
	}
}

func (h SpecificTestCampaignHandler) Handle(
	ctx context.Context,
	qry app.SpecificTestCampaignQuery,
) (app.SpecificTestCampaign, error) {
	tc, err := h.readModel.FindTestCampaign(ctx, qry)

	return tc, errors.Wrap(err, "getting test campaign")
}
