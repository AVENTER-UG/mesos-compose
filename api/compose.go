package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path"
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
func mapComposeServiceToMesosTask(service cfg.Service, data cfg.Compose, vars map[string]string, name string, task mesosutil.Command) {
	var cmd mesosutil.Command

	// if task is set then its not a new task and we have to save old needed parameter
	uuid, _ := util.GenUUID()
	newTaskID := vars["project"] + "_" + name + "." + uuid
	if task.TaskID != "" {
		newTaskID = task.TaskID
		cmd.State = task.State
		cmd.Agent = task.Agent
	}

	cmd.TaskName = config.PrefixTaskName + ":" + vars["project"] + ":" + name
	cmd.CPU = getCPU(service)
	cmd.Memory = getMemory(service)
	cmd.Disk = getDisk(service)
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
	cmd.Executor = getExecutor(service)
	cmd.DockerPortMappings = getDockerPorts(service, cmd.Agent)
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

// Get the Disk value from the compose file, or the default one if it's unset
func getDisk(service cfg.Service) float64 {
	// Currently, onyl default value is supported
	return config.Disk
}

// Get the Hostname value from the compose file, or generate one if it's unset
func getHostname(service cfg.Service) string {
	if service.Hostname != "" {
		return service.Hostname
	}

	if strings.ToLower(service.NetworkMode) == "host" {
		return ""
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
func getRandomHostPort(agent string) int {
	rand.Seed(time.Now().UnixNano())
	// #nosec G404
	v := rand.Intn(framework.PortRangeTo-framework.PortRangeFrom) + framework.PortRangeFrom
	port := uint32(v)
	if portInUse(port, agent) {
		v = getRandomHostPort(agent)
	}
	return v
}

// Check if the port is already in use
func portInUse(port uint32, agent string) bool {
	// get all running services
	logrus.Debug("Check if port is in use: ", port, agent)
	keys := GetAllRedisKeys(framework.FrameworkName + ":*")
	for keys.Next(config.RedisCTX) {
		// get the details of the current running service
		key := GetRedisKey(keys.Val())
		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)

		// check if its the given agent
		if task.Agent == agent {
			// check if the given port is already in use
			ports := task.Discovery.GetPorts()
			if ports == nil {
				for _, hostport := range ports.GetPorts() {
					if hostport.Number == port {
						return true
					}
				}
			}
		}
	}
	return false
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
func getDockerPorts(service cfg.Service, agent string) []mesosproto.ContainerInfo_DockerInfo_PortMapping {
	var ports []mesosproto.ContainerInfo_DockerInfo_PortMapping
	hostport := uint32(getRandomHostPort(agent))
	for i, c := range service.Ports {
		var tmp mesosproto.ContainerInfo_DockerInfo_PortMapping
		var port int
		// "<hostport>:<containerport>"
		// "<hostip>:<hostport>:<containerport>"
		count := strings.Count(c, ":")

		var p []string
		var proto []string
		if count > 0 {
			p = strings.Split(c, ":")
			port, _ = strconv.Atoi(p[count])
			proto = strings.Split(p[count], "/")
		} else {
			port, _ = strconv.Atoi(c)
			proto = strings.Split(c, "/")
		}

		// check if this is a udp protocol
		tmp.Protocol = func() *string { x := "tcp"; return &x }()
		if len(proto) > 1 {
			if strings.ToLower(proto[1]) == "udp" {
				tmp.Protocol = func() *string { x := "udp"; return &x }()
			}
		}

		tmp.ContainerPort = uint32(port)
		tmp.HostPort = hostport + uint32(i)
		ports = append(ports, tmp)
	}
	return ports
}

// Get the discoveryinfo ports of the compose file
func getDiscoveryInfoPorts(service cfg.Service, cmd mesosutil.Command) []mesosproto.Port {
	var disport []mesosproto.Port
	for _, c := range service.Ports {
		var tmpport mesosproto.Port
		var port int
		// "<hostport>:<containerport>"
		// "<hostip>:<hostport>:<containerport>"
		count := strings.Count(c, ":")
		var p []string
		if count > 0 {
			p = strings.Split(c, ":")
			port, _ = strconv.Atoi(p[count])
		} else {
			port, _ = strconv.Atoi(c)
		}

		// create the name of the port
		name := cmd.TaskName + ":" + p[count]

		// get the random hostport
		tmpport.Number, tmpport.Protocol = getHostPortByContainerPort(port, cmd)
		tmpport.Name = func() *string { x := name; return &x }()

		disport = append(disport, tmpport)
	}

	return disport
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
		if len(p) < 2 {
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

// Get custome executer
func getExecutor(service cfg.Service) mesosproto.ExecutorInfo {
	var executorInfo mesosproto.ExecutorInfo
	command := getLabelValueByKey("biz.aventer.mesos_compose.executor", service)
	uri := getLabelValueByKey("biz.aventer.mesos_compose.executor_uri", service)

	if command != "" {
		executorID, _ := util.GenUUID()
		executorInfo = mesosproto.ExecutorInfo{
			Name: func() *string { x := path.Base(command); return &x }(),
			Type: mesosproto.ExecutorInfo_CUSTOM,
			ExecutorID: &mesosproto.ExecutorID{
				Value: executorID,
			},
			FrameworkID: framework.FrameworkInfo.ID,
			Command: &mesosproto.CommandInfo{
				Value: func() *string { x := command; return &x }(),
			},
		}

		if uri != "" {
			executorInfo.Command.URIs = []mesosproto.CommandInfo_URI{
				{
					Value:      uri,
					Extract:    func() *bool { x := false; return &x }(),
					Executable: func() *bool { x := true; return &x }(),
					Cache:      func() *bool { x := false; return &x }(),
					OutputFile: func() *string { x := path.Base(command); return &x }(),
				},
			}
		}

	}
	return executorInfo
}
