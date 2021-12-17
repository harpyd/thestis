package v1

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/harpyd/thestis/internal/app"
)

func unmarshalToCreateTestCampaignCommand(
	w http.ResponseWriter,
	r *http.Request,
) (cmd app.CreateTestCampaignCommand, ok bool) {
	var rb CreateTestCampaignRequest

	if ok = decode(w, r, &rb); !ok {
		return
	}

	var summary string
	if rb.Summary != nil {
		summary = *rb.Summary
	}

	return app.CreateTestCampaignCommand{
		ViewName: rb.ViewName,
		Summary:  summary,
	}, true
}

func unmarshalToSpecificTestCampaignQuery(testCampaignID string) (app.SpecificTestCampaignQuery, bool) {
	return app.SpecificTestCampaignQuery{
		TestCampaignID: testCampaignID,
	}, true
}

func marshalToTestCampaignResponse(w http.ResponseWriter, r *http.Request, tc app.SpecificTestCampaign) {
	response := TestCampaignResponse{
		Id:                    tc.ID,
		ViewName:              tc.ViewName,
		Summary:               &tc.Summary,
		ActiveSpecificationId: &tc.ActiveSpecificationID,
	}

	render.Respond(w, r, response)
}
