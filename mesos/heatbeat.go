package mesos

import (
	"encoding/json"

	api "github.com/AVENTER-UG/mesos-compose/api"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	"github.com/sirupsen/logrus"
)

func Heartbeat() {
	keys := api.GetAllRedisKeys("*")
	for keys.Next(config.RedisCTX) {
		logrus.Info("keys: ", keys.Val())
		// get the values of the current key
		key := api.GetRedisKey(keys.Val())

		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)

		if task.State == "" {
			framework.CommandChan <- task
			logrus.Info("Scheduled Mesos Task")

		}

	}
}
