package v1

import "net/http"

func (h handler) CreateTestCampaign(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetTestCampaigns(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetTestCampaign(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) RemoveTestCampaign(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
