package api

import (
	"testing"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
)

func TestGetShell(t *testing.T) {
	var e API
	e.Service.Command = "test"
	e.Service.Shell = true

	res := e.getShell()

	if !res {
		t.Errorf("getShell was incorrect. Got %t, want %t", res, true)
	}

	e.Service.Command = ""

	res = e.getShell()

	if res {
		t.Errorf("getShell (with empty command) was incorrect. Got %t, want %t", res, false)
	}

	e.Service.Shell = false

	res = e.getShell()

	if res {
		t.Errorf("getShell was incorrect. Got %t, want %t", res, false)
	}
}

func TestGetCommand(t *testing.T) {
	var e API
	e.Service.Command = "sleep"

	res := e.getCommand()

	if res != "sleep" {
		t.Errorf("getCommand was incorrect. Got %s, want %s", res, "sleep")
	}
}

func TestGetHostname(t *testing.T) {
	var e API
	e.Service.Hostname = "hostname"

	res := e.getHostname()

	if res != "hostname" {
		t.Errorf("getHostname (with Hostname) was incorrect. Got %s, want %s", res, "hostname")
	}

	e.Service.ContainerName = "containername"

	res = e.getHostname()

	if res != "hostname" {
		t.Errorf("getHostname (with ContainerName and Hostname) was incorrect. Got %s, want %s", res, "hostname")
	}

	e.Service.Hostname = ""
	res = e.getHostname()

	if res != "containername" {
		t.Errorf("getHostname (with ContainerName) was incorrect. Got %s, want %s", res, "containername")
	}
}

func TestGetNetworkMode(t *testing.T) {
	var e API
	e.Service.NetworkMode = "host"

	res := e.getNetworkMode()

	if res != "host" {
		t.Errorf("getNetworkMode was incorrect. Got %s, want %s", res, "host")
	}
}

func TestGetContainerType(t *testing.T) {
	var e API
	e.Service.ContainerType = "DOCKER"

	res := e.getContainerType()

	if res != "docker" {
		t.Errorf("getContainerType (container type docker) was incorrect. Got %s, want %s", res, "docker")
	}

	e.Service.ContainerType = "MESOS"

	res = e.getContainerType()

	if res != "mesos" {
		t.Errorf("getContainerType (container type mesos) was incorrect. Got %s, want %s", res, "mesos")
	}

	e.Service.ContainerType = ""

	res = e.getContainerType()

	if res != "docker" {
		t.Errorf("getContainerType (empty container type) was incorrect. Got %s, want %s", res, "docker")
	}

	e.Service.ContainerType = "nothing"

	res = e.getContainerType()

	if res != "docker" {
		t.Errorf("getContainerType (wrong container type) was incorrect. Got %s, want %s", res, "docker")
	}
}

func TestGetTaskName(t *testing.T) {
	var e API
	e.Config = &cfg.Config{}
	e.Config.PrefixTaskName = "mc"
	e.Service.Mesos.TaskName = "mc:taskname"

	res := e.getTaskName("project", "name")

	if res != "mc:taskname" {
		t.Errorf("getTaskName was incorrect. Got %s, want %s", res, "mc:taskname")
	}
}

func TestGetCPU(t *testing.T) {
	var e API
	e.Service.Deploy.Resources.Limits.CPUs = 0.1

	res := e.getCPU()

	if res != 0.1 {
		t.Errorf("getCPU was incorrect. Got %f, want %f", res, 0.1)
	}
}

func TestGetMemory(t *testing.T) {
	var e API
	e.Service.Deploy.Resources.Limits.Memory = 1000

	res := e.getMemory()

	if res != 1000.0 {
		t.Errorf("getMemory was incorrect. Got %f, want %f", res, 1000.0)
	}
}

func TestGetDisk(t *testing.T) {
	var e API
	e.Config = &cfg.Config{}
	e.Config.Disk = 1000.0

	res := e.getDisk()

	if res != 1000.0 {
		t.Errorf("getDisk was incorrect. Got %f, want %f", res, 1000.0)
	}
}

func TestGetReplicas(t *testing.T) {
	var e API
	e.Service.Deploy.Replicas = "1"

	res := e.getReplicas()

	if res != 1 {
		t.Errorf("getReplicas was incorrect. Got %d, want %d", res, 1)
	}
}

func TestGetTaskNameWithPrefix(t *testing.T) {
	var e API
	e.Config = &cfg.Config{}
	e.Config.PrefixTaskName = "testprefix"

	// Test scenario: Taskname does not match prefix pattern
	e.Service.Mesos.TaskName = "anotherprefix:something"

	res := e.getTaskName("project", "service")

	if res != "testprefix:project:service" {
		t.Errorf("getTaskName was incorrect. Got %s, want testprefix:project:service", res)
	}
}

func TestGetNetworkModeWithDefault(t *testing.T) {
	var e API
	// Test with default network mode when none set
	e.Service.NetworkMode = ""

	res := e.getNetworkMode()

	if res != "user" {
		t.Errorf("getNetworkMode was incorrect. Got %s, want user", res)
	}
}

func TestGetContainerTypeWithDefault(t *testing.T) {
	var e API
	// Test with default container type when none set
	e.Service.ContainerType = ""

	res := e.getContainerType()

	if res != "docker" {
		t.Errorf("getContainerType was incorrect. Got %s, want docker", res)
	}
}
