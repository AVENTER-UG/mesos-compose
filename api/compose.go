package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	"github.com/AVENTER-UG/util"
	"github.com/sirupsen/logrus"
)

// Map the compose parameters into a mesos task
func (e *API) mapComposeServiceToMesosTask(vars map[string]string, name string, task mesosutil.Command) {
	var cmd mesosutil.Command

	e.Service = e.Compose.Services[name]

	// if task is set then its not a new task and we have to save old needed parameter
	uuid, _ := util.GenUUID()
	newTaskID := vars["project"] + "_" + name + "." + uuid
	if task.TaskID != "" {
		newTaskID = task.TaskID
		cmd.State = task.State
		cmd.Agent = task.Agent
	}

	cmd.TaskName = e.Config.PrefixTaskName + ":" + vars["project"] + ":" + name
	cmd.CPU = e.getCPU()
	cmd.Memory = e.getMemory()
	cmd.Disk = e.getDisk()
	cmd.ContainerType = e.getContainerType()
	cmd.ContainerImage = e.Service.Image
	cmd.NetworkInfo = e.getNetworkInfo()
	cmd.NetworkMode = e.getNetworkMode()
	cmd.TaskID = newTaskID
	cmd.Privileged = e.Service.Privileged
	cmd.Hostname = e.getHostname()
	cmd.Command = e.getCommand()
	cmd.Labels = e.getLabels()
	cmd.Executor = e.getExecutor()
	cmd.DockerPortMappings = e.getDockerPorts(cmd.Agent)
	cmd.Environment.Variables = e.getEnvironment()
	cmd.Volumes = e.getVolumes(cmd.ContainerType)
	cmd.Instances = e.getReplicas()
	cmd.Discovery = e.getDiscoveryInfo(cmd)
	cmd.Shell = e.getShell(cmd)
	cmd.LinuxInfo = e.getLinuxInfo()
	cmd.DockerParameter = e.getDockerParameter(cmd)

	// store/update the mesos task in db
	e.SaveTaskRedis(cmd)
}

// Get the CPU value from the compose file, or the default one if it's unset
func (e *API) getCPU() float64 {
	if e.Service.Deploy.Resources.Limits.CPUs != "" {
		cpu, _ := strconv.ParseFloat(e.Service.Deploy.Resources.Limits.CPUs, 64)
		return cpu
	}
	return e.Config.CPU
}

// Get the Memory value from the compose file, or the default one if it's unset
func (e *API) getMemory() float64 {
	if e.Service.Deploy.Resources.Limits.Memory != "" {
		mem, _ := strconv.ParseFloat(e.Service.Deploy.Resources.Limits.Memory, 64)
		return mem
	}
	return e.Config.Memory
}

// Get the Disk value from the compose file, or the default one if it's unset
func (e *API) getDisk() float64 {
	// Currently, only default value is supported
	return e.Config.Disk
}

// Get the count of Replicas of the tasks
func (e *API) getReplicas() int {
	if e.Service.Deploy.Replicas != "" {
		replicas, _ := strconv.Atoi(e.Service.Deploy.Replicas)
		return replicas
	}
	return 1
}

// Get the Hostname value from the compose file, or generate one if it's unset
func (e *API) getHostname() string {
	if strings.ToLower(e.Service.NetworkMode) == "host" {
		return ""
	}

	if e.Service.Hostname != "" {
		return e.Service.Hostname
	} else if e.Service.ContainerName != "" {
		return e.Service.ContainerName
	}

	uuid, err := util.GenUUID()

	if err != nil {
		logrus.Error("getHostname genUUID Error: ", err)
	}
	return uuid
}

// Get the Command value from the compose file, or generate one if it's unset
func (e *API) getCommand() string {
	if len(e.Service.Command) != 0 {
		return e.Service.Command
	}
	return ""
}

// Get random hostportnumber
func (e *API) getRandomHostPort(agent string) int {
	rand.Seed(time.Now().UnixNano())
	// #nosec G404
	v := rand.Intn(e.Framework.PortRangeTo-e.Framework.PortRangeFrom) + e.Framework.PortRangeFrom
	if v > e.Framework.PortRangeTo {
		v = e.getRandomHostPort(agent)
	} else if v < e.Framework.PortRangeFrom {
		v = e.getRandomHostPort(agent)
	}
	port := uint32(v)
	if e.portInUse(port, agent) {
		v = e.getRandomHostPort(agent)
	}
	return v
}

// Check if the port is already in use
func (e *API) portInUse(port uint32, agent string) bool {
	// get all running services
	logrus.Debug("Check if port is in use: ", port, agent)
	keys := e.GetAllRedisKeys(e.Framework.FrameworkName + ":*")
	for keys.Next(e.Redis.RedisCTX) {
		// get the details of the current running service
		key := e.GetRedisKey(keys.Val())
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
func (e *API) getLabels() []mesosproto.Label {
	var label []mesosproto.Label

	for k, v := range e.Service.Labels {
		var tmp mesosproto.Label
		tmp.Key = k
		tmp.Value = func() *string { x := fmt.Sprint(v); return &x }()
		label = append(label, tmp)
	}
	return label
}

// Return the value of the given key
func (e *API) getLabelValueByKey(label string) string {
	for k, v := range e.Service.Labels {
		if label == k {
			return fmt.Sprint(v)
		}
	}
	return ""
}

// Get the ports of the compose file
func (e *API) getDockerPorts(agent string) []mesosproto.ContainerInfo_DockerInfo_PortMapping {
	var ports []mesosproto.ContainerInfo_DockerInfo_PortMapping
	hostport := uint32(e.getRandomHostPort(agent))
	for i, c := range e.Service.Ports {
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
func (e *API) getDiscoveryInfoPorts(cmd mesosutil.Command) []mesosproto.Port {
	var disport []mesosproto.Port
	for _, c := range e.Service.Ports {
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
		tmpport.Number, tmpport.Protocol = e.getHostPortByContainerPort(port, cmd)
		tmpport.Name = func() *string { x := name; return &x }()

		disport = append(disport, tmpport)
	}

	return disport
}

func (e *API) getDiscoveryInfo(cmd mesosutil.Command) mesosproto.DiscoveryInfo {
	return mesosproto.DiscoveryInfo{
		Visibility: 2,
		Name:       &cmd.TaskName,
		Ports: &mesosproto.Ports{
			Ports: e.getDiscoveryInfoPorts(cmd),
		},
	}
}

// get the random hostport and protcol of the container port
func (e *API) getHostPortByContainerPort(port int, cmd mesosutil.Command) (uint32, *string) {
	for _, v := range cmd.DockerPortMappings {
		ps := v.ContainerPort
		if uint32(port) == ps {
			return v.HostPort, v.Protocol
		}
	}
	return 0, nil
}

// Get the environment of the compose file
func (e *API) getEnvironment() []mesosproto.Environment_Variable {
	var env []mesosproto.Environment_Variable
	for _, c := range e.Service.Environment {
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
func (e *API) getVolumes(containerType string) []mesosproto.Volume {
	var volume []mesosproto.Volume
	for _, c := range e.Service.Volumes {
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

		driver := "local"
		if e.Compose.Volumes[p[0]].Driver != "" {
			driver = e.Compose.Volumes[p[0]].Driver
		}

		switch containerType {
		case "docker":
			tmp.Source = &mesosproto.Volume_Source{
				Type: mesosproto.Volume_Source_DOCKER_VOLUME,
				DockerVolume: &mesosproto.Volume_Source_DockerVolume{
					Name:   p[0],
					Driver: func() *string { x := driver; return &x }(),
				},
			}
		default:
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
func (e *API) getExecutor() mesosproto.ExecutorInfo {
	var executorInfo mesosproto.ExecutorInfo
	command := e.getLabelValueByKey("biz.aventer.mesos_compose.executor")
	uri := e.getLabelValueByKey("biz.aventer.mesos_compose.executor_uri")

	if command != "" {
		command = "exec '" + command + "' " + e.getCommand()
		executorID, _ := util.GenUUID()
		environment := e.getEnvironment()
		executorInfo = mesosproto.ExecutorInfo{
			Name: func() *string { x := path.Base(command); return &x }(),
			Type: *mesosproto.ExecutorInfo_CUSTOM.Enum(),
			ExecutorID: &mesosproto.ExecutorID{
				Value: executorID,
			},
			FrameworkID: e.Framework.FrameworkInfo.ID,
			Command: &mesosproto.CommandInfo{
				Value: func() *string { x := command; return &x }(),
				Environment: &mesosproto.Environment{
					Variables: environment,
				},
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

// get the Network Mode
func (e *API) getNetworkMode() string {
	if len(e.Service.Network) > 0 || len(e.Service.Networks) > 0 {
		// If Network was set, change the network mode to user
		return "user"
	}

	return ""
}

// get the NetworkInfo Name
func (e *API) getNetworkInfo() []mesosproto.NetworkInfo {
	if len(e.Compose.Networks) > 0 {
		var network string

		if len(e.Service.Network) > 0 {
			network = e.Service.Network[0]
		} else if len(e.Service.Networks) > 0 {
			network = e.Service.Networks[0]
		}

		// If Network Info was set, change the network mode to user
		e.Service.NetworkMode = "user"

		return []mesosproto.NetworkInfo{{
			Name: func() *string { x := e.Compose.Networks[network].Name; return &x }(),
		}}
	}

	return []mesosproto.NetworkInfo{}
}

// check if the task command inside of the container have to be executed as shell
func (e *API) getShell(cmd mesosutil.Command) bool {
	return cmd.Command != ""
}

// get linux info like capabilities
func (e *API) getLinuxInfo() mesosproto.LinuxInfo {
	linuxInfo := mesosproto.LinuxInfo{}

	if len(e.Service.CapAdd) > 0 {
		caps, err := json.Marshal(e.Service.CapAdd)
		if err != nil {
			logrus.WithField("func", "getLinuxInfo").Error("Could not marshal cap_add: ", err.Error())
		}
		tmp := "{ \"Capabilities\":" + string(caps) + "}"

		var capability mesosproto.CapabilityInfo

		json.Unmarshal([]byte(tmp), &capability)
		linuxInfo.EffectiveCapabilities = &capability
	}
	return linuxInfo
}

// get the container type
func (e *API) getContainerType() string {
	conType := strings.ToLower(e.getLabelValueByKey("biz.aventer.mesos_compose.container_type"))

	// if contype and custom executor is unset, then set the contype to DOCKER
	if conType == "" && e.getLabelValueByKey("biz.aventer.mesos_compose.executor") == "" {
		conType = "docker"
	}

	return conType
}

func (e *API) getDockerParameter(cmd mesosutil.Command) []mesosproto.Parameter {
	param := cmd.DockerParameter
	if len(param) == 0 {
		param = make([]mesosproto.Parameter, 0)
	}
	if e.Service.NetworkMode != "bridge" && e.getContainerType() == "docker" {
		return e.addDockerParameter(param, mesosproto.Parameter{Key: "net-alias", Value: e.getHostname()})
	}

	return param
}

// Append parameter to the list
func (e *API) addDockerParameter(current []mesosproto.Parameter, newValues mesosproto.Parameter) []mesosproto.Parameter {
	return append(current, newValues)
}
