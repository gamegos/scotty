package storage

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
