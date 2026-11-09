package main

import (
	mesosproto "github.com/m3scluster/mesos-compose/proto"
	"github.com/m3scluster/mesos-compose/redis"
)

type Plugins struct {
	Redis *redis.Redis
}

var plugin *Plugins

func Init(r *redis.Redis) string {
	plugin = &Plugins{
		Redis: r,
	}

	return "example"
}

func Event(event mesosproto.Event) {
}
