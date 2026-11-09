package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"

	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/m3scluster/mesos-compose/api"
	"github.com/m3scluster/mesos-compose/redis"
	"github.com/m3scluster/mesos-compose/scheduler"
	cfg "github.com/m3scluster/mesos-compose/types"
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

// ReconnectBackoffBase is the initial delay before the first resubscription attempt.
const ReconnectBackoffBase = 2 * time.Second

// ReconnectBackoffMax caps the exponential resubscription backoff.
const ReconnectBackoffMax = 60 * time.Second

// reconnectBackoff returns a deterministic exponential backoff delay for the
// given 0-based reconnect attempt: it starts at ReconnectBackoffBase and
// doubles per attempt until it reaches ReconnectBackoffMax, where it plateaus.
// No jitter is added, so the sequence is fully deterministic and testable.
//
// Limitation: main() cannot observe whether the attempt that preceded an
// EventLoop() return ever produced a successful subscription, so the attempt
// counter is never reset during a run; it only restarts from 0 when the
// process is restarted. A caller that *does* know its subscription succeeded
// should pass a fresh (lower) attempt index to reset the backoff.
func reconnectBackoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	delay := ReconnectBackoffBase
	for attempt > 0 && delay < ReconnectBackoffMax {
		delay *= 2
		attempt--
	}
	if delay > ReconnectBackoffMax {
		delay = ReconnectBackoffMax
	}
	return delay
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
		if storedFramework, err := redis.DecodeFramework([]byte(key)); err == nil {
			oldFramework = *storedFramework
		} else {
			logrus.WithField("func", "main").Warn("Could not decode persisted framework: ", err)
		}

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
	backoffAttempt := 0
	for {
		e := scheduler.Subscribe(&config, &framework)
		e.API = a
		e.Vault = v
		ctx, cancel := context.WithCancel(context.Background())
		go e.RunAfterSubscription(ctx, e.HeartbeatLoop)
		go e.RunAfterSubscription(ctx, e.ReconcileLoop)
		e.Redis = r
		e.EventLoop()
		cancel()
		e.Redis.SaveConfig(*e.Config)
		time.Sleep(reconnectBackoff(backoffAttempt))
		backoffAttempt++
	}
}
