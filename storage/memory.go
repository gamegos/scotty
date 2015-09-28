package storage

import "errors"

// MemStorage records and retrieves data from memory.
type MemStorage struct {
	// appid+channelid -> []subscribers
	chans map[string][]string
	// appid -> *App
	apps map[string]*App
	// appid+subscriberId -> Device
	devs map[string][]*Device
	// appid -> [subs1, subs2,...]
	subs map[string][]string
}

// NewMemStorage initializes data structes to store data.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		chans: make(map[string][]string),
		apps:  make(map[string]*App),
		devs:  make(map[string][]*Device),
		subs:  make(map[string][]string),
	}
}

// PutApp creates a new app or updates existing one.
func (stg *MemStorage) PutApp(app *App) error {
	stg.apps[app.ID] = app

	return nil
}

// GetApp gets an app's data.
func (stg *MemStorage) GetApp(appID string) (*App, error) {
	app, ok := stg.apps[appID]

	if !ok {
		return nil, errors.New("not exists")
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

	stg.devs[key] = append(stg.devs[key], device)

	return nil
}

// UpdateDeviceToken updates token of a subscriber's device.
func (stg *MemStorage) UpdateDeviceToken(appID string, subscriberID string, oldDeviceToken string, newDeviceToken string) error {

	key := appID + "." + subscriberID
	devices, ok := stg.devs[key]

	if !ok {
		return errors.New("Device not found.")
	}

	for _, device := range devices {
		if device.Token == oldDeviceToken {
			device.Token = newDeviceToken
		}
	}

	return nil
}

// GetChannelSubscribers gets subscribers of a channel.
func (stg *MemStorage) GetChannelSubscribers(appID string, channelID string) ([]string, error) {

	key := appID + "." + channelID
	return stg.chans[key], nil
}

// GetSubscriberDevices gets devices of a subscriber.
func (stg *MemStorage) GetSubscriberDevices(appID string, subscriberID string) ([]*Device, error) {

	key := appID + "." + subscriberID
	device, _ := stg.devs[key]

	return device, nil
}
