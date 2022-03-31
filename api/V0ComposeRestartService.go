package api

import (
	"encoding/json"
	"net/http"
	"strings"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	util "github.com/AVENTER-UG/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ComposeRestartService will restart a service from a specific project
// example:
// curl -X PUT http://user:password@127.0.0.1:10000/api/compose/v0/{projectname}/{servicename}/restart
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
		var oldTask mesosutil.Command
		json.Unmarshal([]byte(key), &oldTask)
		newTask := oldTask

		// generate new task as copy of old task
		taskName := strings.Split(oldTask.TaskID, ".")
		uuid, _ := util.GenUUID()
		newTask.TaskID = taskName[0] + "." + uuid
		newTask.State = ""
		data, _ := json.Marshal(newTask)
		err := config.RedisClient.Set(config.RedisCTX, newTask.TaskName+":"+newTask.TaskID, data, 0).Err()
		if err != nil {
			d = ErrorMessage(2, "V0ComposeRestartService newTask", err.Error())
			logrus.WithField("func", "V0ComposeRestartService").Error("Redis Error write newTask Data: ", err.Error())
			w.Write(d)
		}

		// set the old task to be killed
		oldTask.State = "__KILL"
		data, _ = json.Marshal(oldTask)
		err = config.RedisClient.Set(config.RedisCTX, oldTask.TaskName+":"+oldTask.TaskID, data, 0).Err()
		if err != nil {
			d = ErrorMessage(2, "V0ComposeRestartService oldTask", err.Error())
			logrus.WithField("func", "V0ComposeRestartService").Error("Redis Error write oldTask Data: ", err.Error())
			w.Write(d)
		}
	}

	d = ErrorMessage(0, "V0ComposeRestartService", "ok")
	w.Write(d)
}
