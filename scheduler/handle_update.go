package scheduler

import (
	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	"github.com/sirupsen/logrus"
)

// HandleUpdate will handle the offers event of mesos
func (e *Scheduler) HandleUpdate(event *mesosproto.Event) {
	update := event.Update

	msg := &mesosproto.Call{
		Type: mesosproto.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &mesosproto.Call_Acknowledge{
			AgentId: update.Status.GetAgentId(),
			TaskId:  update.Status.GetTaskId(),
			Uuid:    update.Status.GetUuid(),
		},
	}

	// get the task of the current event, change the state
	task := e.Redis.GetTaskFromEvent(update)

	// if these object have not TaskID it's currently unknown by these framework.
	if task.TaskID == "" {
		logrus.WithField("func", "scheduler.HandleUpdate").Debug("Could not found Task in Redis: ", update.Status.GetTaskId())
		e.Mesos.Call(msg)
		return
	}

	task.State = update.Status.State.String()

	logrus.WithField("func", "scheduler.HandleUpdate").Tracef("Event %s, State %s, TaskID %s", event.GetType().String(), task.State, task.TaskID)

	switch *update.Status.State {
	case mesosproto.TaskState_TASK_FAILED, mesosproto.TaskState_TASK_KILLED, mesosproto.TaskState_TASK_ERROR, mesosproto.TaskState_TASK_FINISHED:
		if task.TaskID == "" {
			e.Mesos.Call(msg)
			return
		}
		logrus.WithField("func", "scheduler.HandleUpdate").Warn("Task State: " + task.State + " " + task.TaskID + " (" + task.TaskName + ")")

		// check how to handle the event
		switch task.Restart {
		// never restart the task
		case "no":
			e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			e.Mesos.ForceSuppressFramework()
			e.Mesos.Call(msg)
			return
		// only restart the task if it stopped by a failure
		case "on-failure":
			if update.Status.State.String() == mesosproto.TaskState_TASK_FAILED.String() {
				break
			}
			e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			e.Mesos.ForceSuppressFramework()
			e.Mesos.Call(msg)
			return
		// only restart the tasks if it does not stopped
		case "unless-stopped":
			if update.Status.State.String() == mesosproto.TaskState_TASK_FINISHED.String() {
				e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
				e.Mesos.ForceSuppressFramework()
				e.Mesos.Call(msg)
				return
			}
		}
		// all other cases, increase task count and restart task
		e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
		task.TaskID = e.API.IncreaseTaskCount(task.TaskID)
		task.State = ""
		task.Killed = false
		break
	case mesosproto.TaskState_TASK_LOST:
		if task.TaskID == "" {
			e.Mesos.Call(msg)
			return
		}
		logrus.WithField("func", "scheduler.HandleUpdate").Warn("Task State: " + task.State + " " + task.TaskID + " (" + task.TaskName + ")")
		e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
		e.Mesos.ForceSuppressFramework()
		e.Mesos.Call(msg)
		return
	case mesosproto.TaskState_TASK_RUNNING:
		if !e.Mesos.IsSuppress {
			logrus.WithField("func", "scheduler.HandleUpdate").Info("Task State: " + task.State + " " + task.TaskID + " (" + task.TaskName + ")")
		}

		task.MesosAgent = e.Mesos.GetAgentInfo(update.Status.GetAgentId().GetValue())
		task.NetworkInfo = e.Mesos.GetNetworkInfo(task.TaskID)
		task.Agent = update.Status.GetAgentId().GetValue()

		e.Mesos.SuppressFramework()
	}

	// save the new state
	e.Redis.SaveTaskRedis(task)

	e.Mesos.Call(msg)
}
