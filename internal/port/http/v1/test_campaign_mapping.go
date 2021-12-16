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

	var summary string
	if rb.Summary != nil {
		summary = *rb.Summary
	}

	return app.CreateTestCampaignCommand{
		ViewName: rb.ViewName,
		Summary:  summary,
	}, true
}
