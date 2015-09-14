package storage

type Storage interface {
	// App methods

	// AppExists tells whether an app exists or not.
	AppExists(appID string) bool

	// PutApp creates a new app or updates existing one.
	PutApp(app *App) error

	// GetApp gets an app's data.
	GetApp(appID string) (*App, error)

	// Subscriber methods

	// AddSubscriberDevice adds new device to subscriber.
	AddSubscriberDevice(appID string, subscriberID string, device *Device) error

	// UpdateDeviceToken updates token of a subscriber's device.
	UpdateDeviceToken(appID string, subscriberID string, oldDeviceToken string, newDeviceToken string) error

	// GetSubscriberDevices gets devices of a subscriber.
	GetSubscriberDevices(appID string, subscriberID string) ([]*Device, error)

	// Channel methods

	// AddSubscriber adds new subscriber to channel.
	AddSubscriber(appID string, channelID string, subscriberIDs []string) error

	// AddChannel adds new channel to app.
	AddChannel(appID string, channelID string) error

	// DeleteChannel deletes channel and its subscribers from app.
	DeleteChannel(appID string, channelID string) error

	// GetChannelSubscribers gets subscribers of a channel.
	GetChannelSubscribers(appID string, channelID string) ([]string, error)
}
