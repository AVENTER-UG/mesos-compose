package main

import (
	"os"
	"strings"

	util "github.com/AVENTER-UG/util"
	"github.com/Showmax/go-fqdn"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
)

var config cfg.Config

func init() {
	config.FrameworkUser = util.Getenv("FRAMEWORK_USER", "root")
	config.FrameworkName = "mc" + util.Getenv("FRAMEWORK_NAME", "")
	config.FrameworkRole = util.Getenv("FRAMEWORK_ROLE", "mc")
	config.FrameworkPort = util.Getenv("FRAMEWORK_PORT", "10000")
	config.FrameworkHostname = util.Getenv("FRAMEWORK_HOSTNAME", fqdn.Get())
	config.FrameworkInfoFilePath = util.Getenv("FRAMEWORK_STATEFILE_PATH", "/tmp")
	config.Principal = os.Getenv("MESOS_PRINCIPAL")
	config.Username = os.Getenv("MESOS_USERNAME")
	config.Password = os.Getenv("MESOS_PASSWORD")
	config.MesosMasterServer = os.Getenv("MESOS_MASTER")
	config.MesosCNI = util.Getenv("MESOS_CNI", "weave")
	config.LogLevel = util.Getenv("LOGLEVEL", "info")
	config.Domain = util.Getenv("DOMAIN", "local")
	config.Credentials.Username = os.Getenv("AUTH_USERNAME")
	config.Credentials.Password = os.Getenv("AUTH_PASSWORD")
	config.AppName = "Mesos Compose Framework"
	config.PrefixTaskName = util.Getenv("PREFIX_TASKNAME", "mc")
	config.PrefixHostname = util.Getenv("PREFIX_HOSTNAME", "mc")

	// The comunication to the mesos server should be via ssl or not
	if strings.Compare(os.Getenv("MESOS_SSL"), "true") == 0 {
		config.MesosSSL = true
	} else {
		config.MesosSSL = false
	}

}
