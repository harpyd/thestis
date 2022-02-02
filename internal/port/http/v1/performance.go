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

func (h handler) StartPerformance(w http.ResponseWriter, r *http.Request, testCampaignID string) {
	cmd, ok := unmarshalStartPerformanceCommand(w, r, testCampaignID)
	if !ok {
		return
	}

	perfID, messages, err := h.app.Commands.StartPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Location", fmt.Sprintf("/performances/%s", perfID))

		go func(reqID string) {
			logFields := []app.LogField{
				app.BoolLogField("isNew", true),
				app.StringLogField("requestId", reqID),
				app.StringLogField("performanceId", perfID),
			}

			for msg := range messages {
				if msg.Err() == nil || msg.State() == performance.Failed {
					h.logger.Info(msg.String(), logFields...)

					continue
				}

				h.logger.Error(msg.String(), msg.Err(), logFields...)
			}
		}(middleware.GetReqID(r.Context()))

		return
	}

	if user.IsCantSeeTestCampaignError(err) {
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
