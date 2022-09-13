package redis

import (
	"context"
	"encoding/json"

	cfg "github.com/AVENTER-UG/mesos-compose/types"
	mesosutil "github.com/AVENTER-UG/mesos-util"
	mesosproto "github.com/AVENTER-UG/mesos-util/proto"
	goredis "github.com/go-redis/redis/v8"

	"github.com/sirupsen/logrus"
)

// Redis struct about the redis connection
type Redis struct {
	Client   *goredis.Client
	CTX      context.Context
	Server   string
	Password string
	DB       int
	Prefix   string
}

// New will create a new Redis object
func New(server, password, prefix string, db int) *Redis {
	e := &Redis{
		Server:   server,
		Password: password,
		DB:       db,
		Prefix:   prefix,
		CTX:      context.Background(),
	}

	return e
}

// GetAllRedisKeys get out all redis keys to a patter
func (e *Redis) GetAllRedisKeys(pattern string) *goredis.ScanIterator {
	val := e.Client.Scan(e.CTX, 0, pattern, 0).Iterator()
	if err := val.Err(); err != nil {
		logrus.Error("getAllRedisKeys: ", err)
	}

	return val
}

// GetRedisKey get out the data of a key
func (e *Redis) GetRedisKey(key string) string {
	val, err := e.Client.Get(e.CTX, key).Result()
	if err != nil {
		logrus.Error("ge.Key: ", err)
	}

	return val
}

// SetRedisKey store data in redis
func (e *Redis) SetRedisKey(data []byte, key string) {
	err := e.Client.Set(e.CTX, key, data, 0).Err()
	if err != nil {
		logrus.WithField("func", "SetRedisKey").Error("Could not save data in Redis: ", err.Error())
	}
}

// DelRedisKey will delete a redis key
func (e *Redis) DelRedisKey(key string) int64 {
	val, err := e.Client.Del(e.CTX, key).Result()
	if err != nil {
		logrus.Error("de.Key: ", err)
		e.PingRedis()
	}

	return val
}

// GetTaskFromEvent get out the key by a mesos event
func (e *Redis) GetTaskFromEvent(update *mesosproto.Event_Update) mesosutil.Command {
	// search matched taskid in redis and update the status
	keys := e.GetAllRedisKeys(e.Prefix + ":*")
	for keys.Next(e.CTX) {
		// ignore redis keys if they are not mesos tasks
		if e.CheckIfNotTask(keys) {
			continue
		}
		// get the values of the current key
		key := e.GetRedisKey(keys.Val())
		task := mesosutil.DecodeTask(key)

		if task.TaskID == update.Status.TaskID.Value {
			task.State = update.Status.State.String()
			return task
		}
	}

	return mesosutil.Command{}
}

// CountRedisKey will get back the count of the redis key
func (e *Redis) CountRedisKey(pattern string, ignoreState string) int {
	keys := e.GetAllRedisKeys(pattern)
	count := 0
	for keys.Next(e.CTX) {
		// ignore redis keys if they are not mesos tasks
		if e.CheckIfNotTask(keys) {
			continue
		}

		// do not count these key
		if ignoreState != "" {
			// get the values of the current key
			key := e.GetRedisKey(keys.Val())
			task := mesosutil.DecodeTask(key)

			if task.State == ignoreState {
				continue
			}
		}
		count++
	}
	logrus.Debug("CountRedisKey: ", pattern, count)
	return count
}

// SaveConfig store the current framework config
func (e *Redis) SaveConfig(config cfg.Config) {
	data, _ := json.Marshal(config)
	err := e.Client.Set(e.CTX, e.Prefix+":framework_config", data, 0).Err()
	if err != nil {
		logrus.Error("Framework save config state into redis error:", err)
	}
}

// PingRedis to check the health of redis
func (e *Redis) PingRedis() error {
	pong, err := e.Client.Ping(e.CTX).Result()
	if err != nil {
		logrus.WithField("func", "PingRedis").Error("Did not pon Redis: ", pong, err.Error())
	}
	if err != nil {
		return err
	}
	return nil
}

// Connect will connect the redis DB and save the client pointer
func (e *Redis) Connect() bool {
	var redisOptions goredis.Options
	redisOptions.Addr = e.Server
	redisOptions.DB = e.DB
	if e.Password != "" {
		redisOptions.Password = e.Password
	}

	e.Client = goredis.NewClient(&redisOptions)

	err := e.PingRedis()
	if err != nil {
		e.Client.Close()
		e.Connect()
	}

	return true
}

// SaveTaskRedis store mesos task in DB
func (e *Redis) SaveTaskRedis(cmd mesosutil.Command) {
	d, _ := json.Marshal(&cmd)
	e.SetRedisKey(d, cmd.TaskName+":"+cmd.TaskID)
}

// SaveFrameworkRedis store mesos framework in DB
func (e *Redis) SaveFrameworkRedis(framework mesosutil.FrameworkConfig) {
	d, _ := json.Marshal(&framework)
	e.SetRedisKey(d, e.Prefix+":framework")
}

// CheckIfNotTask check if the redis key is a mesos task
func (e *Redis) CheckIfNotTask(keys *goredis.ScanIterator) bool {
	if keys.Val() == e.Prefix+":framework" || keys.Val() == e.Prefix+":framework_config" {
		return true
	}
	return false
}
