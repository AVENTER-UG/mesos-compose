package scheduler

import (
	"testing"

	mesosproto "github.com/m3scluster/mesos-compose/proto"
	cfg "github.com/m3scluster/mesos-compose/types"
)

func TestDefaultResourcesUsesConfiguredDiskFloor(t *testing.T) {
	cpu, memory, disk, gpus := 1.0, 64.0, 10.0, 0.0
	s := &Scheduler{Config: &cfg.Config{Disk: 100.0}}
	resources := s.defaultResources(&cfg.Command{CPU: cpu, Memory: memory, Disk: disk, GPUs: gpus})

	if len(resources) != 4 {
		t.Fatalf("resource count = %d, want 4", len(resources))
	}
	if got := resources[0].GetScalar().GetValue(); got != cpu {
		t.Fatalf("CPU = %v, want %v", got, cpu)
	}
	if got := resources[2].GetScalar().GetValue(); got != 100 {
		t.Fatalf("disk = %v, want configured floor 100", got)
	}
}

func TestPrepareTaskInfoExecuteContainer(t *testing.T) {
	taskID, taskName, image, command := "task-1", "service-1", "busybox:latest", "echo hello"
	agentID := "agent-1"
	s := &Scheduler{Config: &cfg.Config{Disk: 1}}
	cmd := &cfg.Command{
		TaskID: taskID, TaskName: taskName, ContainerType: "docker",
		ContainerImage: image, Command: command, CPU: 1, Memory: 32,
		NetworkMode: "host", PullPolicy: "missing",
	}

	result := s.PrepareTaskInfoExecuteContainer(&mesosproto.AgentID{Value: &agentID}, cmd)
	if len(result) != 1 {
		t.Fatalf("task count = %d, want 1", len(result))
	}
	task := result[0]
	if task.GetName() != taskName || task.GetTaskId().GetValue() != taskID || task.GetAgentId().GetValue() != agentID {
		t.Fatalf("task identity = name %q, task ID %q, agent %q", task.GetName(), task.GetTaskId().GetValue(), task.GetAgentId().GetValue())
	}
	if task.GetCommand().GetValue() != command || task.GetContainer().GetDocker().GetImage() != image {
		t.Fatalf("task command/container not mapped: %#v", task)
	}
	if task.GetContainer().GetDocker().GetNetwork() != mesosproto.ContainerInfo_DockerInfo_HOST {
		t.Fatalf("network = %v, want HOST", task.GetContainer().GetDocker().GetNetwork())
	}
	if task.GetContainer().GetDocker().GetForcePullImage() {
		t.Fatal("missing pull policy unexpectedly enabled force pull")
	}
}
