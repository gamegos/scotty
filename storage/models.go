package storage

import "github.com/gamegos/gcmlib"

// App holds app data.
type App struct {
	ID  string    `json:"id"`
	GCM GCMConfig `json:"gcm"`
}

// GCMConfig holds GCM(Google Cloud Messaging) data.
type GCMConfig struct {
	APIKey    string `json:"apiKey"`
	ProjectID string `json:"projectId"`
}

// Device holds device data.
type Device struct {
	Platform  string
	Token     string
	CreatedAt int
}

// AddDeviceRequest holds the structure of new device request.
type AddDeviceRequest struct {
	SubscriberID string `json:"subscriberId"`
	Platform     string `json:"platform"`
	Token        string `json:"token"`
}

// AddSubscriberRequest holds the structure of new subscriber request.
type AddSubscriberRequest struct {
	SubscriberIds []string `json:"subscribers"`
}

// PublishRequest represents http body of "publish" requests.
type PublishRequest struct {
	Subscribers []string `json:"subscribers"`
	Channels    []string `json:"channels"`
	Message     *gcmlib.Message
}
