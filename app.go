package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/AVENTER-UG/mesos-compose/api"
	"github.com/AVENTER-UG/mesos-compose/redis"
	"github.com/AVENTER-UG/mesos-compose/scheduler"
	cfg "github.com/AVENTER-UG/mesos-compose/types"
	util "github.com/AVENTER-UG/util/util"
	"github.com/AVENTER-UG/util/vault"
	"github.com/sirupsen/logrus"
)

// BuildVersion of m3s
var BuildVersion string

// GitVersion is the revision and commit number
var GitVersion string

// convert Base64 Encodes PEM Certificate to tls object
func decodeBase64Cert(pemCert string) []byte {
	sslPem, err := base64.URLEncoding.DecodeString(pemCert)
	if err != nil {
		logrus.Fatal("Error decoding SSL PEM from Base64: ", err.Error())
	}

	return sslPem
}

func main() {
	// Prints out current version
	var version bool
	flag.BoolVar(&version, "v", false, "Prints current version")
	flag.Parse()
	if version {
		fmt.Print(GitVersion)
		return
	}

	util.SetLogging(config.LogLevel, config.EnableSyslog, config.AppName)
	logrus.Println(config.AppName + " build " + BuildVersion + " git " + GitVersion)

	// Connect the vault if we got a token
	v := vault.New(config.VaultToken, config.VaultURL, config.VaultTimeout)
	if config.VaultToken != "" {
		logrus.Info("Vault Connection: ")
		logrus.Info(v.Connect())
	}

	r := redis.New(&config, &framework)
	if !r.Connect() {
		logrus.WithField("func", "main").Fatal("Could not connect to redis DB")
	}

	// get API
	a := api.New(&config, &framework)
	a.Redis = r

	// load old framework config from database if they exist
	var oldFramework cfg.FrameworkConfig
	key := r.GetRedisKey(framework.FrameworkName + ":framework")
	if key != "" {
		json.Unmarshal([]byte(key), &oldFramework)

		framework.FrameworkInfo.Id = oldFramework.FrameworkInfo.Id
		framework.MesosStreamID = oldFramework.MesosStreamID
	}

	// The Hostname should ever be set after reading the state file.
	framework.FrameworkInfo.Hostname = &framework.FrameworkHostname

	r.SaveConfig(config)
	r.SaveFrameworkRedis(&framework)

	server := &http.Server{
		Addr:              config.Listen,
		Handler:           a.Commands(),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
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

	go loadPlugins(r)

	//	this loop is for resubscribtion purpose
	for {
		e := scheduler.Subscribe(&config, &framework)
		e.API = a
		e.Vault = v
		ctx, cancel := context.WithCancel(context.Background())
		go e.HeartbeatLoop(ctx)
		go e.ReconcileLoop(ctx)
		e.Redis = r
		e.EventLoop()
		cancel()
		e.Redis.SaveConfig(*e.Config)
		time.Sleep(60 * time.Second)
	}
}
