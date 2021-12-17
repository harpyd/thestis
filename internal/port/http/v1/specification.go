package v1

import (
	"fmt"
	"net/http"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/pkg/httperr"
)

func (h handler) LoadSpecification(w http.ResponseWriter, r *http.Request, testCampaignID string) {
	cmd, ok := unmarshalToSpecificationSourceCommand(w, r, testCampaignID)
	if !ok {
		return
	}

	loadedSpecID, err := h.app.Commands.LoadSpecification.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", fmt.Sprintf("/specifications/%s", loadedSpecID))

		return
	}

	if app.IsTestCampaignNotFoundError(err) {
		httperr.NotFound(string(ErrorSlugTestCampaignNotFound), err, w, r)

		return
	}

	if specification.IsBuildSpecificationError(err) {
		httperr.UnprocessableEntity(string(ErrorSlugInvalidSpecificationSource), err, w, r)

		return
	}

	httperr.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetSpecification(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) RemoveSpecification(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
