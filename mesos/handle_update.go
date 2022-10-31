package mesos

import (
	"crypto/tls"
	"encoding/json"
	"net/http"

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
	case mesosproto.TASK_FAILED, mesosproto.TASK_ERROR, mesosproto.TASK_LOST, mesosproto.TASK_FINISHED:
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
			if update.Status.State.String() == mesosproto.TASK_FINISHED.String() {
				e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
				return mesosutil.Call(msg)
			}
			task.State = ""
		// always restart the container
		case "always":
			task.State = ""
		default:
			task.State = ""
		}

	case mesosproto.TASK_KILLED:
		// remove task
		e.Redis.DelRedisKey(task.TaskName + ":" + task.TaskID)
		return mesosutil.Call(msg)

	case mesosproto.TASK_RUNNING:
		task.MesosAgent = mesosutil.GetAgentInfo(update.Status.GetAgentID().Value)
		task.NetworkInfo = e.GetNetworkInfo(task.TaskID)
		task.Agent = update.Status.GetAgentID().Value

		mesosutil.SuppressFramework()
	}

	// save the new state
	e.Redis.SaveTaskRedis(task)

	return mesosutil.Call(msg)
}

// GetNetworkInfo get network info of task
func (e *Scheduler) GetNetworkInfo(taskID string) []mesosproto.NetworkInfo {
	client := &http.Client{}
	// #nosec G402
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	protocol := "https"
	if !e.Framework.MesosSSL {
		protocol = "http"
	}
	req, _ := http.NewRequest("POST", protocol+"://"+e.Framework.MesosMasterServer+"/tasks/?task_id="+taskID+"&framework_id="+e.Framework.FrameworkInfo.ID.GetValue(), nil)
	req.Close = true
	req.SetBasicAuth(e.Framework.Username, e.Framework.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logrus.WithField("func", "getNetworkInfo").Error("Could not connect to agent: ", err.Error())
		return []mesosproto.NetworkInfo{}
	}

	defer res.Body.Close()

	var task mesosutil.MesosTasks
	err = json.NewDecoder(res.Body).Decode(&task)
	if err != nil {
		logrus.WithField("func", "getAgentInfo").Error("Could not encode json result: ", err.Error())
		return []mesosproto.NetworkInfo{}
	}

	if len(task.Tasks) > 0 {
		for _, status := range task.Tasks[0].Statuses {
			if status.State == "TASK_RUNNING" {
				var netw []mesosproto.NetworkInfo
				netw = append(netw, status.ContainerStatus.NetworkInfos[0])
				return netw
			}
		}
	}
	return []mesosproto.NetworkInfo{}
}
