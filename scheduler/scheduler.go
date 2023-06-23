package scheduler

import (
	"bufio"
	"net/http"
	"strings"

	api "github.com/AVENTER-UG/mesos-compose/api"
	"github.com/AVENTER-UG/mesos-compose/mesos"
	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	"github.com/AVENTER-UG/mesos-compose/redis"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/AVENTER-UG/util/vault"
	"github.com/sirupsen/logrus"

	"github.com/gogo/protobuf/jsonpb"
)

// Scheduler include all the current vars and global config
type Scheduler struct {
	Config          *cfg.Config
	Framework       *cfg.FrameworkConfig
	Mesos           mesos.Mesos
	Client          *http.Client
	Req             *http.Request
	API             *api.API
	Vault           *vault.Vault
	Redis           *redis.Redis
	ConnectionError bool
}

// New will create a new Scheduler object
func New(cfg *cfg.Config, frm *cfg.FrameworkConfig) *Scheduler {
	e := &Scheduler{
		Config:    cfg,
		Framework: frm,
		Mesos:     *mesos.New(cfg, frm),
		Redis:     redis.New(cfg, frm),
	}

	e.Client, e.Req = e.Mesos.Subscribe()

	return e
}

// EventLoop is the main loop for the mesos events.
func (e *Scheduler) EventLoop() {
	res, err := e.Client.Do(e.Req)

	if err != nil {
		logrus.WithField("func", "scheduler.EventLoop").Error("Mesos Master connection error: ", err.Error())
		return
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)

	line, _ := reader.ReadString('\n')
	_ = strings.TrimSuffix(line, "\n")

	go e.HeartbeatLoop()
	go e.ReconcileLoop()

	for {
		// Read line from Mesos
		line, err = reader.ReadString('\n')
		if err != nil {
			logrus.WithField("func", "scheduler.EventLoop").Error("Error to read data from Mesos Master: ", err.Error())
			return
		}
		line = strings.TrimSuffix(line, "\n")
		// Read important data
		var event mesosproto.Event // Event as ProtoBuf
		err := jsonpb.UnmarshalString(line, &event)
		if err != nil {
			logrus.WithField("func", "scheduler.EventLoop").Warn("Could not unmarshal Mesos Master data: ", err.Error())
			continue
		}

		switch event.Type {
		case mesosproto.Event_SUBSCRIBED:
			logrus.WithField("func", "scheduler.EventLoop").Info("Subscribed")
			logrus.WithField("func", "scheduler.EventLoop").Debug("FrameworkId: ", event.Subscribed.GetFrameworkID())
			e.Framework.FrameworkInfo.ID = event.Subscribed.GetFrameworkID()
			e.Framework.MesosStreamID = res.Header.Get("Mesos-Stream-Id")

			e.reconcile()
			e.Redis.SaveFrameworkRedis(*e.Framework)
			e.Redis.SaveConfig(*e.Config)
		case mesosproto.Event_UPDATE:
			e.HandleUpdate(&event)
			e.Redis.SaveConfig(*e.Config)
		case mesosproto.Event_OFFERS:
			// Search Failed containers and restart them
			err = e.HandleOffers(event.Offers)
			if err != nil {
				logrus.WithField("func", "scheduler.EventLoop").Warn("Switch Event HandleOffers: ", err)
			}
		}
	}
}

func (e *Scheduler) changeDockerPorts(cmd cfg.Command) []mesosproto.ContainerInfo_DockerInfo_PortMapping {
	var ret []mesosproto.ContainerInfo_DockerInfo_PortMapping
	for _, port := range cmd.DockerPortMappings {
		port.HostPort = e.API.GetRandomHostPort()
		ret = append(ret, port)
	}
	return ret
}

func (e *Scheduler) changeDiscoveryInfo(cmd cfg.Command) mesosproto.DiscoveryInfo {
	for i, port := range cmd.DockerPortMappings {
		cmd.Discovery.Ports.Ports[i].Number = port.HostPort
	}
	return cmd.Discovery
}

// Reconcile will reconcile the task states after the framework was restarted
func (e *Scheduler) reconcile() {
	var oldTasks []mesosproto.Call_Reconcile_Task
	keys := e.Redis.GetAllRedisKeys(e.Framework.FrameworkName + ":*")
	for keys.Next(e.Redis.CTX) {
		// continue if the key is not a mesos task
		if e.Redis.CheckIfNotTask(keys) {
			continue
		}

		key := e.Redis.GetRedisKey(keys.Val())

		task := e.Mesos.DecodeTask(key)

		if task.TaskID == "" || task.Agent == "" {
			continue
		}

		oldTasks = append(oldTasks, mesosproto.Call_Reconcile_Task{
			TaskID: mesosproto.TaskID{
				Value: task.TaskID,
			},
			AgentID: &mesosproto.AgentID{
				Value: task.Agent,
			},
		})
		logrus.WithField("func", "mesos.Reconcile").Debug("Reconcile Task: ", task.TaskID)
	}
	err := e.Mesos.Call(&mesosproto.Call{
		Type:      mesosproto.Call_RECONCILE,
		Reconcile: &mesosproto.Call_Reconcile{Tasks: oldTasks},
	})

	if err != nil {
		logrus.WithField("func", "mesos.Reconcile").Debug("Reconcile Error: ", err)
	}
}

// implicitReconcile will ask Mesos which tasks and there state are registert to this framework
func (e *Scheduler) implicitReconcile() {
	var noTasks []mesosproto.Call_Reconcile_Task
	err := e.Mesos.Call(&mesosproto.Call{
		Type:      mesosproto.Call_RECONCILE,
		Reconcile: &mesosproto.Call_Reconcile{Tasks: noTasks},
	})

	if err != nil {
		logrus.WithField("func", "scheduler.implicitReconcile").Debug("Reconcile Error: ", err)
	}
}
