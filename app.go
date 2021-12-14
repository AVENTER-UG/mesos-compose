package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AVENTER-UG/mesos-compose/api"
	"github.com/AVENTER-UG/mesos-compose/mesos"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	util "github.com/AVENTER-UG/util"
	goredis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// MinVersion is the version number of this program
var MinVersion string

// init the redis cache
func initCache() {
	var redisOptions goredis.Options
	redisOptions.Addr = config.RedisServer
	redisOptions.DB = 1
	if config.RedisPassword != "" {
		redisOptions.Password = config.RedisPassword
	}
	client := goredis.NewClient(&redisOptions)

	config.RedisCTX = context.Background()
	pong, err := client.Ping(config.RedisCTX).Result()
	logrus.Debug("Redis Health: ", pong, err)
	config.RedisClient = client
}

// convert Base64 Encodes PEM Certificate to tls object
func decodeBase64Cert(pemCert string) []byte {
	sslPem, err := base64.URLEncoding.DecodeString(pemCert)
	if err != nil {
		logrus.Fatal("Error decoding SSL PEM from Base64: ", err.Error())
	}
	return sslPem
}

func main() {
	util.SetLogging(config.LogLevel, config.EnableSyslog, config.AppName)
	logrus.Println(config.AppName + " build " + MinVersion)

	listen := fmt.Sprintf(":%s", framework.FrameworkPort)

	failoverTimeout := 5000.0
	checkpoint := true
	webuiurl := fmt.Sprintf("http://%s%s", framework.FrameworkHostname, listen)

	framework.CommandChan = make(chan mesosutil.Command, 100)
	config.Hostname = framework.FrameworkHostname
	config.Listen = listen

	framework.State = map[string]mesosutil.State{}

	framework.FrameworkInfo.User = framework.FrameworkUser
	framework.FrameworkInfo.Name = framework.FrameworkName
	framework.FrameworkInfo.WebUiURL = &webuiurl
	framework.FrameworkInfo.FailoverTimeout = &failoverTimeout
	framework.FrameworkInfo.Checkpoint = &checkpoint
	framework.FrameworkInfo.Principal = &config.Principal
	framework.FrameworkInfo.Role = &framework.FrameworkRole
	//	config.FrameworkInfo.Capabilities = []mesosproto.FrameworkInfo_Capability{
	//		{Type: mesosproto.FrameworkInfo_Capability_RESERVATION_REFINEMENT},
	//	}

	initCache()
	mesosutil.SetConfig(&framework)
	api.SetConfig(&config, &framework)
	mesos.SetConfig(&config, &framework)

	// load framework state from database if they exist
	key := api.GetRedisKey("framework")
	if key != "" {
		json.Unmarshal([]byte(key), &framework)
	}

	// The Hostname should ever be set after reading the state file.
	framework.FrameworkInfo.Hostname = &framework.FrameworkHostname

	server := &http.Server{
		Addr:    listen,
		Handler: api.Commands(),
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
			MinVersion: tls.VersionTLS12,
		},
	}

	if config.SSLCrt != "" && config.SSLKey != "" {
		logrus.Debug("Enable TLS")
		crt := decodeBase64Cert(config.SSLCrt)
		key := decodeBase64Cert(config.SSLKey)
		certs, err := tls.X509KeyPair(crt, key)
		if err != nil {
			logrus.Fatal("TLS Server Error: ", err.Error())
		}
		server.TLSConfig.Certificates = []tls.Certificate{certs}
	}

	go func() {
		if config.SSLCrt != "" && config.SSLKey != "" {
			server.ListenAndServeTLS("", "")
		} else {
			server.ListenAndServe()
		}
	}()
	logrus.Fatal(mesos.Subscribe())
}
