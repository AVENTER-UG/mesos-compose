package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/AVENTER-UG/util/util"
	"github.com/sirupsen/logrus"
)

// Map the compose parameters into a mesos task
func (e *API) mapComposeServiceToMesosTask(vars map[string]string, name string, task cfg.Command) {
	var cmd cfg.Command

	e.Service = e.Compose.Services[name]

	// if task is set then its not a new task and we have to save old needed parameter
	uuid, _ := util.GenUUID()
	newTaskID := vars["project"] + "_" + name + "." + uuid

	if task.TaskID != "" {
		newTaskID = task.TaskID
		cmd.State = task.State
		cmd.Agent = task.Agent
	}

	cmd.TaskName = e.getTaskName(vars["project"], name)
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
	cmd.DockerPortMappings = e.getDockerPorts()
	cmd.Environment.Variables = e.getEnvironment()
	cmd.Volumes = e.getVolumes(cmd.ContainerType)
	cmd.Instances = e.getReplicas()
	cmd.Discovery = e.getDiscoveryInfo(cmd)
	cmd.Shell = e.getShell(cmd)
	cmd.LinuxInfo = e.getLinuxInfo()
	cmd.DockerParameter = e.getDockerParameter(cmd)
	cmd.PullPolicy = e.getPullPolicy()
	cmd.Restart = e.getRestart()

	// set the docker constraints
	e.setConstraints(&cmd)

	// store/update the mesos task in db
	e.Redis.SaveTaskRedis(cmd)
}

// Get the name of the task
func (e *API) getTaskName(project, name string) string {
	taskName := e.getLabelValueByKey("biz.aventer.mesos_compose.taskname")

	if taskName != "" {
		// be sure the taskname is only running under the frameworks prefix
		if strings.Split(taskName, ":")[0] != e.Config.PrefixTaskName {
			return e.Config.PrefixTaskName + ":" + taskName
		}
		return taskName
	}
	return e.Config.PrefixTaskName + ":" + project + ":" + name
}

// Get the Restart value
func (e *API) getRestart() string {
	if e.Service.Restart != "" {
		return e.Service.Restart
	}
	return "unless-stopped"
}

// Get the CPU value from the compose file, or the default one if it's unset
func (e *API) getCPU() float64 {
	if e.Service.Deploy.Resources.Limits.CPUs > 0 {
		return e.Service.Deploy.Resources.Limits.CPUs
	}
	return e.Config.CPU
}

// Get the Memory value from the compose file, or the default one if it's unset
func (e *API) getMemory() float64 {
	if e.Service.Deploy.Resources.Limits.Memory > 0 {
		return e.Service.Deploy.Resources.Limits.Memory
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
	if e.getNetworkMode() == "host" {
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

// GetRandomHostPort Get random hostportnumber
func (e *API) GetRandomHostPort() uint32 {
	rand.Seed(time.Now().UnixNano())
	// #nosec G404
	v := uint32(rand.Intn(e.Framework.PortRangeTo-e.Framework.PortRangeFrom) + e.Framework.PortRangeFrom)

	if v > uint32(e.Framework.PortRangeTo) {
		v = e.GetRandomHostPort()
	} else if v < uint32(e.Framework.PortRangeFrom) {
		v = e.GetRandomHostPort()
	}
	return v
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

// geturn the pullpolicy value
func (e *API) getPullPolicy() string {
	return e.Service.PullPolicy
}

// GetDockerPorts Get the ports of the compose file
func (e *API) getDockerPorts() []mesosproto.ContainerInfo_DockerInfo_PortMapping {
	if e.getNetworkMode() == "host" {
		return nil
	}
	var ports []mesosproto.ContainerInfo_DockerInfo_PortMapping
	for _, c := range e.Service.Ports {
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
			port, _ = strconv.Atoi(proto[0])
			if strings.ToLower(proto[1]) == "udp" {
				tmp.Protocol = func() *string { x := "udp"; return &x }()
			}
		}

		tmp.ContainerPort = uint32(port)
		tmp.HostPort = 0 // the port will be generatet during the schedule
		ports = append(ports, tmp)
	}
	return ports
}

// Get the discoveryinfo ports of the compose file
func (e *API) getDiscoveryInfoPorts(cmd cfg.Command) []mesosproto.Port {
	var disport []mesosproto.Port
	for _, c := range cmd.DockerPortMappings {
		var tmpport mesosproto.Port
		p := func() *string { x := cmd.TaskName + ":" + strconv.FormatUint(uint64(c.ContainerPort), 10); return &x }()
		tmpport.Name = p
		tmpport.Number = c.HostPort
		tmpport.Protocol = c.Protocol

		disport = append(disport, tmpport)
	}

	return disport
}

func (e *API) getDiscoveryInfo(cmd cfg.Command) mesosproto.DiscoveryInfo {
	return mesosproto.DiscoveryInfo{
		Visibility: 2,
		Name:       &cmd.TaskName,
		Ports: &mesosproto.Ports{
			Ports: e.getDiscoveryInfoPorts(cmd),
		},
	}
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
		// check if the value is a secret
		if strings.Contains(p[1], "vault://") {
			p[1] = e.Vault.GetKey(p[1])
		}
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
			var fetch []mesosproto.CommandInfo_URI
			err := json.Unmarshal([]byte(uri), &fetch)

			if err != nil {
				logrus.WithField("func", "getExecutor").Error("Could not unmarchal biz.aventer.mesos_compose.executor_uri")
			}

			for _, uri := range fetch {
				executorInfo.Command.URIs = []mesosproto.CommandInfo_URI{
					{
						Value:      uri.GetValue(),
						Extract:    func() *bool { x := false; return &x }(),
						Executable: func() *bool { x := true; return &x }(),
						Cache:      func() *bool { x := false; return &x }(),
						OutputFile: uri.OutputFile,
					},
				}
			}
		}
	}
	return executorInfo
}

// get the Network Mode
func (e *API) getNetworkMode() string {
	// default network mode
	mode := "user"

	if e.Service.NetworkMode != "" {
		mode = e.Service.NetworkMode
	}

	if len(e.Compose.Networks) > 0 {
		network := e.getNetworkName(0)
		if e.Compose.Networks[network].Driver != "" {
			mode = e.Compose.Networks[network].Driver
		}
	}

	return strings.ToLower(mode)
}

// get the NetworkInfo Name
func (e *API) getNetworkInfo() []mesosproto.NetworkInfo {
	// get network name
	name := e.getNetworkName(0)

	if len(e.Compose.Networks) > 0 {
		name = e.Compose.Networks[name].Name
	}

	return []mesosproto.NetworkInfo{{
		Name: func() *string { x := name; return &x }(),
	}}
}

// get the name of the network parameter
func (e *API) getNetworkName(val int) string {
	// default network name
	network := "default"

	if e.Service.Network != "" {
		network = e.Service.Network
	} else if len(e.Service.Networks) > val {
		keys := reflect.ValueOf(e.Service.Networks).MapKeys()
		network = keys[val].String()
	}

	return network
}

// check if the task command inside of the container have to be executed as shell
func (e *API) getShell(cmd cfg.Command) bool {
	return cmd.Command != ""
}

// get linux info like capabilities
func (e *API) getLinuxInfo() mesosproto.LinuxInfo {
	linuxInfo := mesosproto.LinuxInfo{}

	if len(e.Service.CapAdd) > 0 {
		linuxInfo.EffectiveCapabilities = e.getCapabilities(e.Service.CapAdd)
	}
	if len(e.Service.CapDrop) > 0 {
		linuxInfo.EffectiveCapabilities = e.getCapabilities(e.Service.CapDrop)
	}
	return linuxInfo
}

// get capabilities
func (e *API) getCapabilities(capa []string) *mesosproto.CapabilityInfo {
	caps, err := json.Marshal(capa)
	if err != nil {
		logrus.WithField("func", "getCapabilities").Error("Could not marshal cap_add/drop: ", err.Error())
	}
	tmp := "{ \"Capabilities\":" + string(caps) + "}"

	var capability mesosproto.CapabilityInfo

	json.Unmarshal([]byte(tmp), &capability)
	return &capability
}

// get the container type
func (e *API) getContainerType() string {
	conType := strings.ToLower(e.getLabelValueByKey("biz.aventer.mesos_compose.container_type"))

	// if contype and custom executor is unset, then set the contype to DOCKER
	if conType == "" {
		conType = "docker"
	}

	return conType
}

func (e *API) getDockerParameter(cmd cfg.Command) []mesosproto.Parameter {
	param := cmd.DockerParameter
	if len(param) == 0 {
		param = make([]mesosproto.Parameter, 0)
	}

	if e.getContainerType() == "docker" {
		// set net-alias if its defined
		alias := e.getNetAlias()
		if alias != "" {
			param = e.addDockerParameter(param, mesosproto.Parameter{Key: "net-alias", Value: alias})
		}
		// add default volume driver if there is no defined volume
		if len(e.Service.Volumes) == 0 {
			param = e.addDockerParameter(param, mesosproto.Parameter{Key: "volume-driver", Value: e.Config.DefaultVolumeDriver})
		}
		// configure ulimits
		param = e.getUlimit(param)
	}

	return param
}

func (e *API) getUlimit(param []mesosproto.Parameter) []mesosproto.Parameter {
	if e.Service.Ulimits.Memlock.Hard != 0 {
		param = e.addDockerParameter(param, mesosproto.Parameter{Key: "ulimit", Value: "memlock=" + strconv.Itoa(e.Service.Ulimits.Memlock.Hard) + ":" + strconv.Itoa(e.Service.Ulimits.Memlock.Soft)})
	}

	if e.Service.Ulimits.Nofile.Hard != 0 {
		param = e.addDockerParameter(param, mesosproto.Parameter{Key: "ulimit", Value: "nofile=" + strconv.Itoa(e.Service.Ulimits.Nofile.Hard) + ":" + strconv.Itoa(e.Service.Ulimits.Nofile.Soft)})
	}

	return param
}

func (e *API) getNetAlias() string {
	if len(e.Service.Networks) > 0 {
		network := reflect.ValueOf(e.Service.Networks).MapKeys()[0].String()
		if len(e.Service.Networks[network].Aliases) > 0 {
			return e.Service.Networks[network].Aliases[0]
		}
	}

	if e.getNetworkMode() == "user" {
		return e.getHostname()
	}

	return ""
}

// Append parameter to the list
func (e *API) addDockerParameter(current []mesosproto.Parameter, newValues mesosproto.Parameter) []mesosproto.Parameter {
	return append(current, newValues)
}

// translate the docker-compose placement constraints to labels
func (e *API) setConstraints(cmd *cfg.Command) {
	if len(e.Service.Deploy.Placement.Constraints) > 0 {
		for _, constraint := range e.Service.Deploy.Placement.Constraints {
			cons := strings.Split(constraint, "==")
			if len(cons) >= 2 {
				if cons[0] == "node.hostname" {
					cmd.Labels = append(cmd.Labels, mesosproto.Label{Key: "__mc_placement_node_hostname", Value: &cons[1]})
				}
				if cons[0] == "node.platform.os" {
					cmd.Labels = append(cmd.Labels, mesosproto.Label{Key: "__mc_placement_node_platform_os", Value: &cons[1]})
				}
				if cons[0] == "node.platform.arch" {
					cmd.Labels = append(cmd.Labels, mesosproto.Label{Key: "__mc_placement_node_platform_arch", Value: &cons[1]})
				}
			}
		}
	}
}
