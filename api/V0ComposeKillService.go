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
	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	d := e.ErrorMessage(2, "V0ComposeKillService", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars["project"] == "" || vars["servicename"] == "" {
		w.Write(d)
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

		logrus.WithField("func", "api.V0ComposeKillService").Info("Kill Service: " + e.Config.PrefixTaskName + ":" + project + ":" + servicename)
	}
	d = e.ErrorMessage(0, "V0ComposeKillService", "ok")
	w.Write(d)
}
