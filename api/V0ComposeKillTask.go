package api

import (
	"net/http"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeKillTask will kill a task from a service from a specific project
// example:
// curl -X DELETE http://user:password@127.0.0.1:10000/api/compose/v0/{projectname}/{servicename}/{taskid}
func (e *API) V0ComposeKillTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	logrus.Debug("HTTP DELETE V0ComposeKillTask")
	d := e.ErrorMessage(2, "V0ComposeKillTask", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars["project"] == "" || vars["servicename"] == "" || vars["taskid"] == "" {
		w.Write(d)
		return
	}

	project := vars["project"]
	servicename := vars["servicename"]
	taskID := vars["taskid"]

	key := e.Redis.GetRedisKey(e.Config.PrefixTaskName + ":" + project + ":" + servicename + ":" + taskID)

	task := mesosutil.DecodeTask(key)
	if task.TaskID == taskID {
		// Kill the current task
		mesosutil.Kill(task.TaskID, task.Agent)

		logrus.Debug("V0ComposeKillTask: " + e.Config.PrefixTaskName + ":" + project + ":" + servicename + ":" + taskID)
		d = e.ErrorMessage(0, "V0ComposeKillTask", "ok")
	}
	w.Write(d)
}
