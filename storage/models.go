package storage

// App holds app data.
type App struct {
	ID        string   `json:"id"`
	Platforms Platform `json:"platforms"`
}

// Platform holds platform data.
type Platform struct {
	Apns Apns `json:"apns"`
	Gcm  Gcm  `json:"gcm"`
}

// Apns holds APNS(Apple Push Notification Service) data.
type Apns struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

// Gcm holds GCM(Google Cloud Messaging) data.
type Gcm struct {
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
