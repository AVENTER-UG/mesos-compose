package mesos

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"

	api "github.com/AVENTER-UG/mesos-compose/api"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	"github.com/sirupsen/logrus"

	"github.com/gogo/protobuf/jsonpb"
)

// Scheduler include all the current vars and global config
type Scheduler struct {
	Config    *cfg.Config
	Framework *mesosutil.FrameworkConfig
	Client    *http.Client
	Req       *http.Request
	API       *api.API
}

// Marshaler to serialize Protobuf Message to JSON
var marshaller = jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      " ",
	OrigName:    true,
}

// Subscribe to the mesos backend
func Subscribe(cfg *cfg.Config, frm *mesosutil.FrameworkConfig) *Scheduler {
	e := &Scheduler{
		Config:    cfg,
		Framework: frm,
	}

	subscribeCall := &mesosproto.Call{
		FrameworkID: e.Framework.FrameworkInfo.ID,
		Type:        mesosproto.Call_SUBSCRIBE,
		Subscribe: &mesosproto.Call_Subscribe{
			FrameworkInfo: &e.Framework.FrameworkInfo,
		},
	}
	logrus.Debug(subscribeCall)
	body, _ := marshaller.MarshalToString(subscribeCall)
	logrus.Debug(body)
	client := &http.Client{}
	client.Transport = &http.Transport{
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: e.Config.SkipSSL},
	}

	protocol := "https"
	if !e.Framework.MesosSSL {
		protocol = "http"
	}
	req, _ := http.NewRequest("POST", protocol+"://"+e.Framework.MesosMasterServer+"/api/v1/scheduler", bytes.NewBuffer([]byte(body)))
	req.Close = true
	req.SetBasicAuth(e.Framework.Username, e.Framework.Password)
	req.Header.Set("Content-Type", "application/json")

	e.Req = req
	e.Client = client

	return e
}

// EventLoop is the main loop for the mesos events.
func (e *Scheduler) EventLoop() {
	res, err := e.Client.Do(e.Req)

	if err != nil {
		logrus.Fatal(err)
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)

	line, _ := reader.ReadString('\n')
	_ = strings.TrimSuffix(line, "\n")

	logrus.Debug(line)

	go e.HeartbeatLoop()

	for {
		// Read line from Mesos
		line, _ = reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		// Read important data
		var event mesosproto.Event // Event as ProtoBuf
		err := jsonpb.UnmarshalString(line, &event)
		if err != nil {
			logrus.Error("Mesos Unmarshal Data Error: ", err)
		}

		logrus.Debug(event.GetType())

		switch event.Type {
		case mesosproto.Event_SUBSCRIBED:
			logrus.Debug(event)
			logrus.Info("Subscribed")
			logrus.Info("FrameworkId: ", event.Subscribed.GetFrameworkID())
			e.Framework.FrameworkInfo.ID = event.Subscribed.GetFrameworkID()
			e.Framework.MesosStreamID = res.Header.Get("Mesos-Stream-Id")

			// store framework configuration
			d, _ := json.Marshal(&e.Framework)
			err = e.API.Redis.RedisClient.Set(e.API.Redis.RedisCTX, e.Framework.FrameworkName+":framework", d, 0).Err()
			if err != nil {
				logrus.Error("Framework save config and state into redis Error: ", err)
			}
			e.Reconcile()
			e.API.SaveConfig()
		case mesosproto.Event_UPDATE:
			e.HandleUpdate(&event)
			e.API.SaveConfig()
		case mesosproto.Event_OFFERS:
			// Search Failed containers and restart them
			logrus.Debug("Offer Got")
			err = e.HandleOffers(event.Offers)
			if err != nil {
				logrus.Error("Switch Event HandleOffers: ", err)
			}
		}
	}
}

// Reconcile will reconcile the task states after the framework was restarted
func (e *Scheduler) Reconcile() {
	logrus.Info("Reconcile Tasks")
	var oldTasks []mesosproto.Call_Reconcile_Task
	keys := e.API.GetAllRedisKeys(e.Framework.FrameworkName + ":*")
	for keys.Next(e.API.Redis.RedisCTX) {
		key := e.API.GetRedisKey(keys.Val())

		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)

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
		logrus.Debug("Reconcile Task: ", task.TaskID)
	}
	err := mesosutil.Call(&mesosproto.Call{
		Type:      mesosproto.Call_RECONCILE,
		Reconcile: &mesosproto.Call_Reconcile{Tasks: oldTasks},
	})

	if err != nil {
		e.API.ErrorMessage(3, "Reconcile_Error", err.Error())
		logrus.Debug("Reconcile Error: ", err)
	}
}
