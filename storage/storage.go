package storage

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	// "errors"

	"github.com/garyburd/redigo/redis"
)

// Storage records and retrieves data from data storage.
type Storage struct {
	pool *redis.Pool
}

// AddSubscriber adds new subscriber to channel.
func (stg *Storage) AddSubscriber(appID string, channelID string, subscriberIDs []string) {
	conn := stg.pool.Get()
	defer conn.Close()

	stg.AddChannel(appID, channelID)
	subscribersKey := addPrefix("apps." + appID + ".channels." + channelID + ".subscribers")
	tmpParams := append([]string{subscribersKey}, subscriberIDs...)

	params := make([]interface{}, len(tmpParams))
	for i, v := range tmpParams {
		params[i] = v
	}

	conn.Do("SADD", params...)
}

// AddChannel adds new channel to app.
func (stg *Storage) AddChannel(appID string, channelID string) {
	conn := stg.pool.Get()
	defer conn.Close()

	channelsKey := addPrefix("apps." + appID + ".channels")
	conn.Do("SADD", channelsKey, channelID)
}

// DeleteChannel deletes channel and its subscribers from app.
func (stg *Storage) DeleteChannel(appID string, channelID string) {
	conn := stg.pool.Get()
	defer conn.Close()

	channelsKey := addPrefix("apps." + appID + ".channels")
	conn.Do("SREM", channelsKey, channelID)

	channelKey := channelsKey + "." + channelID + ".subscribers"
	conn.Do("DEL", channelKey)
}

// AddSubscriberDevice adds new device to subscriber.
func (stg *Storage) AddSubscriberDevice(appID string, subscriberID string, device *Device) error {
	conn := stg.pool.Get()
	defer conn.Close()

	subscribersKey := addPrefix("apps." + appID + ".subscribers")
	conn.Do("SADD", subscribersKey, subscriberID)

	devicesKey := subscribersKey + "." + subscriberID + ".devices"
	jstring, _ := json.Marshal(device)
	// todo: multiple devices with same platform and token should not be added
	_, err := conn.Do("HSET", devicesKey, device.Token, jstring)

	if err != nil {
		return err
	}
	return nil
}

// UpdateDeviceToken updates token of a subscriber's device.
func (stg *Storage) UpdateDeviceToken(appID string, subscriberID string, oldDeviceToken string, newDeviceToken string) error {
	conn := stg.pool.Get()
	defer conn.Close()

	key := addPrefix("apps." + appID + ".subscribers." + subscriberID + ".devices")
	deviceData, err := redis.String(conn.Do("HGET", key, oldDeviceToken))
	if err != nil {
		return err
	}

	conn.Do("HDEL", key, oldDeviceToken)

	var device Device
	decoder := json.NewDecoder(strings.NewReader(deviceData))

	if err := decoder.Decode(&device); err != nil {
		return err
	}

	device.Token = newDeviceToken
	err = stg.AddSubscriberDevice(appID, subscriberID, &device)

	if err != nil {
		return err
	}

	return nil
}

// GetChannelSubscribers gets subscribers of a channel.
func (stg *Storage) GetChannelSubscribers(appID string, channelID string) ([]string, error) {
	conn := stg.pool.Get()
	defer conn.Close()

	key := addPrefix("apps." + appID + ".channels." + channelID + ".subscribers")
	subscribers, err := redis.Strings(conn.Do("SMEMBERS", key))

	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

// GetSubscriberDevices gets devices of a subscriber.
func (stg *Storage) GetSubscriberDevices(appID string, subscriberID string) ([]Device, error) {
	conn := stg.pool.Get()
	defer conn.Close()

	key := addPrefix("apps." + appID + ".subscribers." + subscriberID + ".devices")

	var devices map[string]string
	devices, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	var device Device
	var response []Device
	for _, deviceData := range devices {
		decoder := json.NewDecoder(strings.NewReader(deviceData))
		decoder.Decode(&device)
		response = append(response, device)
	}

	return response, nil
}

// AppExists tells whether an app exists or not.
func (stg *Storage) AppExists(appID string) bool {
	conn := stg.pool.Get()
	defer conn.Close()

	status, err := redis.Int(conn.Do("SISMEMBER", addPrefix("apps"), appID))

	if status == 0 || err != nil {
		return false
	}

	return true
}

// CreateApp creates a new app.
func (stg *Storage) CreateApp(appID string, appData string) {
	conn := stg.pool.Get()
	defer conn.Close()

	appsKey := addPrefix("apps")
	conn.Do("SADD", appsKey, appID)
	appKey := appsKey + "." + appID
	conn.Do("SET", appKey, appData)
}

// GetApp gets an app's data.
func (stg *Storage) GetApp(appID string) (string, error) {
	conn := stg.pool.Get()
	defer conn.Close()

	appKey := addPrefix("apps") + "." + appID
	value, err := redis.String(conn.Do("GET", appKey))
	if err != nil {
		return "", err
	}
	return value, nil
}

func addPrefix(key string) string {
	return fmt.Sprintf("srv.push.%s", key)
}

// Init initializes storage with the given config.
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
