package types

import (
	"time"
)

// Config is a struct of the framework configuration
type Config struct {
	Principal         string
	LogLevel          string
	MinVersion        string
	AppName           string
	EnableSyslog      bool
	Hostname          string
	Listen            string
	Domain            string
	Credentials       UserCredentials
	PrefixHostname    string
	PrefixTaskName    string
	CPU               float64
	Memory            float64
	Disk              float64
	RedisServer       string
	RedisPassword     string
	RedisDB           int
	SkipSSL           bool
	SSLKey            string
	SSLCrt            string
	Suppress          bool
	EventLoopTime     time.Duration
	ReconcileLoopTime time.Duration
	VaultToken        string
	VaultURL          string
	VaultTimeout      time.Duration
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
	Environment   []string                `yaml:"environment"`
	DependsOn     []string                `yaml:"depends_on"`
	Ports         []string                `yaml:"ports"`
	Image         string                  `yaml:"image"`
	Labels        map[string]interface{}  `yaml:"labels"`
	NetworkMode   string                  `yaml:"network_mode"`
	Privileged    bool                    `yaml:"privileged"`
	Command       string                  `yaml:"command"`
	Deploy        Deploy                  `yaml:"deploy"`
	Hostname      string                  `yaml:"hostname"`
	ContainerName string                  `yaml:"container_name"`
	CapAdd        []string                `yaml:"cap_add"`
	CapDrop       []string                `yaml:"cap_drop"`
	PullPolicy    string                  `yaml:"pull_policy" default:"always"`
}

// Deploy - The mesos resources to deploy a task
type Deploy struct {
	Placement Placement `yaml:"placement"`
	Replicas  string    `yaml:"replicas"`
	Resources struct {
		Limits struct {
			CPUs   string `yaml:"cpus"`
			Memory string `yaml:"memory"`
		} `yaml:"limits"`
	} `yaml:"resources"`
}

// Placement - The docker-compose placement
type Placement struct {
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
	Driver string `yaml:"driver"`
}

// ErrorMsg hold the structure of error messages
type ErrorMsg struct {
	Message  string
	Number   int
	Function string
}
