package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/user"
	"github.com/harpyd/thestis/internal/port/http/httperr"
)

func (h handler) StartNewPerformance(w http.ResponseWriter, r *http.Request, testCampaignID string) {
	cmd, ok := unmarshalStartNewPerformanceCommand(w, r, testCampaignID)
	if !ok {
		return
	}

	perfID, msg, err := h.app.Commands.StartNewPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Location", fmt.Sprintf("/performances/%s", perfID))

		go func(reqID string) {
			logField := app.StringLogField("requestId", reqID)

			for m := range msg {
				if m.Err() == nil || performance.IsFailedError(err) {
					h.logger.Info(m.String(), logField)

					continue
				}

				h.logger.Error(m.String(), m.Err(), logField)
			}
		}(middleware.GetReqID(r.Context()))

		return
	}

	if user.IsUserCantSeeTestCampaignError(err) {
		httperr.Forbidden(string(ErrorSlugUserCantSeeTestCampaign), err, w, r)

		return
	}

	if app.IsSpecificationNotFoundError(err) {
		httperr.NotFound(string(ErrorSlugSpecificationNotFound), err, w, r)

		return
	}

	if performance.IsAlreadyStartedError(err) {
		httperr.Conflict(string(ErrorSlugPerformanceAlreadyStarted), err, w, r)

		return
	}

	httperr.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetPerformancesHistory(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetPerformance(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) RestartPerformance(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) CancelPerformance(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
