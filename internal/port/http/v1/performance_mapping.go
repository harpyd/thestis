package v1

import (
	"net/http"

	"github.com/harpyd/thestis/internal/app"
)

func unmarshalStartPerformanceCommand(
	w http.ResponseWriter, r *http.Request,
	testCampaignID string,
) (cmd app.StartPerformanceCommand, ok bool) {
	user, ok := unmarshalUser(w, r)
	if !ok {
		return
	}

	return app.StartPerformanceCommand{
		TestCampaignID: testCampaignID,
		StartedByID:    user.UUID,
	}, true
}

func unmarshalRestartPerformanceCommand(
	w http.ResponseWriter,
	r *http.Request,
	performanceID string,
) (cmd app.RestartPerformanceCommand, ok bool) {
	user, ok := unmarshalUser(w, r)
	if !ok {
		return
	}

	return app.RestartPerformanceCommand{
		PerformanceID: performanceID,
		StartedByID:   user.UUID,
	}, true
}

func unmarshalCancelPerformanceCommand(
	w http.ResponseWriter,
	r *http.Request,
	performanceID string,
) (cmd app.CancelPerformanceCommand, ok bool) {
	user, ok := unmarshalUser(w, r)
	if !ok {
		return
	}

	return app.CancelPerformanceCommand{
		PerformanceID: performanceID,
		CanceledByID:  user.UUID,
	}, false
}
