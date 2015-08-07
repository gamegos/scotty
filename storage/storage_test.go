package storage

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"testing"
)

var appID = "fritestapp"
var channelID = "someRandomchannelID"
var subscriberIDs = []string{"sub_foo", "sub_bar"}

var stg Storage

func TestMain(m *testing.M) {
	stg = NewMemStorage()
	os.Exit(m.Run())
}

func TestCreateApp(t *testing.T) {

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

	_, err := stg.GetApp(appID)

	if err != nil {
		t.Error(err)
	}
}

func TestAddChannel(t *testing.T) {

	err := stg.AddChannel(appID, channelID)

	if err != nil {
		t.Error(err)
	}
}

func TestAddSubscriber(t *testing.T) {

	err := stg.AddSubscriber(appID, channelID, subscriberIDs)

	if err != nil {
		t.Error(err)
	}
}

func TestGetChannelSubscribers(t *testing.T) {

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

	for _, subscriberID := range subscriberIDs {
		err := stg.UpdateDeviceToken(appID, subscriberID, "footoken", "bartoken")

		if err != nil {
			t.Error(err)
		}
	}
}

func TestGetSubscriberDevices(t *testing.T) {

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

	if !reflect.DeepEqual(expectedDevice, *device) {
		t.Error("Device data does not match.")
	}
}

func TestDeleteChannel(t *testing.T) {

	err := stg.DeleteChannel(appID, channelID)

	if err != nil {
		t.Error(err)
	}
}
