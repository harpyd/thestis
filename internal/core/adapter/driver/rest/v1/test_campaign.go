package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/adapter/driver/rest"
	"github.com/harpyd/thestis/internal/core/app/service"
)

func (h handler) CreateTestCampaign(w http.ResponseWriter, r *http.Request) {
	cmd, ok := decodeCreateTestCampaignCommand(w, r, uuid.New().String())
	if !ok {
		return
	}

	err := h.app.Commands.CreateTestCampaign.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set(
			"Location",
			fmt.Sprintf("/test-campaigns/%s", cmd.TestCampaignID),
		)
		w.WriteHeader(http.StatusCreated)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetTestCampaigns(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetTestCampaign(w http.ResponseWriter, r *http.Request, testCampaignID string) {
	qry, ok := decodeSpecificTestCampaignQuery(w, r, testCampaignID)
	if !ok {
		return
	}

	tc, err := h.app.Queries.TestCampaign.Handle(r.Context(), qry)
	if err == nil {
		renderTestCampaignResponse(w, r, tc)

		return
	}

	if errors.Is(err, service.ErrTestCampaignNotFound) {
		rest.NotFound(string(ErrorSlugTestCampaignNotFound), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) RemoveTestCampaign(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
