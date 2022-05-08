package v1

import (
	"net/http"

	"github.com/harpyd/thestis/internal/core/app/command"
)

func decodeStartPerformanceCommand(
	w http.ResponseWriter, r *http.Request,
	testCampaignID string,
) (cmd command.StartPerformance, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return command.StartPerformance{
		TestCampaignID: testCampaignID,
		StartedByID:    user.UUID,
	}, true
}

func decodeRestartPerformanceCommand(
	w http.ResponseWriter,
	r *http.Request,
	performanceID string,
) (cmd command.RestartPerformance, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return command.RestartPerformance{
		PerformanceID: performanceID,
		StartedByID:   user.UUID,
	}, true
}

func decodeCancelPerformanceCommand(
	w http.ResponseWriter,
	r *http.Request,
	performanceID string,
) (cmd command.CancelPerformance, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return command.CancelPerformance{
		PerformanceID: performanceID,
		CanceledByID:  user.UUID,
	}, true
}
