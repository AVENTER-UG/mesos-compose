package scheduler

import (
	"encoding/json"
	"strings"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/AVENTER-UG/util/util"
	"github.com/sirupsen/logrus"
)

func (e *Scheduler) defaultResources(cmd *cfg.Command) []*mesosproto.Resource {
	PORT := "ports"
	CPU := "cpus"
	MEM := "mem"
	DISK := "disk"
	GPUs := "gpus"
	cpu := cmd.CPU
	mem := cmd.Memory
	disk := cmd.Disk
	gpus := cmd.GPUs

	// FIX: https://github.com/AVENTER-UG/mesos-compose/issues/8
	// If the task already exists from a prev mesos-compose version, disk
	// would be unset.
	if disk <= e.Config.Disk {
		disk = e.Config.Disk
	}

	res := []*mesosproto.Resource{
		{
			Name:   &CPU,
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: &cpu},
		},
		{
			Name:   &MEM,
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: &mem},
		},
		{
			Name:   &DISK,
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: &disk},
		},
		{
			Name:   &GPUs,
			Type:   mesosproto.Value_SCALAR.Enum(),
			Scalar: &mesosproto.Value_Scalar{Value: &gpus},
		},
	}

	if cmd.DockerPortMappings != nil {
		for _, p := range cmd.DockerPortMappings {
			port := mesosproto.Resource{
				Name: &PORT,
				Type: mesosproto.Value_RANGES.Enum(),
				Ranges: &mesosproto.Value_Ranges{
					Range: []*mesosproto.Value_Range{
						{
							Begin: util.Uint64ToPointer(uint64(*p.HostPort)),
							End:   util.Uint64ToPointer(uint64(*p.HostPort)),
						},
					},
				},
			}
			res = append(res, &port)
		}
	}

	return res
}

// PrepareTaskInfoExecuteContainer will create the TaskInfo Protobuf for Mesos
// nolint: gocyclo
func (e *Scheduler) PrepareTaskInfoExecuteContainer(agent *mesosproto.AgentID, cmd *cfg.Command) []*mesosproto.TaskInfo {
	d, _ := json.Marshal(&cmd)
	logrus.WithField("func", "scheduler.PrepareTaskInfoExecuteContainer").Trace("HandleOffers cmd: ", util.PrettyJSON(d))

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

	msg.Name = &cmd.TaskName
	msg.TaskId = &mesosproto.TaskID{
		Value: &cmd.TaskID,
	}
	msg.AgentId = agent
	msg.Resources = e.defaultResources(cmd)
	// Do not set the cmd.Command if the mesos task is an executor.
	if cmd.Mesos.Executor.Command == "" {
		if cmd.Command == "" {
			msg.Command = &mesosproto.CommandInfo{
				Shell:       &cmd.Shell,
				Uris:        cmd.Uris,
				Environment: cmd.Environment,
			}
		} else {
			msg.Command = &mesosproto.CommandInfo{
				Shell:       &cmd.Shell,
				Value:       &cmd.Command,
				Uris:        cmd.Uris,
				Environment: cmd.Environment,
			}
		}

		if len(cmd.Arguments) > 0 {
			msg.Command.Arguments = cmd.Arguments
		}
	}

	// force to pull the container image
	forcePull := true
	if cmd.PullPolicy == "missing" {
		forcePull = false
	}

	// ExecutorInfo or CommandInfo/Container, both at the same time is not supported
	if contype != nil && cmd.Mesos.Executor.Command == "" {
		msg.Container = &mesosproto.ContainerInfo{}
		msg.Container.Type = contype
		msg.Container.Volumes = cmd.Volumes

		if cmd.ContainerType == "docker" {
			msg.Container.Docker = &mesosproto.ContainerInfo_DockerInfo{
				Image:          util.StringToPointer(cmd.ContainerImage),
				Network:        networkMode,
				PortMappings:   cmd.DockerPortMappings,
				Privileged:     &cmd.Privileged,
				Parameters:     cmd.DockerParameter,
				ForcePullImage: func() *bool { x := forcePull; return &x }(),
			}
		}

		if cmd.ContainerType == "mesos" {
			var cachedImage bool = false
			msg.Container.TtyInfo = &mesosproto.TTYInfo{}
			msg.Container.TtyInfo.WindowSize = &mesosproto.TTYInfo_WindowSize{}
			msg.Container.TtyInfo.WindowSize.Columns = util.Uint32ToPointer(80)
			msg.Container.TtyInfo.WindowSize.Rows = util.Uint32ToPointer(25)
			msg.Container.Mesos = &mesosproto.ContainerInfo_MesosInfo{}
			msg.Container.Mesos.Image = &mesosproto.Image{}
			msg.Container.Mesos.Image.Type = mesosproto.Image_DOCKER.Enum()
			msg.Container.Mesos.Image.Cached = &cachedImage
			msg.Container.Mesos.Image.Docker = &mesosproto.Image_Docker{}
			msg.Container.Mesos.Image.Docker.Name = util.StringToPointer(cmd.ContainerImage)
		}

		msg.Container.NetworkInfos = cmd.NetworkInfo

		if cmd.Hostname != "" {
			msg.Container.Hostname = &cmd.Hostname
		}

		msg.Container.LinuxInfo = cmd.LinuxInfo
	}

	if cmd.Discovery != (&mesosproto.DiscoveryInfo{}) {
		msg.Discovery = cmd.Discovery
	}

	if cmd.Labels != nil {
		msg.Labels = &mesosproto.Labels{
			Labels: cmd.Labels,
		}
	}

	if cmd.Mesos.Executor.Command != "" {
		// FIX: https://github.com/AVENTER-UG/mesos-compose/issues/7
		cmd.Executor.Resources = e.defaultResources(cmd)
		msg.Executor = cmd.Executor
		if cmd.ContainerType != "" {
			msg.Executor.Container = &mesosproto.ContainerInfo{}
			msg.Executor.Container.Type = contype
			msg.Executor.Container.Volumes = cmd.Volumes
			msg.Executor.Container.Docker = &mesosproto.ContainerInfo_DockerInfo{
				Image:          util.StringToPointer(cmd.ContainerImage),
				Network:        mesosproto.ContainerInfo_DockerInfo_HOST.Enum(),
				PortMappings:   cmd.DockerPortMappings,
				Privileged:     &cmd.Privileged,
				ForcePullImage: func() *bool { x := true; return &x }(),
			}
		}
	}

	if cmd.EnableHealthCheck {
		if cmd.Health != nil && cmd.Health.GetType() != 0 {
			msg.HealthCheck = cmd.Health
		} else {
			msg.HealthCheck = nil
		}
	}

	d, _ = json.Marshal(&msg)
	logrus.WithField("func", "scheduler.PrepareTaskInfoExecuteContainer").Trace("HandleOffers msg: ", util.PrettyJSON(d))

	return []*mesosproto.TaskInfo{&msg}
}
