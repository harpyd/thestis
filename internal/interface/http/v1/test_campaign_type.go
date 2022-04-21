package v1

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/harpyd/thestis/internal/app"
)

func decodeCreateTestCampaignCommand(
	w http.ResponseWriter,
	r *http.Request,
) (cmd app.CreateTestCampaignCommand, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

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
		OwnerID:  user.UUID,
	}, true
}

func decodeSpecificTestCampaignQuery(
	w http.ResponseWriter,
	r *http.Request,
	testCampaignID string,
) (qry app.SpecificTestCampaignQuery, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return app.SpecificTestCampaignQuery{
		TestCampaignID: testCampaignID,
		UserID:         user.UUID,
	}, true
}

func renderTestCampaignResponse(w http.ResponseWriter, r *http.Request, tc app.SpecificTestCampaign) {
	response := TestCampaignResponse{
		Id:        tc.ID,
		ViewName:  tc.ViewName,
		Summary:   &tc.Summary,
		CreatedAt: tc.CreatedAt,
	}

	render.Respond(w, r, response)
}
