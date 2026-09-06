package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AVENTER-UG/mesos-compose/redis"
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if !auth {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if vars["project"] == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	composeYAML, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var data cfg.Compose

	err = yaml.Unmarshal(composeYAML, &data)

	if err != nil {
		logrus.WithField("func", "api.V0ComposeUpdate").Error("Error: ", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	e.Compose = data

	logrus.WithField("func", "api.V0ComposeUpdate").Info("Update Mesos Task: " + e.Config.PrefixTaskName + ":" + vars["project"])

	for service := range data.Services {
		taskName := e.Config.PrefixTaskName + ":" + vars["project"] + ":" + service
		// get all keys that start with the taskname
		keys := e.Redis.GetAllRedisKeys(taskName + ":*")

		for keys.Next(e.Redis.CTX) {
			// get the values of the current key
			key := e.Redis.GetRedisKey(keys.Val())
			task := redis.DecodeTaskOrEmpty([]byte(key))
			e.mapComposeServiceToMesosTask(vars, service, task)

			// restore the old MesosAgent info
			key = e.Redis.GetRedisKey(keys.Val())
			updatedTask := redis.DecodeTaskOrEmpty([]byte(key))
			updatedTask.MesosAgent = task.MesosAgent
			updatedTask.State = task.State
			e.Redis.SaveTaskRedis(updatedTask)
		}
	}
	e.Redis.SaveComposeYAML(vars["project"], composeYAML)

	out, _ := json.Marshal(&data)
	w.Write([]byte(util.PrettyJSON(out)))
}
