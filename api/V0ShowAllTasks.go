package api

import (
	"encoding/json"
	"net/http"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	util "github.com/AVENTER-UG/util/util"
	"github.com/sirupsen/logrus"
)

// V0ShowAllTasks will print out all tasks
// example:
// curl -X GET http://user:password@127.0.0.1:10000/api/compose/v0/tasks
func (e *API) V0ShowAllTasks(w http.ResponseWriter, r *http.Request) {
	auth := e.CheckAuth(r, w)

	if !auth {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	logrus.WithField("func", "api.V0ShowAllTasks").Debug("HTTP GET V0ShowAllTasks")

	keys := e.Redis.GetAllRedisKeys(e.Framework.FrameworkName + ":*")

	var list []cfg.Command

	for keys.Next(e.Redis.CTX) {
		// ignore redis keys if they are not mesos tasks
		if e.Redis.CheckIfNotTask(keys) {
			continue
		}

		key := e.Redis.GetRedisKey(keys.Val())
		task := e.Mesos.DecodeTask(key)

		task.Environment = mesosproto.Environment{}

		list = append(list, task)
	}

	out, _ := json.Marshal(&list)
	d := e.ErrorMessage(0, "V0ShowAllTasks", util.PrettyJSON(out))
	w.Write(d)
}
