package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// V0FrameworkReRegister will cleanup framework information in redis to force a register in mesos
// example:
// curl -X PUT http://user:password@127.0.0.1:10000/api/compose/v0/framework/reregister
func (e *API) V0FrameworkReRegister(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "api.V0FrameworkReRegister").Debug("ReRegister Framework")

	auth := e.CheckAuth(r, w)

	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	e.Redis.DelRedisKey(e.Framework.FrameworkName + ":framework")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")
}
