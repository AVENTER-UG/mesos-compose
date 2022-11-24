package api

import (
	"encoding/json"
	"net/http"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	util "github.com/AVENTER-UG/util/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// V0ComposeUpdate will update the docker-compose.yml
// example:
// curl -X GET http://user:password@127.0.0.1:10000/api/compose/v0 --data-binary @docker-compose.yml
func (e *API) V0ComposeUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	logrus.Debug("HTTP PUT V0ComposeUpdate")
	d := e.ErrorMessage(2, "V0ComposeUpdate", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars == nil || !auth {
		w.Write(d)
		return
	}

	var data cfg.Compose

	err := yaml.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		d = e.ErrorMessage(2, "V0ComposeUpdate", err.Error())
		w.Write(d)
		logrus.Error("Error: ", err)
	}

	e.Compose = data

	for service := range data.Services {
		taskName := e.Config.PrefixTaskName + ":" + vars["project"] + ":" + service
		// get all keys that start with the taskname
		keys := e.Redis.GetAllRedisKeys(taskName + ":*")

		for keys.Next(e.Redis.CTX) {
			// get the values of the current key
			key := e.Redis.GetRedisKey(keys.Val())
			task := e.Mesos.DecodeTask(key)
			e.mapComposeServiceToMesosTask(vars, service, task)

			// restore the old MesosAgent info
			key = e.Redis.GetRedisKey(keys.Val())
			updatedTask := e.Mesos.DecodeTask(key)
			updatedTask.MesosAgent = task.MesosAgent
			e.Redis.SaveTaskRedis(updatedTask)
		}
	}

	out, _ := json.Marshal(&data)
	d = e.ErrorMessage(0, "V0ComposeUpdate", util.PrettyJSON(out))
	w.Write(d)
}
