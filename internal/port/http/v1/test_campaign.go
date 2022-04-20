package v1

import (
	"fmt"
	stdhttp "net/http"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/port/http"
)

func (h handler) CreateTestCampaign(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	cmd, ok := decodeCreateTestCampaignCommand(w, r)
	if !ok {
		return
	}

	createdTestCampaignID, err := h.app.Commands.CreateTestCampaign.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Location", fmt.Sprintf("/test-campaigns/%s", createdTestCampaignID))
		w.WriteHeader(stdhttp.StatusCreated)

		return
	}

	http.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetTestCampaigns(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
	w.WriteHeader(stdhttp.StatusNotImplemented)
}

func (h handler) GetTestCampaign(w stdhttp.ResponseWriter, r *stdhttp.Request, testCampaignID string) {
	qry, ok := decodeSpecificTestCampaignQuery(w, r, testCampaignID)
	if !ok {
		return
	}

	tc, err := h.app.Queries.SpecificTestCampaign.Handle(r.Context(), qry)
	if err == nil {
		renderTestCampaignResponse(w, r, tc)

		return
	}

	if errors.Is(err, app.ErrTestCampaignNotFound) {
		http.NotFound(string(ErrorSlugTestCampaignNotFound), err, w, r)

		return
	}

	http.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) RemoveTestCampaign(w stdhttp.ResponseWriter, _ *stdhttp.Request, _ string) {
	w.WriteHeader(stdhttp.StatusNotImplemented)
}
