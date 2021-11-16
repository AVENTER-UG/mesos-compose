package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	"github.com/AVENTER-UG/util"
	"github.com/sirupsen/logrus"
)

// Map the compose parameters into a mesos task
func mapComposeServiceToMesosTask(service cfg.Service, data cfg.Compose, vars map[string]string, name string, taskID string) {
	var cmd mesosutil.Command

	// if taskID is 0, then its a new task and we have to create a new ID
	newTaskID := taskID
	if taskID == "" {
		newTaskID, _ = util.GenUUID()
	}

	cmd.TaskName = config.PrefixTaskName + "_" + vars["project"] + "_" + name
	cmd.CPU = getCPU(service)
	cmd.Memory = getMemory(service)
	cmd.ContainerType = getLabelValueByKey("biz.aventer.mesos_compose.container_type", service)
	cmd.ContainerImage = service.Image
	cmd.NetworkMode = service.NetworkMode
	if len(data.Networks) > 0 {
		cmd.NetworkInfo = []mesosproto.NetworkInfo{{
			Name: func() *string { x := data.Networks[service.Network[0]].Name; return &x }(),
		}}
	}

	cmd.TaskID = newTaskID
	cmd.Privileged = service.Privileged
	cmd.Hostname = getHostname(service)
	cmd.Command = getCommand(service)
	cmd.Labels = getLabels(service)
	cmd.DockerPortMappings = getDockerPorts(service)
	cmd.Environment.Variables = getEnvironment(service)
	cmd.Volumes = getVolumes(service, data)

	cmd.Discovery = mesosproto.DiscoveryInfo{
		Visibility: 2,
		Name:       &cmd.TaskName,
		Ports: &mesosproto.Ports{
			Ports: getDiscoveryInfoPorts(service, cmd),
		},
	}

	if cmd.Command != "" {
		cmd.Shell = true
	} else {
		cmd.Shell = false
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

// Get random hostportnumber
func getRandomHostPort(service cfg.Service) int {
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(framework.PortRangeTo-framework.PortRangeFrom) + framework.PortRangeFrom
	return v
}

// Get the labels of the compose file
func getLabels(service cfg.Service) []mesosproto.Label {
	var label []mesosproto.Label

	for k, v := range service.Labels {
		var tmp mesosproto.Label
		tmp.Key = k
		tmp.Value = func() *string { x := fmt.Sprint(v); return &x }()
		label = append(label, tmp)
	}
	return label
}

// Return the value of the given key
func getLabelValueByKey(label string, service cfg.Service) string {
	for k, v := range service.Labels {
		if label == k {
			return fmt.Sprint(v)
		}
	}
	return ""
}

// Get the ports of the compose file
func getDockerPorts(service cfg.Service) []mesosproto.ContainerInfo_DockerInfo_PortMapping {
	var ports []mesosproto.ContainerInfo_DockerInfo_PortMapping
	hostport := uint32(getRandomHostPort(service))
	for i, c := range service.Ports {
		var tmp mesosproto.ContainerInfo_DockerInfo_PortMapping
		p := strings.Split(c, ":")
		ps, _ := strconv.Atoi(p[1])
		tmp.ContainerPort = uint32(ps)
		tmp.HostPort = uint32(hostport + uint32(i))
		tmp.Protocol = func() *string { x := "tcp"; return &x }()

		// check if this is a udp protocol
		proto := strings.Split(p[1], "/")
		if len(proto) > 1 {
			if strings.ToLower(proto[1]) == "udp" {
				tmp.Protocol = func() *string { x := "udp"; return &x }()
			}
		}

		ports = append(ports, tmp)
	}
	return ports
}

// Get the discoveryinfo ports of the compose file
func getDiscoveryInfoPorts(service cfg.Service, cmd mesosutil.Command) []mesosproto.Port {
	var port []mesosproto.Port
	for _, c := range service.Ports {
		var tmp mesosproto.Port
		p := strings.Split(c, ":")
		if len(p) != 2 {
			continue
		}
		// thats the containerport
		ps, _ := strconv.Atoi(p[1])
		// create the name of the port
		name := cmd.TaskName + "_" + p[1]

		// get the random hostport
		tmp.Number, tmp.Protocol = getHostPortByContainerPort(ps, cmd)
		tmp.Name = func() *string { x := name; return &x }()

		port = append(port, tmp)
	}

	return port
}

// get the random hostport and protcol of the container port
func getHostPortByContainerPort(port int, cmd mesosutil.Command) (uint32, *string) {
	for _, v := range cmd.DockerPortMappings {
		ps := v.ContainerPort
		if uint32(port) == ps {
			return v.HostPort, v.Protocol
		}
	}
	return 0, nil
}

// Get the environment of the compose file
func getEnvironment(service cfg.Service) []mesosproto.Environment_Variable {
	var env []mesosproto.Environment_Variable
	for _, c := range service.Environment {
		var tmp mesosproto.Environment_Variable
		p := strings.Split(c, "=")
		if len(p) != 2 {
			continue
		}
		tmp.Name = p[0]
		tmp.Value = func() *string { x := p[1]; return &x }()
		env = append(env, tmp)
	}
	return env
}

// Get the environment of the compose file
func getVolumes(service cfg.Service, data cfg.Compose) []mesosproto.Volume {
	var volume []mesosproto.Volume
	for _, c := range service.Volumes {
		var tmp mesosproto.Volume
		p := strings.Split(c, ":")
		if len(p) != 2 {
			continue
		}
		tmp.ContainerPath = p[1]
		tmp.Mode = mesosproto.RW.Enum()
		if len(p) == 3 {
			if strings.ToLower(p[2]) == "ro" {
				tmp.Mode = mesosproto.RO.Enum()
			}
		}
		if strings.ToLower(getLabelValueByKey("biz.aventer.mesos_compose.container_type", service)) == "docker" {
			driver := "local"
			if data.Volumes[p[0]].Driver != "" {
				driver = data.Volumes[p[0]].Driver
			}

			tmp.Source = &mesosproto.Volume_Source{
				Type: mesosproto.Volume_Source_DOCKER_VOLUME,
				DockerVolume: &mesosproto.Volume_Source_DockerVolume{
					Name:   p[0],
					Driver: func() *string { x := driver; return &x }(),
				},
			}
		}
		volume = append(volume, tmp)
	}
	return volume
}
