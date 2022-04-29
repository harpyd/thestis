package v1

import (
	"fmt"
	stdhttp "net/http"

	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/domain/specification"
	"github.com/harpyd/thestis/internal/core/domain/user"
	"github.com/harpyd/thestis/internal/core/interface/rest"
)

func (h handler) LoadSpecification(w stdhttp.ResponseWriter, r *stdhttp.Request, testCampaignID string) {
	cmd, ok := decodeSpecificationSourceCommand(w, r, testCampaignID)
	if !ok {
		return
	}

	loadedSpecID, err := h.app.Commands.LoadSpecification.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Location", fmt.Sprintf("/specifications/%s", loadedSpecID))
		w.WriteHeader(stdhttp.StatusCreated)

		return
	}

	if errors.Is(err, app.ErrTestCampaignNotFound) {
		rest.NotFound(string(ErrorSlugTestCampaignNotFound), err, w, r)

		return
	}

	var berr *specification.BuildError

	if errors.As(err, &berr) {
		rest.UnprocessableEntity(string(ErrorSlugInvalidSpecificationSource), err, w, r)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		rest.Forbidden(string(ErrorSlugUserCantSeeTestCampaign), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetSpecification(w stdhttp.ResponseWriter, r *stdhttp.Request, specificationID string) {
	qry, ok := decodeSpecificSpecificationQuery(w, r, specificationID)
	if !ok {
		return
	}

	spec, err := h.app.Queries.SpecificSpecification.Handle(r.Context(), qry)
	if err == nil {
		renderSpecificationResponse(w, r, spec)

		return
	}

	if errors.Is(err, app.ErrSpecificationNotFound) {
		rest.NotFound(string(ErrorSlugSpecificationNotFound), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}
