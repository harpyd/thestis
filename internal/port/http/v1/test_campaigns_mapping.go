package v1

import (
	"net/http"

	"github.com/harpyd/thestis/internal/app"
)

func unmarshalCreateTestCampaignCommand(
	w http.ResponseWriter,
	r *http.Request,
) (cmd app.CreateTestCampaignCommand, ok bool) {
	var rb CreateTestCampaignRequest

	if ok = decode(w, r, &rb); !ok {
		return
	}

	return app.CreateTestCampaignCommand{
		ViewName: rb.ViewName,
	}, true
}
