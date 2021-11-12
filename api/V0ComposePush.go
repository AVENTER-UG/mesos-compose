package api

import (
	"net/http"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// V0ComposePush will read and interpret the docker-compose.yml
// example:
// curl -X GET http://user:password@127.0.0.1:10000/v0/compose --data-binary @docker-compose.yml
func V0ComposePush(w http.ResponseWriter, r *http.Request) {
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
	logrus.Println(data)

	for service := range data.Services {
		mapComposeServiceToMesosTask(data.Services[service], data, vars, service, "")
	}
}
