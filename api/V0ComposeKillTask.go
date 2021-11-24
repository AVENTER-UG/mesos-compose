package api

import (
	"encoding/json"
	"net/http"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeKillTask will kill a task from a service from a specific project
// example:
// curl -X DELETE http://user:password@127.0.0.1:10000/v0/compose/{projectname}/{servicename}/{taskid}
func V0ComposeKillTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	logrus.Debug("HTTP DELETE V0ComposeKillTask")
	d := ErrorMessage(2, "V0ComposeKillTask", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars["project"] == "" || vars["servicename"] == "" || vars["taskid"] == "" {
		w.Write(d)
		return
	}

	project := vars["project"]
	servicename := vars["servicename"]
	taskID := vars["taskid"]

	key := GetRedisKey(config.PrefixTaskName + "_" + project + "_" + servicename + ":" + taskID)

	var task mesosutil.Command
	json.Unmarshal([]byte(key), &task)
	logrus.Debug(task)
	if task.TaskID == taskID {
		err := mesosutil.Kill(task.TaskID, task.Agent)
		if err != nil {
			d = ErrorMessage(2, "V0ComposeKillTask", err.Error())
			logrus.Error("V0ComposeKillTask Error during kill: ", err)
			w.Write(d)
			return
		}
		logrus.Debug("V0ComposeKillTask: " + config.PrefixTaskName + "_" + project + "_" + servicename + ":" + taskID)
		d = ErrorMessage(0, "V0ComposeKillTask", "ok")
	}
	w.Write(d)
	return
}
