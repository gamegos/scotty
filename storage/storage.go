package storage

import (
	"time"
	"fmt"
	"encoding/json"
	// "errors"

	"github.com/garyburd/redigo/redis"
)

type Storage struct {
	pool *redis.Pool
}

func (stg *Storage) AddSubscriber(appId string, channelId string, subscriberIds []string) {
	conn := stg.pool.Get()
	// todo: defer close
	stg.AddChannel(appId, channelId)
	subscribersKey := addPrefix("apps."+appId+".channels."+channelId+".subscribers")
	tmpParams := append([]string{subscribersKey}, subscriberIds...)

	params := make([]interface{}, len(tmpParams))
	for i, v := range tmpParams {
		params[i] = v
	}

	conn.Do("SADD", params...)
}

func (stg *Storage) AddChannel(appId string, channelId string) {
	conn := stg.pool.Get()
	channelsKey := addPrefix("apps."+appId+".channels")
	conn.Do("SADD", channelsKey, channelId)
}

func (stg *Storage) DeleteChannel(appId string, channelId string) {
	conn := stg.pool.Get()
	channelsKey := addPrefix("apps."+appId+".channels")
	conn.Do("SREM", channelsKey, channelId)

	channelKey := channelsKey+"."+channelId+".subscribers"
	conn.Do("DEL", channelKey)
}

func (stg *Storage) AddSubscriberDevice(appId string, subscriberId string, device *Device) {
	conn := stg.pool.Get()
	subscribersKey := addPrefix("apps."+appId+".subscribers")
	conn.Do("SADD", subscribersKey, subscriberId)

	devicesKey := subscribersKey+"."+subscriberId+".devices"
	jstring, _ := json.Marshal(device)
	// todo: multiple devices with same platform and token should not be added
	conn.Do("SADD", devicesKey, jstring)
}

func (stg *Storage) AppExists(appId string) bool {
	conn := stg.pool.Get()
	status, err := redis.Int(conn.Do("SISMEMBER", addPrefix("apps"), appId))
	_ = err
	if status == 0 {
		return false
	} else {
		return true
	}
}

func (stg *Storage) CreateApp(appId string, appData string) {
	conn := stg.pool.Get()
	appsKey := addPrefix("apps")
	conn.Do("SADD", appsKey, appId)
	appKey := appsKey + "." + appId
	conn.Do("SET", appKey, appData)
}

func (stg *Storage) GetApp(appId string) (string,error) {
	conn := stg.pool.Get()
	appKey := addPrefix("apps") + "." + appId
	value, err := redis.String(conn.Do("GET", appKey))
	if err != nil {
		return "",err
	}
	return value,nil
}

func addPrefix(key string) string {
	return fmt.Sprintf("srv.push.%s", key)
}

func Init(conf *RedisConfig) *Storage {

	stg := new(Storage)
	pool := &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Wait:        conf.Wait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(conf.Network, conf.Addr)
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	stg.pool = pool
	return stg
}
