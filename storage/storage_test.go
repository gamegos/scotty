package storage

import (
	"encoding/json"
	"flag"
	"reflect"
	"sort"
	"testing"
)

var appID = "fritestapp"
var channelID = "someRandomchannelID"
var subscriberIDs = []string{"sub_foo", "sub_bar"}
var confFile = flag.String("config", "", "Config file")

var stg *RedisStorage

func setup() {

	var conf = InitConfig(*confFile)
	stg = Init(&conf.Redis)
}

func TestCreateApp(t *testing.T) {
	setup()

	app := App{
		ID: appID,
		Platforms: Platform{
			Apns: Apns{
				Certificate: "apnscertificate",
				PrivateKey:  "privatekey",
			},
			Gcm: Gcm{
				APIKey:    "apikey",
				ProjectID: "projectid",
			},
		},
	}

	str, _ := json.Marshal(app)
	err := stg.CreateApp(appID, string(str))

	if err != nil {
		t.Error(err)
	}
}

func TestGetApp(t *testing.T) {
	setup()

	_, err := stg.GetApp(appID)

	if err != nil {
		t.Error(err)
	}
}

func TestAddChannel(t *testing.T) {
	setup()

	err := stg.AddChannel(appID, channelID)

	if err != nil {
		t.Error(err)
	}
}

func TestAddSubscriber(t *testing.T) {
	setup()

	err := stg.AddSubscriber(appID, channelID, subscriberIDs)

	if err != nil {
		t.Error(err)
	}
}

func TestGetChannelSubscribers(t *testing.T) {
	setup()

	receivedSubscribers, err := stg.GetChannelSubscribers(appID, channelID)

	if err != nil {
		t.Error(err)
	}

	sort.Strings(receivedSubscribers)
	sort.Strings(subscriberIDs)

	if !reflect.DeepEqual(receivedSubscribers, subscriberIDs) {
		t.Error("Subscriber IDs does not match.")
	}
}

func TestAddSubscriberDevice(t *testing.T) {
	setup()

	device := Device{
		Platform:  "apns",
		Token:     "footoken",
		CreatedAt: 1436546411,
	}

	for _, subscriberID := range subscriberIDs {
		err := stg.AddSubscriberDevice(appID, subscriberID, &device)

		if err != nil {
			t.Error(err)
		}
	}
}

func TestUpdateDeviceToken(t *testing.T) {
	setup()

	for _, subscriberID := range subscriberIDs {
		err := stg.UpdateDeviceToken(appID, subscriberID, "footoken", "bartoken")

		if err != nil {
			t.Error(err)
		}
	}
}

func TestGetSubscriberDevices(t *testing.T) {
	setup()

	expectedDevice := Device{
		Platform:  "apns",
		Token:     "bartoken",
		CreatedAt: 1436546411,
	}

	devices, err := stg.GetSubscriberDevices(appID, subscriberIDs[0])

	if err != nil {
		t.Error(err)
	}

	device := devices[0]
	if !reflect.DeepEqual(expectedDevice, device) {
		t.Error("Device data does not match.")
	}
}

func TestDeleteChannel(t *testing.T) {
	setup()

	err := stg.DeleteChannel(appID, channelID)

	if err != nil {
		t.Error(err)
	}
}
