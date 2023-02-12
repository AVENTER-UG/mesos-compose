package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeKillTask will kill a task from a service from a specific project
// example:
// curl -X DELETE http://user:password@127.0.0.1:10000/api/compose/v0/{projectname}/{servicename}/{taskid}
func (e *API) V0ComposeKillTask(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "api.V0ComposeKillTask").Debug("Kill Task")

	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if vars["project"] == "" || vars["servicename"] == "" || vars["taskid"] == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	project := vars["project"]
	servicename := vars["servicename"]
	taskID := vars["taskid"]

	key := e.Redis.GetRedisKey(e.Config.PrefixTaskName + ":" + project + ":" + servicename + ":" + taskID)

	task := e.Mesos.DecodeTask(key)
	if task.TaskID == taskID {
		task.State = "__KILL"
		task.Restart = "no"
		e.Redis.SaveTaskRedis(task)
		e.Mesos.SuppressFramework()

		logrus.WithField("func", "api.V0ComposeKillTask").Info("Kill Task: " + e.Config.PrefixTaskName + ":" + project + ":" + servicename + ":" + taskID)
	}
}
