package v1

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/harpyd/thestis/internal/core/app/command"
	"github.com/harpyd/thestis/internal/core/app/query"
)

func decodeCreateTestCampaignCommand(
	w http.ResponseWriter,
	r *http.Request,
	testCampaignID string,
) (cmd command.CreateTestCampaign, ok bool) {
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

	return command.CreateTestCampaign{
		TestCampaignID: testCampaignID,
		ViewName:       rb.ViewName,
		Summary:        summary,
		OwnerID:        user.UUID,
	}, true
}

func decodeSpecificTestCampaignQuery(
	w http.ResponseWriter,
	r *http.Request,
	testCampaignID string,
) (qry query.SpecificTestCampaign, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return query.SpecificTestCampaign{
		TestCampaignID: testCampaignID,
		UserID:         user.UUID,
	}, true
}

func renderTestCampaignResponse(
	w http.ResponseWriter,
	r *http.Request,
	tc query.SpecificTestCampaignView,
) {
	response := TestCampaignResponse{
		Id:        tc.ID,
		ViewName:  tc.ViewName,
		Summary:   &tc.Summary,
		CreatedAt: tc.CreatedAt,
	}

	render.Respond(w, r, response)
}
