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
	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	d := e.ErrorMessage(2, "V0ComposePush", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars["project"] == "" || vars["servicename"] == "" {
		w.Write(d)
		return
	}

	project := vars["project"]
	servicename := vars["servicename"]

	keys := e.Redis.GetAllRedisKeys(e.Config.PrefixTaskName + ":" + project + ":" + servicename + ":*")

	if keys == nil {
		logrus.WithField("func", "api.V0ComposeRestartService").Error("Could not find Project or Servicename")

		d = e.ErrorMessage(0, "V0ComposeRestartService", "Could not find Project or Servicename")
		w.Write(d)
		return
	}

	logrus.WithField("func", "api.V0ComposeRestartService").Info("Restart Task" + e.Config.PrefixTaskName + ":" + project + ":" + servicename)

	for keys.Next(e.Redis.CTX) {
		key := e.Redis.GetRedisKey(keys.Val())
		task := e.Mesos.DecodeTask(key)
		task.State = "__KILL"
		task.Restart = "no"
		e.Redis.SaveTaskRedis(task)

		// generate new task as copy of old task
		task.TaskID = e.IncreaseTaskCount(task.TaskID)
		task.State = ""
		e.Redis.SaveTaskRedis(task)
	}

	d = e.ErrorMessage(0, "V0ComposeRestartService", "ok")
	w.Write(d)
}
