package mesos

import (
	"encoding/json"
	"time"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	"github.com/sirupsen/logrus"
)

// Heartbeat - The Apache Mesos heatbeat function
func (e *Scheduler) Heartbeat() {
	// Check Connection state of Redis
	err := e.API.PingRedis()
	if err != nil {
		e.API.ConnectRedis()
	}

	keys := e.API.GetAllRedisKeys(e.Framework.FrameworkName + ":*")
	suppress := true
	for keys.Next(e.API.Redis.RedisCTX) {
		// get the values of the current key
		key := e.API.GetRedisKey(keys.Val())

		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)

		if task.TaskID == "" || task.TaskName == "" {
			continue
		}

		if task.State == "" && e.API.CountRedisKey(task.TaskName+":*") <= task.Instances {
			mesosutil.Revive()
			task.State = "__NEW"
			// these will save the current time at the task. we need it to check
			// if the state will change in the next 'n min. if not, we have to
			// give these task a recall.
			task.StateTime = time.Now()

			// add task to communication channel
			e.Framework.CommandChan <- task

			data, _ := json.Marshal(task)
			err := e.API.Redis.RedisClient.Set(e.API.Redis.RedisCTX, task.TaskName+":"+task.TaskID, data, 0).Err()
			if err != nil {
				logrus.Error("HandleUpdate Redis set Error: ", err)
			}

			logrus.Info("Scheduled Mesos Task: ", task.TaskName)
		}

		if task.State == "__NEW" {
			suppress = false
			e.Config.Suppress = false
		}

		if task.State == "__KILL" {
			mesosutil.Kill(task.TaskID, task.Agent)
		}
	}

	if suppress && !e.Config.Suppress {
		mesosutil.SuppressFramework()
		e.Config.Suppress = true
	}
}

// HeartbeatLoop - The main loop for the hearbeat
func (e *Scheduler) HeartbeatLoop() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		e.Heartbeat()
	}
}
