package api

import (
	"encoding/json"
	"strconv"
	"strings"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	"github.com/AVENTER-UG/util"
	"github.com/sirupsen/logrus"
)

// Map the compose parameters into a mesos task
func mapComposeServiceToMesosTask(service cfg.Service, network map[string]cfg.Networks, vars map[string]string, name string, taskID string) {
	var cmd mesosutil.Command

	// if taskID is 0, then its a new task and we have to create a new ID
	newTaskID := taskID
	if taskID == "" {
		newTaskID, _ = util.GenUUID()
	}

	cmd.TaskName = config.PrefixTaskName + "_" + vars["project"] + "_" + name
	cmd.CPU = getCPU(service)
	cmd.Memory = getMemory(service)
	cmd.ContainerType = service.Labels.ContainerType
	cmd.ContainerImage = service.Image
	cmd.NetworkMode = service.NetworkMode
	cmd.NetworkInfo = []mesosproto.NetworkInfo{{
		Name: func() *string { x := network[service.Network[0]].Name; return &x }(),
	}}

	cmd.TaskID = newTaskID
	cmd.Privileged = service.Privileged
	cmd.Hostname = getHostname(service)
	cmd.Command = getCommand(service)

	if cmd.Command != "" {
		cmd.Shell = true
	} else {
		cmd.Shell = false
	}

	var hostport uint32
	hostport = 31859
	protocol := "tcp"
	cmd.DockerPortMappings = []mesosproto.ContainerInfo_DockerInfo_PortMapping{
		{
			HostPort:      hostport,
			ContainerPort: 10422,
			Protocol:      &protocol,
		},
	}

	// store mesos task in db
	d, _ := json.Marshal(&cmd)
	logrus.Debug("Scheduled Mesos Task: ", util.PrettyJSON(d))
	err := config.RedisClient.Set(config.RedisCTX, cmd.TaskName+":"+newTaskID, d, 0).Err()

	if err != nil {
		logrus.Error("Cloud not store Mesos Task in Redis: ", err)
	}

}

// Get the CPU value from the compose file, or the default one if it's unset
func getCPU(service cfg.Service) float64 {
	if service.Deploy.Resources.Limits.CPUs != "" {
		cpu, _ := strconv.ParseFloat(service.Deploy.Resources.Limits.CPUs, 64)
		return cpu
	}
	return config.CPU
}

// Get the Memory value from the compose file, or the default one if it's unset
func getMemory(service cfg.Service) float64 {
	if service.Deploy.Resources.Limits.Memory != "" {
		mem, _ := strconv.ParseFloat(service.Deploy.Resources.Limits.Memory, 64)
		return mem
	}
	return config.Memory
}

// Get the Hostname value from the compose file, or generate one if it's unset
func getHostname(service cfg.Service) string {
	if service.Hostname != "" {
		return service.Hostname
	}
	uuid, err := util.GenUUID()

	if err != nil {
		logrus.Error("getHostname genUUID Error: ", err)
	}
	return uuid
}

// Get the Command value from the compose file, or generate one if it's unset
func getCommand(service cfg.Service) string {
	if len(service.Command) != 0 {
		comm := strings.Join(service.Command, " ")
		return comm
	}

	return ""
}
