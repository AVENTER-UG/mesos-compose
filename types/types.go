package types

import (
	"context"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	goredis "github.com/go-redis/redis/v8"
)

// Config is a struct of the framework configuration
type Config struct {
	Principal       string
	LogLevel        string
	MinVersion      string
	AppName         string
	EnableSyslog    bool
	Hostname        string
	Listen          string
	Domain          string
	Credentials     UserCredentials
	PrefixHostname  string
	PrefixTaskName  string
	CPU             float64
	Memory          float64
	FrameworkConfig mesosutil.FrameworkConfig
	RedisServer     string
	RedisClient     *goredis.Client
	RedisCTX        context.Context
}

// UserCredentials - The Username and Password to authenticate against this framework
type UserCredentials struct {
	Username string
	Password string
}

// Yaml2Go
type Compose struct {
	Version  string              `yaml:"version"`
	Services map[string]Service  `yaml:"services"`
	Networks map[string]Networks `yaml:"networks"`
}

// Web
type Service struct {
	Network     []string `yaml:"network"`
	Build       string   `yaml:"build"`
	Restart     string   `yaml:"restart"`
	Volumes     []string `yaml:"volumes"`
	Environment []string `yaml:"environment"`
	DependsOn   []string `yaml:"depends_on"`
	Ports       []string `yaml:"ports"`
	Image       string   `yaml:"image"`
	Labels      Labels   `yaml:"labels"`
	NetworkMode string   `yaml:"network_mode"`
	Privileged  bool     `yaml:"priviliged"`
	Command     []string `yaml:"command"`
	Deploy      Deploy   `yaml:"deploy"`
	Hostname    string   `yaml:"hostname"`
}

type Deploy struct {
	Resources struct {
		Limits struct {
			CPUs   string `yaml:"cpus"`
			Memory string `yaml:"memory"`
		} `yaml:"limits"`
	} `yaml:"resources"`
}

type Networks struct {
	External bool   `yaml:"external"`
	Name     string `yaml:"name"`
	Driver   string `yaml:"driver"`
}

type Labels struct {
	ContainerType string `yaml:"biz.aventer.mesos_compose.container_type"`
}

// Command is a chan which include all the Information about the started tasks
type Command struct {
	ContainerImage     string                                            `json:"container_image,omitempty"`
	ContainerType      string                                            `json:"container_type,omitempty"`
	TaskName           string                                            `json:"task_name,omitempty"`
	Command            string                                            `json:"command,omitempty"`
	Hostname           string                                            `json:"hostname,omitempty"`
	Domain             string                                            `json:"domain,omitempty"`
	Privileged         bool                                              `json:"privileged,omitempty"`
	NetworkMode        string                                            `json:"network_mode,omitempty"`
	Volumes            []mesosproto.Volume                               `protobuf:"bytes,2,rep,name=volumes" json:"volumes,omitempty"`
	Shell              bool                                              `protobuf:"varint,6,opt,name=shell,def=1" json:"shell,omitempty"`
	Uris               []mesosproto.CommandInfo_URI                      `protobuf:"bytes,1,rep,name=uris" json:"uris,omitempty"`
	Environment        mesosproto.Environment                            `protobuf:"bytes,2,opt,name=environment" json:"environment,omitempty"`
	NetworkInfo        []mesosproto.NetworkInfo                          `protobuf:"bytes,2,opt,name=networkinfo" json:"networkinfo,omitempty"`
	DockerPortMappings []mesosproto.ContainerInfo_DockerInfo_PortMapping `protobuf:"bytes,3,rep,name=port_mappings,json=portMappings" json:"port_mappings,omitempty"`
	DockerParameter    []mesosproto.Parameter                            `protobuf:"bytes,5,rep,name=parameters" json:"parameters,omitempty"`
	Arguments          []string                                          `protobuf:"bytes,7,rep,name=arguments" json:"arguments,omitempty"`
	Discovery          mesosproto.DiscoveryInfo                          `protobuf:"bytes,12,opt,name=discovery" json:"discovery,omitempty"`
	Executor           mesosproto.ExecutorInfo
	TaskID             uint64
	Memory             float64
	CPU                float64
	Agent              string
	Labels             []mesosproto.Label
}
