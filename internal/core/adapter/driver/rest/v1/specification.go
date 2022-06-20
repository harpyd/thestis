package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/adapter/driver/rest"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/specification"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

func (h handler) LoadSpecification(w http.ResponseWriter, r *http.Request, testCampaignID string) {
	cmd, ok := decodeSpecificationSourceCommand(w, r, uuid.New().String(), testCampaignID)
	if !ok {
		return
	}

	err := h.app.Commands.LoadSpecification.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Location", fmt.Sprintf("/specifications/%s", cmd.SpecificationID))
		w.WriteHeader(http.StatusCreated)

		return
	}

	if errors.Is(err, service.ErrTestCampaignNotFound) {
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

func (h handler) GetSpecification(w http.ResponseWriter, r *http.Request, specificationID string) {
	qry, ok := decodeSpecificSpecificationQuery(w, r, specificationID)
	if !ok {
		return
	}

	spec, err := h.app.Queries.Specification.Handle(r.Context(), qry)
	if err == nil {
		renderSpecificationResponse(w, r, spec)

		return
	}

	if errors.Is(err, service.ErrSpecificationNotFound) {
		rest.NotFound(string(ErrorSlugSpecificationNotFound), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}
