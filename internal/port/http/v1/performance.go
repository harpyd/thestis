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
		w.Header().Set("Location", fmt.Sprintf("/performances/%s", perfID))
		w.WriteHeader(http.StatusAccepted)

		go h.logMessages(
			r,
			messages,
			app.StringLogField("performanceId", perfID),
			app.BoolLogField("restarted", false),
		)

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

	httperr.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) RestartPerformance(w http.ResponseWriter, r *http.Request, performanceID string) {
	cmd, ok := unmarshalRestartPerformanceCommand(w, r, performanceID)
	if !ok {
		return
	}

	messages, err := h.app.Commands.RestartPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)

		go h.logMessages(r, messages, app.BoolLogField("restarted", true))

		return
	}

	if user.IsCantSeePerformanceError(err) {
		httperr.Forbidden(string(ErrorSlugUserCantSeePerformance), err, w, r)

		return
	}

	if app.IsPerformanceNotFoundError(err) {
		httperr.NotFound(string(ErrorSlugPerformanceNotFound), err, w, r)

		return
	}

	if performance.IsAlreadyStartedError(err) {
		httperr.Conflict(string(ErrorSlugPerformanceAlreadyStarted), err, w, r)

		return
	}

	httperr.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) CancelPerformance(w http.ResponseWriter, r *http.Request, performanceID string) {
	cmd, ok := unmarshalCancelPerformanceCommand(w, r, performanceID)
	if !ok {
		return
	}

	err := h.app.Commands.CancelPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if user.IsCantSeePerformanceError(err) {
		httperr.Forbidden(string(ErrorSlugUserCantSeePerformance), err, w, r)

		return
	}

	if app.IsPerformanceNotFoundError(err) {
		httperr.NotFound(string(ErrorSlugPerformanceNotFound), err, w, r)

		return
	}

	if performance.IsNotStartedError(err) {
		httperr.Conflict(string(ErrorSlugPerformanceNotStarted), err, w, r)

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

func (h handler) logMessages(r *http.Request, messages <-chan app.Message, extraFields ...app.LogField) {
	extraFields = append(extraFields, app.StringLogField("requestId", middleware.GetReqID(r.Context())))

	for msg := range messages {
		if msg.Err() == nil || msg.State() == performance.Failed {
			h.logger.Info(msg.String(), extraFields...)

			continue
		}

		h.logger.Error(msg.String(), msg.Err(), extraFields...)
	}
}
