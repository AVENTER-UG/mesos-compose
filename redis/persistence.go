package redis

import (
	"encoding/json"
	"fmt"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/sirupsen/logrus"
)

// EncodeTask preserves the existing Redis JSON representation while keeping
// persistence serialization behind one stable boundary.
func EncodeTask(task *cfg.Command) ([]byte, error) {
	if task == nil {
		return nil, fmt.Errorf("task must not be nil")
	}
	return json.Marshal(task)
}

// DecodeTask preserves legacy task JSON handling, including Command's custom
// CSI volume unmarshaling.
func DecodeTask(data []byte) (*cfg.Command, error) {
	var task cfg.Command
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// DecodeTaskOrEmpty preserves the legacy best-effort task lookup behavior.
func DecodeTaskOrEmpty(data []byte) *cfg.Command {
	task, err := DecodeTask(data)
	if err != nil {
		logrus.WithField("func", "redis.DecodeTaskOrEmpty").Error("Could not decode task: ", err)
		return &cfg.Command{}
	}
	return task
}

// EncodeFramework preserves the existing framework JSON representation.
func EncodeFramework(framework *cfg.FrameworkConfig) ([]byte, error) {
	if framework == nil {
		return nil, fmt.Errorf("framework must not be nil")
	}
	return json.Marshal(framework)
}

// DecodeFramework decodes the legacy framework record without exposing
// clusterd-go wire types to the persistence layer.
func DecodeFramework(data []byte) (*cfg.FrameworkConfig, error) {
	var framework cfg.FrameworkConfig
	if err := json.Unmarshal(data, &framework); err != nil {
		return nil, err
	}
	return &framework, nil
}
