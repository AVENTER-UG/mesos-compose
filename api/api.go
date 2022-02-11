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

// Service include all the current vars and global config
var config *cfg.Config
var framework *mesosutil.FrameworkConfig

// SetConfig set the global config
func SetConfig(cfg *cfg.Config, frm *mesosutil.FrameworkConfig) {
	config = cfg
	framework = frm
}

// Commands is the main function of this package
func Commands() *mux.Router {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/api/compose/versions", Versions).Methods("GET")
	rtr.HandleFunc("/api/compose/v0/tasks", V0ShowAllTasks).Methods("GET")
	rtr.HandleFunc("/api/compose/v0/{project}", V0ComposePush).Methods("PUT")
	rtr.HandleFunc("/api/compose/v0/{project}", V0ComposeUpdate).Methods("UPDATE")
	rtr.HandleFunc("/api/compose/v0/{project}/{servicename}", V0ComposeKillService).Methods("DELETE")
	rtr.HandleFunc("/api/compose/v0/{project}/{servicename}/restart", V0ComposeRestartService).Methods("PUT")
	rtr.HandleFunc("/api/compose/v0/{project}/{servicename}/{taskid}", V0ComposeKillTask).Methods("DELETE")

	return rtr
}

// Versions give out a list of Versions
func Versions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "-")
	w.Write([]byte("/api/compose/v0"))
}

// CheckAuth will check if the token is valid
func CheckAuth(r *http.Request, w http.ResponseWriter) bool {
	// if no credentials are configured, then we dont have to check
	if config.Credentials.Username == "" || config.Credentials.Password == "" {
		return true
	}

	username, password, ok := r.BasicAuth()

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	if username == config.Credentials.Username && password == config.Credentials.Password {
		w.WriteHeader(http.StatusOK)
		return true
	}

	w.WriteHeader(http.StatusUnauthorized)
	return false
}

// ErrorMessage will create a message json
func ErrorMessage(number int, function string, msg string) []byte {
	var err cfg.ErrorMsg
	err.Function = function
	err.Number = number
	err.Message = msg

	data, _ := json.Marshal(err)
	return data
}
