package scheduler

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// Heartbeat - The Apache Mesos heatbeat function
// nolint:gocyclo
func (e *Scheduler) Heartbeat() {
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

		// kill task
		if task.State == "__KILL" {
			// if agent is unset, the task is not running we can just delete the DB key
			if task.Agent == "" {
				e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
				continue
			} else {
				task.Restart = "no"
				e.Redis.SaveTaskRedis(task)
				e.Mesos.Kill(task.TaskID, task.Agent)
			}
			continue
		}

		// there are lesser instances are running as it should be
		if e.Redis.CountRedisKey(task.TaskName+":*", "") < task.Instances {
			// do not schedule task if the placement is unique and the amount of running tasks
			// are same like the amount of mesos agents.
			if e.getLabelValue("__mc_placement", task) == "unique" && e.Redis.CountRedisKeyState(task.TaskName+":*", "TASK_RUNNING") >= e.Mesos.CountAgent {
				continue
			}

			logrus.WithField("func", "scheduler.CheckState").Info("Scale up Mesos Task: ", task.TaskName)
			e.Mesos.Revive()

			taskNew := *task
			taskNew.State = ""
			taskNew.TaskID = e.API.IncreaseTaskCount(taskNew.TaskID)
			e.Redis.SaveTaskRedis(&taskNew)
			continue
		}

		// there are more instances are running as it should be
		if e.Redis.CountRedisKey(task.TaskName+":*", "__KILL") > task.Instances {
			logrus.WithField("func", "scheduler.CheckState").Info("Scale down Mesos Task: ", task.TaskName)
			e.Mesos.Revive()
			task.State = "__KILL"
			e.Redis.SaveTaskRedis(task)
			continue
		}

		// Schedule task
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

			e.Redis.SaveTaskRedis(task)

			// add task to communication channel
			e.Framework.CommandChan <- *task

			logrus.WithField("func", "scheduler.CheckState").Info("Scheduled Mesos Task: ", task.TaskName)
			continue
		}

		// Revieve Mesos to get new offers
		if task.State == "__NEW" {
			e.Mesos.IsRevive = false
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

// Check Connection state of Redis
func (e *Scheduler) checkRedis() {
	err := e.Redis.PingRedis()
	if err != nil {
		e.Redis.Connect()
	}
}

// HeartbeatLoop - The main loop for the hearbeat
func (e *Scheduler) HeartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(e.Config.EventLoopTime)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			e.checkRedis()
			e.Heartbeat()
		}
	}
}

// ReconcileLoop - The reconcile loop to check periodicly the task state
func (e *Scheduler) ReconcileLoop(ctx context.Context) {
	ticker := time.NewTicker(e.Config.ReconcileLoopTime)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		select {
		case <-ctx.Done():
			logrus.WithField("Reconcileloop", e.Framework.FrameworkName).Info("Reconcile Stop")
			return
		default:
			logrus.WithField("ReconcileLoop", e.Framework.FrameworkName).Info("Reconcile")
			e.reconcile()
			e.implicitReconcile()
		}
	}
}
