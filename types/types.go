package types

import (
	"encoding/json"
	"plugin"
	"time"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	"gopkg.in/yaml.v3"
)

// Config is a struct of the framework configuration
type Config struct {
	Principal                  string
	LogLevel                   string
	MinVersion                 string
	AppName                    string
	EnableSyslog               bool
	Hostname                   string
	Listen                     string
	Domain                     string
	Credentials                UserCredentials
	PrefixHostname             string
	PrefixTaskName             string
	CPU                        float64
	Memory                     float64
	Disk                       float64
	RedisServer                string
	RedisPassword              string
	RedisDB                    int
	RedisPoolSize              int
	SkipSSL                    bool
	SSLKey                     string
	SSLCrt                     string
	Suppress                   bool
	EventLoopTime              time.Duration
	ReconcileLoopTime          time.Duration
	VaultToken                 string
	VaultURL                   string
	VaultTimeout               time.Duration
	DefaultVolumeDriver        string
	DiscoveryInfoNameDelimiter string
	DiscoveryPortNameDelimiter string
	Plugins                    map[string]*plugin.Plugin
	PluginsEnable              bool

	EnableGPUAllocation bool
	HostConstraintsList []string
	TaskLostRemovesTask bool
	CORSAllowedOrigins  []string
}

// UserCredentials - The Username and Password to authenticate against this framework
type UserCredentials struct {
	Username string
	Password string
}

// Compose - The main structure of the supported docker-compose syntax
type Compose struct {
	Version  string              `yaml:"version"`
	Services map[string]Service  `yaml:"services"`
	Networks map[string]Networks `yaml:"networks"`
	Volumes  map[string]Volumes  `yaml:"volumes"`
}

// Service - The docker-compose service parameters
type Service struct {
	Network       string                  `yaml:"network"`
	Networks      map[string]NetworksLong `yaml:"networks"`
	Build         string                  `yaml:"build"`
	Restart       string                  `yaml:"restart" default:"unless-stopped"`
	Volumes       []string                `yaml:"volumes"`
	Environment   map[string]string       `yaml:"environment"`
	Arguments     []string                `yaml:"arguments"`
	DependsOn     []string                `yaml:"depends_on"`
	Ports         []string                `yaml:"ports"`
	Image         string                  `yaml:"image"`
	Labels        map[string]interface{}  `yaml:"labels"`
	NetworkMode   string                  `yaml:"network_mode"`
	Privileged    bool                    `yaml:"privileged"`
	Command       string                  `yaml:"command"`
	Shell         bool                    `yaml:"shell"`
	Deploy        Deploy                  `yaml:"deploy"`
	Hostname      string                  `yaml:"hostname"`
	ContainerName string                  `yaml:"container_name"`
	ContainerType string                  `yaml:"container_type" default:"docker"`
	CapAdd        []string                `yaml:"cap_add"`
	CapDrop       []string                `yaml:"cap_drop"`
	PullPolicy    string                  `yaml:"pull_policy" default:"always"`
	Ulimits       Ulimits                 `yaml:"ulimits"`
	Mesos         Mesos                   `yaml:"mesos"`
	HealthCheck   *mesosproto.HealthCheck `yaml:"healthcheck"`
	GPUs          GPUs                    `yaml:"gpus"`
}

// Gpus holds the config for gpu usage
type GPUs struct {
	Driver string `yaml:"driver"`
	Device int    `yaml:"device"`
}

// Mesos custom mesos task configuration
type Mesos struct {
	TaskName string                        `yaml:"task_name"`
	Executor Executor                      `yaml:"executor"`
	Fetch    []*mesosproto.CommandInfo_URI `yaml:"fetch"`
}

// Executor to configure the to use executor
type Executor struct {
	Command string `yaml:"command"`
}

// Ulimits - Configure ulimits of a mesos task
type Ulimits struct {
	Memlock struct {
		Soft int `yaml:"soft"`
		Hard int `yaml:"hard"`
	} `yaml:"memlock"`
	Nofile struct {
		Soft int `yaml:"soft"`
		Hard int `yaml:"hard"`
	} `yaml:"nofile"`
}

// Deploy - Deploy information of a mesos task
type Deploy struct {
	Placement Placement `yaml:"placement"`
	Replicas  string    `yaml:"replicas"`
	Resources struct {
		Limits struct {
			CPUs   float64 `yaml:"cpus"`
			Memory float64 `yaml:"memory"`
		} `yaml:"limits"`
	} `yaml:"resources"`
	Runtime string `yaml:"runtime"`
}

// Placement - The docker-compose placement
type Placement struct {
	Attributes  []string `yaml:"attributes"`
	Constraints []string `yaml:"constraints"`
}

// Networks - The docker-compose network syntax
type Networks struct {
	External bool   `yaml:"external"`
	Name     string `yaml:"name"`
	Driver   string `yaml:"driver"`
}

// NetworksLong - Supportet structure for Networks
type NetworksLong struct {
	Aliases []string `yaml:"aliases"`
}

// Volumes - The docker-compose volumes syntax
type Volumes struct {
	Driver             string                `yaml:"driver"`
	PluginName         string                `yaml:"plugin_name"`
	StaticProvisioning CSIStaticProvisioning `yaml:"static_provisioning"`
}

// CSIStaticProvisioning describes a pre-provisioned CSI volume.
type CSIStaticProvisioning struct {
	VolumeID         string              `yaml:"volume_id"`
	VolumeCapability CSIVolumeCapability `yaml:"volume_capability"`
	VolumeContext    map[string]string   `yaml:"volume_context"`
}

// CSIVolumeCapability describes how a CSI volume is accessed.
type CSIVolumeCapability struct {
	Mount      *CSIVolumeMount `yaml:"mount"`
	AccessMode CSIAccessMode   `yaml:"access_mode"`
}

// CSIVolumeMount contains filesystem mount options for a CSI volume.
type CSIVolumeMount struct {
	FsType     string   `yaml:"fs_type"`
	MountFlags []string `yaml:"mount_flags"`
}

// CSIAccessMode accepts both the concise scalar form and Mesos' object form.
type CSIAccessMode struct {
	Mode string `yaml:"mode"`
}

func (m *CSIAccessMode) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		m.Mode = value.Value
		return nil
	}
	type accessMode CSIAccessMode
	var decoded accessMode
	if err := value.Decode(&decoded); err != nil {
		return err
	}
	*m = CSIAccessMode(decoded)
	return nil
}

// ErrorMsg hold the structure of error messages
type ErrorMsg struct {
	Message  string
	Number   int
	Function string
}

type FrameworkConfig struct {
	FrameworkHostname     string
	FrameworkPort         string
	FrameworkBind         string
	FrameworkUser         string
	FrameworkName         string
	FrameworkRole         string
	FrameworkInfo         mesosproto.FrameworkInfo
	FrameworkInfoFile     string
	FrameworkInfoFilePath string
	PortRangeFrom         int
	PortRangeTo           int
	CommandChan           chan Command `json:"-"`
	Username              string
	Password              string
	MesosMasterServer     string
	MesosSSL              bool
	MesosStreamID         string
	TaskID                string
	SSL                   bool
	State                 map[string]State
}

// Command is a chan which include all the Information about the started tasks
type Command struct {
	ContainerImage     string                                             `json:"container_image,omitempty"`
	ContainerType      string                                             `json:"container_type,omitempty"`
	TaskName           string                                             `json:"task_name,omitempty"`
	Command            string                                             `json:"command,omitempty"`
	Hostname           string                                             `json:"hostname,omitempty"`
	Domain             string                                             `json:"domain,omitempty"`
	Privileged         bool                                               `json:"privileged,omitempty"`
	NetworkMode        string                                             `json:"network_mode,omitempty"`
	Volumes            []*mesosproto.Volume                               `protobuf:"bytes,1,rep,name=volumes" json:"volumes,omitempty"`
	Shell              bool                                               `protobuf:"varint,2,opt,name=shell,def=1" json:"shell,omitempty"`
	Uris               []*mesosproto.CommandInfo_URI                      `protobuf:"bytes,3,rep,name=uris" json:"uris,omitempty"`
	Environment        *mesosproto.Environment                            `protobuf:"bytes,4,opt,name=environment" json:"environment,omitempty"`
	NetworkInfo        []*mesosproto.NetworkInfo                          `protobuf:"bytes,5,opt,name=networkinfo" json:"networkinfo,omitempty"`
	DockerPortMappings []*mesosproto.ContainerInfo_DockerInfo_PortMapping `protobuf:"bytes,6,rep,name=port_mappings,json=portMappings" json:"port_mappings,omitempty"`
	DockerParameter    []*mesosproto.Parameter                            `protobuf:"bytes,7,rep,name=parameters" json:"parameters,omitempty"`
	Arguments          []string                                           `protobuf:"bytes,8,rep,name=arguments" json:"arguments,omitempty"`
	Discovery          *mesosproto.DiscoveryInfo                          `protobuf:"bytes,9,opt,name=discovery" json:"discovery,omitempty"`
	Executor           *mesosproto.ExecutorInfo                           `protobuf:"bytes,10,opt,name=executor" json:"executor,omitempty"`
	Restart            string
	InternalID         int
	TaskID             string
	Memory             float64
	Mesos              Mesos
	CPU                float64
	GPUs               float64
	Disk               float64
	Agent              string
	Labels             []*mesosproto.Label
	State              string
	StateTime          time.Time
	Instances          int
	LinuxInfo          *mesosproto.LinuxInfo `protobuf:"bytes,11,opt,name=linux_info,json=linuxInfo" json:"linux_info,omitempty"`
	PullPolicy         string
	EnableHealthCheck  bool
	Health             *mesosproto.HealthCheck
	MesosAgent         MesosSlaves
	Attributes         []*mesosproto.Label
}

// UnmarshalJSON keeps Redis task recovery compatible with CSI volumes. The
// generated protobuf JSON representation of the CSI access_type oneof cannot
// be unmarshaled by encoding/json into its interface field.
func (c *Command) UnmarshalJSON(data []byte) error {
	type commandAlias Command
	if err := json.Unmarshal(data, (*commandAlias)(c)); err == nil {
		return nil
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	volumesData := fields["volumes"]
	delete(fields, "volumes")
	withoutVolumes, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	*c = Command{}
	if err := json.Unmarshal(withoutVolumes, (*commandAlias)(c)); err != nil {
		return err
	}
	if volumesData == nil {
		return nil
	}

	var volumeData []map[string]json.RawMessage
	if err := json.Unmarshal(volumesData, &volumeData); err != nil {
		return err
	}
	c.Volumes = make([]*mesosproto.Volume, 0, len(volumeData))
	for _, fields := range volumeData {
		sourceData := fields["source"]
		delete(fields, "source")
		volumeBytes, err := json.Marshal(fields)
		if err != nil {
			return err
		}
		var volume mesosproto.Volume
		if err := json.Unmarshal(volumeBytes, &volume); err != nil {
			return err
		}
		if sourceData != nil {
			if err := unmarshalVolumeSource(sourceData, &volume); err != nil {
				return err
			}
		}
		c.Volumes = append(c.Volumes, &volume)
	}
	return nil
}

func unmarshalVolumeSource(data []byte, volume *mesosproto.Volume) error {
	var source struct {
		Type         *mesosproto.Volume_Source_Type         `json:"type"`
		DockerVolume *mesosproto.Volume_Source_DockerVolume `json:"docker_volume"`
		CSIVolume    json.RawMessage                        `json:"csi_volume"`
	}
	if err := json.Unmarshal(data, &source); err != nil {
		return err
	}
	volume.Source = &mesosproto.Volume_Source{
		Type:         source.Type,
		DockerVolume: source.DockerVolume,
	}
	if source.CSIVolume == nil {
		return nil
	}

	var csi struct {
		PluginName         *string         `json:"plugin_name"`
		StaticProvisioning json.RawMessage `json:"static_provisioning"`
	}
	if err := json.Unmarshal(source.CSIVolume, &csi); err != nil {
		return err
	}
	var provisioning struct {
		VolumeID         *string           `json:"volume_id"`
		VolumeCapability json.RawMessage   `json:"volume_capability"`
		VolumeContext    map[string]string `json:"volume_context"`
	}
	if err := json.Unmarshal(csi.StaticProvisioning, &provisioning); err != nil {
		return err
	}
	var capability struct {
		AccessType      json.RawMessage                                                 `json:"AccessType"`
		ProtoAccessType json.RawMessage                                                 `json:"access_type"`
		AccessMode      *mesosproto.Volume_Source_CSIVolume_VolumeCapability_AccessMode `json:"access_mode"`
	}
	if err := json.Unmarshal(provisioning.VolumeCapability, &capability); err != nil {
		return err
	}
	mesosCapability := &mesosproto.Volume_Source_CSIVolume_VolumeCapability{
		AccessMode: capability.AccessMode,
	}
	if len(capability.AccessType) == 0 {
		capability.AccessType = capability.ProtoAccessType
	}
	if len(capability.AccessType) > 0 {
		var accessType map[string]json.RawMessage
		if err := json.Unmarshal(capability.AccessType, &accessType); err != nil {
			return err
		}
		mountData := accessType["Mount"]
		if mountData == nil {
			mountData = accessType["mount"]
		}
		if mountData != nil {
			originalMountData := mountData
			var mountWrapper map[string]json.RawMessage
			if err := json.Unmarshal(mountData, &mountWrapper); err != nil {
				return err
			}
			if mountWrapper["Mount"] != nil {
				mountData = mountWrapper["Mount"]
			} else if mountWrapper["mount"] != nil {
				mountData = mountWrapper["mount"]
			} else {
				mountData = originalMountData
			}
			var mount mesosproto.Volume_Source_CSIVolume_VolumeCapability_MountVolume
			if err := json.Unmarshal(mountData, &mount); err != nil {
				return err
			}
			mesosCapability.AccessType = &mesosproto.Volume_Source_CSIVolume_VolumeCapability_Mount{Mount: &mount}
		}
	}
	volume.Source.CsiVolume = &mesosproto.Volume_Source_CSIVolume{
		PluginName: csi.PluginName,
		StaticProvisioning: &mesosproto.Volume_Source_CSIVolume_StaticProvisioning{
			VolumeId:         provisioning.VolumeID,
			VolumeCapability: mesosCapability,
			VolumeContext:    provisioning.VolumeContext,
		},
	}
	return nil
}

// State will have the state of all tasks stated by this framework
type State struct {
	Command Command                `json:"command"`
	Status  *mesosproto.TaskStatus `json:"status"`
}

// MesosAgents
type MesosAgent struct {
	Slaves          []MesosSlaves `json:"slaves"`
	RecoveredSlaves []interface{} `json:"recovered_slaves"`
}

// MesosSlaves ..
type MesosSlaves struct {
	ID         string `json:"id"`
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	Attributes struct {
	} `json:"attributes"`
	Pid              string  `json:"pid"`
	RegisteredTime   float64 `json:"registered_time"`
	ReregisteredTime float64 `json:"reregistered_time"`
	Resources        struct {
		Disk  float64 `json:"disk"`
		Mem   float64 `json:"mem"`
		Gpus  float64 `json:"gpus"`
		Cpus  float64 `json:"cpus"`
		Ports string  `json:"ports"`
	} `json:"resources"`
	UsedResources struct {
		Disk  float64 `json:"disk"`
		Mem   float64 `json:"mem"`
		Gpus  float64 `json:"gpus"`
		Cpus  float64 `json:"cpus"`
		Ports string  `json:"ports"`
	} `json:"used_resources"`
	OfferedResources struct {
		Disk float64 `json:"disk"`
		Mem  float64 `json:"mem"`
		Gpus float64 `json:"gpus"`
		Cpus float64 `json:"cpus"`
	} `json:"offered_resources"`
	ReservedResources struct {
	} `json:"reserved_resources"`
	UnreservedResources struct {
		Disk  float64 `json:"disk"`
		Mem   float64 `json:"mem"`
		Gpus  float64 `json:"gpus"`
		Cpus  float64 `json:"cpus"`
		Ports string  `json:"ports"`
	} `json:"unreserved_resources"`
	Active                bool     `json:"active"`
	Deactivated           bool     `json:"deactivated"`
	Version               string   `json:"version"`
	Capabilities          []string `json:"capabilities"`
	ReservedResourcesFull struct {
	} `json:"reserved_resources_full"`
	UnreservedResourcesFull []struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Scalar struct {
			Value float64 `json:"value"`
		} `json:"scalar,omitempty"`
		Role   string `json:"role"`
		Ranges struct {
			Range []struct {
				Begin int `json:"begin"`
				End   int `json:"end"`
			} `json:"range"`
		} `json:"ranges,omitempty"`
	} `json:"unreserved_resources_full"`
	UsedResourcesFull []struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Scalar struct {
			Value float64 `json:"value"`
		} `json:"scalar,omitempty"`
		Role           string `json:"role"`
		AllocationInfo struct {
			Role string `json:"role"`
		} `json:"allocation_info"`
		Ranges struct {
			Range []struct {
				Begin int `json:"begin"`
				End   int `json:"end"`
			} `json:"range"`
		} `json:"ranges,omitempty"`
	} `json:"used_resources_full"`
	OfferedResourcesFull []interface{} `json:"offered_resources_full"`
}

// MesosTasks hold the information of the task
type MesosTasks struct {
	Tasks []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		FrameworkID string `json:"framework_id"`
		ExecutorID  string `json:"executor_id"`
		SlaveID     string `json:"slave_id"`
		AgentID     string `json:"agent_id"`
		State       string `json:"state"`
		Resources   struct {
			Disk float64 `json:"disk"`
			Mem  float64 `json:"mem"`
			Gpus float64 `json:"gpus"`
			Cpus float64 `json:"cpus"`
		} `json:"resources"`
		Role     string `json:"role"`
		Statuses []struct {
			State           string  `json:"state"`
			Timestamp       float64 `json:"timestamp"`
			ContainerStatus struct {
				ContainerID struct {
					Value string `json:"value"`
				} `json:"container_id"`
				NetworkInfos []*mesosproto.NetworkInfo `json:"network_infos"`
			} `json:"container_status,omitempty"`
		} `json:"statuses"`
		Discovery mesosproto.DiscoveryInfo `json:"discovery"`
		Container mesosproto.ContainerInfo `json:"container"`
	} `json:"tasks"`
}
