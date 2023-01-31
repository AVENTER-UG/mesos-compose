package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// V0FrameworkSupress suppress the framework
// example:
// curl -X PUT http://user:password@127.0.0.1:10000/api/compose/v0/framework/supress
func (e *API) V0FrameworkSupress(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "api.V0FrameworkSupress").Debug("Suppress Framework")

	auth := e.CheckAuth(r, w)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	e.Mesos.SuppressFramework()
}
