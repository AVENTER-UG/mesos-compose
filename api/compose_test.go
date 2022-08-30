package api

import (
	"testing"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
)

func TestGetShell(t *testing.T) {
	var e API
	cmd := mesosutil.Command{}
	cmd.Command = "/bin/bash"

	res := e.getShell(cmd)

	if !res {
		t.Errorf("GetShell was incorrect")
	}
}

func TestGetCPU(t *testing.T) {
	var e API
	e.Service.Deploy.Resources.Limits.CPUs = "0.1"

	res := e.getCPU()

	if res != 0.1 {
		t.Errorf("GetCPU was incorrect. Got %f, want %f ", res, 0.1)
	}
}

func TestGetMemory(t *testing.T) {
	var e API
	e.Service.Deploy.Resources.Limits.Memory = "1000"

	res := e.getMemory()

	if res != 1000.0 {
		t.Errorf("GetMemory was incorrect. Got %f, want %f ", res, 1000.0)
	}
}

func TestGetDisk(t *testing.T) {
	var e API
	e.Config = &cfg.Config{}
	e.Config.Disk = 1000.0

	res := e.getDisk()

	if res != 1000.0 {
		t.Errorf("GetDisk was incorrect. Got %f, want %f ", res, 1000.0)
	}
}

func TestGetReplicas(t *testing.T) {
	var e API
	e.Service.Deploy.Replicas = "1"

	res := e.getReplicas()

	if res != 1 {
		t.Errorf("GetDisk was incorrect. Got %d, want %d ", res, 1)
	}
}
