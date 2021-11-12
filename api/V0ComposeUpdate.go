package api

import (
	"encoding/json"
	"net/http"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// V0ComposeUpdate will update the docker-compose.yml
// example:
// curl -X GET http://user:password@127.0.0.1:10000/v0/compose --data-binary @docker-compose.yml
func V0ComposeUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	var data cfg.Compose

	err := yaml.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		logrus.Error("Error: ", err)
	}

	for service := range data.Services {
		taskName := config.PrefixTaskName + "_" + vars["project"] + "_" + service
		// get all keys that start with the taskname
		keys := GetAllRedisKeys(taskName + ":*")

		for keys.Next(config.RedisCTX) {
			logrus.Info("keys: ", keys.Val())
			// get the values of the current key
			key := GetRedisKey(keys.Val())
			var task mesosutil.Command
			json.Unmarshal([]byte(key), &task)
			mapComposeServiceToMesosTask(data.Services[service], data, vars, service, task.TaskID)
		}
	}
}
