package types

import (
	"context"

	goredis "github.com/go-redis/redis/v8"
)

// Config is a struct of the framework configuration
type Config struct {
	Principal      string
	LogLevel       string
	MinVersion     string
	AppName        string
	EnableSyslog   bool
	Hostname       string
	Listen         string
	Domain         string
	Credentials    UserCredentials
	PrefixHostname string
	PrefixTaskName string
	CPU            float64
	Memory         float64
	RedisServer    string
	RedisClient    *goredis.Client
	RedisCTX       context.Context
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
	Volumes  map[string]Volumes  `yaml:"volumes"`
}

// Web
type Service struct {
	Network     []string               `yaml:"network"`
	Build       string                 `yaml:"build"`
	Restart     string                 `yaml:"restart"`
	Volumes     []string               `yaml:"volumes"`
	Environment []string               `yaml:"environment"`
	DependsOn   []string               `yaml:"depends_on"`
	Ports       []string               `yaml:"ports"`
	Image       string                 `yaml:"image"`
	Labels      map[string]interface{} `yaml:"labels"`
	NetworkMode string                 `yaml:"network_mode"`
	Privileged  bool                   `yaml:"priviliged"`
	Command     []string               `yaml:"command"`
	Deploy      Deploy                 `yaml:"deploy"`
	Hostname    string                 `yaml:"hostname"`
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

type Volumes struct {
	Driver string `yaml:"driver"`
}

// ErrorMsg hold the structure of error messages
type ErrorMsg struct {
	Message  string
	Number   int
	Function string
}
