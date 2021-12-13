package v1

import (
	"fmt"
	"net/http"

	"github.com/harpyd/thestis/pkg/httperr"
)

func (h handler) CreateTestCampaign(w http.ResponseWriter, r *http.Request) {
	cmd, ok := unmarshalCreateTestCampaignCommand(w, r)
	if !ok {
		return
	}

	createdTestCampaignID, err := h.app.Commands.CreateTestCampaign.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Location", fmt.Sprintf("/test-campaigns/%s", createdTestCampaignID))
		w.WriteHeader(http.StatusCreated)

		return
	}

	httperr.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetTestCampaigns(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetTestCampaign(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) RemoveTestCampaign(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
