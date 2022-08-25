package api

import (
	"net/http"
	"strings"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	util "github.com/AVENTER-UG/util/util"
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

	logrus.Debug("HTTP PUT V0ComposeRestartService")
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
		logrus.WithField("func", "V0ComposeRestartService").Error("Could not find Project or Servicename")

		d = e.ErrorMessage(0, "V0ComposeRestartService", "Could not find Project or Servicename")
		w.Write(d)
	}

	for keys.Next(e.Redis.CTX) {
		key := e.Redis.GetRedisKey(keys.Val())
		oldTask := mesosutil.DecodeTask(key)
		newTask := oldTask

		// generate new task as copy of old task
		taskName := strings.Split(oldTask.TaskID, ".")
		uuid, _ := util.GenUUID()
		newTask.TaskID = taskName[0] + "." + uuid
		newTask.State = ""
		e.Redis.SaveTaskRedis(newTask)

		// set the old task to be killed
		oldTask.State = "__KILL"
		e.Redis.SaveTaskRedis(oldTask)
	}

	d = e.ErrorMessage(0, "V0ComposeRestartService", "ok")
	w.Write(d)
}
