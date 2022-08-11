package mesos

import (
	"time"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"

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
		// continue if the key is not a mesos task
		if e.API.CheckIfNotTask(keys) {
			continue
		}
		// get the values of the current key
		key := e.API.GetRedisKey(keys.Val())

		task := mesosutil.DecodeTask(key)

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

			// Change the Dynamic Host Ports
			task.DockerPortMappings = e.changeDockerPorts(task)
			task.Discovery = e.changeDiscoveryInfo(task)

			// add task to communication channel
			e.Framework.CommandChan <- task

			e.API.SaveTaskRedis(task)

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

func (e *Scheduler) changeDockerPorts(cmd mesosutil.Command) []mesosproto.ContainerInfo_DockerInfo_PortMapping {
	var ret []mesosproto.ContainerInfo_DockerInfo_PortMapping
	for _, port := range cmd.DockerPortMappings {
		port.HostPort = e.API.GetRandomHostPort()
		ret = append(ret, port)
	}
	return ret
}

func (e *Scheduler) changeDiscoveryInfo(cmd mesosutil.Command) mesosproto.DiscoveryInfo {
	for i, port := range cmd.DockerPortMappings {
		cmd.Discovery.Ports.Ports[i].Number = port.HostPort
	}
	return cmd.Discovery
}
