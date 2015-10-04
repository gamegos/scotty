package redis

import (
	"encoding/json"
	"strings"

	"github.com/gamegos/scotty/storage"
	redigo "github.com/garyburd/redigo/redis"
)

// RedisStorage records and retrieves data from Redis storage.
type RedisStorage struct {
	pool *redigo.Pool
}

// AddSubscriber adds new subscriber to channel.
func (stg *RedisStorage) AddSubscriber(appID string, channelID string, subscriberIDs []string) error {
	conn := stg.pool.Get()
	defer conn.Close()

	err := stg.AddChannel(appID, channelID)

	if err != nil {
		return err
	}

	subscribersKey := keyChannelSubscribers(appID, channelID)
	tmpParams := append([]string{subscribersKey}, subscriberIDs...)

	params := make([]interface{}, len(tmpParams))
	for i, v := range tmpParams {
		params[i] = v
	}

	_, err = conn.Do("SADD", params...)

	if err != nil {
		return err
	}

	return nil
}

// AddChannel adds new channel to app.
func (stg *RedisStorage) AddChannel(appID string, channelID string) error {
	conn := stg.pool.Get()
	defer conn.Close()

	channelsKey := keyAppChannels(appID)
	_, err := conn.Do("SADD", channelsKey, channelID)

	if err != nil {
		return err
	}

	return nil
}

// DeleteChannel deletes channel and its subscribers from app.
func (stg *RedisStorage) DeleteChannel(appID string, channelID string) error {
	conn := stg.pool.Get()
	defer conn.Close()

	channelsKey := keyAppChannels(appID)
	_, err := conn.Do("SREM", channelsKey, channelID)

	if err != nil {
		return err
	}

	channelKey := keyChannelSubscribers(appID, channelID)
	_, err = conn.Do("DEL", channelKey)

	if err != nil {
		return err
	}

	return nil
}

// AddSubscriberDevice adds new device to subscriber.
func (stg *RedisStorage) AddSubscriberDevice(appID string, subscriberID string, device *storage.Device) error {
	conn := stg.pool.Get()
	defer conn.Close()

	subscribersKey := keyAppSubscribers(appID)
	_, err := conn.Do("SADD", subscribersKey, subscriberID)

	if err != nil {
		return err
	}

	devicesKey := keySubscriberDevices(appID, subscriberID)
	jstring, _ := json.Marshal(device)
	// todo: multiple devices with same platform and token should not be added
	_, err = conn.Do("HSET", devicesKey, device.Token, jstring)

	if err != nil {
		return err
	}

	return nil
}

// UpdateDeviceToken updates token of a subscriber's device.
func (stg *RedisStorage) UpdateDeviceToken(appID string, subscriberID string, oldDeviceToken string, newDeviceToken string) error {
	conn := stg.pool.Get()
	defer conn.Close()

	key := keySubscriberDevices(appID, subscriberID)
	deviceData, err := redigo.String(conn.Do("HGET", key, oldDeviceToken))

	if err != nil {
		return err
	}

	_, err = conn.Do("HDEL", key, oldDeviceToken)

	if err != nil {
		return err
	}

	var device storage.Device
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
func (stg *RedisStorage) GetChannelSubscribers(appID string, channelID string) ([]string, error) {
	conn := stg.pool.Get()
	defer conn.Close()

	key := keyChannelSubscribers(appID, channelID)
	subscribers, err := redigo.Strings(conn.Do("SMEMBERS", key))

	if err != nil {
		return nil, err
	}

	return subscribers, nil
}

// GetSubscriberDevices gets devices of a subscriber.
func (stg *RedisStorage) GetSubscriberDevices(appID string, subscriberID string) ([]*storage.Device, error) {
	conn := stg.pool.Get()
	defer conn.Close()

	key := keySubscriberDevices(appID, subscriberID)

	var devices map[string]string
	devices, err := redigo.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	var device storage.Device
	var response []*storage.Device
	for _, deviceData := range devices {
		decoder := json.NewDecoder(strings.NewReader(deviceData))
		decoder.Decode(&device)
		response = append(response, &device)
	}

	return response, nil
}

// PutApp creates a new app or updates existing one.
func (stg *RedisStorage) PutApp(app *storage.App) error {
	conn := stg.pool.Get()
	defer conn.Close()

	appID := app.ID
	appData, err := json.Marshal(app)
	if err != nil {
		return err
	}

	if _, err := conn.Do("HSET", keyApps(), appID, appData); err != nil {
		return err
	}

	return nil
}

// GetApp gets an app's data.
func (stg *RedisStorage) GetApp(appID string) (*storage.App, error) {
	conn := stg.pool.Get()
	defer conn.Close()

	value, err := redigo.Bytes(conn.Do("HGET", keyApps(), appID))
	if err != nil {
		return nil, err
	}

	var app *storage.App
	if err := json.Unmarshal(value, &app); err != nil {
		return nil, err
	}

	return app, nil
}

const redisPrefix = "scotty"

func buildKey(part ...string) string {
	return redisPrefix + ":" + strings.Join(part, ".")
}

func keyApps() string {
	return buildKey("apps")
}

func keyAppSubscribers(appID string) string {
	return buildKey("apps", appID, "subs")
}

func keyAppChannels(appID string) string {
	return buildKey("apps", appID, "chans")
}

func keyChannelSubscribers(appID, channelID string) string {
	return buildKey("apps", appID, "chans", channelID, "subs")
}

func keySubscriberDevices(appID, subscriberID string) string {
	return buildKey("apps", appID, "subs", subscriberID, "devs")
}
