package v1

import (
	"fmt"
	stdhttp "net/http"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/user"
	"github.com/harpyd/thestis/internal/port/http"
)

func (h handler) LoadSpecification(w stdhttp.ResponseWriter, r *stdhttp.Request, testCampaignID string) {
	cmd, ok := unmarshalToSpecificationSourceCommand(w, r, testCampaignID)
	if !ok {
		return
	}

	loadedSpecID, err := h.app.Commands.LoadSpecification.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(stdhttp.StatusCreated)
		w.Header().Set("Location", fmt.Sprintf("/specifications/%s", loadedSpecID))

		return
	}

	if errors.Is(err, app.ErrTestCampaignNotFound) {
		http.NotFound(string(ErrorSlugTestCampaignNotFound), err, w, r)

		return
	}

	var berr *specification.BuildError

	if errors.As(err, &berr) {
		http.UnprocessableEntity(string(ErrorSlugInvalidSpecificationSource), err, w, r)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		http.Forbidden(string(ErrorSlugUserCantSeeTestCampaign), err, w, r)

		return
	}

	http.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetSpecification(w stdhttp.ResponseWriter, r *stdhttp.Request, specificationID string) {
	qry, ok := unmarshalToSpecificSpecificationQuery(w, r, specificationID)
	if !ok {
		return
	}

	spec, err := h.app.Queries.SpecificSpecification.Handle(r.Context(), qry)
	if err == nil {
		marshalToSpecificationResponse(w, r, spec)

		return
	}

	if errors.Is(err, app.ErrSpecificationNotFound) {
		http.NotFound(string(ErrorSlugSpecificationNotFound), err, w, r)

		return
	}

	http.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}
