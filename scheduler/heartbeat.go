package scheduler

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// Heartbeat - The Apache Mesos heatbeat function
// nolint:gocyclo
func (e *Scheduler) Heartbeat() {
	// Check Connection state of Redis
	err := e.Redis.PingRedis()
	if err != nil {
		e.Redis.Connect()
	}

	keys := e.Redis.GetAllRedisKeys(e.Framework.FrameworkName + ":*")
	suppress := true
	for keys.Next(e.Redis.CTX) {
		// continue if the key is not a mesos task
		if e.Redis.CheckIfNotTask(keys) {
			continue
		}
		// get the values of the current key
		key := e.Redis.GetRedisKey(keys.Val())

		task := e.Mesos.DecodeTask(key)

		if task.TaskID == "" || task.TaskName == "" {
			continue
		}

		if task.State == "__KILL" {
			// if agent is unset, the task is not running we can just delete the DB key
			if task.Agent == "" {
				e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			} else {
				e.Mesos.Kill(task.TaskID, task.Agent)
				e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			}
			continue
		}

		if e.Redis.CountRedisKey(task.TaskName+":*", "__KILL") < task.Instances {
			e.Mesos.Revive()
			task.TaskID = task.TaskID + "." + strconv.Itoa(e.Redis.CountRedisKey(task.TaskName+":*", "__KILL")+1)
			task.State = ""
			e.Redis.SaveTaskRedis(task)
			continue
		}

		if task.State == "" && e.Redis.CountRedisKey(task.TaskName+":*", "__KILL") <= task.Instances {
			e.Mesos.Revive()
			task.State = "__NEW"
			// these will save the current time at the task. we need it to check
			// if the state will change in the next 'n min. if not, we have to
			// give these task a recall.
			task.StateTime = time.Now()

			// Change the Dynamic Host Ports
			task.DockerPortMappings = e.changeDockerPorts(task)
			task.Discovery = e.changeDiscoveryInfo(task)

			// add task to communication channel
			e.Framework.CommandChan <- task

			e.Redis.SaveTaskRedis(task)

			logrus.Info("Scheduled Mesos Task: ", task.TaskName)
			continue
		}

		if task.State == "__NEW" {
			e.Mesos.Revive()
			suppress = false
			e.Config.Suppress = false
			continue
		}

		// Remove corrupt tasks
		if task.State == "" && task.StateTime.Year() == 1 {
			e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			continue
		}
	}

	if suppress && !e.Config.Suppress {
		e.Mesos.SuppressFramework()
		e.Config.Suppress = true
	}
}

// HeartbeatLoop - The main loop for the hearbeat
func (e *Scheduler) HeartbeatLoop() {
	ticker := time.NewTicker(e.Config.EventLoopTime)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		e.Heartbeat()
	}
}

// ReconcileLoop - The reconcile loop to check periodicly the task state
func (e *Scheduler) ReconcileLoop() {
	ticker := time.NewTicker(e.Config.ReconcileLoopTime)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		go e.reconcile()
		go e.implicitReconcile()
	}
}
