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
	vars := mux.Vars(r)
	auth := e.CheckAuth(r, w)

	logrus.Debug("HTTP PUT V0ComposePush")
	d := e.ErrorMessage(2, "V0ComposePush", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars == nil || !auth {
		w.Write(d)
		return
	}

	var data cfg.Compose

	err := yaml.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		logrus.Error("Error: ", err.Error())
		d = e.ErrorMessage(2, "V0ComposePush", err.Error())
		w.Write(d)
		return
	}

	for service := range data.Services {
		// only schedule if the max instances is not reached
		TaskName := e.Config.PrefixTaskName + ":" + vars["project"] + ":" + service
		Instances := e.getReplicas()
		if e.Redis.CountRedisKey(TaskName+":*", "") < Instances {
			e.Compose = data
			e.mapComposeServiceToMesosTask(vars, service, cfg.Command{})
		}
	}

	out, _ := json.Marshal(&data)
	d = e.ErrorMessage(0, "V0ComposePush", util.PrettyJSON(out))
	w.Write(d)
}
