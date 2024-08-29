package api

import (
	"net/http"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeKillTask will kill a task
// example:
// curl -X PUT http://user:password@127.0.0.1:10000/api/compose/v0/tasks/{taskid}/restart
func (e *API) V0ComposeRestartTask(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "api.V0ComposeRestartTask").Debug("Restart Mesos Task")

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
		logrus.WithField("func", "api.V0ComposeRestartTask").Errorf("Could not find TaskID (%s)", taskID)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	logrus.WithField("func", "api.V0ComposeRestartTask").Infof("Restart Task (%s)", taskID)

	// copy task
	newTask := new(cfg.Command)
	*newTask = *task
	task.State = "__KILL"
	task.Restart = "no"
	e.Redis.SaveTaskRedis(task)

	// generate new task as copy of old task
	newTask.TaskID = e.IncreaseTaskCount(newTask.TaskID)
	newTask.State = ""
	e.Redis.SaveTaskRedis(newTask)
}
