package v1

import "net/http"

func (h handler) LoadSpecification(w http.ResponseWriter, r *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) GetSpecification(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h handler) RemoveSpecification(w http.ResponseWriter, _ *http.Request, _ string) {
	w.WriteHeader(http.StatusNotImplemented)
}
