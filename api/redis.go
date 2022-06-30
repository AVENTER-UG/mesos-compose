package api

import (
	"context"
	"encoding/json"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	goredis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// GetAllRedisKeys get out all redis keys to a patter
func (e *API) GetAllRedisKeys(pattern string) *goredis.ScanIterator {
	val := e.Config.RedisClient.Scan(e.Config.RedisCTX, 0, pattern, 0).Iterator()
	if err := val.Err(); err != nil {
		logrus.Error("getAllRedisKeys: ", err)
	}

	return val
}

// GetRedisKey get out the data of a key
func (e *API) GetRedisKey(key string) string {
	val, err := e.Config.RedisClient.Get(e.Config.RedisCTX, key).Result()
	if err != nil {
		logrus.Error("getRedisKey: ", err)
	}

	return val
}

// DelRedisKey will delete a redis key
func (e *API) DelRedisKey(key string) int64 {
	val, err := e.Config.RedisClient.Del(e.Config.RedisCTX, key).Result()
	if err != nil {
		logrus.Error("delRedisKey: ", err)
		e.PingRedis()
	}

	return val
}

// GetTaskFromEvent get out the key by a mesos event
func (e *API) GetTaskFromEvent(update *mesosproto.Event_Update) mesosutil.Command {
	// search matched taskid in redis and update the status
	keys := e.GetAllRedisKeys(e.Framework.FrameworkName + ":*")
	for keys.Next(e.Config.RedisCTX) {
		// get the values of the current key
		key := e.GetRedisKey(keys.Val())

		// update the status of the matches task
		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)
		if task.TaskID == update.Status.TaskID.Value {
			return task
		}
	}

	return mesosutil.Command{}
}

// CountRedisKey will get back the count of the redis key
func (e *API) CountRedisKey(pattern string) int {
	keys := e.GetAllRedisKeys(pattern)
	count := 0
	for keys.Next(e.Config.RedisCTX) {
		count++
	}
	logrus.Debug("CountRedisKey: ", pattern, count)
	return count
}

// SaveConfig store the current framework config
func (e *API) SaveConfig() {
	data, _ := json.Marshal(e.Config)
	err := e.Config.RedisClient.Set(e.Config.RedisCTX, e.Framework.FrameworkName+":framework_config", data, 0).Err()
	if err != nil {
		logrus.Error("Framework save config state into redis error:", err)
	}
}

// PingRedis to check the health of redis
func (e *API) PingRedis() error {
	pong, err := e.Config.RedisClient.Ping(e.Config.RedisCTX).Result()
	logrus.Debug("Redis Health: ", pong, err)
	if err != nil {
		return err
	}
	return nil
}

// ConnectRedis will connect the redis DB and save the client pointer
func (e *API) ConnectRedis() {
	var redisOptions goredis.Options
	redisOptions.Addr = e.Config.RedisServer
	redisOptions.DB = e.Config.RedisDB
	if e.Config.RedisPassword != "" {
		redisOptions.Password = e.Config.RedisPassword
	}

	e.Config.RedisClient = goredis.NewClient(&redisOptions)
	e.Config.RedisCTX = context.Background()

	err := e.PingRedis()
	if err != nil {
		e.Config.RedisClient.Close()
		e.ConnectRedis()
	}
}
