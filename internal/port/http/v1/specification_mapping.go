package v1

import (
	"io"
	"net/http"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/pkg/httperr"
)

func unmarshalSpecificationSourceCommand(
	w http.ResponseWriter,
	r *http.Request,
	testCampaignID string,
) (cmd app.LoadSpecificationCommand, ok bool) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		httperr.BadRequest(string(ErrorSlugBadRequest), err, w, r)

		return
	}

	return app.LoadSpecificationCommand{
		Content:        content,
		TestCampaignID: testCampaignID,
	}, true
}
