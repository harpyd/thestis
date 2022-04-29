package v1

import (
	"fmt"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/domain/performance"
	"github.com/harpyd/thestis/internal/core/domain/user"
	"github.com/harpyd/thestis/internal/core/interface/http"
)

func (h handler) StartPerformance(w stdhttp.ResponseWriter, r *stdhttp.Request, testCampaignID string) {
	cmd, ok := decodeStartPerformanceCommand(w, r, testCampaignID)
	if !ok {
		return
	}

	perfID, messages, err := h.app.Commands.StartPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Location", fmt.Sprintf("/performances/%s", perfID))
		w.WriteHeader(stdhttp.StatusAccepted)

		go h.logMessages(
			r,
			messages,
			app.StringLogField("performanceId", perfID),
			app.BoolLogField("restarted", false),
		)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		http.Forbidden(string(ErrorSlugUserCantSeeTestCampaign), err, w, r)

		return
	}

	if errors.Is(err, app.ErrSpecificationNotFound) {
		http.NotFound(string(ErrorSlugSpecificationNotFound), err, w, r)

		return
	}

	http.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) RestartPerformance(w stdhttp.ResponseWriter, r *stdhttp.Request, performanceID string) {
	cmd, ok := decodeRestartPerformanceCommand(w, r, performanceID)
	if !ok {
		return
	}

	messages, err := h.app.Commands.RestartPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(stdhttp.StatusNoContent)

		go h.logMessages(r, messages, app.BoolLogField("restarted", true))

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		http.Forbidden(string(ErrorSlugUserCantSeePerformance), err, w, r)

		return
	}

	if errors.Is(err, app.ErrPerformanceNotFound) {
		http.NotFound(string(ErrorSlugPerformanceNotFound), err, w, r)

		return
	}

	if errors.Is(err, performance.ErrAlreadyStarted) {
		http.Conflict(string(ErrorSlugPerformanceAlreadyStarted), err, w, r)

		return
	}

	http.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) CancelPerformance(w stdhttp.ResponseWriter, r *stdhttp.Request, performanceID string) {
	cmd, ok := decodeCancelPerformanceCommand(w, r, performanceID)
	if !ok {
		return
	}

	err := h.app.Commands.CancelPerformance.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(stdhttp.StatusNoContent)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		http.Forbidden(string(ErrorSlugUserCantSeePerformance), err, w, r)

		return
	}

	if errors.Is(err, app.ErrPerformanceNotFound) {
		http.NotFound(string(ErrorSlugPerformanceNotFound), err, w, r)

		return
	}

	if errors.Is(err, performance.ErrNotStarted) {
		http.Conflict(string(ErrorSlugPerformanceNotStarted), err, w, r)

		return
	}

	http.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetPerformanceHistory(w stdhttp.ResponseWriter, _ *stdhttp.Request, _ string) {
	w.WriteHeader(stdhttp.StatusNotImplemented)
}

func (h handler) GetPerformance(w stdhttp.ResponseWriter, _ *stdhttp.Request, _ string) {
	w.WriteHeader(stdhttp.StatusNotImplemented)
}

func (h handler) logMessages(r *stdhttp.Request, messages <-chan app.Message, extraFields ...app.LogField) {
	extraFields = append(extraFields, app.StringLogField("requestId", middleware.GetReqID(r.Context())))

	for msg := range messages {
		if msg.Err() == nil || msg.Event() == performance.FiredFail {
			h.logger.Info(msg.String(), extraFields...)

			continue
		}

		h.logger.Error(msg.String(), msg.Err(), extraFields...)
	}
}
