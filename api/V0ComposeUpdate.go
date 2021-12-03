package api

import (
	"encoding/json"
	"net/http"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	util "github.com/AVENTER-UG/util"
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

	logrus.Debug("HTTP PUT V0ComposeUpdate")
	d := ErrorMessage(2, "V0ComposeUpdate", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars == nil || !auth {
		w.Write(d)
		return
	}

	var data cfg.Compose

	err := yaml.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		d := ErrorMessage(2, "V0ComposeUpdate", err.Error())
		w.Write(d)
		logrus.Error("Error: ", err)
	}

	for service := range data.Services {
		taskName := config.PrefixTaskName + ":" + vars["project"] + ":" + service
		// get all keys that start with the taskname
		keys := GetAllRedisKeys(taskName + ":*")

		for keys.Next(config.RedisCTX) {
			// get the values of the current key
			key := GetRedisKey(keys.Val())
			var task mesosutil.Command
			json.Unmarshal([]byte(key), &task)
			mapComposeServiceToMesosTask(data.Services[service], data, vars, service, task)
		}
	}

	out, _ := json.Marshal(&data)
	d = ErrorMessage(0, "V0ComposeUpdate", util.PrettyJSON(out))
	w.Write(d)
}
