package scheduler

import (
	"testing"

	mesosclient "github.com/AVENTER-UG/mesos-compose/mesos"
	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
)

func TestResetFrameworkIdentityClearsOnlySessionIdentity(t *testing.T) {
	frameworkID := "framework-1"
	framework := &cfg.FrameworkConfig{
		FrameworkInfo: mesosproto.FrameworkInfo{Id: &mesosproto.FrameworkID{Value: &frameworkID}},
		MesosStreamID: "stream-1",
		FrameworkName: "mc",
		FrameworkRole: "mc",
	}
	scheduler := &Scheduler{
		Framework: framework,
		Mesos:     *mesosclient.New(&cfg.Config{}, framework),
	}

	scheduler.resetFrameworkIdentity()
	if got := framework.FrameworkInfo.GetId().GetValue(); got != "" {
		t.Fatalf("framework ID = %q, want empty", got)
	}
	if got := framework.MesosStreamID; got != "" {
		t.Fatalf("stream ID = %q, want empty", got)
	}
	if framework.FrameworkName != "mc" || framework.FrameworkRole != "mc" {
		t.Fatal("framework configuration was modified")
	}
}
