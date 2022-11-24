package scheduler

import (
	"encoding/json"
	"strings"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/AVENTER-UG/util/util"
	"github.com/sirupsen/logrus"
)

func (e *Scheduler) defaultResources(cmd cfg.Command) []mesosproto.Resource {
	PORT := "ports"
	CPU := "cpus"
	MEM := "mem"
	DISK := "disk"
	cpu := cmd.CPU
	mem := cmd.Memory
	disk := cmd.Disk

	// FIX: https://github.com/AVENTER-UG/mesos-compose/issues/8
	// If the task already exists from a prev mesos-compose version, disk
	// would be unset.
	if disk <= e.Config.Disk {
		disk = e.Config.Disk
	}

	res := []mesosproto.Resource{
		{
			Name:   CPU,
			Type:   mesosproto.SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: cpu},
		},
		{
			Name:   MEM,
			Type:   mesosproto.SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: mem},
		},
		{
			Name:   DISK,
			Type:   mesosproto.SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: disk},
		},
	}

	if cmd.DockerPortMappings != nil {
		for _, p := range cmd.DockerPortMappings {
			port := mesosproto.Resource{
				Name: PORT,
				Type: mesosproto.RANGES.Enum(),
				Ranges: &mesosproto.Value_Ranges{
					Range: []mesosproto.Value_Range{
						{
							Begin: uint64(p.HostPort),
							End:   uint64(p.HostPort),
						},
					},
				},
			}
			res = append(res, port)
		}
	}

	return res
}

// PrepareTaskInfoExecuteContainer will create the TaskInfo Protobuf for Mesos
// nolint: gocyclo
func (e *Scheduler) PrepareTaskInfoExecuteContainer(agent mesosproto.AgentID, cmd cfg.Command) ([]mesosproto.TaskInfo, error) {
	d, _ := json.Marshal(&cmd)
	logrus.Debug("HandleOffers cmd: ", util.PrettyJSON(d))

	// Set Container Type
	var contype *mesosproto.ContainerInfo_Type
	switch cmd.ContainerType {
	case "mesos":
		contype = mesosproto.ContainerInfo_MESOS.Enum()
	case "docker":
		contype = mesosproto.ContainerInfo_DOCKER.Enum()
	default:
		contype = nil
	}

	// Set Container Network Mode
	networkMode := mesosproto.ContainerInfo_DockerInfo_BRIDGE.Enum()
	switch strings.ToLower(cmd.NetworkMode) {
	case "host":
		networkMode = mesosproto.ContainerInfo_DockerInfo_HOST.Enum()
	case "none":
		networkMode = mesosproto.ContainerInfo_DockerInfo_NONE.Enum()
	case "user":
		networkMode = mesosproto.ContainerInfo_DockerInfo_USER.Enum()
	case "bridge":
		networkMode = mesosproto.ContainerInfo_DockerInfo_BRIDGE.Enum()
	}

	var msg mesosproto.TaskInfo

	msg.Name = cmd.TaskName
	msg.TaskID = mesosproto.TaskID{
		Value: cmd.TaskID,
	}
	msg.AgentID = agent
	msg.Resources = e.defaultResources(cmd)

	if e.getLabelValue("biz.aventer.mesos_compose.executor", cmd) == "" {
		if cmd.Command == "" {
			msg.Command = &mesosproto.CommandInfo{
				Shell:       &cmd.Shell,
				URIs:        cmd.Uris,
				Environment: &cmd.Environment,
			}
		} else {
			msg.Command = &mesosproto.CommandInfo{
				Shell:       &cmd.Shell,
				Value:       &cmd.Command,
				URIs:        cmd.Uris,
				Environment: &cmd.Environment,
			}
		}
	}

	// force to pull the container image
	forcePull := true
	if cmd.PullPolicy == "missing" {
		forcePull = false
	}

	// ExecutorInfo or CommandInfo/Container, both is not supportet
	if contype != nil && e.getLabelValue("biz.aventer.mesos_compose.executor", cmd) == "" {
		msg.Container = &mesosproto.ContainerInfo{}
		msg.Container.Type = contype
		msg.Container.Volumes = cmd.Volumes
		msg.Container.Docker = &mesosproto.ContainerInfo_DockerInfo{
			Image:          cmd.ContainerImage,
			Network:        networkMode,
			PortMappings:   cmd.DockerPortMappings,
			Privileged:     &cmd.Privileged,
			Parameters:     cmd.DockerParameter,
			ForcePullImage: func() *bool { x := forcePull; return &x }(),
		}
		msg.Container.NetworkInfos = cmd.NetworkInfo

		if cmd.Hostname != "" {
			msg.Container.Hostname = &cmd.Hostname
		}

		msg.Container.LinuxInfo = &cmd.LinuxInfo
	}

	if cmd.Discovery != (mesosproto.DiscoveryInfo{}) {
		msg.Discovery = &cmd.Discovery
	}

	if cmd.Labels != nil {
		msg.Labels = &mesosproto.Labels{
			Labels: cmd.Labels,
		}
	}

	if e.getLabelValue("biz.aventer.mesos_compose.executor", cmd) != "" {
		// FIX: https://github.com/AVENTER-UG/mesos-compose/issues/7
		cmd.Executor.Resources = e.defaultResources(cmd)
		msg.Executor = &cmd.Executor
		if cmd.ContainerType != "" {
			msg.Executor.Container = &mesosproto.ContainerInfo{}
			msg.Executor.Container.Type = contype
			msg.Executor.Container.Volumes = cmd.Volumes
			msg.Executor.Container.Docker = &mesosproto.ContainerInfo_DockerInfo{
				Image:          cmd.ContainerImage,
				Network:        mesosproto.ContainerInfo_DockerInfo_HOST.Enum(),
				PortMappings:   cmd.DockerPortMappings,
				Privileged:     &cmd.Privileged,
				ForcePullImage: func() *bool { x := true; return &x }(),
			}
		}
	}

	return []mesosproto.TaskInfo{msg}, nil
}
