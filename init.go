package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
	"time"

	util "github.com/AVENTER-UG/util/util"
	"github.com/Showmax/go-fqdn"
	"github.com/sirupsen/logrus"

	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	"github.com/AVENTER-UG/mesos-compose/redis"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
)

var config cfg.Config
var framework cfg.FrameworkConfig

func init() {
	framework.FrameworkUser = util.Getenv("FRAMEWORK_USER", "root")
	framework.FrameworkName = util.Getenv("FRAMEWORK_NAME", "mc")
	framework.FrameworkRole = util.Getenv("FRAMEWORK_ROLE", "mc")
	framework.FrameworkPort = util.Getenv("FRAMEWORK_PORT", "10000")
	framework.FrameworkHostname = util.Getenv("FRAMEWORK_HOSTNAME", fqdn.Get())
	framework.FrameworkInfoFilePath = util.Getenv("FRAMEWORK_STATEFILE_PATH", "/tmp")
	framework.Username = util.Getenv("MESOS_USERNAME", "")
	framework.Password = util.Getenv("MESOS_PASSWORD", "")
	framework.MesosMasterServer = util.Getenv("MESOS_MASTER", "localhost:5050")
	framework.PortRangeFrom, _ = strconv.Atoi(util.Getenv("PORTRANGE_FROM", "31000"))
	framework.PortRangeTo, _ = strconv.Atoi(util.Getenv("PORTRANGE_TO", "32000"))
	config.Principal = util.Getenv("MESOS_PRINCIPAL", "")
	config.CPU, _ = strconv.ParseFloat(util.Getenv("DEFAULT_CPU", "0.1"), 64)
	config.Memory, _ = strconv.ParseFloat(util.Getenv("DEFAULT_MEMORY", "50"), 64)
	config.Disk, _ = strconv.ParseFloat(util.Getenv("DEFAULT_DISK", "1000"), 64)
	config.DefaultVolumeDriver = util.Getenv("DEFAULT_VOLUME_DRIVER", "local")
	config.LogLevel = util.Getenv("LOGLEVEL", "info")
	config.Domain = util.Getenv("DOMAIN", "local")
	config.Credentials.Username = util.Getenv("AUTH_USERNAME", "")
	config.Credentials.Password = util.Getenv("AUTH_PASSWORD", "")
	config.AppName = "Mesos Compose Framework"
	config.ReconcileLoopTime, _ = time.ParseDuration(util.Getenv("RECONCILE_WAIT", "30m"))
	config.RedisServer = util.Getenv("REDIS_SERVER", "127.0.0.1:6379")
	config.RedisPassword = util.Getenv("REDIS_PASSWORD", "")
	config.RedisDB, _ = strconv.Atoi(util.Getenv("REDIS_DB", "1"))
	config.RedisPoolSize, _ = strconv.Atoi(util.Getenv("REDIS_POOLSIZE", "0"))
	config.SSLKey = util.Getenv("SSL_KEY_BASE64", "")
	config.SSLCrt = util.Getenv("SSL_CRT_BASE64", "")
	config.PrefixTaskName = util.Getenv("PREFIX_TASKNAME", framework.FrameworkName)
	config.PrefixHostname = util.Getenv("PREFIX_HOSTNAME", framework.FrameworkName)
	config.EventLoopTime, _ = time.ParseDuration(util.Getenv("HEARTBEAT_INTERVAL", "15s"))
	config.VaultToken = util.Getenv("VAULT_TOKEN", "")
	config.VaultURL = util.Getenv("VAULT_URL", "http://127.0.0.1:8200")
	config.VaultTimeout, _ = time.ParseDuration(util.Getenv("VAULT_TIMEOUT", "10s"))
	config.DiscoveryInfoNameDelimiter = util.Getenv("DISCOVERY_INFONAME_DELIMITER", ".")
	config.DiscoveryPortNameDelimiter = util.Getenv("DISCOVERY_PORTNAME_DELIMITER", "_")

	// Enable Threads
	if strings.Compare(util.Getenv("THREAD_ENABLE", "false"), "false") == 0 {
		config.ThreadEnable = false
	} else {
		config.ThreadEnable = true
	}

	// Enable plugins
	if strings.Compare(util.Getenv("COMPOSE_PLUGINS_ENABLE", "false"), "false") == 0 {
		config.PluginsEnable = false
	} else {
		config.PluginsEnable = true
	}

	// The comunication to the mesos server should be via ssl or not
	if util.Getenv("MESOS_SSL", "false") == "true" {
		framework.MesosSSL = true
	} else {
		framework.MesosSSL = false
	}

	// Skip SSL Verification
	if util.Getenv("SKIP_SSL", "true") == "true" {
		config.SkipSSL = true
	} else {
		config.SkipSSL = false
	}

	// Enable GPU Allocation in Mesos.. If false, GPU can still be used but allocation wont be impacted in mesos
	if strings.Compare(util.Getenv("ENABLE_GPU_ALLOCATION", "true"), "true") == 0 {
		config.EnableGPUAllocation = true
	} else {
		config.EnableGPUAllocation = false
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

	framework.CommandChan = make(chan cfg.Command, 100)
	config.Hostname = framework.FrameworkHostname
	config.Listen = listen
	config.Suppress = false

	// tell mesos that we can handle gpu offers
	capabilities := &mesosproto.FrameworkInfo_Capability{
		Type: mesosproto.FrameworkInfo_Capability_GPU_RESOURCES.Enum(),
	}
	framework.State = map[string]cfg.State{}

	framework.FrameworkInfo.User = util.StringToPointer(framework.FrameworkUser)
	framework.FrameworkInfo.Name = util.StringToPointer(framework.FrameworkName)
	framework.FrameworkInfo.WebuiUrl = &webuiurl
	framework.FrameworkInfo.FailoverTimeout = &failoverTimeout
	framework.FrameworkInfo.Checkpoint = &checkpoint
	framework.FrameworkInfo.Principal = &config.Principal
	framework.FrameworkInfo.Role = util.StringToPointer(framework.FrameworkRole)
	framework.FrameworkInfo.Capabilities = append(framework.FrameworkInfo.Capabilities, capabilities)
}

func loadPlugins(r *redis.Redis) {
	if config.PluginsEnable {
		config.Plugins = map[string]*plugin.Plugin{}

		plugins, err := filepath.Glob("plugins/*.so")
		if err != nil {
			logrus.WithField("func", "main.loadPlugins").Info("No Plugins found")
			return
		}

		for _, filename := range plugins {
			p, err := plugin.Open(filename)
			if err != nil {
				logrus.WithField("func", "main.initPlugins").Error("Error during loading plugin: ", err.Error())
				continue
			}

			symbol, err := p.Lookup("Init")
			if err != nil {
				logrus.WithField("func", "main.initPlugins").Error("Error lookup init plugin: ", err.Error())
				continue
			}

			initPluginFunc, ok := symbol.(func(*redis.Redis) string)

			if !ok {
				logrus.WithField("func", "main.initPlugins").Error("Error plugin does not have init function")
				continue
			}

			name := initPluginFunc(r)
			config.Plugins[name] = p
		}
	}
}
