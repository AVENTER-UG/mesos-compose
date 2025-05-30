package api

import (

	//"encoding/json"

	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	//"io/ioutil"
	"net/http"

	"github.com/AVENTER-UG/mesos-compose/mesos"
	"github.com/AVENTER-UG/mesos-compose/redis"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	"github.com/AVENTER-UG/util/vault"
)

// API Service include all the current vars and global config
type API struct {
	Config    *cfg.Config
	Framework *cfg.FrameworkConfig
	Service   cfg.Service
	Compose   cfg.Compose
	Redis     *redis.Redis
	Vault     *vault.Vault
	Mesos     mesos.Mesos
}

// New will create a new API object
func New(cfg *cfg.Config, frm *cfg.FrameworkConfig) *API {
	e := &API{
		Config:    cfg,
		Framework: frm,
		Mesos:     *mesos.New(cfg, frm),
	}

	// Connect the vault if we got a token
	e.Vault = vault.New(cfg.VaultToken, cfg.VaultURL, cfg.VaultTimeout)
	if e.Config.VaultToken != "" {
		logrus.WithField("func", "api.New").Info("Vault Connection: ", e.Vault.Connect())
	}

	return e
}

// Commands is the main function of this package
func (e *API) Commands() *mux.Router {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/api/compose/versions", e.Versions).Methods("GET")
	rtr.HandleFunc("/api/compose/v0/tasks", e.V0ShowAllTasks).Methods("GET")
	rtr.HandleFunc("/api/compose/v0/tasks/{taskid}", e.V0ComposeKillTask).Methods("DELETE")
	rtr.HandleFunc("/api/compose/v0/framework/reregister", e.V0FrameworkReRegister).Methods("PUT")
	rtr.HandleFunc("/api/compose/v0/framework/suppress", e.V0FrameworkSuppress).Methods("PUT")
	rtr.HandleFunc("/api/compose/v0/framework", e.V0FrameworkRemoveID).Methods("DELETE")
	rtr.HandleFunc("/api/compose/v0/{project}", e.V0ComposePush).Methods("PUT")
	rtr.HandleFunc("/api/compose/v0/{project}", e.V0ComposeUpdate).Methods("UPDATE")
	rtr.HandleFunc("/api/compose/v0/{project}/{servicename}", e.V0ComposeKillService).Methods("DELETE")
	rtr.HandleFunc("/api/compose/v0/{project}/{servicename}/restart", e.V0ComposeRestartService).Methods("PUT")

	return rtr
}

// Versions give out a list of Versions
func (e *API) Versions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "-")
	w.Write([]byte("/api/compose/v0"))
}

// CheckAuth will check if the token is valid
func (e *API) CheckAuth(r *http.Request, w http.ResponseWriter) bool {
	// if no credentials are configured, then we dont have to check
	if e.Config.Credentials.Username == "" || e.Config.Credentials.Password == "" {
		return true
	}

	username, password, ok := r.BasicAuth()

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	if username == e.Config.Credentials.Username && password == e.Config.Credentials.Password {
		w.WriteHeader(http.StatusOK)
		return true
	}

	w.WriteHeader(http.StatusUnauthorized)
	return false
}

// IncreaseTaskCount split the taskID and increase the app count digit
func (e *API) IncreaseTaskCount(taskID string) string {
	tp := strings.Split(taskID, ".")
	if len(tp) == 3 {
		taskNr, _ := strconv.Atoi(tp[2])
		return tp[0] + "." + tp[1] + "." + strconv.Itoa(taskNr+1)
	}
	return taskID + ".0"
}
