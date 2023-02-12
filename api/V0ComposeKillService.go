package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeKillService will kill a task from a service from a specific project
// example:
// curl -X DELETE http://user:password@127.0.0.1:10000/api/compose/v0/{projectname}/{servicename}/{taskid}
func (e *API) V0ComposeKillService(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "api.V0ComposeKillService").Debug("Kill Service")

	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if vars["project"] == "" || vars["servicename"] == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	project := vars["project"]
	servicename := vars["servicename"]

	keys := e.Redis.GetAllRedisKeys(e.Config.PrefixTaskName + ":" + project + ":" + servicename + ":*")

	for keys.Next(e.Redis.CTX) {
		key := e.Redis.GetRedisKey(keys.Val())

		task := e.Mesos.DecodeTask(key)
		task.State = "__KILL"
		task.Restart = "no"
		e.Redis.SaveTaskRedis(task)
		e.Mesos.SuppressFramework()

		logrus.WithField("func", "api.V0ComposeKillService").Info("Kill Service: " + e.Config.PrefixTaskName + ":" + project + ":" + servicename)
	}
}
