package scheduler

import (
	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	"github.com/sirupsen/logrus"
)

// HandleUpdate will handle the offers event of mesos
func (e *Scheduler) HandleUpdate(event *mesosproto.Event) error {
	update := event.Update

	msg := &mesosproto.Call{
		Type: mesosproto.Call_ACKNOWLEDGE,
		Acknowledge: &mesosproto.Call_Acknowledge{
			AgentID: *update.Status.AgentID,
			TaskID:  update.Status.TaskID,
			UUID:    update.Status.UUID,
		},
	}

	// get the task of the current event, change the state
	task := e.Redis.GetTaskFromEvent(update)

	// if these object have not TaskID it's currently unknown by these framework.
	if task.TaskID == "" {
		logrus.WithField("func", "scheduler.HandleUpdate").Debug("Could not found Task in Redis: ", update.Status.TaskID.Value)

		if *update.Status.State != mesosproto.TASK_LOST {
			e.Mesos.Kill(update.Status.TaskID.Value, update.Status.AgentID.Value)
		}
	}

	task.State = update.Status.State.String()

	switch *update.Status.State {
	case mesosproto.TASK_FAILED, mesosproto.TASK_ERROR, mesosproto.TASK_FINISHED, mesosproto.TASK_KILLED:
		if task.TaskID == "" {
			return e.Mesos.Call(msg)
		}
		logrus.WithField("func", "scheduler.HandleUpdate").Warn("Task State: " + task.State + " " + task.TaskID + " (" + task.TaskName + ")")

		// check how to handle the event
		switch task.Restart {
		// never restart the task
		case "no":
			e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			e.Mesos.ForceSuppressFramework()
			return e.Mesos.Call(msg)
		// only restart the task if it stopped by a failure
		case "on-failure":
			if update.Status.State.String() == mesosproto.TASK_FAILED.String() {
				break
			}
			e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			e.Mesos.ForceSuppressFramework()
			return e.Mesos.Call(msg)
		// only restart the tasks if it does not stopped
		case "unless-stopped":
			if update.Status.State.String() == mesosproto.TASK_FINISHED.String() {
				e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
				e.Mesos.ForceSuppressFramework()
				return e.Mesos.Call(msg)
			}
		}
		// all other cases, increase task count and restart task
		e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
		task.TaskID = e.API.IncreaseTaskCount(task.TaskID)
		task.State = ""
	case mesosproto.TASK_LOST:
		if task.TaskID == "" {
			return e.Mesos.Call(msg)
		}
		logrus.WithField("func", "scheduler.HandleUpdate").Warn("Task State: " + task.State + " " + task.TaskID + " (" + task.TaskName + ")")
		// all other cases, increase task count and restart task
		e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
		task.TaskID = e.API.IncreaseTaskCount(task.TaskID)
		task.State = ""
	case mesosproto.TASK_RUNNING:
		if !e.Mesos.IsSuppress {
			logrus.WithField("func", "scheduler.HandleUpdate").Info("Task State: " + task.State + " " + task.TaskID + " (" + task.TaskName + ")")
		}

		task.MesosAgent = e.Mesos.GetAgentInfo(update.Status.GetAgentID().Value)
		task.NetworkInfo = e.Mesos.GetNetworkInfo(task.TaskID)
		task.Agent = update.Status.GetAgentID().Value

		e.Mesos.SuppressFramework()
	}

	// save the new state
	e.Redis.SaveTaskRedis(task)

	return e.Mesos.Call(msg)
}
