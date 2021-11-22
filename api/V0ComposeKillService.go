package api

import (
	"encoding/json"
	"net/http"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeKillService will kill a task from a service from a specific project
// example:
// curl -X DELETE http://user:password@127.0.0.1:10000/v0/compose/{projectname}/{servicename}/{taskid}
func V0ComposeKillService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	logrus.Debug("HTTP DELETE V0ComposeKillService")
	d := []byte("nok")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars["project"] == "" || vars["servicename"] == "" || vars["taskid"] == "" {
		w.Write(d)
		return
	}

	project := vars["project"]
	servicename := vars["servicename"]
	taskID := vars["taskid"]

	keys := GetAllRedisKeys(config.PrefixTaskName + "_" + project + "_" + servicename + ":*")

	for keys.Next(config.RedisCTX) {
		logrus.Info("keys: ", keys.Val())
		key := GetRedisKey(keys.Val())

		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)
		logrus.Debug(task)
		mesosutil.Kill(task.TaskID, task.Agent)
		logrus.Debug("V0ComposeKillService: " + config.PrefixTaskName + "_" + project + "_" + servicename + ":" + taskID)
	}
	d = []byte("ok")
	w.Write(d)
}
