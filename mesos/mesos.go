package mesos

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	api "github.com/AVENTER-UG/mesos-compose/api"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	"github.com/sirupsen/logrus"

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
func SetConfig(cfg *cfg.Config, frm *mesosutil.FrameworkConfig) {
	config = cfg
	framework = frm
}

// Subscribe to the mesos backend
func Subscribe() error {
	subscribeCall := &mesosproto.Call{
		FrameworkID: framework.FrameworkInfo.ID,
		Type:        mesosproto.Call_SUBSCRIBE,
		Subscribe: &mesosproto.Call_Subscribe{
			FrameworkInfo: &framework.FrameworkInfo,
		},
	}
	logrus.Debug(subscribeCall)
	body, _ := marshaller.MarshalToString(subscribeCall)
	logrus.Debug(body)
	client := &http.Client{}
	client.Transport = &http.Transport{
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipSSL},
	}

	protocol := "https"
	if !framework.MesosSSL {
		protocol = "http"
	}
	req, _ := http.NewRequest("POST", protocol+"://"+framework.MesosMasterServer+"/api/v1/scheduler", bytes.NewBuffer([]byte(body)))
	req.Close = true
	req.SetBasicAuth(framework.Username, framework.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logrus.Fatal(err)
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)

	line, _ := reader.ReadString('\n')
	bytesCount, _ := strconv.Atoi(strings.Trim(line, "\n"))

	for {
		// Read line from Mesos
		line, _ = reader.ReadString('\n')
		line = strings.Trim(line, "\n")
		// Read important data
		data := line[:bytesCount]
		// Rest data will be bytes of next message
		bytesCount, _ = strconv.Atoi((line[bytesCount:]))
		var event mesosproto.Event // Event as ProtoBuf
		err := jsonpb.UnmarshalString(data, &event)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Debug("Subscribe Got: ", event.GetType())

		switch event.Type {
		case mesosproto.Event_SUBSCRIBED:
			logrus.Debug(event)
			logrus.Info("Subscribed")
			logrus.Info("FrameworkId: ", event.Subscribed.GetFrameworkID())
			framework.FrameworkInfo.ID = event.Subscribed.GetFrameworkID()
			framework.MesosStreamID = res.Header.Get("Mesos-Stream-Id")

			// store framework configuration
			d, _ := json.Marshal(&framework)
			err = config.RedisClient.Set(config.RedisCTX, "framework", d, 0).Err()
			if err != nil {
				logrus.Error("Framework save config and state into redis Error: ", err)
			}
			Reconcile()
		case mesosproto.Event_UPDATE:
			logrus.Debug("Update", HandleUpdate(&event))
			err = api.SaveConfig()
			if err != nil {
				api.ErrorMessage(1, "Event_UPDATE", "Could not save config data")
			}
		case mesosproto.Event_HEARTBEAT:
			Heartbeat()
		case mesosproto.Event_OFFERS:
			// Search Failed containers and restart them
			logrus.Debug("Offer Got")
			err = HandleOffers(event.Offers)
			if err != nil {
				logrus.Error("Switch Event HandleOffers: ", err)
			}
		default:
			logrus.Debug("DEFAULT EVENT: ", event.Offers)
		}
	}
}

// Reconcile will reconcile the task states after the framework was restarted
func Reconcile() {
	logrus.Info("Reconcile Tasks")
	var oldTasks []mesosproto.Call_Reconcile_Task
	keys := api.GetAllRedisKeys("*")
	for keys.Next(config.RedisCTX) {
		key := api.GetRedisKey(keys.Val())

		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)

		if task.TaskID == "" {
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
		api.ErrorMessage(3, "Reconcile_Error", err.Error())
		logrus.Debug("Reconcile Error: ", err)
	}
}
