package mesos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	clusterdcalls "github.com/m3scluster/clusterd-go/api/v1/lib/scheduler/calls"
)

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

func TestNewReusesConfiguredHTTPTransport(t *testing.T) {
	e := New(&cfg.Config{SkipSSL: true}, &cfg.FrameworkConfig{})
	if e.Client == nil || e.Transport == nil || e.Client.Transport != e.Transport {
		t.Fatal("Mesos did not initialize a reusable HTTP client and transport")
	}
	if !e.Transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("Mesos transport ignored SkipSSL")
	}
}

func TestCallReturnsMesosHTTPErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("framework is not subscribed"))
	}))
	defer server.Close()

	frameworkID := "framework-1"
	e := New(&cfg.Config{}, &cfg.FrameworkConfig{
		FrameworkInfo:     mesosproto.FrameworkInfo{Id: &mesosproto.FrameworkID{Value: &frameworkID}},
		MesosMasterServer: strings.TrimPrefix(server.URL, "http://"),
		MesosStreamID:     "stream-1",
	})

	err := e.CallClusterd(clusterdcalls.Revive())
	if err == nil || !strings.Contains(err.Error(), "framework is not subscribed") {
		t.Fatalf("Call error = %v, want Mesos response diagnostic", err)
	}
}

func TestCallAcceptsOKAndAcceptedResponses(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusAccepted} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if status == http.StatusOK {
				w.Header().Set("Content-Type", "application/json")
			}
			w.WriteHeader(status)
		}))

		frameworkID := "framework-1"
		e := New(&cfg.Config{}, &cfg.FrameworkConfig{
			FrameworkInfo:     mesosproto.FrameworkInfo{Id: &mesosproto.FrameworkID{Value: &frameworkID}},
			MesosMasterServer: strings.TrimPrefix(server.URL, "http://"),
		})
		if err := e.CallClusterd(clusterdcalls.Revive()); err != nil {
			t.Fatalf("Call status %d returned error: %v", status, err)
		}
		server.Close()
	}
}

func TestSubscribeBuildsClusterdCall(t *testing.T) {
	frameworkID := "framework-1"
	e := New(&cfg.Config{}, &cfg.FrameworkConfig{
		FrameworkInfo:     mesosproto.FrameworkInfo{Id: &mesosproto.FrameworkID{Value: &frameworkID}},
		MesosMasterServer: "master.mesos:5050",
	})

	e.Subscribe()
	if e.ClusterdSubscribe == nil {
		t.Fatal("Subscribe did not create a clusterd call")
	}
	if got := e.ClusterdSubscribe.GetType().String(); got != "SUBSCRIBE" {
		t.Fatalf("call type = %q, want SUBSCRIBE", got)
	}
	if got := e.ClusterdSubscribe.GetFrameworkID().GetValue(); got != frameworkID {
		t.Fatalf("framework ID = %q, want %q", got, frameworkID)
	}
}

func TestGetAgentInfoUsesTTLCache(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		_, _ = w.Write([]byte(`{"slaves":[{"id":"agent-1","hostname":"node-1"}]}`))
	}))
	defer server.Close()

	e := New(&cfg.Config{}, &cfg.FrameworkConfig{
		MesosMasterServer: strings.TrimPrefix(server.URL, "http://"),
	})
	first := e.GetAgentInfo("agent-1")
	second := e.GetAgentInfo("agent-1")
	if first.ID != "agent-1" || second.Hostname != "node-1" {
		t.Fatalf("cached agent data = %#v, %#v", first, second)
	}
	if got := atomic.LoadInt32(&requests); got != 1 {
		t.Fatalf("request count = %d, want 1 cache miss", got)
	}

	e.agentCacheMu.Lock()
	e.agentCache["agent-1"] = agentCacheEntry{agent: first, expires: time.Now().Add(-time.Second)}
	e.agentCacheMu.Unlock()
	_ = e.GetAgentInfo("agent-1")
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Fatalf("request count after expiry = %d, want 2", got)
	}
}

func TestCallTimesOutWithoutAffectingStreamClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	frameworkID := "framework-1"
	e := New(&cfg.Config{}, &cfg.FrameworkConfig{
		FrameworkInfo:     mesosproto.FrameworkInfo{Id: &mesosproto.FrameworkID{Value: &frameworkID}},
		MesosMasterServer: strings.TrimPrefix(server.URL, "http://"),
	})
	e.callTimeout = 10 * time.Millisecond

	if err := e.CallClusterd(clusterdcalls.Revive()); err == nil {
		t.Fatal("Call succeeded despite an expired request context")
	}
}

func TestCallSerializesSchedulerRequests(t *testing.T) {
	var active, maxActive int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&active, 1)
		for {
			old := atomic.LoadInt32(&maxActive)
			if current <= old || atomic.CompareAndSwapInt32(&maxActive, old, current) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		atomic.AddInt32(&active, -1)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	frameworkID := "framework-1"
	e := New(&cfg.Config{}, &cfg.FrameworkConfig{
		FrameworkInfo:     mesosproto.FrameworkInfo{Id: &mesosproto.FrameworkID{Value: &frameworkID}},
		MesosMasterServer: strings.TrimPrefix(server.URL, "http://"),
	})

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := e.CallClusterd(clusterdcalls.Revive()); err != nil {
				t.Errorf("Call returned error: %v", err)
			}
		}()
	}
	wg.Wait()
	if got := atomic.LoadInt32(&maxActive); got != 1 {
		t.Fatalf("maximum concurrent scheduler calls = %d, want 1", got)
	}
}

func TestClusterdSchedulerCallPreservesWireFields(t *testing.T) {
	frameworkID := "framework-1"
	call := &mesosproto.Call{
		FrameworkId: &mesosproto.FrameworkID{Value: &frameworkID},
		Type:        mesosproto.Call_REVIVE.Enum(),
	}
	converted, err := clusterdSchedulerCall(call)
	if err != nil {
		t.Fatalf("convert scheduler call: %v", err)
	}
	data, err := json.Marshal(&converted)
	if err != nil {
		t.Fatalf("marshal clusterd scheduler call: %v", err)
	}
	var wire map[string]any
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("decode wire JSON: %v", err)
	}
	if wire["type"] != "REVIVE" {
		t.Fatalf("wire type = %v, want REVIVE", wire["type"])
	}
	framework, ok := wire["framework_id"].(map[string]any)
	if !ok || framework["value"] != frameworkID {
		t.Fatalf("wire framework_id = %v, want %q", wire["framework_id"], frameworkID)
	}
}

func stringPointer(value string) *string { return &value }
