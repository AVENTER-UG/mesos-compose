package main

import (
	"os"
	"strconv"
	"strings"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	util "github.com/AVENTER-UG/util"
	"github.com/Showmax/go-fqdn"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
)

var config cfg.Config
var framework mesosutil.FrameworkConfig

func init() {
	framework.FrameworkUser = util.Getenv("FRAMEWORK_USER", "root")
	framework.FrameworkName = "mc" + util.Getenv("FRAMEWORK_NAME", "")
	framework.FrameworkRole = util.Getenv("FRAMEWORK_ROLE", "mc")
	framework.FrameworkPort = util.Getenv("FRAMEWORK_PORT", "10000")
	framework.FrameworkHostname = util.Getenv("FRAMEWORK_HOSTNAME", fqdn.Get())
	framework.FrameworkInfoFilePath = util.Getenv("FRAMEWORK_STATEFILE_PATH", "/tmp")
	framework.Username = os.Getenv("MESOS_USERNAME")
	framework.Password = os.Getenv("MESOS_PASSWORD")
	framework.MesosMasterServer = os.Getenv("MESOS_MASTER")
	framework.MesosCNI = util.Getenv("MESOS_CNI", "weave")
	config.CPU, _ = strconv.ParseFloat(util.Getenv("DEFAULT_CPU", "0.001"), 64)
	config.Memory, _ = strconv.ParseFloat(util.Getenv("DEFAULT_CONTAINER", "50"), 64)
	config.Principal = os.Getenv("MESOS_PRINCIPAL")
	config.LogLevel = util.Getenv("LOGLEVEL", "info")
	config.Domain = util.Getenv("DOMAIN", "local")
	config.Credentials.Username = os.Getenv("AUTH_USERNAME")
	config.Credentials.Password = os.Getenv("AUTH_PASSWORD")
	config.AppName = "Mesos Compose Framework"
	config.PrefixTaskName = util.Getenv("PREFIX_TASKNAME", "mc")
	config.PrefixHostname = util.Getenv("PREFIX_HOSTNAME", "mc")
	config.RedisServer = util.Getenv("REDIS_SERVER", "127.0.0.1:6379")

	// The comunication to the mesos server should be via ssl or not
	if strings.Compare(os.Getenv("MESOS_SSL"), "true") == 0 {
		framework.MesosSSL = true
	} else {
		framework.MesosSSL = false
	}

}
