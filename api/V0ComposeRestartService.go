package api

import (
	"encoding/json"
	"net/http"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	util "github.com/AVENTER-UG/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeRestartService will restart a service from a specific project
// example:
// curl -X PUT http://user:password@127.0.0.1:10000/v0/compose/{projectname}/{servicename}/restart
func V0ComposeRestartService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	logrus.Debug("HTTP PUT V0ComposeRestartService")
	d := ErrorMessage(2, "V0ComposePush", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars["project"] == "" || vars["servicename"] == "" {
		w.Write(d)
		return
	}

	project := vars["project"]
	servicename := vars["servicename"]

	keys := GetAllRedisKeys(config.PrefixTaskName + ":" + project + ":" + servicename + ":*")

	for keys.Next(config.RedisCTX) {
		key := GetRedisKey(keys.Val())
		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)

		// start a copy of this task and kill the old one
		// TODO: green blue
		mesosutil.Kill(task.TaskID, task.Agent)
		task.State = ""
		task.TaskID, _ = util.GenUUID()
		data, _ := json.Marshal(task)
		err := config.RedisClient.Set(config.RedisCTX, task.TaskName+":"+task.TaskID, data, 0).Err()
		if err != nil {
			d = ErrorMessage(2, "V0ComposeRestartService", err.Error())
			logrus.Error("V0ComposeRestartService Redis set Error: ", err)
			w.Write(d)
		}
	}

	d = ErrorMessage(0, "V0ComposeRestartService", "ok")
	w.Write(d)
}
