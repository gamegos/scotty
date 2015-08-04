package storage

type Storage interface {
	// AddSubscriber adds new subscriber to channel.
	AddSubscriber(appID string, channelID string, subscriberIDs []string) error

	// AddChannel adds new channel to app.
	AddChannel(appID string, channelID string) error

	// DeleteChannel deletes channel and its subscribers from app.
	DeleteChannel(appID string, channelID string) error

	// AddSubscriberDevice adds new device to subscriber.
	AddSubscriberDevice(appID string, subscriberID string, device *Device) error

	// UpdateDeviceToken updates token of a subscriber's device.
	UpdateDeviceToken(appID string, subscriberID string, oldDeviceToken string, newDeviceToken string) error

	// GetChannelSubscribers gets subscribers of a channel.
	GetChannelSubscribers(appID string, channelID string) ([]string, error)

	// GetSubscriberDevices gets devices of a subscriber.
	GetSubscriberDevices(appID string, subscriberID string) ([]Device, error)

	// AppExists tells whether an app exists or not.
	AppExists(appID string) bool

	// CreateApp creates a new app.
	CreateApp(appID string, appData string) error

	// GetApp gets an app's data.
	GetApp(appID string) (*App, error)
}
