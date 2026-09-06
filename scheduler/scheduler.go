package scheduler

import (
	"context"
	stdjson "encoding/json"
	"errors"
	"io"
	"strings"
	"sync/atomic"
	"time"

	api "github.com/AVENTER-UG/mesos-compose/api"
	"github.com/AVENTER-UG/mesos-compose/mesos"
	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	"github.com/AVENTER-UG/mesos-compose/redis"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/AVENTER-UG/util/vault"
	clusterdclient "github.com/m3scluster/clusterd-go/api/v1/lib/client"

	clusterdscheduler "github.com/m3scluster/clusterd-go/api/v1/lib/scheduler"
	clusterdcalls "github.com/m3scluster/clusterd-go/api/v1/lib/scheduler/calls"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

// Scheduler include all the current vars and global config
type Scheduler struct {
	Config          *cfg.Config
	Framework       *cfg.FrameworkConfig
	Mesos           mesos.Mesos
	API             *api.API
	Vault           *vault.Vault
	Redis           *redis.Redis
	ConnectionError bool
	Subscribed      chan struct{}
}

const streamLivenessTimeout = 2 * time.Minute

// Subscribe to the mesos backend
func Subscribe(cfg *cfg.Config, frm *cfg.FrameworkConfig) *Scheduler {
	e := &Scheduler{
		Config:     cfg,
		Framework:  frm,
		Mesos:      *mesos.New(cfg, frm),
		Subscribed: make(chan struct{}),
	}

	e.Mesos.Subscribe()

	return e
}

// EventLoop is the main loop for the mesos events.
func (e *Scheduler) EventLoop() {
	response, err := e.Mesos.ClusterdClient.Send(
		clusterdclient.RequestSingleton(e.Mesos.ClusterdSubscribe),
		clusterdclient.ResponseClassAuto,
	)
	if err != nil {
		logrus.WithField("func", "scheduler.EventLoop").Error("Mesos subscription failed: ", err)
		return
	}
	defer response.Close()
	watchdogDone := make(chan struct{})
	defer close(watchdogDone)
	lastRecord := time.Now().UnixNano()
	go func() {
		ticker := time.NewTicker(streamLivenessTimeout / 2)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				last := time.Unix(0, atomic.LoadInt64(&lastRecord))
				if time.Since(last) >= streamLivenessTimeout {
					_ = response.Close()
					return
				}
			case <-watchdogDone:
				return
			}
		}
	}()

	for {
		var wireEvent clusterdscheduler.Event
		if err := response.Decode(&wireEvent); err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				logrus.WithField("func", "scheduler.EventLoop").Debug("Mesos Master stream closed")
			} else {
				logrus.WithField("func", "scheduler.EventLoop").Error("Mesos stream decode failed: ", err)
			}
			return
		}
		atomic.StoreInt64(&lastRecord, time.Now().UnixNano())
		if wireEvent.GetType() == clusterdscheduler.Event_SUBSCRIBED {
			wireFrameworkID := wireEvent.GetSubscribed().GetFrameworkID()
			frameworkID := wireFrameworkID.GetValue()
			e.Framework.FrameworkInfo.Id = &mesosproto.FrameworkID{Value: &frameworkID}
			select {
			case <-e.Subscribed:
			default:
				close(e.Subscribed)
			}
			continue
		}
		if wireEvent.GetType() == clusterdscheduler.Event_HEARTBEAT {
			continue
		}
		if wireEvent.GetType() == clusterdscheduler.Event_ERROR {
			message := wireEvent.GetError().GetMessage()
			logrus.WithField("func", "scheduler.EventLoop").Error("Mesos scheduler error: ", message)
			if strings.Contains(strings.ToLower(message), "framework failed over") {
				e.resetFrameworkIdentity()
			}
			return
		}
		data, err := stdjson.Marshal(&wireEvent)
		if err != nil {
			logrus.WithField("func", "scheduler.EventLoop").Error("Mesos event conversion failed: ", err)
			return
		}
		var event mesosproto.Event
		if err := protojson.Unmarshal(data, &event); err != nil {
			logrus.WithField("func", "scheduler.EventLoop").Warn("Mesos event conversion failed: ", err)
			continue
		}
		e.dispatchDomainEvent(&event)
		if event.GetType() == mesosproto.Event_UPDATE {
			go e.callPluginEvent(&event)
		}
	}
}

func (e *Scheduler) resetFrameworkIdentity() {
	emptyID := ""
	e.Framework.FrameworkInfo.Id = &mesosproto.FrameworkID{Value: &emptyID}
	e.Mesos.SetStreamID("")
	if e.Redis != nil {
		e.Redis.SaveFrameworkRedis(e.Framework)
	}
}

func (e *Scheduler) dispatchDomainEvent(event *mesosproto.Event) {
	logrus.WithField("func", "scheduler.EventLoop").Tracef("Event %s", event.GetType().String())
	switch event.Type.Number() {
	case mesosproto.Event_UPDATE.Number():
		e.HandleUpdate(event)
	case mesosproto.Event_OFFERS.Number():
		if err := e.HandleOffers(event.Offers); err != nil {
			logrus.WithField("func", "scheduler.dispatchDomainEvent").Warn("HandleOffers: ", err)
		}
	}
}

func (e *Scheduler) changeDockerPorts(cmd *cfg.Command) []*mesosproto.ContainerInfo_DockerInfo_PortMapping {
	var ret []*mesosproto.ContainerInfo_DockerInfo_PortMapping
	for _, port := range cmd.DockerPortMappings {
		port.HostPort = e.API.GetRandomHostPort()
		ret = append(ret, port)
	}
	return ret
}

func (e *Scheduler) changeDiscoveryInfo(cmd *cfg.Command) *mesosproto.DiscoveryInfo {
	for i, port := range cmd.DockerPortMappings {
		cmd.Discovery.Ports.Ports[i].Number = port.HostPort
	}
	return cmd.Discovery
}

// reconcile will ask Mesos about the current state of the given tasks
func (e *Scheduler) reconcile() {
	tasks := make(map[string]string)
	keys := e.Redis.GetAllRedisKeys(e.Framework.FrameworkName + ":*")
	for keys.Next(e.Redis.CTX) {
		// continue if the key is not a mesos task
		if e.Redis.CheckIfNotTask(keys) {
			continue
		}

		keys.Val()

		key := e.Redis.GetRedisKey(keys.Val())

		task := redis.DecodeTaskOrEmpty([]byte(key))

		if task.TaskID == "" || task.Agent == "" || task.State == "__NEW" || task.State == "__KILL" || task.State == "" {
			continue
		}

		tasks[task.TaskID] = task.MesosAgent.ID
		logrus.WithField("func", "mesos.Reconcile").Debug("Reconcile Task: ", task.TaskID)
	}
	err := e.Mesos.CallClusterd(clusterdcalls.Reconcile(clusterdcalls.ReconcileTasks(tasks)))

	if err != nil {
		logrus.WithField("func", "scheduler.reconcile").Debug("Reconcile Error: ", err)
	}
}

// implicitReconcile will ask Mesos which tasks and there state are registert to this framework
func (e *Scheduler) implicitReconcile() {
	err := e.Mesos.CallClusterd(clusterdcalls.Reconcile())

	if err != nil {
		logrus.WithField("func", "scheduler.implicitReconcile").Debug("Reconcile Error: ", err)
	}
}

func (e *Scheduler) callPluginEvent(event *mesosproto.Event) {
	if e.Config.PluginsEnable {
		for _, p := range e.Config.Plugins {
			symbol, err := p.Lookup("Event")
			if err != nil {
				logrus.WithField("func", "scheduler.callPluginEvent").Error("Error lookup event function in plugin: ", err.Error())
				continue
			}

			eventPluginFunc, ok := symbol.(func(*mesosproto.Event))
			if !ok {
				logrus.WithField("func", "main.initPlugins").Error("Error plugin does not have init function")
				continue
			}

			eventPluginFunc(event)
		}
	}
}
