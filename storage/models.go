package storage

type App struct {
	Id        string   `json:"id"`
	Platforms Platform `json:"platforms"`
}

type Platform struct {
	Apns Apns `json:"apns"`
	Gcm  Gcm  `json:"gcm"`
}

type Apns struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

type Gcm struct {
	ApiKey    string `json:"apiKey"`
	ProjectId string `json:"projectId"`
}

type Device struct {
	Platform  string
	Token     string
	CreatedAt int
}

type AddDeviceRequest struct {
	SubscriberId string `json:"subscriberId"`
	Platform     string `json:"platform"`
	Token        string `json:"token"`
}

type AddSubscriberRequest struct {
	SubscriberIds []string `json:"subscribers"`
}
