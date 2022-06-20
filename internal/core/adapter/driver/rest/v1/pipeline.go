package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/adapter/driver/rest"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

func (h handler) StartPipeline(w http.ResponseWriter, r *http.Request, testCampaignID string) {
	cmd, ok := decodeStartPipelineCommand(w, r, uuid.New().String(), testCampaignID)
	if !ok {
		return
	}

	err := h.app.Commands.StartPipeline.Handle(r.Context(), cmd)
	if err == nil {
		w.Header().Set("Location", fmt.Sprintf("/pipelines/%s", cmd.PipelineID))
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

func (h handler) RestartPipeline(w http.ResponseWriter, r *http.Request, pipelineID string) {
	cmd, ok := decodeRestartPipelineCommand(w, r, pipelineID)
	if !ok {
		return
	}

	err := h.app.Commands.RestartPipeline.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		rest.Forbidden(string(ErrorSlugUserCantSeePipeline), err, w, r)

		return
	}

	if errors.Is(err, service.ErrPipelineNotFound) {
		rest.NotFound(string(ErrorSlugPipelineNotFound), err, w, r)

		return
	}

	if errors.Is(err, pipeline.ErrAlreadyStarted) {
		rest.Conflict(string(ErrorSlugPipelineAlreadyStarted), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) CancelPipeline(w http.ResponseWriter, r *http.Request, pipelineID string) {
	cmd, ok := decodeCancelPipelineCommand(w, r, pipelineID)
	if !ok {
		return
	}

	err := h.app.Commands.CancelPipeline.Handle(r.Context(), cmd)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	var aerr *user.AccessError

	if errors.As(err, &aerr) {
		rest.Forbidden(string(ErrorSlugUserCantSeePipeline), err, w, r)

		return
	}

	if errors.Is(err, service.ErrPipelineNotFound) {
		rest.NotFound(string(ErrorSlugPipelineNotFound), err, w, r)

		return
	}

	if errors.Is(err, pipeline.ErrNotStarted) {
		rest.Conflict(string(ErrorSlugPipelineNotStarted), err, w, r)

		return
	}

	rest.InternalServerError(string(ErrorSlugUnexpectedError), err, w, r)
}

func (h handler) GetPipelineHistory(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetPipeline(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
