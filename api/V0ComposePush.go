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

// V0ComposePush will read and interpret the docker-compose.yml
// example:
// curl -X GET http://user:password@127.0.0.1:10000/api/compose/v0 --data-binary @docker-compose.yml
func (e *API) V0ComposePush(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("func", "api.V0ComposePush").Debug("Launch Mesos Task")

	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if vars["project"] == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var data cfg.Compose

	err := yaml.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		logrus.WithField("func", "V0ComposePush").Error("Error: ", err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	for service := range data.Services {
		// only schedule if the max instances is not reached
		taskName := e.Config.PrefixTaskName + ":" + vars["project"] + ":" + service
		logrus.WithField("func", "api.V0ComposePush").Info("Push Task: " + taskName)

		Instances := e.getReplicas()
		if e.Redis.CountRedisKey(taskName+":*", "") < Instances {
			e.Compose = data
			e.mapComposeServiceToMesosTask(vars, service, cfg.Command{})
		}
	}

	out, _ := json.Marshal(&data)
	w.Write([]byte(util.PrettyJSON(out)))
}
