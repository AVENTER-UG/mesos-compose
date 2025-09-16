package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeRestartTask will restart a task
// example:
// curl -X DELETE http://user:password@127.0.0.1:10000/api/compose/v0/tasks/{taskid}
func (e *API) V0ComposeKillTask(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "api.V0ComposeKillTask").Debug("Kill Mesos Task")

	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if vars["taskid"] == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	taskID := vars["taskid"]

	task := e.Redis.GetTaskFromTaskID(taskID)

	if task.TaskID == "" {
		logrus.WithField("func", "api.V0ComposeKillTask").Errorf("Could not find TaskID (%s)", taskID)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	logrus.WithField("func", "api.V0ComposeKillTask").Infof("Kill Task (%s)", taskID)

	task.State = "__KILL"
	task.ExpectedState = "__KILL"
	task.Restart = "no"
	task.Instances = 0
	e.Redis.SaveTaskRedis(task)
}
