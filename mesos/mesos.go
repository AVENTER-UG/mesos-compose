package mesos

import (
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"

	"github.com/gogo/protobuf/jsonpb"
)

// Service include all the current vars and global config
var config *cfg.Config
var framework *mesosutil.FrameworkConfig

// Marshaler to serialize Protobuf Message to JSON
var marshaller = jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      " ",
	OrigName:    true,
}

// SetConfig set the global config
func SetConfig(cfg *cfg.Config) {
	config = cfg
}

// Restart failed container
func RestartFailedContainer() {
	if framework.State != nil {
		for _, element := range framework.State {
			if element.Status != nil {
				switch *element.Status.State {
				case mesosproto.TASK_FAILED, mesosproto.TASK_ERROR:
					mesosutil.DeleteOldTask(element.Status.TaskID)
				case mesosproto.TASK_KILLED:
					mesosutil.DeleteOldTask(element.Status.TaskID)
				case mesosproto.TASK_LOST:
					mesosutil.DeleteOldTask(element.Status.TaskID)
				}
			}
		}
	}
}

// Heartbeat
func Heartbeat() {
}
