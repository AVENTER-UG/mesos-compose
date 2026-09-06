package mesos

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"time"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	clusterd "github.com/m3scluster/clusterd-go/api/v1/lib"
	"github.com/m3scluster/clusterd-go/api/v1/lib/encoding/codecs"
	"github.com/m3scluster/clusterd-go/api/v1/lib/httpcli"
	clusterdscheduler "github.com/m3scluster/clusterd-go/api/v1/lib/scheduler"
	clusterdcalls "github.com/m3scluster/clusterd-go/api/v1/lib/scheduler/calls"
	"github.com/sirupsen/logrus"
)

// Mesos include all the current vars and global config
type Mesos struct {
	Config     *cfg.Config
	Framework  *cfg.FrameworkConfig
	IsSuppress bool
	IsRevive   bool
	CountAgent int

	// Client and Transport are created once in New and reused for the whole
	// lifetime of the framework, so non-stream Calls can keep connections
	// alive and Subscribe applies the same TLS policy as everything else.
	Client            *http.Client
	Transport         *http.Transport
	ClusterdClient    *httpcli.Client
	ClusterdSubscribe *clusterdscheduler.Call
	callTimeout       time.Duration
	callMu            sync.Mutex
	streamIDMu        sync.RWMutex

	agentCacheMu sync.Mutex
	agentCache   map[string]agentCacheEntry
}

const agentInfoCacheTTL = 30 * time.Second

func mesosEndpoint(framework *cfg.FrameworkConfig, path string) string {
	protocol := "https"
	if !framework.MesosSSL {
		protocol = "http"
	}
	return protocol + "://" + framework.MesosMasterServer + path
}

type agentCacheEntry struct {
	agent   cfg.MesosSlaves
	expires time.Time
}

type streamIDRoundTripper struct {
	next  http.RoundTripper
	owner *Mesos
}

func (rt *streamIDRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	response, err := rt.next.RoundTrip(req)
	if err == nil {
		if streamID := response.Header.Get("Mesos-Stream-Id"); streamID != "" {
			rt.owner.SetStreamID(streamID)
		}
	}
	return response, err
}

func clusterdSchedulerCall(call *mesosproto.Call) (clusterdscheduler.Call, error) {
	if call == nil {
		return clusterdscheduler.Call{}, fmt.Errorf("scheduler call must not be nil")
	}
	data, err := json.Marshal(call)
	if err != nil {
		return clusterdscheduler.Call{}, fmt.Errorf("marshal scheduler call: %w", err)
	}
	var result clusterdscheduler.Call
	if err := json.Unmarshal(data, &result); err != nil {
		return clusterdscheduler.Call{}, fmt.Errorf("unmarshal clusterd scheduler call: %w", err)
	}
	return result, nil
}

func (e *Mesos) SetStreamID(streamID string) {
	e.streamIDMu.Lock()
	e.Framework.MesosStreamID = streamID
	e.streamIDMu.Unlock()
}

func (e *Mesos) StreamID() string {
	e.streamIDMu.RLock()
	defer e.streamIDMu.RUnlock()
	return e.Framework.MesosStreamID
}

// New will create a new API object
func New(cfg *cfg.Config, frm *cfg.FrameworkConfig) *Mesos {
	e := &Mesos{
		Config:     cfg,
		Framework:  frm,
		IsSuppress: false,
		IsRevive:   false,
	}

	// One Transport (and one Client on top of it) is created for the whole
	// lifetime of the framework and reused by Subscribe and every Call.
	// Idle connections stay open, so non-stream Calls can be served over
	// keep-alive connections instead of one TLS handshake per request.
	e.Transport = &http.Transport{
		// #nosec G402
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: cfg.SkipSSL},
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
	}
	e.Client = &http.Client{
		Transport: e.Transport,
	}
	e.callTimeout = 15 * time.Second
	e.ClusterdClient = httpcli.New(
		httpcli.Endpoint(mesosEndpoint(frm, "/api/v1/scheduler")),
		httpcli.Codec(codecs.ByMediaType[codecs.MediaTypeJSON]),
		httpcli.Do(httpcli.With(
			httpcli.BasicAuth(frm.Username, frm.Password),
			httpcli.RoundTripper(e.Transport),
			httpcli.WrapRoundTripper(func(next http.RoundTripper) http.RoundTripper {
				return &streamIDRoundTripper{next: next, owner: e}
			}),
			// #nosec G402 -- SkipSSL is an explicit framework configuration option.
			httpcli.TLSConfig(&tls.Config{InsecureSkipVerify: cfg.SkipSSL}),
		)),
	)
	e.agentCache = make(map[string]agentCacheEntry)

	return e
}

func (e *Mesos) Subscribe() {
	subscribeCall := &mesosproto.Call{
		FrameworkId: e.Framework.FrameworkInfo.Id,
		Type:        mesosproto.Call_SUBSCRIBE.Enum(),
		Subscribe: &mesosproto.Call_Subscribe{
			FrameworkInfo: &e.Framework.FrameworkInfo,
		},
	}

	if len(e.Config.HostConstraintsList) > 0 {

		offerConstraintGroups := []*mesosproto.OfferConstraints_RoleConstraints_Group{}
		for _, hostname := range e.Config.HostConstraintsList {
			offerConstraint := mesosproto.OfferConstraints_RoleConstraints_Group{
				AttributeConstraints: []*mesosproto.AttributeConstraint{
					{
						Selector: &mesosproto.AttributeConstraint_Selector{
							Selector: &mesosproto.AttributeConstraint_Selector_PseudoattributeType_{
								PseudoattributeType: mesosproto.AttributeConstraint_Selector_HOSTNAME,
							},
						},
						Predicate: &mesosproto.AttributeConstraint_Predicate{
							Predicate: &mesosproto.AttributeConstraint_Predicate_TextEquals_{
								TextEquals: &mesosproto.AttributeConstraint_Predicate_TextEquals{
									Value: &hostname,
								},
							},
						},
					},
				},
			}
			offerConstraintGroups = append(offerConstraintGroups, &offerConstraint)
		}

		offerConstraints := &mesosproto.OfferConstraints{
			RoleConstraints: map[string]*mesosproto.OfferConstraints_RoleConstraints{
				e.Framework.FrameworkRole: {
					Groups: offerConstraintGroups,
				},
			},
		}

		subscribeCall.Subscribe.OfferConstraints = offerConstraints
	}

	logrus.WithField("func", "scheduler.Subscribe").Debug(subscribeCall)
	if wireCall, err := clusterdSchedulerCall(subscribeCall); err == nil {
		e.ClusterdSubscribe = &wireCall
	} else {
		logrus.WithField("func", "scheduler.Subscribe").Error("Could not create clusterd subscription call: ", err)
	}
}

// Revive will revive the mesos tasks to clean up
func (e *Mesos) Revive() {
	if !e.IsRevive {
		logrus.WithField("func", "mesos.Revive").Info("Framework Revive")
		e.IsSuppress = false
		e.IsRevive = true
		err := e.CallClusterd(clusterdcalls.Revive())
		if err != nil {
			logrus.WithField("func", "mesos.Revive").Error("Call Revive: ", err)
		}
	}
}

// ForceSuppressFramework if all Tasks are running, suppress framework offers
func (e *Mesos) ForceSuppressFramework() {
	logrus.WithField("func", "mesos.ForceSuppressFramework").Info("Framework Suppress")
	e.IsSuppress = false
	e.SuppressFramework()
}

// SuppressFramework if all Tasks are running, suppress framework offers
func (e *Mesos) SuppressFramework() {
	if !e.IsSuppress {
		logrus.WithField("func", "mesos.SuppressFramework").Info("Framework Suppress")
		e.IsSuppress = true
		e.IsRevive = false
		err := e.CallClusterd(clusterdcalls.Suppress())
		if err != nil {
			logrus.WithField("func", "mesos.SupressFramework").Error("Suppress Framework Call: ")
		}
	}
}

// Kill a Task with the given taskID
func (e *Mesos) Kill(taskID string, agentID string) error {
	logrus.WithField("func", "mesos.Kill").Info("Kill task ", taskID)
	// tell mesos to shutdonw the given task
	return e.CallClusterd(clusterdcalls.Kill(taskID, agentID))
}

func (e *Mesos) CallClusterd(message *clusterdscheduler.Call) error {
	if message == nil {
		return fmt.Errorf("scheduler call must not be nil")
	}
	e.callMu.Lock()
	defer e.callMu.Unlock()
	frameworkID := e.Framework.FrameworkInfo.Id.GetValue()
	message.FrameworkID = &clusterd.FrameworkID{Value: frameworkID}
	ctx, cancel := context.WithTimeout(context.Background(), e.callTimeout)
	defer cancel()
	response, err := e.ClusterdClient.Do(message,
		httpcli.Context(ctx),
		httpcli.Header("Mesos-Stream-Id", e.StreamID()),
	)
	if response != nil {
		_ = response.Close()
	}
	return err
}

func (e *Mesos) AcceptOffer(offerID string, tasks []*mesosproto.TaskInfo, refuse time.Duration) error {
	clusterdTasks := make([]clusterd.TaskInfo, 0, len(tasks))
	for _, task := range tasks {
		data, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("marshal task for clusterd ACCEPT: %w", err)
		}
		var converted clusterd.TaskInfo
		if err := json.Unmarshal(data, &converted); err != nil {
			return fmt.Errorf("unmarshal task for clusterd ACCEPT: %w", err)
		}
		clusterdTasks = append(clusterdTasks, converted)
	}
	call := clusterdcalls.Accept(
		clusterdcalls.OfferOperations{
			clusterdcalls.OpLaunch(clusterdTasks...),
		}.WithOffers(clusterd.OfferID{Value: offerID}),
	).With(clusterdcalls.RefuseSeconds(refuse))
	return e.CallClusterd(call)
}

func (e *Mesos) DeclineOffers(offerIDs []*mesosproto.OfferID, refuse time.Duration) error {
	ids := make([]clusterd.OfferID, 0, len(offerIDs))
	for _, offerID := range offerIDs {
		if offerID != nil {
			ids = append(ids, clusterd.OfferID{Value: offerID.GetValue()})
		}
	}
	return e.CallClusterd(clusterdcalls.Decline(ids...).With(clusterdcalls.RefuseSeconds(refuse)))
}

func (e *Mesos) AcknowledgeUpdate(status *mesosproto.TaskStatus) error {
	if status == nil {
		return nil
	}
	return e.CallClusterd(clusterdcalls.Acknowledge(
		status.GetAgentId().GetValue(),
		status.GetTaskId().GetValue(),
		status.GetUuid(),
	))
}

// IsRessourceMatched - check if the ressources of the offer are matching the needs of the cmd
// nolint:gocyclo
func (e *Mesos) IsRessourceMatched(ressource []*mesosproto.Resource, cmd *cfg.Command) bool {
	mem := false
	cpu := false
	ports := true

	for _, v := range ressource {
		if v.GetName() == "cpus" && v.Scalar.GetValue() >= cmd.CPU {
			logrus.WithField("func", "mesos.IsRessourceMatched").Debug("Matched Offer CPU: ", cmd.CPU)
			cpu = true
		}
		if v.GetName() == "mem" && v.Scalar.GetValue() >= cmd.Memory {
			logrus.WithField("func", "mesos.IsRessourceMatched").Debug("Matched Offer Memory: ", cmd.Memory)
			mem = true
		}
		if v.GetName() == "gpus" && v.Scalar.GetValue() >= cmd.GPUs {
			logrus.WithField("func", "mesos.IsRessourceMatched").Debug("Matched Offer GPU: ", cmd.GPUs)
			mem = true
		}
		if len(cmd.DockerPortMappings) > 0 {
			if v.GetName() == "ports" {
				for _, taskPort := range cmd.DockerPortMappings {
					for _, portRange := range v.GetRanges().Range {
						portBegin := uint32(portRange.GetBegin())
						portEnd := uint32(portRange.GetEnd())
						if *taskPort.HostPort >= portBegin && *taskPort.HostPort <= portEnd {
							logrus.WithField("func", "mesos.IsRessourceMatched").Debug("Matched Offer TaskPort: ", taskPort.GetHostPort())
							logrus.WithField("func", "mesos.IsRessourceMatched").Debug("Matched Offer RangePort: ", portRange)
							ports = ports || true
							break
						}
						ports = ports || false
					}
				}
			}
		}
	}

	return mem && cpu && ports
}

// GetAgentInfo get information about the agent
func (e *Mesos) GetAgentInfo(agentID string) cfg.MesosSlaves {
	now := time.Now()
	e.agentCacheMu.Lock()
	if e.agentCache == nil {
		e.agentCache = make(map[string]agentCacheEntry)
	}
	if cached, ok := e.agentCache[agentID]; ok && now.Before(cached.expires) {
		e.agentCacheMu.Unlock()
		return cached.agent
	}
	e.agentCacheMu.Unlock()

	req, err := http.NewRequest("POST", mesosEndpoint(e.Framework, "/slaves/"+agentID), nil)
	if err != nil {
		logrus.WithField("func", "mesos.getAgentInfo").Error("Could not create agent request: ", err)
		return cfg.MesosSlaves{}
	}
	req.SetBasicAuth(e.Framework.Username, e.Framework.Password)
	req.Header.Set("Mesos-Stream-Id", e.StreamID())
	req.Header.Set("Content-Type", "application/json")
	res, err := e.Client.Do(req)

	if err != nil {
		logrus.WithField("func", "mesos.getAgentInfo").Error("Could not connect to master: ", err.Error())
		return cfg.MesosSlaves{}
	}

	if res.StatusCode == http.StatusOK {
		defer res.Body.Close()

		var agent cfg.MesosAgent
		err = json.NewDecoder(res.Body).Decode(&agent)
		if err != nil {
			logrus.WithField("func", "getAgentInfo").Error("Could not encode json result: ", err.Error())
			// if there is an error, dump out the res.Body as debug
			bodyBytes, err := io.ReadAll(res.Body)
			if err == nil {
				logrus.WithField("func", "getAgentInfo").Debug("response Body Dump: ", string(bodyBytes))
			}
			return cfg.MesosSlaves{}
		}

		// save how many agents the cluster has
		e.CountAgent = len(agent.Slaves)

		// get the used agent info
		for _, a := range agent.Slaves {
			if a.ID == agentID {
				e.agentCacheMu.Lock()
				e.agentCache[agentID] = agentCacheEntry{agent: a, expires: time.Now().Add(agentInfoCacheTTL)}
				e.agentCacheMu.Unlock()
				return a
			}
		}
	}

	return cfg.MesosSlaves{}
}

func (e *Mesos) GetNetworkInfo(taskID string) []*mesosproto.NetworkInfo {
	task := e.GetTaskInfo(taskID)

	if len(task.Tasks) > 0 {
		for _, status := range task.Tasks[0].Statuses {
			if status.State == "TASK_RUNNING" {
				var netw []*mesosproto.NetworkInfo
				netw = append(netw, status.ContainerStatus.NetworkInfos[0])
				return netw
			}
		}
	}
	return []*mesosproto.NetworkInfo{}
}

// GetTaskInfo get the task object to the given ID
func (e *Mesos) GetTaskInfo(taskID string) cfg.MesosTasks {
	req, err := http.NewRequest("POST", mesosEndpoint(e.Framework, "/tasks/?task_id="+taskID+"&framework_id="+e.Framework.FrameworkInfo.Id.GetValue()), nil)
	if err != nil {
		logrus.WithField("func", "mesos.GetTaskInfo").Error("Could not create task request: ", err)
		return cfg.MesosTasks{}
	}
	req.SetBasicAuth(e.Framework.Username, e.Framework.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := e.Client.Do(req)

	if err != nil {
		logrus.WithField("func", "mesos.GetTaskInfo").Error("Could not connect to mesos-master: ", err.Error())
		return cfg.MesosTasks{}
	}

	defer res.Body.Close()

	var task cfg.MesosTasks
	err = json.NewDecoder(res.Body).Decode(&task)
	if err != nil {
		logrus.WithField("func", "mesos.GetTaskInfo").Error("Could not encode json result: ", err.Error())
		return cfg.MesosTasks{}
	}

	return task
}
