package api

import (
	"encoding/json"

	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	goredis "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func GetAllRedisKeys(pattern string) *goredis.ScanIterator {
	val := config.RedisClient.Scan(config.RedisCTX, 0, pattern, 0).Iterator()
	if err := val.Err(); err != nil {
		logrus.Error("getAllRedisKeys: ", err)
	}

	return val
}

func GetRedisKey(key string) string {
	val, err := config.RedisClient.Get(config.RedisCTX, key).Result()
	if err != nil {
		logrus.Error("getRedisKey: ", err)
	}

	return val
}

func GetTaskFromEvent(update *mesosproto.Event_Update) mesosutil.Command {
	// search matched taskid in redis and update the status
	keys := GetAllRedisKeys("*")
	for keys.Next(config.RedisCTX) {
		// get the values of the current key
		key := GetRedisKey(keys.Val())

		// update the status of the matches task
		var task mesosutil.Command
		json.Unmarshal([]byte(key), &task)
		if task.TaskID == update.Status.TaskID.Value {
			return task
		}
	}

	return mesosutil.Command{}
}

// SaveConfig store the current framework config
func SaveConfig() error {
	data, _ := json.Marshal(config)
	err := config.RedisClient.Set(config.RedisCTX, "framework_config", data, 0).Err()
	if err != nil {
		logrus.Error("getRedisKey: ", err)
		return err
	}
	return nil
}
