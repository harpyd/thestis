package v1

import (
	"net/http"

	"github.com/harpyd/thestis/internal/core/app/command"
)

func decodeStartPipelineCommand(
	w http.ResponseWriter, r *http.Request,
	pipelineID, testCampaignID string,
) (cmd command.StartPipeline, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return command.StartPipeline{
		PipelineID:     pipelineID,
		TestCampaignID: testCampaignID,
		StartedByID:    user.UUID,
	}, true
}

func decodeRestartPipelineCommand(
	w http.ResponseWriter,
	r *http.Request,
	pipelineID string,
) (cmd command.RestartPipeline, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return command.RestartPipeline{
		PipelineID:  pipelineID,
		StartedByID: user.UUID,
	}, true
}

func decodeCancelPipelineCommand(
	w http.ResponseWriter,
	r *http.Request,
	pipelineID string,
) (cmd command.CancelPipeline, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return command.CancelPipeline{
		PipelineID:   pipelineID,
		CanceledByID: user.UUID,
	}, true
}
