package storage

import (
	"encoding/json"
	"errors"
)

type MemStorage struct {
	// appid+channelid -> []subsctibers
	chans map[string][]string
	// appid -> *App
	apps map[string]*App
	// appid+subscriberId -> Device
	devs map[string]*Device
	// appid -> [subs1, subs2,...]
	subs map[string][]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		chans: make(map[string][]string),
		apps:  make(map[string]*App),
		devs:  make(map[string]*Device),
		subs:  make(map[string][]string),
	}
}

// AppExists tells whether an app exists or not.
func (stg *MemStorage) AppExists(appID string) bool {
	_, ok := stg.apps[appID]

	return ok
}

// CreateApp creates a new app.
func (stg *MemStorage) CreateApp(appID string, appData string) error {

	var app *App
	err := json.Unmarshal([]byte(appData), &app)
	if err != nil {
		stg.apps[appID] = app
	}

	return err
}

// GetApp gets an app's data.
func (stg *MemStorage) GetApp(appID string) (*App, error) {
	app, ok := stg.apps[appID]

	if !ok {
		return nil, errors.New("not exits")
	}

	return app, nil
}

// AddSubscriber adds new subscriber to channel.
func (stg *MemStorage) AddSubscriber(appID string, channelID string, subscriberIDs []string) error {
	key := appID + "." + channelID
	_, ok := stg.chans[key]
	if !ok {
		stg.chans[key] = []string{}
	}

	stg.chans[key] = append(stg.chans[key], subscriberIDs...)

	return nil
}

// AddChannel adds new channel to app.
func (stg *MemStorage) AddChannel(appID string, channelID string) error {
	key := appID + "." + channelID
	_, ok := stg.chans[key]
	if !ok {
		stg.chans[key] = []string{}
	}

	return nil
}

// DeleteChannel deletes channel and its subscribers from app.
func (stg *MemStorage) DeleteChannel(appID string, channelID string) error {
	key := appID + "." + channelID
	delete(stg.chans, key)

	return nil
}

// AddSubscriberDevice adds new device to subscriber.
func (stg *MemStorage) AddSubscriberDevice(appID string, subscriberID string, device *Device) error {
	key := appID + "." + subscriberID

	stg.devs[key] = device

	return nil
}

// UpdateDeviceToken updates token of a subscriber's device.
func (stg *MemStorage) UpdateDeviceToken(appID string, subscriberID string, oldDeviceToken string, newDeviceToken string) error {

	return nil
}

// GetChannelSubscribers gets subscribers of a channel.
func (stg *MemStorage) GetChannelSubscribers(appID string, channelID string) ([]string, error) {

	return []string{}, nil
}

// GetSubscriberDevices gets devices of a subscriber.
func (stg *MemStorage) GetSubscriberDevices(appID string, subscriberID string) ([]Device, error) {

	return []Device{}, nil
}
