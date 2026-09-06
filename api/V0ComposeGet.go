package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// V0ComposeGet returns the raw compose YAML stored for a project.
func (e *API) V0ComposeGet(w http.ResponseWriter, r *http.Request) {
	if !e.CheckAuth(r, w) {
		return
	}

	project := mux.Vars(r)["project"]
	if project == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	composeYAML, ok := e.Redis.GetComposeYAML(project)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
	w.Header().Set("Api-Service", "v0")
	_, _ = w.Write([]byte(composeYAML))
}
