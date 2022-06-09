package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/performance"
	"github.com/harpyd/thestis/internal/core/entity/user"
	"github.com/harpyd/thestis/internal/core/interface/rest"
)

func (h handler) StartPerformance(w http.ResponseWriter, r *http.Request, testCampaignID string) {
	cmd, ok := decodeStartPerformanceCommand(w, r, uuid.New().String(), testCampaignID)
	if !ok {
		return
	}

	err := h.app.Commands.StartPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Location", fmt.Sprintf("/performances/%s", cmd.PerformanceID))
		w.WriteHeader(http.StatusAccepted)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		rest.Forbidden(string(ErrorSlugUserCantSeeTestCampaign), err, w, r)

		return
	}

	if errors.Is(err, service.ErrSpecificationNotFound) {
		rest.NotFound(string(ErrorSlugSpecificationNotFound), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) RestartPerformance(w http.ResponseWriter, r *http.Request, performanceID string) {
	cmd, ok := decodeRestartPerformanceCommand(w, r, performanceID)
	if !ok {
		return
	}

	err := h.app.Commands.RestartPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		rest.Forbidden(string(ErrorSlugUserCantSeePerformance), err, w, r)

		return
	}

	if errors.Is(err, service.ErrPerformanceNotFound) {
		rest.NotFound(string(ErrorSlugPerformanceNotFound), err, w, r)

		return
	}

	if errors.Is(err, performance.ErrAlreadyStarted) {
		rest.Conflict(string(ErrorSlugPerformanceAlreadyStarted), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) CancelPerformance(w http.ResponseWriter, r *http.Request, performanceID string) {
	cmd, ok := decodeCancelPerformanceCommand(w, r, performanceID)
	if !ok {
		return
	}

	err := h.app.Commands.CancelPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		rest.Forbidden(string(ErrorSlugUserCantSeePerformance), err, w, r)

		return
	}

	if errors.Is(err, service.ErrPerformanceNotFound) {
		rest.NotFound(string(ErrorSlugPerformanceNotFound), err, w, r)

		return
	}

	if errors.Is(err, performance.ErrNotStarted) {
		rest.Conflict(string(ErrorSlugPerformanceNotStarted), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetPerformanceHistory(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetPerformance(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
