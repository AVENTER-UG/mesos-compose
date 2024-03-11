package main

import (
	mesosproto "github.com/AVENTER-UG/mesos-compose/proto"
	"github.com/AVENTER-UG/mesos-compose/redis"
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
