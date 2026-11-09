package redis

import (
	"testing"

	mesosproto "github.com/m3scluster/mesos-compose/proto"
	cfg "github.com/m3scluster/mesos-compose/types"
)

func TestFrameworkPersistenceRoundTrip(t *testing.T) {
	frameworkID := "framework-1"
	framework := &cfg.FrameworkConfig{
		FrameworkName:     "mc",
		FrameworkRole:     "mc",
		MesosStreamID:     "stream-1",
		MesosMasterServer: "master.example:5050",
		FrameworkInfo: mesosproto.FrameworkInfo{
			Id: &mesosproto.FrameworkID{Value: &frameworkID},
		},
	}

	data, err := EncodeFramework(framework)
	if err != nil {
		t.Fatalf("encode framework: %v", err)
	}
	decoded, err := DecodeFramework(data)
	if err != nil {
		t.Fatalf("decode framework: %v", err)
	}
	if decoded.FrameworkInfo.GetId().GetValue() != frameworkID || decoded.MesosStreamID != framework.MesosStreamID {
		t.Fatalf("framework round trip = %#v", decoded)
	}
}

func TestTaskPersistenceRoundTrip(t *testing.T) {
	task := &cfg.Command{
		TaskID:         "task-1",
		TaskName:       "web",
		ContainerImage: "busybox:latest",
		NetworkMode:    "host",
		Restart:        "unless-stopped",
		CPU:            1,
		Memory:         64,
		State:          "TASK_RUNNING",
	}

	data, err := EncodeTask(task)
	if err != nil {
		t.Fatalf("encode task: %v", err)
	}
	decoded, err := DecodeTask(data)
	if err != nil {
		t.Fatalf("decode task: %v", err)
	}
	if decoded.TaskID != task.TaskID || decoded.TaskName != task.TaskName || decoded.ContainerImage != task.ContainerImage || decoded.NetworkMode != task.NetworkMode || decoded.State != task.State {
		t.Fatalf("task round trip = %#v", decoded)
	}
}
