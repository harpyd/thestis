package v1

import (
	"net/http"

	"github.com/harpyd/thestis/internal/app"
)

func unmarshalStartNewPerformanceCommand(
	w http.ResponseWriter, r *http.Request,
	testCampaignID string,
) (cmd app.StartNewPerformanceCommand, ok bool) {
	user, ok := unmarshalUser(w, r)
	if !ok {
		return
	}

	return app.StartNewPerformanceCommand{
		TestCampaignID: testCampaignID,
		StartedByID:    user.UUID,
	}, true
}
