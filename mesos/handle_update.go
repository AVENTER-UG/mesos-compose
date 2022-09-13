package mesos

import (
	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
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

	if task.TaskID == "" {
		return nil
	}

	task.State = update.Status.State.String()

	logrus.WithField("func", "HandleUpdate").Debug("Task State: ", task.State)

	switch *update.Status.State {
	case mesosproto.TASK_FAILED, mesosproto.TASK_KILLED, mesosproto.TASK_ERROR, mesosproto.TASK_LOST, mesosproto.TASK_FINISHED:
		// check how to handle the event
		switch task.Restart {
		// never restart the task
		case "no":
			e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			return mesosutil.Call(msg)
		// only restart the task if it stopped by a failure
		case "on-failure":
			if update.Status.State.String() == mesosproto.TASK_FAILED.String() {
				task.State = ""
			}
			e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
			return mesosutil.Call(msg)
		// only restart the tasks if it does not stopped
		case "unless-stopped":
			if update.Status.State.String() == mesosproto.TASK_FINISHED.String() ||
				update.Status.State.String() == mesosproto.TASK_KILLED.String() {
				e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
				return mesosutil.Call(msg)
			}
			task.State = ""
		// always restart the container
		case "always":
			task.State = ""
		}
	case mesosproto.TASK_RUNNING:
		task.MesosAgent = mesosutil.GetAgentInfo(update.Status.GetAgentID().Value)
		task.NetworkInfo = mesosutil.GetNetworkInfo(task.TaskID)
		task.Agent = update.Status.GetAgentID().Value

		mesosutil.SuppressFramework()
	}

	// save the new state
	e.Redis.SaveTaskRedis(task)

	return mesosutil.Call(msg)
}
