package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"path"
	"reflect"
	"strconv"
	"strings"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/AVENTER-UG/util/util"
	"github.com/sirupsen/logrus"
)

// Map the compose parameters into a mesos task
func (e *API) mapComposeServiceToMesosTask(vars map[string]string, name string, task *cfg.Command) {
	cmd := new(cfg.Command)

	e.Service = e.Compose.Services[name]

	// if task is set then its not a new task and we have to save old needed parameter
	uuid, err := util.GenUUID()
	if err != nil {
		logrus.WithField("func", "api.mapComposeServiceToMesosTask").Error("Error during create uuid: ", err.Error())
		return
	}
	newTaskID := e.IncreaseTaskCount(vars["project"] + "_" + name + "." + uuid)

	if task.TaskID != "" {
		newTaskID = task.TaskID
		cmd.State = task.State
		cmd.Agent = task.Agent
	}

	cmd.TaskID = newTaskID
	cmd.TaskName = e.getTaskName(vars["project"], name)
	cmd.CPU = e.getCPU()
	cmd.Memory = e.getMemory()
	cmd.Disk = e.getDisk()
	cmd.ContainerType = e.getContainerType()
	cmd.ContainerImage = e.Service.Image
	cmd.NetworkInfo = e.getNetworkInfo()
	cmd.NetworkMode = e.getNetworkMode()
	cmd.Privileged = e.Service.Privileged
	cmd.Hostname = e.getHostname()
	cmd.Command = e.getCommand()
	cmd.Labels = e.getLabels()
	cmd.Executor = e.getExecutor()
	cmd.DockerPortMappings = e.getDockerPorts()
	cmd.Environment = new(mesosproto.Environment)
	cmd.Environment.Variables = e.getEnvironment()
	cmd.Volumes = e.getVolumes()
	cmd.Instances = e.getReplicas()
	cmd.Discovery = e.getDiscoveryInfo(cmd)
	cmd.Shell = e.getShell()
	cmd.LinuxInfo = e.getLinuxInfo()
	cmd.DockerParameter = e.getDockerParameter(cmd)
	cmd.DockerParameter = e.getGPUs(cmd)
	cmd.PullPolicy = e.getPullPolicy()
	cmd.Restart = e.getRestart()
	cmd.Mesos = e.Service.Mesos
	cmd.Uris = e.getURIs()
	cmd.Arguments = e.getArguments()


	// set healthcheck if it's configured
	e.setHealthCheck(cmd)

	// set the docker constraints
	e.setConstraints(cmd)

	// store/update the mesos task in db
	e.Redis.SaveTaskRedis(cmd)
}

// Set GPU config for the docker container. Mesos container not supportet right now.
func (e *API) getGPUs(cmd *cfg.Command) []*mesosproto.Parameter {
	param := cmd.DockerParameter

	if len(param) == 0 {
		param = make([]*mesosproto.Parameter, 0)
	}

	if e.getContainerType() == "docker" && len(e.Service.GPUs.Driver) > 0 {
    if e.Service.GPUs.Driver == "amd" {
      param = e.addDockerParameter(param, "device", "/dev/kfd")
      param = e.addDockerParameter(param, "device", "/dev/dri")
      param = e.addDockerParameter(param, "security-opt", "seccomp=unconfined")
    }

    if e.Service.GPUs.Driver == "nvidia" && e.Service.GPUs.Device >= 0 {
      i := strconv.Itoa(e.Service.GPUs.Device)
      param = e.addDockerParameter(param, "gpus", "device="+i)
    }
	}

	return param
}

// Set the healtchek configuration
func (e *API) setHealthCheck(cmd *cfg.Command) {
	if e.Service.HealthCheck != nil {
		if e.Service.HealthCheck.Command.GetValue() != "" {
			e.Service.HealthCheck.Type = func() *mesosproto.HealthCheck_Type { x := mesosproto.HealthCheck_COMMAND; return &x }()
		} else if e.Service.HealthCheck.Http.GetPort() != 0 {
			e.Service.HealthCheck.Type = func() *mesosproto.HealthCheck_Type { x := mesosproto.HealthCheck_HTTP; return &x }()
		} else if e.Service.HealthCheck.Tcp.GetPort() != 0 {
			e.Service.HealthCheck.Type = func() *mesosproto.HealthCheck_Type { x := mesosproto.HealthCheck_TCP; return &x }()
		} else {
			// if no command, http port or tcp port is set, then disable the healthcheck
			cmd.EnableHealthCheck = false
			return
		}
		cmd.EnableHealthCheck = true
		cmd.Health = e.Service.HealthCheck
		return
	}
	cmd.EnableHealthCheck = false
}

// Get the Command value from the compose file, or generate one if it's unset
func (e *API) getArguments() []string {
	if len(e.Service.Arguments) != 0 {
		return e.Service.Arguments
	}
	return []string{}
}

// Get the URIS to fetch
func (e *API) getURIs() []*mesosproto.CommandInfo_URI {
	if len(e.Service.Mesos.Fetch) > 0 {
		return e.Service.Mesos.Fetch
	}

	var res []*mesosproto.CommandInfo_URI
	return res
}

// Get the name of the task
func (e *API) getTaskName(project, name string) string {
	taskName := e.Service.Mesos.TaskName
	// be sure the taskname is only running under the frameworks prefix
	if strings.Split(taskName, ":")[0] != e.Config.PrefixTaskName {
		return e.Config.PrefixTaskName + ":" + project + ":" + name
	}
	return taskName
}

// Get shell
func (e *API) getShell() bool {
	if e.Service.Shell && e.Service.Command != "" {
		return true
	}
	return false
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
func (e *API) GetRandomHostPort() *uint32 {
	// #nosec G404
	v := util.Uint32ToPointer(uint32(rand.Intn(e.Framework.PortRangeTo-e.Framework.PortRangeFrom) + e.Framework.PortRangeFrom))

	if *v > uint32(e.Framework.PortRangeTo) {
		v = e.GetRandomHostPort()
	} else if *v < uint32(e.Framework.PortRangeFrom) {
		v = e.GetRandomHostPort()
	}
	return v
}

// Get the labels of the compose file
func (e *API) getLabels() []*mesosproto.Label {
	var label []*mesosproto.Label

	for k, v := range e.Service.Labels {
		var tmp mesosproto.Label
		tmp.Key = util.StringToPointer(k)
		tmp.Value = func() *string { x := fmt.Sprint(v); return &x }()
		label = append(label, &tmp)
	}
	return label
}

// geturn the pullpolicy value
func (e *API) getPullPolicy() string {
	return e.Service.PullPolicy
}

// GetDockerPorts Get the ports of the compose file
func (e *API) getDockerPorts() []*mesosproto.ContainerInfo_DockerInfo_PortMapping {
	if e.getNetworkMode() == "host" {
		return nil
	}
	ports := make([]*mesosproto.ContainerInfo_DockerInfo_PortMapping, 0)
	for _, c := range e.Service.Ports {
		tmp := new(mesosproto.ContainerInfo_DockerInfo_PortMapping)
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

		// set protocol. default is tcp.
		tmp.Protocol = util.StringToPointer("tcp")
		if len(proto) > 1 {
			port, _ = strconv.Atoi(proto[0])
			if strings.ToLower(proto[1]) == "udp" {
				tmp.Protocol = util.StringToPointer("udp")
			}
			if strings.ToLower(proto[1]) == "wss" {
				tmp.Protocol = util.StringToPointer("wss")
			}
			if strings.ToLower(proto[1]) == "http" {
				tmp.Protocol = util.StringToPointer("http")
			}
			if strings.ToLower(proto[1]) == "https" {
				tmp.Protocol = util.StringToPointer("https")
			}
			if strings.ToLower(proto[1]) == "h2c" {
				tmp.Protocol = util.StringToPointer("h2c")
			}
		}

		tmp.ContainerPort = util.Uint32ToPointer(uint32(port))
		tmp.HostPort = util.Uint32ToPointer(0) // the port will be generatet during the schedule
		ports = append(ports, tmp)
	}
	return ports
}

// Get the discoveryinfo ports of the compose file
func (e *API) getDiscoveryInfoPorts(cmd *cfg.Command) []*mesosproto.Port {
	disport := make([]*mesosproto.Port, 0)
	for i, c := range cmd.DockerPortMappings {
		tmpport := new(mesosproto.Port)
		p := util.StringToPointer(cmd.TaskName + ":" + strconv.FormatUint(uint64(*c.ContainerPort), 10))
		tmpport.Name = func() *string { x := strings.ReplaceAll(*p, ":", e.Config.DiscoveryPortNameDelimiter); return &x }()
		tmpport.Number = c.HostPort
		tmpport.Protocol = c.Protocol

		// Docker understand only tcp and udp.
		if *c.Protocol != "udp" && *c.Protocol != "tcp" {
			cmd.DockerPortMappings[i].Protocol = func() *string { x := "tcp"; return &x }()
		}

		disport = append(disport, tmpport)
	}

	return disport
}

func (e *API) getDiscoveryInfo(cmd *cfg.Command) *mesosproto.DiscoveryInfo {
	name := strings.ReplaceAll(cmd.TaskName, ":", e.Config.DiscoveryInfoNameDelimiter)
	return &mesosproto.DiscoveryInfo{
		Visibility: mesosproto.DiscoveryInfo_EXTERNAL.Enum(),
		Name:       &name,
		Ports: &mesosproto.Ports{
			Ports: e.getDiscoveryInfoPorts(cmd),
		},
	}
}

// Get the environment of the compose file
func (e *API) getEnvironment() []*mesosproto.Environment_Variable {
	env := make([]*mesosproto.Environment_Variable, 0)
	for k, v := range e.Service.Environment {
		tmp := new(mesosproto.Environment_Variable)
		tmp.Name = util.StringToPointer(k) //p[0]
		// check if the value is a secret
		if strings.Contains(v, "vault://") {
			v = e.Vault.GetKey(v)
		}
		tmp.Value = util.StringToPointer(v)
		env = append(env, tmp)
	}
	return env
}

// Get the environment of the compose file
func (e *API) getVolumes() []*mesosproto.Volume {
	volume := make([]*mesosproto.Volume, 0)

	for _, c := range e.Service.Volumes {
		tmp := new(mesosproto.Volume)
		p := strings.Split(c, ":")
		if len(p) < 2 {
			continue
		}
		tmp.ContainerPath = util.StringToPointer(p[1])
		tmp.Mode = mesosproto.Volume_RW.Enum()
		if len(p) == 3 {
			if strings.ToLower(p[2]) == "ro" {
				tmp.Mode = mesosproto.Volume_RO.Enum()
			}
		}

		driver := "local"
		if e.Compose.Volumes[p[0]].Driver != "" {
			driver = e.Compose.Volumes[p[0]].Driver
		}

		tmp.Source = &mesosproto.Volume_Source{
			Type: mesosproto.Volume_Source_DOCKER_VOLUME.Enum(),
			DockerVolume: &mesosproto.Volume_Source_DockerVolume{
				Name:   util.StringToPointer(p[0]),
				Driver: util.StringToPointer(driver),
			},
		}
		volume = append(volume, tmp)
	}

	return volume
}

// Get custome executer
func (e *API) getExecutor() *mesosproto.ExecutorInfo {
	var executorInfo *mesosproto.ExecutorInfo
	var command string

	if (e.Service.Mesos.Executor != cfg.Executor{}) {
		if e.Service.Mesos.Executor.Command != "" {
			command = e.Service.Mesos.Executor.Command
		}
	}

	if command != "" {
		command = "exec '" + command + "' " + e.getCommand()
		executorID, _ := util.GenUUID()
		environment := e.getEnvironment()
		executorInfo = &mesosproto.ExecutorInfo{
			Name: func() *string { x := path.Base(command); return &x }(),
			Type: mesosproto.ExecutorInfo_CUSTOM.Enum(),
			ExecutorId: &mesosproto.ExecutorID{
				Value: util.StringToPointer(executorID),
			},
			FrameworkId: e.Framework.FrameworkInfo.Id,
			Command: &mesosproto.CommandInfo{
				Value: func() *string { x := command; return &x }(),
				Environment: &mesosproto.Environment{
					Variables: environment,
				},
			},
		}
		executorInfo.Command.Uris = e.getURIs()
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
		// the name of the network driver (CNI) will be configured by getNetworkInfo
		// and it would be a part of the mesos networkinfo.

		// to use CNI or docker network plugins, we have to set the mode to user.
		// it can be overwritten by the driver name.
		if e.Compose.Networks[network].Name != "" {
		  mode = "user"
		}
		if e.Compose.Networks[network].Driver != "" {
		  mode = e.Compose.Networks[network].Driver
		}
	}

	return strings.ToLower(mode)
}

// get the NetworkInfo Name
func (e *API) getNetworkInfo() []*mesosproto.NetworkInfo {
	// get network name
	name := e.getNetworkName(0)

	if len(e.Compose.Networks) > 0 {
		name = e.Compose.Networks[name].Name
	}

	return []*mesosproto.NetworkInfo{{
		Name: util.StringToPointer(name),
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

// get linux info like capabilities
func (e *API) getLinuxInfo() *mesosproto.LinuxInfo {
	linuxInfo := &mesosproto.LinuxInfo{}

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
	var conType string

	if e.Service.ContainerType != "" {
		conType = strings.ToLower(e.Service.ContainerType)
	}

	// if contype and custom executor is unset, then set the contype to DOCKER
	if conType == "" {
		conType = "docker"
	}

	if conType != "mesos" && conType != "docker" {
		conType = "docker"
	}

	return conType
}

func (e *API) getDockerParameter(cmd *cfg.Command) []*mesosproto.Parameter {
	param := cmd.DockerParameter
	if len(param) == 0 {
		param = make([]*mesosproto.Parameter, 0)
	}

	if e.getContainerType() == "docker" {
		// set net-alias if its defined
		alias := e.getNetAlias()
		if alias != "" {
			param = e.addDockerParameter(param, "net-alias", alias)
		}
		// add default volume driver if there is no defined volume
		if len(e.Service.Volumes) == 0 {
			param = e.addDockerParameter(param, "volume-driver", e.Config.DefaultVolumeDriver)
		}
		// configure ulimits
		param = e.getUlimit(param)
	}

	return param
}

func (e *API) getUlimit(param []*mesosproto.Parameter) []*mesosproto.Parameter {
	if e.Service.Ulimits.Memlock.Hard != 0 {
		param = e.addDockerParameter(param, "ulimit", "memlock="+strconv.Itoa(e.Service.Ulimits.Memlock.Hard)+":"+strconv.Itoa(e.Service.Ulimits.Memlock.Soft))
	}

	if e.Service.Ulimits.Nofile.Hard != 0 {
		param = e.addDockerParameter(param, "ulimit", "nofile="+strconv.Itoa(e.Service.Ulimits.Nofile.Hard)+":"+strconv.Itoa(e.Service.Ulimits.Nofile.Soft))
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

	return ""
}

func (e *API) addDockerParameter(current []*mesosproto.Parameter, key string, values string) []*mesosproto.Parameter {
	newValues := &mesosproto.Parameter{
		Key:   func() *string { x := key; return &x }(),
		Value: func() *string { x := values; return &x }(),
	}

	return append(current, newValues)
}

// translate the docker-compose placement constraints to labels
func (e *API) setConstraints(cmd *cfg.Command) {
	if len(e.Service.Deploy.Placement.Constraints) > 0 {
		for _, constraint := range e.Service.Deploy.Placement.Constraints {
			if strings.Contains(constraint, "==") {
				cons := strings.Split(constraint, "==")
				if len(cons) >= 2 {
					if cons[0] == "node.hostname" {
						cmd.Labels = append(cmd.Labels, &mesosproto.Label{Key: util.StringToPointer("__mc_placement_node_hostname"), Value: &cons[1]})
					}
					if cons[0] == "node.platform.os" {
						cmd.Labels = append(cmd.Labels, &mesosproto.Label{Key: util.StringToPointer("__mc_placement_node_platform_os"), Value: &cons[1]})
					}
					if cons[0] == "node.platform.arch" {
						cmd.Labels = append(cmd.Labels, &mesosproto.Label{Key: util.StringToPointer("__mc_placement_node_platform_arch"), Value: &cons[1]})
					}
				}
			}
			if strings.ToLower(constraint) == "unique" {
				val := func() *string { x := "unique"; return &x }()
				cmd.Labels = append(cmd.Labels, &mesosproto.Label{Key: util.StringToPointer("__mc_placement"), Value: val})
			}
		}
	}
}
