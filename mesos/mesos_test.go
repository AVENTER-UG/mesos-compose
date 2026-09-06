package mesos

import (
	"encoding/json"
	"testing"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
)

func TestDecodeTask(t *testing.T) {
	e := &Mesos{}
	decoded := e.DecodeTask(`{"TaskID":"task-1","State":"TASK_RUNNING","CPU":1.5}`)
	if decoded.TaskID != "task-1" || decoded.State != "TASK_RUNNING" || decoded.CPU != 1.5 {
		t.Fatalf("decoded task = %#v", decoded)
	}

	if got := e.DecodeTask("not-json"); got.TaskID != "" {
		t.Fatalf("invalid JSON returned task %#v, want empty task", got)
	}
	if got := e.DecodeTask(""); got.TaskID != "" {
		t.Fatalf("empty input returned task %#v, want empty task", got)
	}
}

func TestDecodeTaskWithCSIVolume(t *testing.T) {
	fsType := "cifs"
	pluginName := "org.apache.mesos.csi.smb"
	volumeID := "mvs-csi-proof"
	accessMode := mesosproto.Volume_Source_CSIVolume_VolumeCapability_AccessMode_SINGLE_NODE_WRITER
	task := &cfg.Command{
		TaskID: "csi-task",
		Volumes: []*mesosproto.Volume{{
			ContainerPath: stringPointer("/mnt/mvs"),
			Mode:          mesosproto.Volume_RW.Enum(),
			Source: &mesosproto.Volume_Source{
				Type: mesosproto.Volume_Source_CSI_VOLUME.Enum(),
				CsiVolume: &mesosproto.Volume_Source_CSIVolume{
					PluginName: &pluginName,
					StaticProvisioning: &mesosproto.Volume_Source_CSIVolume_StaticProvisioning{
						VolumeId: &volumeID,
						VolumeCapability: &mesosproto.Volume_Source_CSIVolume_VolumeCapability{
							AccessType: &mesosproto.Volume_Source_CSIVolume_VolumeCapability_Mount{
								Mount: &mesosproto.Volume_Source_CSIVolume_VolumeCapability_MountVolume{
									FsType:     &fsType,
									MountFlags: []string{"vers=3.0"},
								},
							},
							AccessMode: &mesosproto.Volume_Source_CSIVolume_VolumeCapability_AccessMode{Mode: &accessMode},
						},
						VolumeContext: map[string]string{"source": "//192.168.150.82/mvs"},
					},
				},
			},
		}},
	}
	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("marshal task: %v", err)
	}

	decoded := (&Mesos{}).DecodeTask(string(data))
	if len(decoded.Volumes) != 1 {
		t.Fatalf("decoded volumes = %d, want 1", len(decoded.Volumes))
	}
	csi := decoded.Volumes[0].GetSource().GetCsiVolume()
	if csi.GetPluginName() != pluginName || csi.GetStaticProvisioning().GetVolumeId() != volumeID {
		t.Fatalf("decoded CSI volume = %v", csi)
	}
	capability := csi.GetStaticProvisioning().GetVolumeCapability()
	if capability.GetMount().GetFsType() != fsType || capability.GetAccessMode().GetMode() != accessMode {
		t.Fatalf("decoded capability = %v", capability)
	}
}

func TestIsRessourceMatched(t *testing.T) {
	nameCPU, nameMem := "cpus", "mem"
	cpu, mem := 2.0, 128.0
	resources := []*mesosproto.Resource{
		{Name: &nameCPU, Scalar: &mesosproto.Value_Scalar{Value: &cpu}},
		{Name: &nameMem, Scalar: &mesosproto.Value_Scalar{Value: &mem}},
	}
	e := &Mesos{}
	if !e.IsRessourceMatched(resources, &cfg.Command{CPU: 1, Memory: 64}) {
		t.Fatal("matching resources reported as insufficient")
	}
	if e.IsRessourceMatched(resources, &cfg.Command{CPU: 3, Memory: 64}) {
		t.Fatal("insufficient CPU reported as sufficient")
	}
}

func TestDeclineOffer(t *testing.T) {
	offerID := &mesosproto.OfferID{Value: stringPointer("offer-1")}
	call := (&Mesos{}).DeclineOffer([]*mesosproto.OfferID{offerID})

	if call.GetType() != mesosproto.Call_DECLINE {
		t.Fatalf("call type = %v, want DECLINE", call.GetType())
	}
	if len(call.GetDecline().GetOfferIds()) != 1 || call.GetDecline().GetOfferIds()[0].GetValue() != "offer-1" {
		t.Fatalf("declined offers = %#v", call.GetDecline().GetOfferIds())
	}
	if got := call.GetDecline().GetFilters().GetRefuseSeconds(); got != 120 {
		t.Fatalf("refuse seconds = %v, want 120", got)
	}
}

func stringPointer(value string) *string { return &value }
