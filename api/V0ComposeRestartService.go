package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeRestartService will restart a service from a specific project
// example:
// curl -X PUT http://user:password@127.0.0.1:10000/api/compose/v0/{projectname}/{servicename}/restart
func (e *API) V0ComposeRestartService(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "api.V0ComposeRestartService").Debug("Restart Mesos Task")

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

	if keys == nil {
		logrus.WithField("func", "api.V0ComposeRestartService").Error("Could not find Project or Servicename")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	logrus.WithField("func", "api.V0ComposeRestartService").Info("Restart Task" + e.Config.PrefixTaskName + ":" + project + ":" + servicename)

	for keys.Next(e.Redis.CTX) {
		key := e.Redis.GetRedisKey(keys.Val())
		task := e.Mesos.DecodeTask(key)
		newTask := task
		task.State = "__KILL"
		task.Restart = "no"
		e.Redis.SaveTaskRedis(task)

		// generate new task as copy of old task
		newTask.TaskID = e.IncreaseTaskCount(newTask.TaskID)
		newTask.State = ""
		e.Redis.SaveTaskRedis(newTask)
	}
}
