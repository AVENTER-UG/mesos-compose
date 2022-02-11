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

// V0ComposePush will read and interpret the docker-compose.yml
// example:
// curl -X GET http://user:password@127.0.0.1:10000/api/compose/v0 --data-binary @docker-compose.yml
func V0ComposePush(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := CheckAuth(r, w)

	logrus.Debug("HTTP PUT V0ComposePush")
	d := ErrorMessage(2, "V0ComposePush", "nok")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	if vars == nil || !auth {
		w.Write(d)
		return
	}

	var data cfg.Compose

	err := yaml.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		d = ErrorMessage(2, "V0ComposePush", err.Error())
		logrus.Error("Error: ", err)
		w.Write(d)
		return
	}

	for service := range data.Services {
		mapComposeServiceToMesosTask(data.Services[service], data, vars, service, mesosutil.Command{})
	}

	out, _ := json.Marshal(&data)
	d = ErrorMessage(0, "V0ComposePush", util.PrettyJSON(out))
	w.Write(d)
}
