package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	util "github.com/AVENTER-UG/util/util"
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
	framework.Username = util.Getenv("MESOS_USERNAME", "")
	framework.Password = util.Getenv("MESOS_PASSWORD", "")
	framework.MesosMasterServer = util.Getenv("MESOS_MASTER", "localhost:5050")
	framework.MesosCNI = util.Getenv("MESOS_CNI", "weave")
	framework.PortRangeFrom, _ = strconv.Atoi(util.Getenv("PORTRANGE_FROM", "31000"))
	framework.PortRangeTo, _ = strconv.Atoi(util.Getenv("PORTRANGE_TO", "32000"))
	config.Principal = util.Getenv("MESOS_PRINCIPAL", "")
	config.CPU, _ = strconv.ParseFloat(util.Getenv("DEFAULT_CPU", "0.001"), 64)
	config.Memory, _ = strconv.ParseFloat(util.Getenv("DEFAULT_MEMORY", "50"), 64)
	config.Disk, _ = strconv.ParseFloat(util.Getenv("DEFAULT_DISK", "1000"), 64)
	config.LogLevel = util.Getenv("LOGLEVEL", "info")
	config.Domain = util.Getenv("DOMAIN", "local")
	config.Credentials.Username = util.Getenv("AUTH_USERNAME", "")
	config.Credentials.Password = util.Getenv("AUTH_PASSWORD", "")
	config.AppName = "Mesos Compose Framework"
	config.ReconcileLoopTime, _ = time.ParseDuration(util.Getenv("RECONCILE_WAIT", "10m"))
	config.RedisServer = util.Getenv("REDIS_SERVER", "127.0.0.1:6379")
	config.RedisPassword = util.Getenv("REDIS_PASSWORD", "")
	config.RedisDB, _ = strconv.Atoi(util.Getenv("REDIS_DB", "1"))
	config.SSLKey = util.Getenv("SSL_KEY_BASE64", "")
	config.SSLCrt = util.Getenv("SSL_CRT_BASE64", "")
	config.PrefixTaskName = util.Getenv("PREFIX_TASKNAME", framework.FrameworkName)
	config.PrefixHostname = util.Getenv("PREFIX_HOSTNAME", framework.FrameworkName)
	config.EventLoopTime, _ = time.ParseDuration(util.Getenv("HEARTBEAT_INTERVAL", "15s"))
	config.VaultToken = util.Getenv("VAULT_TOKEN", "")
	config.VaultURL = util.Getenv("VAULT_URL", "http://127.0.0.1:8200")
	config.VaultTimeout, _ = time.ParseDuration(util.Getenv("VAULT_TIMEOUT", "10s"))

	// The comunication to the mesos server should be via ssl or not
	if strings.Compare(os.Getenv("MESOS_SSL"), "true") == 0 {
		framework.MesosSSL = true
	} else {
		framework.MesosSSL = false
	}

	// Skip SSL Verification
	if strings.Compare(os.Getenv("SKIP_SSL"), "true") == 0 {
		config.SkipSSL = true
	} else {
		config.SkipSSL = false
	}

	listen := fmt.Sprintf(":%s", framework.FrameworkPort)

	failoverTimeout := 5000.0
	checkpoint := true
	webuiurl := fmt.Sprintf("http://%s%s", framework.FrameworkHostname, listen)
	if config.SSLCrt != "" && config.SSLKey != "" {
		webuiurl = fmt.Sprintf("https://%s%s", framework.FrameworkHostname, listen)
	}

	// overwrite the webui URL
	if os.Getenv("FRAMEWORK_WEBUIURL") != "" {
		webuiurl = os.Getenv("FRAMEWORK_WEBUIURL")
	}

	framework.CommandChan = make(chan mesosutil.Command, 100)
	config.Hostname = framework.FrameworkHostname
	config.Listen = listen
	config.Suppress = false

	framework.State = map[string]mesosutil.State{}

	framework.FrameworkInfo.User = framework.FrameworkUser
	framework.FrameworkInfo.Name = framework.FrameworkName
	framework.FrameworkInfo.WebUiURL = &webuiurl
	framework.FrameworkInfo.FailoverTimeout = &failoverTimeout
	framework.FrameworkInfo.Checkpoint = &checkpoint
	framework.FrameworkInfo.Principal = &config.Principal
	framework.FrameworkInfo.Role = &framework.FrameworkRole
}
