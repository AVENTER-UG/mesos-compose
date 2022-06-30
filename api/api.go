package api

import (

	//"encoding/json"

	"encoding/json"

	"github.com/gorilla/mux"

	//"io/ioutil"
	"net/http"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
)

// API Service include all the current vars and global config
type API struct {
	Config    *cfg.Config
	Framework *mesosutil.FrameworkConfig
	Service   cfg.Service
	Compose   cfg.Compose
	Redis     Redis
}

// New will create a new API object
func New(cfg *cfg.Config, frm *mesosutil.FrameworkConfig) *API {
	e := &API{
		Config:    cfg,
		Framework: frm,
	}

	return e
}

// Commands is the main function of this package
func (e *API) Commands() *mux.Router {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/api/compose/versions", e.Versions).Methods("GET")
	rtr.HandleFunc("/api/compose/v0/tasks", e.V0ShowAllTasks).Methods("GET")
	rtr.HandleFunc("/api/compose/v0/{project}", e.V0ComposePush).Methods("PUT")
	rtr.HandleFunc("/api/compose/v0/{project}", e.V0ComposeUpdate).Methods("UPDATE")
	rtr.HandleFunc("/api/compose/v0/{project}/{servicename}", e.V0ComposeKillService).Methods("DELETE")
	rtr.HandleFunc("/api/compose/v0/{project}/{servicename}/restart", e.V0ComposeRestartService).Methods("PUT")
	rtr.HandleFunc("/api/compose/v0/{project}/{servicename}/{taskid}", e.V0ComposeKillTask).Methods("DELETE")

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

// ErrorMessage will create a message json
func (e *API) ErrorMessage(number int, function string, msg string) []byte {
	var err cfg.ErrorMsg
	err.Function = function
	err.Number = number
	err.Message = msg

	data, _ := json.Marshal(err)
	return data
}
