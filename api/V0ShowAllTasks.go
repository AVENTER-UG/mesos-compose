package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// V0ShowAllTasks will print out all tasks
// example:
// curl -X GET http://user:password@127.0.0.1:10000/api/compose/v0/tasks
func V0ShowAllTasks(w http.ResponseWriter, r *http.Request) {
	auth := CheckAuth(r, w)

	if !auth {
		return
	}

	logrus.Debug("HTTP GET V0ShowAllTasks")

	keys := GetAllRedisKeys("*")

	for keys.Next(config.RedisCTX) {
		logrus.Info("keys: ", keys.Val())
	}
}
