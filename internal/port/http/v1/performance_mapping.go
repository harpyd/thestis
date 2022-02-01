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
