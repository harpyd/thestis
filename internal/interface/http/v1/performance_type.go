package v1

import (
	"net/http"

	"github.com/harpyd/thestis/internal/app"
)

func decodeStartPerformanceCommand(
	w http.ResponseWriter, r *http.Request,
	testCampaignID string,
) (cmd app.StartPerformanceCommand, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return app.StartPerformanceCommand{
		TestCampaignID: testCampaignID,
		StartedByID:    user.UUID,
	}, true
}

func decodeRestartPerformanceCommand(
	w http.ResponseWriter,
	r *http.Request,
	performanceID string,
) (cmd app.RestartPerformanceCommand, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return app.RestartPerformanceCommand{
		PerformanceID: performanceID,
		StartedByID:   user.UUID,
	}, true
}

func decodeCancelPerformanceCommand(
	w http.ResponseWriter,
	r *http.Request,
	performanceID string,
) (cmd app.CancelPerformanceCommand, ok bool) {
	user, ok := authorize(w, r)
	if !ok {
		return
	}

	return app.CancelPerformanceCommand{
		PerformanceID: performanceID,
		CanceledByID:  user.UUID,
	}, true
}
